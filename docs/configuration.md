# Configuration Guide

This guide explains all the configuration options available in the go-tabbyapi client library. Understanding these options will help you optimize the client for your specific use case.

## Client Configuration Options

The `tabby.NewClient()` function accepts multiple configuration options using the functional options pattern:

```go
client := tabby.NewClient(
    tabby.WithBaseURL("http://localhost:8080"),
    tabby.WithAPIKey("your-api-key"),
    tabby.WithTimeout(30*time.Second),
    // Add more options as needed
)
```

Let's explore each configuration option in detail:

## Core Options

### WithBaseURL

Sets the base URL for the TabbyAPI server:

```go
tabby.WithBaseURL("http://localhost:8080")
```

- **Default**: `"http://localhost:8080"`
- **Purpose**: Specifies the endpoint where the TabbyAPI server is running
- **Format**: Should include protocol (http/https) and host, with optional port
- **Examples**:
  - Local server: `"http://localhost:8080"`
  - Remote server: `"https://tabby-api.example.com"`
  - Custom port: `"http://localhost:9000"`

**Note**: Do not include API version or endpoint paths in the base URL. The library automatically appends them.

### WithTimeout

Sets the timeout duration for all API requests:

```go
tabby.WithTimeout(30*time.Second)
```

- **Default**: `30*time.Second`
- **Purpose**: Controls how long to wait for API responses before timing out
- **Considerations**:
  - Set shorter timeouts for quick operations like health checks
  - Set longer timeouts for model loading or large generation requests
  - For streaming operations, this timeout applies to establishing the connection, not the entire stream duration

**Example with different timeouts**:

```go
// Quick operations
quickClient := tabby.NewClient(
    tabby.WithBaseURL("http://localhost:8080"),
    tabby.WithTimeout(5*time.Second),
)

// Long-running operations
longRunningClient := tabby.NewClient(
    tabby.WithBaseURL("http://localhost:8080"),
    tabby.WithTimeout(2*time.Minute),
)
```

### WithHTTPClient

Provides a custom HTTP client for making API requests:

```go
customClient := &http.Client{
    // Custom configuration
}
tabby.WithHTTPClient(customClient)
```

- **Default**: Standard HTTP client with 30-second timeout
- **Purpose**: Allows advanced HTTP configuration such as:
  - Custom transport settings
  - Connection pooling
  - Proxy configuration
  - TLS settings
  - Custom redirect policy

**Example with advanced configuration**:

```go
transport := &http.Transport{
    MaxIdleConns:        100,
    MaxIdleConnsPerHost: 100,
    IdleConnTimeout:     90 * time.Second,
    TLSClientConfig: &tls.Config{
        MinVersion: tls.VersionTLS12,
    },
    DisableCompression: false,
    Proxy:              http.ProxyFromEnvironment,
}

customClient := &http.Client{
    Transport: transport,
    Timeout:   60 * time.Second,
}

client := tabby.NewClient(
    tabby.WithBaseURL("http://localhost:8080"),
    tabby.WithHTTPClient(customClient),
)
```

## Authentication Options

### WithAPIKey

Sets the API key for standard authentication:

```go
tabby.WithAPIKey("your-api-key")
```

- **Default**: No authentication
- **Purpose**: Authenticates requests using the X-API-Key header
- **Permissions**: Typically provides read and write permissions but not administrative access
- **Usage**: Standard choice for most applications

### WithAdminKey

Sets the admin key for administrative operations:

```go
tabby.WithAdminKey("your-admin-key")
```

- **Default**: No authentication
- **Purpose**: Authenticates requests using the X-Admin-Key header
- **Permissions**: Provides full administrative access, including model management
- **Usage**: Use only when administrative operations are required

### WithBearerToken

Sets a bearer token for OAuth or JWT authentication:

```go
tabby.WithBearerToken("your-bearer-token")
```

- **Default**: No authentication
- **Purpose**: Authenticates requests using the Authorization header with Bearer scheme
- **Permissions**: Depends on the token's claims and server configuration
- **Usage**: Useful when integrating with OAuth/JWT systems

**Note**: Only use one authentication method at a time. If multiple are provided, the last one specified will take precedence.

## Retry Policy Options

### WithRetryPolicy

Sets the retry policy for failed requests:

```go
tabby.WithRetryPolicy(tabby.DefaultRetryPolicy())
```

- **Default**: No retry policy
- **Purpose**: Automatically retries failed requests based on configurable conditions
- **Implementation**: Uses the `RetryPolicy` interface with these methods:
  - `ShouldRetry(resp *http.Response, err error) bool`: Determines if a request should be retried
  - `RetryDelay(attempts int) time.Duration`: Returns the delay before the next retry
  - `MaxRetries() int`: Returns the maximum number of retry attempts

### Default Retry Policy

The library provides a default retry policy that can be used:

```go
tabby.WithRetryPolicy(tabby.DefaultRetryPolicy())
```

The default policy provides:
- Maximum of 3 retry attempts
- Exponential backoff with jitter:
  - 1st retry: ~100ms
  - 2nd retry: ~1.2s 
  - 3rd retry: ~2.3s
- Retries on:
  - Any network or connection error
  - Any HTTP status code >= 500 (server errors)

### Custom Retry Policy

For more control, you can create a custom retry policy:

```go
customPolicy := &tabby.SimpleRetryPolicy{
    MaxRetryCount: 5,  // Try up to 5 times
    
    // Exponential backoff with jitter
    RetryDelayFunc: func(attempts int) time.Duration {
        baseDelay := time.Duration(1<<uint(attempts-1)) * time.Second  // 1s, 2s, 4s, 8s, 16s
        jitter := time.Duration(rand.Int63n(int64(500 * time.Millisecond)))
        return baseDelay + jitter
    },
    
    // Retry conditions
    RetryableFunc: func(resp *http.Response, err error) bool {
        // Retry on network errors
        if err != nil {
            return true
        }
        
        // Retry on server errors (5xx) and rate limiting (429)
        statusCode := resp.StatusCode
        return statusCode >= 500 || statusCode == 429
    },
}

client := tabby.NewClient(
    tabby.WithBaseURL("http://localhost:8080"),
    tabby.WithAPIKey("your-api-key"),
    tabby.WithRetryPolicy(customPolicy),
)
```

## Complete Configuration Example

Here's a comprehensive example showing all configuration options together:

```go
package main

import (
    "crypto/tls"
    "math/rand"
    "net"
    "net/http"
    "time"
    
    "github.com/pixelsquared/go-tabbyapi/tabby"
)

func main() {
    // Seed the random number generator (for jitter in retry policy)
    rand.Seed(time.Now().UnixNano())
    
    // Create a custom transport
    transport := &http.Transport{
        Proxy: http.ProxyFromEnvironment,
        DialContext: (&net.Dialer{
            Timeout:   30 * time.Second,
            KeepAlive: 30 * time.Second,
        }).DialContext,
        ForceAttemptHTTP2:     true,
        MaxIdleConns:          100,
        IdleConnTimeout:       90 * time.Second,
        TLSHandshakeTimeout:   10 * time.Second,
        ExpectContinueTimeout: 1 * time.Second,
        TLSClientConfig: &tls.Config{
            MinVersion: tls.VersionTLS12,
        },
    }
    
    // Create a custom HTTP client
    httpClient := &http.Client{
        Transport: transport,
        Timeout:   60 * time.Second,
    }
    
    // Create a custom retry policy
    retryPolicy := &tabby.SimpleRetryPolicy{
        MaxRetryCount: 5,
        RetryDelayFunc: func(attempts int) time.Duration {
            baseDelay := time.Duration(1<<uint(attempts-1)) * time.Second
            jitter := time.Duration(rand.Int63n(int64(500 * time.Millisecond)))
            return baseDelay + jitter
        },
        RetryableFunc: func(resp *http.Response, err error) bool {
            if err != nil {
                return true
            }
            statusCode := resp.StatusCode
            return statusCode >= 500 || statusCode == 429
        },
    }
    
    // Create the TabbyAPI client with all configuration options
    client := tabby.NewClient(
        // Core options
        tabby.WithBaseURL("https://tabby-api.example.com"),
        tabby.WithHTTPClient(httpClient),
        
        // Authentication option (choose one)
        tabby.WithAPIKey("your-api-key"),
        // OR tabby.WithAdminKey("your-admin-key"),
        // OR tabby.WithBearerToken("your-bearer-token"),
        
        // Retry policy
        tabby.WithRetryPolicy(retryPolicy),
    )
    defer client.Close()
    
    // Use the client to interact with TabbyAPI
    // ...
}
```

## Client Methods

Once configured, the client provides access to all TabbyAPI services:

```go
// Create the client
client := tabby.NewClient(/* options */)

// Access services
completions := client.Completions()  // CompletionsService
chat := client.Chat()                // ChatService
models := client.Models()            // ModelsService
embeddings := client.Embeddings()    // EmbeddingsService
lora := client.Lora()                // LoraService
templates := client.Templates()      // TemplatesService
tokens := client.Tokens()            // TokensService
sampling := client.Sampling()        // SamplingService
health := client.Health()            // HealthService
auth := client.Auth()                // AuthService
```

## Configuration Best Practices

### Timeouts

- **Set appropriate timeouts** based on operation type:
  - Health checks: 5-10 seconds
  - Text generation: 30-60 seconds
  - Model loading: 2-5 minutes
  - Embedding generation: 10-30 seconds

```go
// Different clients for different operation types
healthClient := tabby.NewClient(
    tabby.WithBaseURL(baseURL),
    tabby.WithAPIKey(apiKey),
    tabby.WithTimeout(5*time.Second),
)

generationClient := tabby.NewClient(
    tabby.WithBaseURL(baseURL),
    tabby.WithAPIKey(apiKey),
    tabby.WithTimeout(60*time.Second),
)

adminClient := tabby.NewClient(
    tabby.WithBaseURL(baseURL),
    tabby.WithAdminKey(adminKey),
    tabby.WithTimeout(5*time.Minute),
)
```

### Authentication

- **Use the least privileged authentication** required for your operations:
  - Regular operations (completions, chat, embeddings): Use API key
  - Administrative operations (model management): Use admin key
  - Never mix admin and non-admin operations in the same client instance

```go
// Client for standard operations
standardClient := tabby.NewClient(
    tabby.WithBaseURL(baseURL),
    tabby.WithAPIKey(apiKey),
)

// Separate client for admin operations
adminClient := tabby.NewClient(
    tabby.WithBaseURL(baseURL),
    tabby.WithAdminKey(adminKey),
)
```

### Retry Policies

- **Use retries for transient errors** only:
  - Network issues
  - Server errors (5xx)
  - Rate limiting (429)
  - Do not retry client errors (4xx) except for rate limiting

- **Implement exponential backoff** to avoid overwhelming the server:
  ```go
  baseDelay := time.Duration(1<<uint(attempts-1)) * time.Second  // 1s, 2s, 4s, 8s...
  ```

- **Add jitter** to prevent retry storms:
  ```go
  jitter := time.Duration(rand.Int63n(int64(500 * time.Millisecond)))
  delay := baseDelay + jitter
  ```

### HTTP Client Tuning

- **Increase MaxIdleConnsPerHost** for high-throughput applications:
  ```go
  transport.MaxIdleConnsPerHost = 100
  ```

- **Set reasonable connection timeouts**:
  ```go
  transport.IdleConnTimeout = 90 * time.Second
  ```

- **Enable HTTP/2** for better performance:
  ```go
  transport.ForceAttemptHTTP2 = true
  ```

## Environment-Specific Configuration

### Production

```go
client := tabby.NewClient(
    tabby.WithBaseURL(os.Getenv("TABBY_API_URL")),
    tabby.WithAPIKey(os.Getenv("TABBY_API_KEY")),
    tabby.WithTimeout(60*time.Second),
    tabby.WithRetryPolicy(tabby.DefaultRetryPolicy()),
)
```

### Development

```go
client := tabby.NewClient(
    tabby.WithBaseURL("http://localhost:8080"),
    // Often no authentication needed in development
    tabby.WithTimeout(30*time.Second),
    // No retry policy for faster debugging
)
```

### Testing

```go
// Create a test server
testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Mock responses here
}))
defer testServer.Close()

// Configure client to use test server
client := tabby.NewClient(
    tabby.WithBaseURL(testServer.URL),
    tabby.WithTimeout(5*time.Second),
)
```

## Creating and Managing Multiple Clients

For different operation types, consider creating specialized clients:

```go
type TabbyClients struct {
    Standard *tabby.Client // For regular operations
    Admin    *tabby.Client // For admin operations
    Streaming *tabby.Client // For streaming operations with longer timeouts
}

func NewTabbyClients(config Config) *TabbyClients {
    standard := tabby.NewClient(
        tabby.WithBaseURL(config.BaseURL),
        tabby.WithAPIKey(config.APIKey),
        tabby.WithTimeout(30*time.Second),
        tabby.WithRetryPolicy(tabby.DefaultRetryPolicy()),
    )
    
    admin := tabby.NewClient(
        tabby.WithBaseURL(config.BaseURL),
        tabby.WithAdminKey(config.AdminKey),
        tabby.WithTimeout(2*time.Minute),
        tabby.WithRetryPolicy(createAdminRetryPolicy()),
    )
    
    streaming := tabby.NewClient(
        tabby.WithBaseURL(config.BaseURL),
        tabby.WithAPIKey(config.APIKey),
        tabby.WithTimeout(5*time.Minute),
        // No retry policy for streaming
    )
    
    return &TabbyClients{
        Standard: standard,
        Admin:    admin,
        Streaming: streaming,
    }
}

// Remember to close all clients
func (c *TabbyClients) Close() {
    c.Standard.Close()
    c.Admin.Close()
    c.Streaming.Close()
}
```

## Conclusion

Properly configuring your TabbyAPI client can significantly impact the performance, reliability, and security of your application. Choose configuration options that match your specific requirements, and remember to always close your client when you're done using it to release resources.

For more detailed information about using the configured client, refer to the other documentation pages:
- [Getting Started](./getting-started.md)
- [Error Handling](./error-handling.md)
- [Services Documentation](./services/README.md)