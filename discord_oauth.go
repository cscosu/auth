package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func (r *Router) DiscordSignin(w http.ResponseWriter, req *http.Request) {
	state := generateStateOauthCookie(w)
	redirectUri := "http://localhost:3000/discord/callback"
	url := fmt.Sprintf("https://discord.com/oauth2/authorize?client_id=1172564200007155903&response_type=code&redirect_uri=%v&scope=identify&state=%v", url.QueryEscape(redirectUri), state)
	http.Redirect(w, req, url, http.StatusTemporaryRedirect)
}

func (r *Router) DiscordCallback(w http.ResponseWriter, req *http.Request) {
	stateCookie, err := req.Cookie("oauthstate")
	if err != nil {
		http.Error(w, "Missing oauthstate cookie", http.StatusBadRequest)
		return
	}
	stateParam := req.URL.Query().Get("state")
	if stateParam == "" {
		http.Error(w, "Missing state url parameter", http.StatusBadRequest)
		return
	}
	if stateCookie.Value != stateParam {
		http.Error(w, "State cookie and state parameter don't match", http.StatusBadRequest)
		return
	}

	code := req.URL.Query().Get("code")
	authToken, err := getDiscordAuthToken(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
	}

	discordUser, err := getDiscordUser(authToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
	}

	fmt.Fprintf(w, "You are %s, id = %s", discordUser.Username, discordUser.ID)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(2 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)
	return state
}

func getDiscordAuthToken(code string) (string, error) {
	data := url.Values{
		"client_id":     {"1172564200007155903"},
		"client_secret": {"wOp5aVQOsnO0I3_Nt6a3hQyVz-LTqLaz"},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {"http://localhost:3000/discord/callback"},
	}

	req, err := http.NewRequest("POST", "https://discord.com/api/oauth2/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("discord token endpoint responded with %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("discord token endpoint response did not have access_token")
	}

	return accessToken, nil
}

type DiscordUser struct {
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
	Email         string `json:"email"`
	Flags         int    `json:"flags"`
	ID            string `json:"id"`
	Username      string `json:"username"`
}

func getDiscordUser(authToken string) (DiscordUser, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	resp, err := client.Do(req)
	if err != nil {
		return DiscordUser{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return DiscordUser{}, fmt.Errorf("failed to get user info")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return DiscordUser{}, err
	}

	var user DiscordUser
	if err := json.Unmarshal(body, &user); err != nil {
		return DiscordUser{}, err
	}

	return user, nil
}
