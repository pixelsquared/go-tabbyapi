# Error Handling in go-tabbyapi

This guide explains the error handling system in go-tabbyapi and provides best practices for dealing with errors in your applications.

## Error Types

The go-tabbyapi library uses a structured error system with several specialized error types that implement the common `tabby.Error` interface:

```go
// Error is the interface implemented by all errors in the TabbyAPI client library.
type Error interface {
	// error provides the standard Error() string method
	error

	// Code returns a string identifier for the error type
	Code() string

	// HTTPStatusCode returns the HTTP status code associated with the error
	HTTPStatusCode() int
}
```

### APIError

`APIError` represents errors returned directly by the TabbyAPI server:

```go
type APIError struct {
	StatusCode int         `json:"status_code"` // HTTP status code
	Message    string      `json:"message"`     // Error message
	Details    interface{} `json:"details,omitempty"` // Additional details
	RequestID  string      `json:"request_id,omitempty"` // For debugging
}
```

Common APIError codes include:
- `"invalid_request"` (400) - Malformed request or invalid parameters
- `"authentication_error"` (401) - Invalid or missing authentication credentials
- `"permission_error"` (403) - Insufficient permissions for the operation
- `"not_found"` (404) - Requested resource does not exist
- `"rate_limit_exceeded"` (429) - Too many requests
- `"server_error"` (500+) - Internal server error

### ValidationError

`ValidationError` represents client-side validation errors:

```go
type ValidationError struct {
	Field   string `json:"field"`   // Field that failed validation
	Message string `json:"message"` // Validation error message
	Type    string `json:"type,omitempty"` // Type of validation error
}
```

### RequestError

`RequestError` represents errors that occur during the HTTP request lifecycle:

```go
type RequestError struct {
	Message    string // Error description
	StatusCode int    // HTTP status code (if applicable)
	Err        error  // Underlying error
}
```

### StreamError

`StreamError` represents errors that occur during streaming operations:

```go
type StreamError struct {
	Message string // Error description
	Err     error  // Underlying error
}
```

## Predefined Error Variables

The library provides predefined error variables for common error conditions:

```go
var (
	ErrInvalidRequest = &RequestError{Message: "invalid request parameters", StatusCode: 400}
	ErrAuthentication = &APIError{StatusCode: 401, Message: "authentication failed"}
	ErrPermission     = &APIError{StatusCode: 403, Message: "permission denied"}
	ErrNotFound       = &APIError{StatusCode: 404, Message: "resource not found"}
	ErrServerError    = &APIError{StatusCode: 500, Message: "server error"}
	ErrTimeout        = &RequestError{Message: "request timed out", StatusCode: 504}
	ErrCanceled       = &RequestError{Message: "request canceled"}
	ErrStreamClosed   = &StreamError{Message: "stream closed"}
)
```

You can use these variables with `errors.Is()` for simple error checking.

## Error Handling Patterns

### Basic Error Handling

For simple cases, you can handle errors with a standard Go error check:

```go
resp, err := client.Completions().Create(ctx, req)
if err != nil {
    log.Fatalf("Error generating completion: %v", err)
}
```

### Type Assertions

For more detailed handling, use type assertions to check specific error types:

```go
resp, err := client.Completions().Create(ctx, req)
if err != nil {
    var apiErr *tabby.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("API Error (Status %d): %s\n", apiErr.HTTPStatusCode(), apiErr.Error())
        // Handle based on status code or error code
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
    return
}
```

### Error Code Handling

Handle errors based on their error code:

```go
resp, err := client.Chat().Create(ctx, req)
if err != nil {
    var apiErr *tabby.APIError
    if errors.As(err, &apiErr) {
        switch apiErr.Code() {
        case "authentication_error":
            fmt.Println("Authentication failed. Check your API key.")
        case "permission_error":
            fmt.Println("You don't have permission to perform this operation.")
        case "invalid_request":
            fmt.Println("Invalid request parameters:", apiErr.Error())
        case "not_found":
            fmt.Println("Resource not found:", apiErr.Error())
        case "rate_limit_exceeded":
            fmt.Println("Rate limit exceeded. Try again later.")
        case "server_error":
            fmt.Println("Server error. Try again later.")
        default:
            fmt.Printf("API error (%s): %s\n", apiErr.Code(), apiErr.Error())
        }
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
    return
}
```

### HTTP Status Code Handling

Handle errors based on their HTTP status code:

```go
resp, err := client.Models().Load(ctx, req)
if err != nil {
    var apiErr *tabby.APIError
    if errors.As(err, &apiErr) {
        switch apiErr.HTTPStatusCode() {
        case http.StatusBadRequest: // 400
            fmt.Println("Bad request:", apiErr.Error())
        case http.StatusUnauthorized: // 401
            fmt.Println("Unauthorized. Check your API key.")
        case http.StatusForbidden: // 403
            fmt.Println("Forbidden. You don't have permission.")
        case http.StatusNotFound: // 404
            fmt.Println("Resource not found:", apiErr.Error())
        case http.StatusTooManyRequests: // 429
            fmt.Println("Too many requests. Try again later.")
        case http.StatusInternalServerError: // 500
            fmt.Println("Server error. Try again later.")
        default:
            fmt.Printf("HTTP error %d: %s\n", apiErr.HTTPStatusCode(), apiErr.Error())
        }
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
    return
}
```

### Handling Request Errors

Handle network or client-side errors:

```go
resp, err := client.Embeddings().Create(ctx, req)
if err != nil {
    var reqErr *tabby.RequestError
    if errors.As(err, &reqErr) {
        // Check the underlying error
        switch {
        case errors.Is(reqErr.Unwrap(), context.DeadlineExceeded):
            fmt.Println("Request timed out. Try increasing the timeout.")
        case errors.Is(reqErr.Unwrap(), context.Canceled):
            fmt.Println("Request was canceled.")
        default:
            fmt.Printf("Request error: %v\n", reqErr)
        }
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
    return
}
```

### Handling Stream Errors

Handle errors during streaming operations:

```go
stream, err := client.Completions().CreateStream(ctx, req)
if err != nil {
    log.Fatalf("Error creating stream: %v", err)
}
defer stream.Close()

for {
    resp, err := stream.Recv()
    if err == io.EOF {
        break // Stream completed
    }
    
    if err != nil {
        var streamErr *tabby.StreamError
        if errors.As(err, &streamErr) {
            fmt.Printf("Stream error: %v\n", streamErr)
            if errors.Is(streamErr.Unwrap(), tabby.ErrStreamClosed) {
                fmt.Println("Stream was closed.")
            }
        } else {
            fmt.Printf("Other error: %v\n", err)
        }
        break
    }
    
    // Process response
    fmt.Print(resp.Choices[0].Text)
}
```

## Comprehensive Error Handling Example

Here's a comprehensive example that handles various error scenarios:

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "io"
    "log"
    "net/http"
    "time"
    
    "github.com/pixelsquared/go-tabbyapi/tabby"
)

func main() {
    client := tabby.NewClient(
        tabby.WithBaseURL("http://localhost:8080"),
        tabby.WithAPIKey("your-api-key"),
        tabby.WithTimeout(30*time.Second),
    )
    defer client.Close()
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Example operation
    req := &tabby.CompletionRequest{
        Prompt:    "Write a function to calculate the factorial of a number:",
        MaxTokens: 100,
    }
    
    resp, err := client.Completions().Create(ctx, req)
    if err != nil {
        handleError(err)
        return
    }
    
    // Process successful response
    if len(resp.Choices) > 0 {
        fmt.Println("Generated code:")
        fmt.Println(resp.Choices[0].Text)
    }
}

// handleError provides comprehensive error handling for all error types
func handleError(err error) {
    // Check for predefined errors with errors.Is
    switch {
    case errors.Is(err, tabby.ErrInvalidRequest):
        fmt.Println("Invalid request parameters. Please check your request.")
        return
    case errors.Is(err, tabby.ErrAuthentication):
        fmt.Println("Authentication failed. Please check your API key.")
        return
    case errors.Is(err, tabby.ErrPermission):
        fmt.Println("Permission denied. Your API key doesn't have sufficient permissions.")
        return
    case errors.Is(err, tabby.ErrNotFound):
        fmt.Println("Resource not found. Please check your request.")
        return
    case errors.Is(err, tabby.ErrTimeout):
        fmt.Println("Request timed out. Please try again later.")
        return
    case errors.Is(err, tabby.ErrCanceled):
        fmt.Println("Request was canceled.")
        return
    case errors.Is(err, tabby.ErrStreamClosed):
        fmt.Println("Stream was closed unexpectedly.")
        return
    }
    
    // Check for specific error types with errors.As
    var apiErr *tabby.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("API Error (Code: %s, Status: %d): %s\n", 
            apiErr.Code(), apiErr.HTTPStatusCode(), apiErr.Error())
        
        // Handle specific status codes
        switch apiErr.HTTPStatusCode() {
        case http.StatusTooManyRequests: // 429
            fmt.Println("Rate limit exceeded. Please try again later.")
        case http.StatusServiceUnavailable: // 503
            fmt.Println("Service unavailable. Please try again later.")
        }
        return
    }
    
    var validationErr *tabby.ValidationError
    if errors.As(err, &validationErr) {
        fmt.Printf("Validation error in field '%s': %s\n", 
            validationErr.Field, validationErr.Error())
        return
    }
    
    var reqErr *tabby.RequestError
    if errors.As(err, &reqErr) {
        fmt.Printf("Request error: %s\n", reqErr.Error())
        
        // Check underlying error
        if reqErr.Err != nil {
            if errors.Is(reqErr.Err, context.DeadlineExceeded) {
                fmt.Println("The request exceeded the configured timeout.")
            } else if errors.Is(reqErr.Err, context.Canceled) {
                fmt.Println("The request was canceled.")
            }
        }
        return
    }
    
    var streamErr *tabby.StreamError
    if errors.As(err, &streamErr) {
        fmt.Printf("Stream error: %s\n", streamErr.Error())
        return
    }
    
    // Handle any other errors
    fmt.Printf("Unknown error: %v\n", err)
}
```

## Retry Strategies

For transient errors, implement a retry strategy:

```go
// retryWithBackoff attempts an operation with exponential backoff
func retryWithBackoff(operation func() error, maxRetries int) error {
    var err error
    
    for attempt := 0; attempt <= maxRetries; attempt++ {
        // First attempt (attempt=0) has no delay
        if attempt > 0 {
            // Calculate backoff with jitter
            backoff := time.Duration(1<<uint(attempt-1)) * time.Second
            jitter := time.Duration(rand.Int63n(int64(time.Second)))
            delay := backoff + jitter
            
            fmt.Printf("Retrying after %v (attempt %d of %d)...\n", 
                delay, attempt, maxRetries)
            time.Sleep(delay)
        }
        
        // Attempt the operation
        err = operation()
        
        // Check if the error is retryable
        if err == nil || !isRetryableError(err) {
            return err
        }
    }
    
    return fmt.Errorf("operation failed after %d attempts: %w", maxRetries, err)
}

// isRetryableError determines if an error should be retried
func isRetryableError(err error) bool {
    // Check for network errors and server errors
    var apiErr *tabby.APIError
    if errors.As(err, &apiErr) {
        // Retry on server errors (5xx) and rate limiting (429)
        statusCode := apiErr.HTTPStatusCode()
        return statusCode >= 500 || statusCode == 429
    }
    
    var reqErr *tabby.RequestError
    if errors.As(err, &reqErr) {
        // Retry on most request errors except cancellation
        return !errors.Is(reqErr.Unwrap(), context.Canceled)
    }
    
    // Don't retry other types of errors
    return false
}

// Example usage
func generateWithRetry(client tabby.Client, ctx context.Context, req *tabby.CompletionRequest) (*tabby.CompletionResponse, error) {
    var resp *tabby.CompletionResponse
    
    err := retryWithBackoff(func() error {
        var err error
        resp, err = client.Completions().Create(ctx, req)
        return err
    }, 3) // Maximum 3 retries
    
    return resp, err
}
```

## Error Logging Best Practices

When logging errors, include relevant context:

```go
func logError(operation string, err error, requestDetails interface{}) {
    // Basic error information
    errorInfo := map[string]interface{}{
        "operation": operation,
        "error":     err.Error(),
        "timestamp": time.Now().Format(time.RFC3339),
        "details":   requestDetails,
    }
    
    // Add error-specific information
    var apiErr *tabby.APIError
    if errors.As(err, &apiErr) {
        errorInfo["error_code"] = apiErr.Code()
        errorInfo["status_code"] = apiErr.HTTPStatusCode()
        errorInfo["request_id"] = apiErr.RequestID
    }
    
    // Log as JSON for structured logging
    jsonBytes, _ := json.Marshal(errorInfo)
    log.Println(string(jsonBytes))
}

// Example usage
func generateCompletion(client tabby.Client, ctx context.Context, prompt string) {
    req := &tabby.CompletionRequest{
        Prompt:    prompt,
        MaxTokens: 100,
    }
    
    resp, err := client.Completions().Create(ctx, req)
    if err != nil {
        // Log the error with context
        logError("generate_completion", err, map[string]interface{}{
            "prompt":     prompt,
            "max_tokens": 100,
        })
        return
    }
    
    // Process the response
    // ...
}
```

## Summary

- Use the `errors.As()` function to check for specific error types (`APIError`, `RequestError`, etc.)
- Use the `errors.Is()` function to check for predefined error variables
- Implement retry strategies for transient errors like server errors (5xx) and rate limiting (429)
- Include relevant context when logging errors
- Consider creating a dedicated error handling package for large applications