package main

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"net/url"

	"github.com/crewjam/saml/samlsp"
)

type SamlAuthProvider struct {
	samlSP *samlsp.Middleware
}

func (s SamlAuthProvider) attributesFromContext(ctx context.Context) *OSUAttributes {
	session := samlsp.SessionFromContext(ctx).(samlsp.SessionWithAttributes)
	attributes := session.GetAttributes()

	return &OSUAttributes{
		Surname:      attributes.Get("sn"),
		IDMUID:       attributes.Get("IDMUID"),
		DisplayName:  attributes.Get("displayName"),
		Affiliations: attributes["eduPersonScopedAffiliation"],
		BuckID:       attributes.Get("employeeNumber"),
		GivenName:    attributes.Get("givenName"),
		Email:        attributes.Get("mail"),
		SessionIndex: attributes.Get("SessionIndex"),
	}
}

func (s SamlAuthProvider) requireAuth(handler http.Handler) http.Handler {
	return s.samlSP.RequireAccount(handler)
}

func (s *SamlAuthProvider) globalLogout(w http.ResponseWriter, r *http.Request) {
	err := s.samlSP.Session.DeleteSession(w, r)
	if err != nil {
		panic(err) // TODO handle error
	}

	w.Header().Add("Location", "https://webauth.service.ohio-state.edu/cgi-bin/logout.cgi")
	w.WriteHeader(http.StatusFound)
}

func samlAuthProvider(mux *http.ServeMux, rootURL *url.URL, keyPair *tls.Certificate) (*SamlAuthProvider, error) {
	var err error

	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		return nil, err
	}

	idpMetadataURL, err := url.Parse("https://webauth.service.ohio-state.edu/OSU-idp-metadata.xml")
	if err != nil {
		return nil, err
	}
	idpMetadata, err := samlsp.FetchMetadata(context.Background(), http.DefaultClient, *idpMetadataURL)
	if err != nil {
		return nil, err
	}

	samlSP, err := samlsp.New(samlsp.Options{
		URL:         *rootURL,
		EntityID:    "https://auth.osucyber.club/shibboleth",
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
		IDPMetadata: idpMetadata,
		CookieName:  "_saml",
	})
	if err != nil {
		return nil, err
	}

	acsUrl := *rootURL
	acsUrl.Path = "/Shibboleth.sso/SAML2/POST"
	samlSP.ServiceProvider.AcsURL = acsUrl
	samlSP.OnError = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Println("SAML error:", err)
	}

	mux.Handle("/Shibboleth.sso/", samlSP)
	mux.Handle("/metadata.xml", http.HandlerFunc(samlSP.ServeMetadata))

	return &SamlAuthProvider{
		samlSP,
	}, nil
}
