// Package tabby provides a Go client for TabbyAPI.
package tabby

import (
	"net/http"
	"time"
)

// Option configures a Client using the functional options pattern.
// Use these option functions with NewClient to customize client behavior.
//
// Example:
//
//	client := tabby.NewClient(
//	    tabby.WithBaseURL("http://your-server:8080"),
//	    tabby.WithAPIKey("your-api-key"),
//	    tabby.WithTimeout(30 * time.Second),
//	)
type Option func(*clientImpl)

// WithBaseURL sets the base URL for the TabbyAPI server.
//
// The URL should include the protocol and host (with optional port),
// but should not include API version or endpoint paths.
//
// Example: "http://localhost:8080" or "https://api.example.com"
func WithBaseURL(url string) Option {
	return func(c *clientImpl) {
		c.baseURL = url
	}
}

// WithHTTPClient sets a custom HTTP client for making API requests.
//
// This can be used to configure advanced HTTP client behaviors such as
// custom transport settings, timeouts, or middleware.
func WithHTTPClient(client *http.Client) Option {
	return func(c *clientImpl) {
		c.httpClient = client
	}
}

// WithAPIKey sets the API key for standard API authentication.
//
// The API key is sent with each request in the X-API-Key header.
// API keys typically provide read and write access but not administrative
// capabilities.
func WithAPIKey(key string) Option {
	return func(c *clientImpl) {
		c.auth = &APIKeyAuthenticator{Key: key}
	}
}

// WithAdminKey sets the admin key for administrative API authentication.
//
// The admin key is sent with each request in the X-Admin-Key header.
// Admin keys provide full access to all API operations, including
// administrative functions like model and LoRA management.
func WithAdminKey(key string) Option {
	return func(c *clientImpl) {
		c.auth = &AdminKeyAuthenticator{Key: key}
	}
}

// WithBearerToken sets a bearer token for OAuth or JWT authentication.
//
// The token is sent with each request in the Authorization header
// with the Bearer scheme.
func WithBearerToken(token string) Option {
	return func(c *clientImpl) {
		c.auth = &BearerTokenAuthenticator{Token: token}
	}
}

// WithTimeout sets the timeout duration for all API requests.
//
// The timeout includes connection time, any redirects, and reading
// the response body. For streaming operations, the timeout applies
// to the initial connection, not the duration of the stream.
func WithTimeout(timeout time.Duration) Option {
	return func(c *clientImpl) {
		c.httpClient.Timeout = timeout
	}
}

// RetryPolicy defines how the client should retry failed requests.
// This interface allows for customizable retry behavior, including
// determining which requests should be retried, how long to wait between
// retries, and how many retries to attempt.
type RetryPolicy interface {
	// ShouldRetry determines if a request should be retried based on
	// the HTTP response and/or error. Typical conditions for retrying
	// include network errors, server errors (5xx status codes), and
	// rate limiting (429 status code).
	//
	// Parameters:
	//   - resp: The HTTP response, may be nil if a response was not received
	//   - err: The error that occurred, may be nil if the request succeeded
	//          but the response indicates an error condition
	//
	// Returns true if the request should be retried, false otherwise.
	ShouldRetry(resp *http.Response, err error) bool

	// RetryDelay returns the duration to wait before the next retry attempt.
	// This can implement various backoff strategies like exponential backoff
	// or fixed intervals.
	//
	// Parameters:
	//   - attempts: The number of retry attempts that have been made so far
	//
	// Returns the duration to wait before the next retry.
	RetryDelay(attempts int) time.Duration

	// MaxRetries returns the maximum number of retry attempts that should be made.
	// After this many retries, the request will fail with the last error received.
	//
	// Returns the maximum number of retry attempts.
	MaxRetries() int
}

// WithRetryPolicy sets the retry policy for the client.
// This allows customizing how and when failed requests are retried.
// For most use cases, DefaultRetryPolicy() provides a reasonable default,
// but custom policies can be implemented for specific requirements.
func WithRetryPolicy(policy RetryPolicy) Option {
	return func(c *clientImpl) {
		c.retryPolicy = policy
	}
}

// SimpleRetryPolicy provides a basic retry policy implementation that can be
// configured with custom functions for determining retry conditions, delays
// between retries, and the maximum number of retry attempts.
//
// This struct provides a flexible way to create custom retry policies
// without implementing the entire RetryPolicy interface.
type SimpleRetryPolicy struct {
	// MaxRetryCount specifies the maximum number of retry attempts
	MaxRetryCount int

	// RetryDelayFunc returns the delay duration for a given retry attempt
	// It receives the attempt number (starting at 1) and returns a duration
	RetryDelayFunc func(attempts int) time.Duration

	// RetryableFunc determines if a request should be retried based on
	// the HTTP response and/or error
	RetryableFunc func(resp *http.Response, err error) bool
}

// ShouldRetry implements the RetryPolicy interface by delegating to the
// RetryableFunc function provided in the SimpleRetryPolicy struct.
func (p *SimpleRetryPolicy) ShouldRetry(resp *http.Response, err error) bool {
	return p.RetryableFunc(resp, err)
}

// RetryDelay implements the RetryPolicy interface by delegating to the
// RetryDelayFunc function provided in the SimpleRetryPolicy struct.
func (p *SimpleRetryPolicy) RetryDelay(attempts int) time.Duration {
	return p.RetryDelayFunc(attempts)
}

// MaxRetries implements the RetryPolicy interface by returning the
// MaxRetryCount value from the SimpleRetryPolicy struct.
func (p *SimpleRetryPolicy) MaxRetries() int {
	return p.MaxRetryCount
}

// DefaultRetryPolicy returns a reasonable default retry policy with the following behavior:
//
// - Maximum of 3 retry attempts
// - Exponential backoff delay with jitter:
//   - 1st retry: ~100ms
//   - 2nd retry: ~1.2s
//   - 3rd retry: ~2.3s
//
// - Retries on the following conditions:
//   - Any network or connection error
//   - Any HTTP status code >= 500 (server errors)
//
// This policy is suitable for most use cases and provides a balance between
// reliability and responsiveness. For more specific requirements, create a
// custom RetryPolicy implementation or use SimpleRetryPolicy with tailored
// parameters.
func DefaultRetryPolicy() RetryPolicy {
	return &SimpleRetryPolicy{
		MaxRetryCount: 3,
		RetryDelayFunc: func(attempts int) time.Duration {
			// Exponential backoff with jitter
			return time.Duration(1<<uint(attempts-1))*time.Second + time.Duration(100*attempts)*time.Millisecond
		},
		RetryableFunc: func(resp *http.Response, err error) bool {
			// Retry on connection errors or 5xx status codes
			if err != nil {
				return true
			}
			return resp.StatusCode >= 500
		},
	}
}
