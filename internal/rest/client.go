// Package rest provides a REST client for making HTTP requests to the Tabby API.
package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pixelsquared/go-tabbyapi/internal/auth"
	"github.com/pixelsquared/go-tabbyapi/internal/errors"
)

// ClientOption is a function that configures a Client.
type ClientOption func(*Client)

// Client is a REST client for making HTTP requests to the Tabby API.
type Client struct {
	baseURL     string
	httpClient  *http.Client
	auth        auth.Authenticator
	contentType string
}

// New creates a new REST client.
func New(baseURL string, options ...ClientOption) *Client {
	client := &Client{
		baseURL:     strings.TrimRight(baseURL, "/"),
		httpClient:  http.DefaultClient,
		contentType: "application/json",
	}

	for _, option := range options {
		option(client)
	}

	return client
}

// WithHTTPClient sets the HTTP client for the REST client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithAuth sets the authenticator for the REST client.
func WithAuth(auth auth.Authenticator) ClientOption {
	return func(c *Client) {
		c.auth = auth
	}
}

// WithContentType sets the content type for requests.
func WithContentType(contentType string) ClientOption {
	return func(c *Client) {
		c.contentType = contentType
	}
}

// Get sends a GET request to the specified endpoint.
func (c *Client) Get(ctx context.Context, endpoint string, params url.Values, result interface{}) error {
	url := c.buildURL(endpoint, params)
	return c.Do(ctx, http.MethodGet, url, nil, result)
}

// Post sends a POST request to the specified endpoint.
func (c *Client) Post(ctx context.Context, endpoint string, body, result interface{}) error {
	url := c.buildURL(endpoint, nil)
	return c.Do(ctx, http.MethodPost, url, body, result)
}

// Put sends a PUT request to the specified endpoint.
func (c *Client) Put(ctx context.Context, endpoint string, body, result interface{}) error {
	url := c.buildURL(endpoint, nil)
	return c.Do(ctx, http.MethodPut, url, body, result)
}

// Delete sends a DELETE request to the specified endpoint.
func (c *Client) Delete(ctx context.Context, endpoint string, params url.Values) error {
	url := c.buildURL(endpoint, params)
	return c.Do(ctx, http.MethodDelete, url, nil, nil)
}

// Do sends an HTTP request and returns the response.
func (c *Client) Do(ctx context.Context, method, url string, body, result interface{}) error {
	req, err := c.createRequest(ctx, method, url, body)
	if err != nil {
		return &errors.RequestError{
			Message: "failed to create request",
			Err:     err,
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &errors.RequestError{
			Message: "failed to execute request",
			Err:     err,
		}
	}
	defer resp.Body.Close()

	return c.handleResponse(resp, result)
}

// DoRaw sends an HTTP request and returns the raw response for streaming.
func (c *Client) DoRaw(ctx context.Context, method, url string, body interface{}) (*http.Response, error) {
	req, err := c.createRequest(ctx, method, url, body)
	if err != nil {
		return nil, &errors.RequestError{
			Message: "failed to create request",
			Err:     err,
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &errors.RequestError{
			Message: "failed to execute request",
			Err:     err,
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		return nil, c.parseErrorResponse(resp)
	}

	return resp, nil
}

// createRequest creates a new HTTP request.
func (c *Client) createRequest(ctx context.Context, method, url string, body interface{}) (*http.Request, error) {
	var bodyReader io.Reader

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", c.contentType)
	}
	req.Header.Set("Accept", "application/json")

	if c.auth != nil {
		c.auth.Apply(req)
	}

	return req, nil
}

// handleResponse processes the HTTP response.
func (c *Client) handleResponse(resp *http.Response, result interface{}) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return c.parseErrorResponse(resp)
	}

	if result == nil || resp.StatusCode == http.StatusNoContent {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &errors.RequestError{
			Message:    "failed to read response body",
			StatusCode: resp.StatusCode,
			Err:        err,
		}
	}

	if len(body) == 0 {
		return nil
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return &errors.RequestError{
			Message:    "failed to unmarshal response body",
			StatusCode: resp.StatusCode,
			Err:        err,
		}
	}

	return nil
}

// parseErrorResponse extracts error information from an error response.
func (c *Client) parseErrorResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &errors.APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("failed to read error response body: %v", err),
		}
	}

	if len(body) == 0 {
		return &errors.APIError{
			StatusCode: resp.StatusCode,
			Message:    http.StatusText(resp.StatusCode),
		}
	}

	var apiError errors.APIError
	if err := json.Unmarshal(body, &apiError); err != nil {
		return &errors.APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	apiError.StatusCode = resp.StatusCode
	return &apiError
}

// buildURL constructs the full URL for a request.
// BuildURL constructs the full URL for a request.
func (c *Client) BuildURL(endpoint string, params url.Values) string {
	endpoint = strings.TrimLeft(endpoint, "/")
	fullURL := fmt.Sprintf("%s/%s", c.baseURL, endpoint)

	if params != nil && len(params) > 0 {
		fullURL = fmt.Sprintf("%s?%s", fullURL, params.Encode())
	}

	return fullURL
}

// Private method to maintain backward compatibility
func (c *Client) buildURL(endpoint string, params url.Values) string {
	return c.BuildURL(endpoint, params)
}
