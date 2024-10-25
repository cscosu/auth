package main

import (
	"context"
	"net/http"
)

// Full list of attributes is here:
// https://webauth.service.ohio-state.edu/~shibboleth/user-attribute-reference.html (https://archive.is/H9bOB)
//
// Note that OSUCyber is only authorized for a subset of these attributes. As of 2024-10-23, this is:
// sn, IDMUID, displayName, eduPersonScopedAffiliation, employeeNumber, givenName, mail, SessionIndex, eduPersonPrincipalName
type OSUAttributes struct {
	// This attribute is `sn` in the OSU Shibboleth user attribute reference.
	Surname string
	// This attribute is `IDM ID` in the OSU Shibboleth user attribute reference.
	IDMUID string
	// This attribute is `displayName` in the OSU Shibboleth user attribute reference.
	DisplayName string
	// This attribute is `eduPersonScopedAffiliation` in the OSU Shibboleth user attribute reference.
	Affiliations []string
	// This attribute is `employeeNumber` in the OSU Shibboleth user attribute reference.
	BuckID string
	// This attribute is `givenName` in the OSU Shibboleth user attribute reference.
	GivenName string
	// This attribute is `mail` in the OSU Shibboleth user attribute reference.
	Email string
	// This attribute is `SessionIndex` in the OSU Shibboleth user attribute reference.
	SessionIndex string
}

type AuthProvider interface {
	attributesFromContext(ctx context.Context) *OSUAttributes
	requireAuth(handler http.Handler) http.Handler
	globalLogout(w http.ResponseWriter, r *http.Request)
	logout(w http.ResponseWriter, r *http.Request)
}

type MockAuthProvider struct{}

func (m MockAuthProvider) attributesFromContext(ctx context.Context) *OSUAttributes {
	return &OSUAttributes{
		GivenName:    "Brutus",
		Surname:      "Buckeye",
		DisplayName:  "Brutus Buckeye",
		BuckID:       "500123456",
		IDMUID:       "IDM123456789",
		Email:        "buckeye.1@osu.edu",
		Affiliations: []string{"member@osu.edu", "student@osu.edu"},
		SessionIndex: "_0123456789abcdef01234566890abcde",
	}
}

func (m MockAuthProvider) requireAuth(handler http.Handler) http.Handler {
	return handler
}

func (m MockAuthProvider) globalLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Location", "https://webauth.service.ohio-state.edu/cgi-bin/logout.cgi")
	w.WriteHeader(http.StatusFound)
}

func (m MockAuthProvider) logout(w http.ResponseWriter, r *http.Request) {
}

func mockAuthProvider() *MockAuthProvider {
	return &MockAuthProvider{}
}
