// Package auth provides authentication mechanisms for the Tabby API client.
package auth

import (
	"net/http"
)

// Authenticator provides authentication for API requests.
type Authenticator interface {
	// Apply adds authentication to the provided request.
	Apply(req *http.Request)
}

// NoAuth is an authenticator that doesn't add any authentication.
var NoAuth Authenticator = &noAuth{}

// noAuth implements Authenticator but doesn't apply any authentication.
type noAuth struct{}

// Apply implements the Authenticator interface but does nothing.
func (a *noAuth) Apply(req *http.Request) {}

// APIKeyAuthenticator uses the X-API-Key header for authentication.
type APIKeyAuthenticator struct {
	Key string
}

// Apply implements the Authenticator interface.
func (a *APIKeyAuthenticator) Apply(req *http.Request) {
	req.Header.Set("X-API-Key", a.Key)
}

// AdminKeyAuthenticator uses the X-Admin-Key header for authentication.
type AdminKeyAuthenticator struct {
	Key string
}

// Apply implements the Authenticator interface.
func (a *AdminKeyAuthenticator) Apply(req *http.Request) {
	req.Header.Set("X-Admin-Key", a.Key)
}

// BearerTokenAuthenticator uses the Authorization header with Bearer token.
type BearerTokenAuthenticator struct {
	Token string
}

// Apply implements the Authenticator interface.
func (a *BearerTokenAuthenticator) Apply(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+a.Token)
}

// NewAPIKeyAuthenticator creates a new APIKeyAuthenticator.
func NewAPIKeyAuthenticator(key string) Authenticator {
	return &APIKeyAuthenticator{Key: key}
}

// NewAdminKeyAuthenticator creates a new AdminKeyAuthenticator.
func NewAdminKeyAuthenticator(key string) Authenticator {
	return &AdminKeyAuthenticator{Key: key}
}

// NewBearerTokenAuthenticator creates a new BearerTokenAuthenticator.
func NewBearerTokenAuthenticator(token string) Authenticator {
	return &BearerTokenAuthenticator{Token: token}
}
