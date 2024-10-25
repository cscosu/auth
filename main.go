package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/k0kubun/pp/v3"
	_ "modernc.org/sqlite"
)

type Router struct {
	db           *sql.DB
	authProvider AuthProvider
	jwtSecret    []byte
}

const AUTH_COOKIE_NAME string = "csc-auth"

func (r *Router) signin(w http.ResponseWriter, req *http.Request) {
	attributes := r.authProvider.attributesFromContext(req.Context())

	pp.Println(attributes)
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"idm_id": attributes.IDMUID,
		"iat":    now.Unix(),
		"exp":    now.AddDate(1, 0, 0).Unix(),
	})

	signedTokenString, err := token.SignedString(r.jwtSecret)

	if err != nil {
		log.Fatalln("Failed to sign JWT:", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     AUTH_COOKIE_NAME,
		HttpOnly: true,
		Value:    signedTokenString,
		MaxAge:   365 * 24 * 60 * 60, // 1 year
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	nameNum := strings.TrimSuffix(attributes.Email, "@osu.edu")
	nameNum = strings.TrimSuffix(nameNum, "@buckeyemail.osu.edu")

	student := false
	alum := false
	employee := false
	faculty := false

	for _, affiliation := range attributes.Affiliations {
		if affiliation == "student@osu.edu" {
			student = true
		} else if affiliation == "alum@osu.edu" {
			alum = true
		} else if affiliation == "employee@osu.edu" {
			employee = true
		} else if affiliation == "faculty@osu.edu" {
			faculty = true
		}
	}

	r.db.Exec(`
		INSERT OR REPLACE INTO users (idm_id, buck_id, name_num, display_name, student, alum, employee, faculty)
		VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8)
	`, attributes.IDMUID, attributes.BuckID, nameNum, attributes.DisplayName, student, alum, employee, faculty)

	redirect := req.URL.Query().Get("redirect")
	if redirect != "" {
		http.Redirect(w, req, redirect, http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "Hello, %s!", attributes.GivenName)
}

func getUserIDFromContext(ctx context.Context) (string, bool) {
	userId, ok := ctx.Value(CONTEXT_USER_ID_KEY).(string)

	return userId, ok
}

func (r *Router) hello(w http.ResponseWriter, req *http.Request) {
	userId, hasUserId := getUserIDFromContext(req.Context())

	if hasUserId {
		row := r.db.QueryRow("SELECT display_name FROM users WHERE idm_id = ?", userId)
		var displayName string
		row.Scan(&displayName)
		fmt.Fprintf(w, "Hello, %s!", displayName)
	} else {
		fmt.Fprintln(w, "Hello, unknown user!")
	}
}

type contextUserIdType int

const CONTEXT_USER_ID_KEY contextUserIdType = iota

func (r *Router) InjectJwtMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie(AUTH_COOKIE_NAME)
		if err != nil {
			handler.ServeHTTP(w, req)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return r.jwtSecret, nil
		})
		if err != nil {
			log.Println(err)
			http.Redirect(w, req, fmt.Sprintf("/signin?redirect=%v", req.URL.Path), http.StatusTemporaryRedirect)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			log.Println("Invalid token", token)
			http.Redirect(w, req, fmt.Sprintf("/signin?redirect=%v", req.URL.Path), http.StatusTemporaryRedirect)
			return
		}

		idm_id := claims["idm_id"].(string)

		req = req.WithContext(context.WithValue(req.Context(), CONTEXT_USER_ID_KEY, idm_id))
		handler.ServeHTTP(w, req)
	})
}

func (r *Router) EnforceJwtMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_, hasUserId := getUserIDFromContext(req.Context())
		if !hasUserId {
			http.Redirect(w, req, fmt.Sprintf("/signin?redirect=%v", req.URL.Path), http.StatusTemporaryRedirect)
			return
		}

		handler.ServeHTTP(w, req)
	})
}

//go:embed migrations/*
var migrations embed.FS

func main() {
	mux := http.NewServeMux()

	db, err := sql.Open("sqlite", "./auth.db")

	if err != nil {
		log.Fatalln("Failed to load the database:", err)
	}

	dirs, err := migrations.ReadDir("migrations")

	if err != nil {
		log.Fatalln("Failed to read migrations directory:", err)
	}

	slices.SortStableFunc(dirs, func(a fs.DirEntry, b fs.DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})

	for _, entry := range dirs {
		data, err := migrations.ReadFile(fmt.Sprintf("migrations/%v", entry.Name()))
		if err != nil {
			log.Fatalln("Failed to read", entry.Name(), err)
		}
		sql := string(data)
		_, err = db.Exec(sql)
		if err != nil {
			log.Fatalln("Failed to run", entry.Name(), err)
		}
	}

	authEnvironment := os.Getenv("ENV")
	var authProvider AuthProvider

	if authEnvironment == "" {
		authProvider = mockAuthProvider()
	} else {
		rootURL, _ := url.Parse("https://auth.osucyber.club")
		if authEnvironment == "saml" {
			rootURL, _ = url.Parse("https://auth-test.osucyber.club")
		}

		keyPair, err := tls.LoadX509KeyPair("keys/sp-cert.pem", "keys/sp-key.pem")
		if err != nil {
			panic(err)
		}

		authProvider, _ = samlAuthProvider(mux, rootURL, &keyPair)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		if authEnvironment != "" && authEnvironment != "saml" {
			log.Fatalln("JWT_SECRET not set")
		}

		log.Println("DEFAULTING JWT_SECRET TO `secret` DO NOT RUN IN PRODUCTION")
		jwtSecret = "secret"
	}

	router := &Router{
		db:           db,
		authProvider: authProvider,
		jwtSecret:    []byte(jwtSecret),
	}

	mux.Handle("/hello", router.InjectJwtMiddleware(router.EnforceJwtMiddleware(http.HandlerFunc(router.hello))))
	// mux.Handle("/hello", router.InjectJwtMiddleware(http.HandlerFunc(router.hello)))
	mux.Handle("/signin", authProvider.requireAuth(http.HandlerFunc(router.signin)))
	mux.Handle("/logout", authProvider.requireAuth(http.HandlerFunc(authProvider.globalLogout)))

	if authEnvironment == "saml" {
		log.Println("Starting server on :443. Visit https://auth-test.osucyber.club and accept the self-signed certificate")
		keyPair, err := getTlsCert()
		if err != nil {
			panic(err)
		}
		server := &http.Server{
			Addr:              ":443",
			ReadHeaderTimeout: time.Second * 10,
			Handler:           mux,
			TLSConfig: &tls.Config{
				MinVersion:   tls.VersionTLS12,
				Certificates: []tls.Certificate{*keyPair},
			},
		}
		_ = server.ListenAndServeTLS("", "")
	} else {
		if authEnvironment == "" {
			log.Println("Starting server on :3000. Visit http://localhost:3000")
		} else {
			log.Println("Starting server on :3000")
		}

		server := &http.Server{
			Addr:              ":3000",
			ReadHeaderTimeout: time.Second * 10,
			Handler:           mux,
		}
		_ = server.ListenAndServe()
	}
}
