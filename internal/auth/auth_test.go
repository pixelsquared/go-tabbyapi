package auth

import (
	"net/http"
	"testing"
)

func TestAPIKeyAuthenticator_Apply(t *testing.T) {
	// Create an authenticator with a test API key
	apiKey := "test-api-key"
	auth := NewAPIKeyAuthenticator(apiKey)

	// Create a test request
	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Apply the authenticator
	auth.Apply(req)

	// Check that the header was set correctly
	if got := req.Header.Get("X-API-Key"); got != apiKey {
		t.Errorf("Expected X-API-Key header to be %q, got %q", apiKey, got)
	}
}

func TestAdminKeyAuthenticator_Apply(t *testing.T) {
	// Create an authenticator with a test admin key
	adminKey := "test-admin-key"
	auth := NewAdminKeyAuthenticator(adminKey)

	// Create a test request
	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Apply the authenticator
	auth.Apply(req)

	// Check that the header was set correctly
	if got := req.Header.Get("X-Admin-Key"); got != adminKey {
		t.Errorf("Expected X-Admin-Key header to be %q, got %q", adminKey, got)
	}
}

func TestBearerTokenAuthenticator_Apply(t *testing.T) {
	// Create an authenticator with a test bearer token
	token := "test-token"
	auth := NewBearerTokenAuthenticator(token)

	// Create a test request
	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Apply the authenticator
	auth.Apply(req)

	// Check that the header was set correctly
	expectedAuth := "Bearer " + token
	if got := req.Header.Get("Authorization"); got != expectedAuth {
		t.Errorf("Expected Authorization header to be %q, got %q", expectedAuth, got)
	}
}
