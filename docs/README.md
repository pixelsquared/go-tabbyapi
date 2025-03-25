# Go TabbyAPI Client

[![Go Reference](https://pkg.go.dev/badge/github.com/pixelsquared/go-tabbyapi.svg)](https://pkg.go.dev/github.com/pixelsquared/go-tabbyapi)

A comprehensive Go client library for interacting with TabbyAPI, an open-source self-hosted AI coding assistant.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Authentication](#authentication)
- [Error Handling](#error-handling)
- [Services](#services)
- [Examples](#examples)
- [Advanced Usage](#advanced-usage)
- [API Reference](#api-reference)

## Installation

```bash
go get github.com/pixelsquared/go-tabbyapi
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/pixelsquared/go-tabbyapi/tabby"
)

func main() {
    // Create a new TabbyAPI client
    client := tabby.NewClient(
        tabby.WithBaseURL("http://localhost:8080"),
        tabby.WithAPIKey("your-api-key"),
        tabby.WithTimeout(30*time.Second),
    )
    defer client.Close()
    
    // Create a context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Generate a completion
    resp, err := client.Completions().Create(ctx, &tabby.CompletionRequest{
        Prompt:      "func fibonacci(n int) int {",
        MaxTokens:   100,
        Temperature: 0.7,
    })
    
    if err != nil {
        log.Fatalf("Error generating completion: %v", err)
    }
    
    // Print the generated text
    fmt.Println(resp.Choices[0].Text)
}
```

## Authentication

The library supports multiple authentication methods:

### API Key

```go
client := tabby.NewClient(
    tabby.WithAPIKey("your-api-key"),
)
```

### Admin Key

```go
client := tabby.NewClient(
    tabby.WithAdminKey("your-admin-key"),
)
```

### Bearer Token

```go
client := tabby.NewClient(
    tabby.WithBearerToken("your-bearer-token"),
)
```

## Error Handling

The library provides detailed error types to help you handle different kinds of errors:

```go
resp, err := client.Completions().Create(ctx, req)
if err != nil {
    // Check for specific error types
    var apiErr *tabby.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("API Error (code %s): %s\n", apiErr.Code(), apiErr.Error())
        
        // Check specific error types
        if apiErr.HTTPStatusCode() == http.StatusUnauthorized {
            fmt.Println("Authentication failed")
        }
    }
    
    // Or use predefined error variables
    if errors.Is(err, tabby.ErrAuthentication) {
        fmt.Println("Authentication failed")
    } else if errors.Is(err, tabby.ErrPermission) {
        fmt.Println("Permission denied")
    }
    
    log.Fatalf("Error: %v", err)
}
```

## Services

The client provides access to all TabbyAPI services:

- **Completions**: Generate text completions
- **Chat**: Chat-based interactions
- **Embeddings**: Generate vector embeddings
- **Models**: Model management
- **Lora**: LoRA adapter management
- **Templates**: Manage prompt templates
- **Tokens**: Token encoding and decoding
- **Sampling**: Sampling parameter management
- **Health**: Health checks
- **Auth**: Authentication permissions

See the [Services documentation](services/README.md) for detailed information about each service.

## Examples

Check out the [examples](../examples) directory for comprehensive examples of each service.

## Advanced Usage

### Streaming Responses

For long-running requests, you can use streaming to get incremental results:

```go
stream, err := client.Completions().CreateStream(ctx, &tabby.CompletionRequest{
    Prompt:      "Explain quantum computing in simple terms:",
    MaxTokens:   500,
    Temperature: 0.7,
    Stream:      true,
})
if err != nil {
    log.Fatalf("Error creating stream: %v", err)
}
defer stream.Close()

// Process streaming responses
for {
    resp, err := stream.Recv()
    if err == io.EOF {
        break  // Stream completed
    }
    if err != nil {
        log.Fatalf("Error receiving from stream: %v", err)
    }
    
    // Process each chunk
    fmt.Print(resp.Choices[0].Text)
}
```

### Custom HTTP Client

You can provide a custom HTTP client for advanced configuration:

```go
httpClient := &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 100,
        IdleConnTimeout:     90 * time.Second,
    },
}

client := tabby.NewClient(
    tabby.WithHTTPClient(httpClient),
)
```

### Custom Retry Policy

Configure retries for failed requests:

```go
retryPolicy := &tabby.SimpleRetryPolicy{
    MaxRetryCount: 5,
    RetryDelayFunc: func(attempts int) time.Duration {
        return time.Duration(1<<uint(attempts-1)) * time.Second
    },
    RetryableFunc: func(resp *http.Response, err error) bool {
        // Retry on connection errors or server errors (5xx)
        return err != nil || (resp != nil && resp.StatusCode >= 500)
    },
}

client := tabby.NewClient(
    tabby.WithRetryPolicy(retryPolicy),
)
```

## API Reference

For complete API documentation, see the [pkg.go.dev reference](https://pkg.go.dev/github.com/pixelsquared/go-tabbyapi).