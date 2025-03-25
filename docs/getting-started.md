# Getting Started with go-tabbyapi

This guide will help you get started with the go-tabbyapi client library, providing step-by-step instructions for common use cases.

## Prerequisites

Before you begin, make sure you have:

- Go 1.18 or higher installed
- A running TabbyAPI server (local or remote)
- API key or admin key for authentication (if required by your TabbyAPI server)

## Installation

Install the library using the standard Go toolchain:

```bash
go get github.com/pixelsquared/go-tabbyapi
```

## Basic Setup

All interactions with TabbyAPI begin with creating a client instance:

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
        tabby.WithBaseURL("http://localhost:8080"), // TabbyAPI server address
        tabby.WithAPIKey("your-api-key"),           // Optional API key
        tabby.WithTimeout(30*time.Second),          // Request timeout
    )
    defer client.Close() // Always close the client when done
    
    // Create a context with timeout for operations
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Now you can use the client to interact with TabbyAPI
    // ...
}
```

## Common Use Cases

### 1. Generating Text Completions

Generate text based on a prompt:

```go
func generateCompletion(client tabby.Client, ctx context.Context) {
    // Create a completion request
    req := &tabby.CompletionRequest{
        Prompt:      "Once upon a time in a distant galaxy,",
        MaxTokens:   100,
        Temperature: 0.7,
    }
    
    // Call the API
    resp, err := client.Completions().Create(ctx, req)
    if err != nil {
        log.Fatalf("Error generating completion: %v", err)
    }
    
    // Print the generated text
    if len(resp.Choices) > 0 {
        fmt.Printf("Generated text: %s\n", resp.Choices[0].Text)
    }
}
```

### 2. Having a Chat Conversation

Interact with a model using multi-turn conversations:

```go
func chatConversation(client tabby.Client, ctx context.Context) {
    // Define the conversation
    messages := []tabby.ChatMessage{
        {
            Role:    tabby.ChatMessageRoleSystem,
            Content: "You are a helpful assistant specialized in Python programming.",
        },
        {
            Role:    tabby.ChatMessageRoleUser,
            Content: "How do I create a web server with Flask?",
        },
    }
    
    // Create a chat completion request
    req := &tabby.ChatCompletionRequest{
        Messages:    messages,
        MaxTokens:   150,
        Temperature: 0.7,
    }
    
    // Call the API
    resp, err := client.Chat().Create(ctx, req)
    if err != nil {
        log.Fatalf("Error generating chat completion: %v", err)
    }
    
    // Print the assistant's response
    if len(resp.Choices) > 0 {
        fmt.Printf("Assistant: %s\n", resp.Choices[0].Message.Content)
        
        // Add the response to the conversation history
        messages = append(messages, tabby.ChatMessage{
            Role:    tabby.ChatMessageRoleAssistant,
            Content: resp.Choices[0].Message.Content,
        })
        
        // Add a follow-up question
        messages = append(messages, tabby.ChatMessage{
            Role:    tabby.ChatMessageRoleUser,
            Content: "How do I add a route that returns JSON data?",
        })
        
        // Update the request with the new messages
        req.Messages = messages
        
        // Get the follow-up response
        resp, err = client.Chat().Create(ctx, req)
        if err != nil {
            log.Fatalf("Error generating follow-up response: %v", err)
        }
        
        // Print the follow-up response
        if len(resp.Choices) > 0 {
            fmt.Printf("\nFollow-up response:\n%s\n", resp.Choices[0].Message.Content)
        }
    }
}
```

### 3. Generating Embeddings

Convert text into vector embeddings:

```go
func generateEmbeddings(client tabby.Client, ctx context.Context) {
    // Create an embeddings request
    req := &tabby.EmbeddingsRequest{
        Input: []string{
            "The quick brown fox jumps over the lazy dog",
            "Hello world",
        },
    }
    
    // Call the API
    resp, err := client.Embeddings().Create(ctx, req)
    if err != nil {
        log.Fatalf("Error generating embeddings: %v", err)
    }
    
    // Process the embeddings
    fmt.Printf("Generated %d embeddings\n", len(resp.Data))
    
    for i, embedding := range resp.Data {
        // Type assertion to get the embedding vector
        if values, ok := embedding.Embedding.([]float32); ok {
            fmt.Printf("Embedding %d: dimension=%d\n", i+1, len(values))
        }
    }
}
```

### 4. Managing Models

List, load, and unload models:

```go
func manageModels(client tabby.Client, ctx context.Context) {
    // List available models
    models, err := client.Models().List(ctx)
    if err != nil {
        log.Fatalf("Error listing models: %v", err)
    }
    
    fmt.Printf("Available models: %d\n", len(models.Data))
    for i, model := range models.Data {
        fmt.Printf("%d. %s\n", i+1, model.ID)
    }
    
    // Choose a model to load (if any models are available)
    if len(models.Data) > 0 {
        modelName := models.Data[0].ID
        
        // Load the model
        loadReq := &tabby.ModelLoadRequest{
            ModelName: modelName,
            MaxSeqLen: 4096,  // Context window size
        }
        
        fmt.Printf("Loading model: %s\n", modelName)
        loadResp, err := client.Models().Load(ctx, loadReq)
        if err != nil {
            log.Fatalf("Error loading model: %v", err)
        }
        
        fmt.Printf("Model loaded successfully. Status: %s\n", loadResp.Status)
        
        // Get the current model
        current, err := client.Models().Get(ctx)
        if err != nil {
            log.Fatalf("Error getting current model: %v", err)
        }
        
        fmt.Printf("Current model: %s\n", current.ID)
    }
}
```

### 5. Streaming Text Generation

Generate text with streaming responses:

```go
func streamingCompletion(client tabby.Client, ctx context.Context) {
    // Create a streaming completion request
    req := &tabby.CompletionRequest{
        Prompt:      "Write a short story about space exploration:",
        MaxTokens:   200,
        Temperature: 0.7,
        Stream:      true,  // Enable streaming
    }
    
    // Create a stream
    fmt.Println("Starting streaming completion...")
    stream, err := client.Completions().CreateStream(ctx, req)
    if err != nil {
        log.Fatalf("Error creating stream: %v", err)
    }
    defer stream.Close()  // Always close the stream when done
    
    // Process the streaming responses
    fmt.Println("\nGenerated story:")
    for {
        resp, err := stream.Recv()
        if err == io.EOF {
            break  // Stream completed
        }
        if err != nil {
            log.Fatalf("Error receiving from stream: %v", err)
        }
        
        // Print each token as it's generated
        if len(resp.Choices) > 0 {
            fmt.Print(resp.Choices[0].Text)
        }
    }
    fmt.Println("\nStreaming completed")
}
```

## Complete Example

Here's a complete example that demonstrates how to set up a client and perform some common operations:

```go
package main

import (
    "context"
    "fmt"
    "io"
    "log"
    "os"
    "time"
    
    "github.com/pixelsquared/go-tabbyapi/tabby"
)

func main() {
    // Get configuration from environment variables
    endpoint := getEnvOrDefault("TABBY_API_ENDPOINT", "http://localhost:8080")
    apiKey := os.Getenv("TABBY_API_KEY")
    
    // Create a new TabbyAPI client
    client := tabby.NewClient(
        tabby.WithBaseURL(endpoint),
        tabby.WithAPIKey(apiKey),
        tabby.WithTimeout(30*time.Second),
    )
    defer client.Close()
    
    // Create a context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()
    
    // Check the health of the server
    fmt.Println("Checking TabbyAPI server health...")
    health, err := client.Health().Check(ctx)
    if err != nil {
        log.Fatalf("Error checking health: %v", err)
    }
    
    if health.Status != "ok" {
        log.Fatalf("TabbyAPI server is not healthy. Status: %s", health.Status)
    }
    
    fmt.Println("TabbyAPI server is healthy!")
    
    // Check authentication permissions
    authResp, err := client.Auth().GetPermission(ctx)
    if err != nil {
        log.Fatalf("Error checking permissions: %v", err)
    }
    
    fmt.Printf("Current permission level: %s\n", authResp.Permission)
    
    // Generate a completion
    fmt.Println("\n=== Text Completion Example ===")
    req := &tabby.CompletionRequest{
        Prompt:      "Write a function in Go that calculates the fibonacci sequence:",
        MaxTokens:   150,
        Temperature: 0.7,
    }
    
    resp, err := client.Completions().Create(ctx, req)
    if err != nil {
        log.Fatalf("Error generating completion: %v", err)
    }
    
    if len(resp.Choices) > 0 {
        fmt.Println("\nGenerated Code:")
        fmt.Println(resp.Choices[0].Text)
    }
    
    // Generate a streaming completion
    fmt.Println("\n=== Streaming Completion Example ===")
    streamReq := &tabby.CompletionRequest{
        Prompt:      "Explain quantum computing in simple terms:",
        MaxTokens:   200,
        Temperature: 0.7,
        Stream:      true,
    }
    
    stream, err := client.Completions().CreateStream(ctx, streamReq)
    if err != nil {
        log.Fatalf("Error creating stream: %v", err)
    }
    defer stream.Close()
    
    fmt.Println("\nStreaming response:")
    for {
        streamResp, err := stream.Recv()
        if err == io.EOF {
            break  // Stream completed
        }
        if err != nil {
            log.Fatalf("Error receiving from stream: %v", err)
        }
        
        if len(streamResp.Choices) > 0 {
            fmt.Print(streamResp.Choices[0].Text)
        }
    }
    
    fmt.Println("\n\nExecution completed successfully!")
}

// getEnvOrDefault returns the value of the environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}
```

## Error Handling

Implement robust error handling to manage various failure scenarios:

```go
resp, err := client.Completions().Create(ctx, req)
if err != nil {
    // Check for specific error types
    var apiErr *tabby.APIError
    if errors.As(err, &apiErr) {
        switch apiErr.Code() {
        case "authentication_error":
            log.Fatalf("Authentication failed. Check your API key.")
        case "permission_error":
            log.Fatalf("Permission denied. This operation requires higher permissions.")
        case "invalid_request":
            log.Fatalf("Invalid request: %s", apiErr.Error())
        case "not_found":
            log.Fatalf("Resource not found: %s", apiErr.Error())
        case "server_error":
            log.Fatalf("Server error: %s", apiErr.Error())
        default:
            log.Fatalf("API error (%s): %s", apiErr.Code(), apiErr.Error())
        }
    }
    
    // Check for request errors (network, timeout, etc.)
    var reqErr *tabby.RequestError
    if errors.As(err, &reqErr) {
        if errors.Is(reqErr.Unwrap(), context.DeadlineExceeded) {
            log.Fatalf("Request timed out. Try increasing the timeout.")
        } else {
            log.Fatalf("Request error: %v", reqErr)
        }
    }
    
    // Handle other errors
    log.Fatalf("Unknown error: %v", err)
}
```

## Advanced Configuration

### Custom HTTP Client

Configure a custom HTTP client with specific settings:

```go
// Create a custom HTTP client with advanced settings
customHTTPClient := &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 100,
        IdleConnTimeout:     90 * time.Second,
        TLSHandshakeTimeout: 10 * time.Second,
        DialContext: (&net.Dialer{
            Timeout:   30 * time.Second,
            KeepAlive: 30 * time.Second,
        }).DialContext,
    },
    Timeout: 60 * time.Second,
}

// Use the custom HTTP client with the TabbyAPI client
client := tabby.NewClient(
    tabby.WithBaseURL("http://localhost:8080"),
    tabby.WithHTTPClient(customHTTPClient),
    tabby.WithAPIKey("your-api-key"),
)
```

### Custom Retry Policy

Implement a custom retry policy for handling failed requests:

```go
// Define a custom retry policy
customRetryPolicy := &tabby.SimpleRetryPolicy{
    MaxRetryCount: 5,  // Maximum number of retries
    
    // Exponential backoff with jitter
    RetryDelayFunc: func(attempts int) time.Duration {
        delay := time.Duration(1<<uint(attempts-1)) * time.Second
        jitter := time.Duration(rand.Int63n(int64(time.Second)))
        return delay + jitter
    },
    
    // Custom retry conditions
    RetryableFunc: func(resp *http.Response, err error) bool {
        // Retry on network errors
        if err != nil {
            return true
        }
        
        // Retry on server errors and rate limiting
        return resp.StatusCode >= 500 || resp.StatusCode == 429
    },
}

// Use the custom retry policy with the TabbyAPI client
client := tabby.NewClient(
    tabby.WithBaseURL("http://localhost:8080"),
    tabby.WithAPIKey("your-api-key"),
    tabby.WithRetryPolicy(customRetryPolicy),
)
```

## Next Steps

After getting familiar with the basics, explore the detailed documentation for each service:

- [Completions Service](services/completions.md)
- [Chat Service](services/chat.md)
- [Embeddings Service](services/embeddings.md)
- [Models Service](services/models.md)
- [Lora Service](services/lora.md)
- [Templates Service](services/templates.md)
- [Tokens Service](services/tokens.md)
- [Sampling Service](services/sampling.md)
- [Health Service](services/health.md)
- [Auth Service](services/auth.md)