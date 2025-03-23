package rest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pixelsquared/go-tabbyapi/internal/auth"
	"github.com/pixelsquared/go-tabbyapi/internal/errors"
)

type testResponse struct {
	Message string `json:"message"`
}

func TestClient_Get(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the request method is GET
		if r.Method != http.MethodGet {
			t.Errorf("Expected method %s, got %s", http.MethodGet, r.Method)
		}

		// Check the path
		if r.URL.Path != "/test" {
			t.Errorf("Expected path /test, got %s", r.URL.Path)
		}

		// Check query parameters
		if r.URL.Query().Get("param") != "value" {
			t.Errorf("Expected query param 'param' to be 'value', got %s", r.URL.Query().Get("param"))
		}

		// Check headers
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept header to be application/json, got %s", r.Header.Get("Accept"))
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	// Create client
	client := New(server.URL)

	// Create parameters
	params := url.Values{}
	params.Set("param", "value")

	// Send GET request
	var result testResponse
	err := client.Get(context.Background(), "/test", params, &result)
	if err != nil {
		t.Fatalf("Get returned an error: %v", err)
	}

	// Check response
	if result.Message != "success" {
		t.Errorf("Expected message 'success', got %s", result.Message)
	}
}

func TestClient_Post(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the request method is POST
		if r.Method != http.MethodPost {
			t.Errorf("Expected method %s, got %s", http.MethodPost, r.Method)
		}

		// Check the path
		if r.URL.Path != "/test" {
			t.Errorf("Expected path /test, got %s", r.URL.Path)
		}

		// Check headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header to be application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Check body
		var requestBody map[string]string
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &requestBody)
		if requestBody["key"] != "value" {
			t.Errorf("Expected body key 'key' to be 'value', got %s", requestBody["key"])
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"created"}`))
	}))
	defer server.Close()

	// Create client
	client := New(server.URL)

	// Create request body
	body := map[string]string{"key": "value"}

	// Send POST request
	var result testResponse
	err := client.Post(context.Background(), "/test", body, &result)
	if err != nil {
		t.Fatalf("Post returned an error: %v", err)
	}

	// Check response
	if result.Message != "created" {
		t.Errorf("Expected message 'created', got %s", result.Message)
	}
}

func TestClient_Error(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"invalid request", "status_code":400}`))
	}))
	defer server.Close()

	// Create client
	client := New(server.URL)

	// Send GET request
	var result testResponse
	err := client.Get(context.Background(), "/test", nil, &result)

	// Check that we got the expected error
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}

	// Check that the error is of the expected type
	apiError, ok := err.(*errors.APIError)
	if !ok {
		t.Fatalf("Expected error of type *errors.APIError, got %T", err)
	}

	// Check error details
	if apiError.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, apiError.StatusCode)
	}
	if apiError.Message != "invalid request" {
		t.Errorf("Expected message 'invalid request', got %s", apiError.Message)
	}
}

func TestClient_WithAuth(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the API key header is set
		if r.Header.Get("X-API-Key") != "test-key" {
			t.Errorf("Expected X-API-Key header to be 'test-key', got %s", r.Header.Get("X-API-Key"))
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"authenticated"}`))
	}))
	defer server.Close()

	// Create client with API key auth
	apiKey := "test-key"
	authenticator := auth.NewAPIKeyAuthenticator(apiKey)
	client := New(server.URL, WithAuth(authenticator))

	// Send GET request
	var result testResponse
	err := client.Get(context.Background(), "/test", nil, &result)
	if err != nil {
		t.Fatalf("Get returned an error: %v", err)
	}

	// Check response
	if result.Message != "authenticated" {
		t.Errorf("Expected message 'authenticated', got %s", result.Message)
	}
}

func TestClient_DoRaw(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the request method is GET
		if r.Method != http.MethodGet {
			t.Errorf("Expected method %s, got %s", http.MethodGet, r.Method)
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	// Create client
	client := New(server.URL)

	// Send raw GET request
	resp, err := client.DoRaw(context.Background(), http.MethodGet, server.URL+"/test", nil)
	if err != nil {
		t.Fatalf("DoRaw returned an error: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Read and check response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var result testResponse
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if result.Message != "success" {
		t.Errorf("Expected message 'success', got %s", result.Message)
	}
}
