// Package errors provides error types for the TabbyAPI client.
package errors

import (
	"fmt"
	"net/http"
)

// Error is the interface implemented by all errors in the library
type Error interface {
	error
	Code() string
	HTTPStatusCode() int
}

// APIError represents an error returned by the TabbyAPI
type APIError struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
	RequestID  string      `json:"request_id,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("tabby API error (status %d, request ID: %s): %s", e.StatusCode, e.RequestID, e.Message)
	}
	return fmt.Sprintf("tabby API error (status %d): %s", e.StatusCode, e.Message)
}

// Code returns a string code for the error
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

// HTTPStatusCode returns the HTTP status code
func (e *APIError) HTTPStatusCode() int {
	return e.StatusCode
}

// RequestError represents a client request error
type RequestError struct {
	Message    string
	StatusCode int
	Err        error
}

// Error implements the error interface
func (e *RequestError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("request error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("request error: %s", e.Message)
}

// Code returns a string code for the error
func (e *RequestError) Code() string {
	return "request_error"
}

// HTTPStatusCode returns the HTTP status code
func (e *RequestError) HTTPStatusCode() int {
	if e.StatusCode != 0 {
		return e.StatusCode
	}
	return http.StatusInternalServerError
}

// Unwrap returns the wrapped error
func (e *RequestError) Unwrap() error {
	return e.Err
}

// StreamError represents an error in a streaming response
type StreamError struct {
	Message string
	Err     error
}

// Error implements the error interface
func (e *StreamError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("stream error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("stream error: %s", e.Message)
}

// Code returns a string code for the error
func (e *StreamError) Code() string {
	return "stream_error"
}

// HTTPStatusCode returns the HTTP status code
func (e *StreamError) HTTPStatusCode() int {
	return http.StatusInternalServerError
}

// Unwrap returns the wrapped error
func (e *StreamError) Unwrap() error {
	return e.Err
}

// Predefined error variables
var (
	ErrInvalidRequest = &RequestError{Message: "invalid request parameters", StatusCode: http.StatusBadRequest}
	ErrAuthentication = &APIError{StatusCode: http.StatusUnauthorized, Message: "authentication failed"}
	ErrPermission     = &APIError{StatusCode: http.StatusForbidden, Message: "permission denied"}
	ErrNotFound       = &APIError{StatusCode: http.StatusNotFound, Message: "resource not found"}
	ErrServerError    = &APIError{StatusCode: http.StatusInternalServerError, Message: "server error"}
	ErrTimeout        = &RequestError{Message: "request timed out", StatusCode: http.StatusGatewayTimeout}
	ErrCanceled       = &RequestError{Message: "request canceled"}
	ErrStreamClosed   = &StreamError{Message: "stream closed"}
)
