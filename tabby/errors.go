package tabby

import (
	"fmt"
	"net/http"
)

// Error is the interface implemented by all errors in the TabbyAPI client library.
// It extends the standard error interface with methods to get error codes and HTTP status codes,
// making it easier to programmatically handle different types of errors.
type Error interface {
	// error provides the standard Error() string method
	error

	// Code returns a string identifier for the error type, such as "authentication_error"
	// or "invalid_request", which can be used for programmatic error handling.
	Code() string

	// HTTPStatusCode returns the HTTP status code associated with the error, such as
	// 400 for bad requests or 401 for authentication errors.
	HTTPStatusCode() int
}

// APIError represents an error returned directly by the TabbyAPI server.
// These errors come from the server's error response JSON and may include
// additional details about what went wrong.
type APIError struct {
	// StatusCode is the HTTP status code of the error response
	StatusCode int `json:"status_code"`

	// Message is the human-readable error message
	Message string `json:"message"`

	// Details contains additional error information which may be type-specific
	Details interface{} `json:"details,omitempty"`

	// RequestID uniquely identifies the request for server-side debugging
	RequestID string `json:"request_id,omitempty"`
}

// Error implements the error interface by returning a formatted error message
// that includes the status code, request ID (if available), and error message.
func (e *APIError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("tabby API error (status %d, request ID: %s): %s", e.StatusCode, e.RequestID, e.Message)
	}
	return fmt.Sprintf("tabby API error (status %d): %s", e.StatusCode, e.Message)
}

// Code returns a string code for the error based on its HTTP status code.
// This provides programmatic error classification without parsing error messages,
// making it easier to handle specific error types in application code.
func (e *APIError) Code() string {
	switch e.StatusCode {
	case http.StatusBadRequest:
		return "invalid_request"
	case http.StatusUnauthorized:
		return "authentication_error"
	case http.StatusForbidden:
		return "permission_error"
	case http.StatusNotFound:
		return "not_found"
	case http.StatusTooManyRequests:
		return "rate_limit_exceeded"
	default:
		if e.StatusCode >= 500 {
			return "server_error"
		}
		return "api_error"
	}
}

// HTTPStatusCode returns the HTTP status code of the error.
func (e *APIError) HTTPStatusCode() int {
	return e.StatusCode
}

// ValidationError represents a validation error that occurs when request
// parameters or input values do not meet the required format or constraints.
// These errors typically indicate client-side issues that need to be fixed
// before retrying the request.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
}

// Error implements the error interface by returning a formatted message
// that includes the field name, error message, and optionally the
// validation error type.
func (e *ValidationError) Error() string {
	if e.Type != "" {
		return fmt.Sprintf("validation error in %s: %s (type: %s)", e.Field, e.Message, e.Type)
	}
	return fmt.Sprintf("validation error in %s: %s", e.Field, e.Message)
}

// Code returns the string code "validation_error" to identify validation errors.
func (e *ValidationError) Code() string {
	return "validation_error"
}

// HTTPStatusCode returns http.StatusBadRequest (400) for validation errors,
// as they typically represent client-side input errors.
func (e *ValidationError) HTTPStatusCode() int {
	return http.StatusBadRequest
}

// RequestError represents errors that occur during the HTTP request lifecycle,
// such as connection problems, timeouts, or malformed responses. These errors
// typically occur on the client side before reaching the server or when
// processing the server's response.
type RequestError struct {
	// Message describes the error in human-readable form
	Message string

	// StatusCode is the HTTP status code, if applicable
	StatusCode int

	// Err is the underlying error that caused this RequestError,
	// such as a network error or context cancellation
	Err error
}

// Error implements the error interface by returning a formatted message
// that includes the error description and the underlying error, if any.
func (e *RequestError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("request error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("request error: %s", e.Message)
}

// Code returns the string code "request_error" to identify request-related errors.
func (e *RequestError) Code() string {
	return "request_error"
}

// HTTPStatusCode returns the HTTP status code associated with this error.
// If a specific status code was set, it returns that code.
// Otherwise, it returns http.StatusInternalServerError (500) as a default.
func (e *RequestError) HTTPStatusCode() int {
	if e.StatusCode != 0 {
		return e.StatusCode
	}
	return http.StatusInternalServerError
}

// Unwrap returns the underlying error that caused this RequestError.
// This method enables using the errors.Is and errors.As functions
// to inspect and compare the underlying error.
func (e *RequestError) Unwrap() error {
	return e.Err
}

// StreamError represents an error that occurs during streaming operations,
// such as errors reading from an SSE stream, deserializing stream data,
// or when a stream is unexpectedly closed.
type StreamError struct {
	// Message describes the error in human-readable form
	Message string

	// Err is the underlying error that caused this StreamError
	Err error
}

// Error implements the error interface by returning a formatted message
// that includes the stream error description and the underlying error, if any.
func (e *StreamError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("stream error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("stream error: %s", e.Message)
}

// Code returns the string code "stream_error" to identify streaming-related errors.
func (e *StreamError) Code() string {
	return "stream_error"
}

// HTTPStatusCode returns http.StatusInternalServerError (500) for stream errors.
func (e *StreamError) HTTPStatusCode() int {
	return http.StatusInternalServerError
}

// Unwrap returns the underlying error that caused this StreamError.
// This method enables using the errors.Is and errors.As functions
// to inspect and compare the underlying error.
func (e *StreamError) Unwrap() error {
	return e.Err
}

// Predefined error variables provide common error instances that can be
// returned or checked against using errors.Is() for specific error conditions.
var (
	// ErrInvalidRequest is returned when request parameters are invalid or malformed.
	// Status code: 400 Bad Request
	ErrInvalidRequest = &RequestError{Message: "invalid request parameters", StatusCode: http.StatusBadRequest}

	// ErrAuthentication is returned when authentication fails.
	// Status code: 401 Unauthorized
	ErrAuthentication = &APIError{StatusCode: http.StatusUnauthorized, Message: "authentication failed"}

	// ErrPermission is returned when the authenticated user lacks permission for the requested operation.
	// Status code: 403 Forbidden
	ErrPermission = &APIError{StatusCode: http.StatusForbidden, Message: "permission denied"}

	// ErrNotFound is returned when the requested resource does not exist.
	// Status code: 404 Not Found
	ErrNotFound = &APIError{StatusCode: http.StatusNotFound, Message: "resource not found"}

	// ErrServerError is returned for general server-side errors.
	// Status code: 500 Internal Server Error
	ErrServerError = &APIError{StatusCode: http.StatusInternalServerError, Message: "server error"}

	// ErrTimeout is returned when a request exceeds the timeout duration.
	// Status code: 504 Gateway Timeout
	ErrTimeout = &RequestError{Message: "request timed out", StatusCode: http.StatusGatewayTimeout}

	// ErrCanceled is returned when a request is canceled, typically by the context being canceled.
	ErrCanceled = &RequestError{Message: "request canceled"}

	// ErrStreamClosed is returned when attempting to read from a closed stream.
	// This typically happens if Recv() is called after Close() or after the stream ends.
	ErrStreamClosed = &StreamError{Message: "stream closed"}
)
