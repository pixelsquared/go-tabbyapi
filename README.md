# Go TabbyAPI Client

A Go client library for interacting with [TabbyAPI](https://github.com/TabbyML/tabby), providing a simple and idiomatic way to use TabbyAPI services from Go applications.

## Overview

This client library provides a comprehensive interface to TabbyAPI's features, including:

- Text completions with support for both basic and streaming responses
- Chat completions with conversation history management
- Embeddings generation
- Model management (loading, unloading, listing)
- LoRA adapter management
- Prompt template management
- Token encoding and decoding
- Sampling parameter management
- Health checks
- Authentication and permission management

## Installation

```bash
go get github.com/pixelsquared/go-tabbyapi
```

## Quick Start

### Creating a Client

```go
import (
    "github.com/pixelsquared/go-tabbyapi/tabby"
)

// Create a client with default options (localhost:8080)
client := tabby.NewClient()

// Or customize with options
client := tabby.NewClient(
    tabby.WithBaseURL("http://your-tabby-server:8080"),
    tabby.WithAPIKey("your-api-key"),
    tabby.WithTimeout(30 * time.Second),
)

// Always close the client when you're done
defer client.Close()
```

### Text Completions

```go
// Basic completion
req := &tabby.CompletionRequest{
    Prompt:      "Once upon a time, there was a programmer who",
    MaxTokens:   100,
    Temperature: 0.7,
}

resp, err := client.Completions().Create(ctx, req)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Println(resp.Choices[0].Text)

// Streaming completion
req.Stream = true
stream, err := client.Completions().CreateStream(ctx, req)
if err != nil {
    log.Fatalf("Error: %v", err)
}
defer stream.Close()

for {
    response, err := stream.Recv()
    if err == io.EOF {
        break // End of stream
    }
    if err != nil {
        log.Fatalf("Stream error: %v", err)
    }
    
    fmt.Print(response.Choices[0].Text) // Print each chunk as it arrives
}
```

### Chat Completions

```go
req := &tabby.ChatCompletionRequest{
    Messages: []tabby.ChatMessage{
        {
            Role:    tabby.ChatMessageRoleSystem,
            Content: "You are a helpful assistant.",
        },
        {
            Role:    tabby.ChatMessageRoleUser,
            Content: "Hello, who are you?",
        },
    },
    MaxTokens:   150,
    Temperature: 0.7,
}

resp, err := client.Chat().Create(ctx, req)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Println(resp.Choices[0].Message.Content)
```

### Embeddings

```go
req := &tabby.EmbeddingsRequest{
    Input: "The quick brown fox jumps over the lazy dog.",
}

resp, err := client.Embeddings().Create(ctx, req)
if err != nil {
    log.Fatalf("Error: %v", err)
}

// Process embeddings
for _, embedding := range resp.Data {
    // Access the embedding vectors
    switch e := embedding.Embedding.(type) {
    case []interface{}:
        // Convert to []float64 for actual use
        floatEmbed := make([]float64, len(e))
        for i, v := range e {
            if f, ok := v.(float64); ok {
                floatEmbed[i] = f
            }
        }
        // Use floatEmbed...
    }
}
```

## Authentication

The client supports three authentication methods:

```go
// API Key authentication
client := tabby.NewClient(tabby.WithAPIKey("your-api-key"))

// Admin Key authentication
client := tabby.NewClient(tabby.WithAdminKey("your-admin-key"))

// Bearer Token authentication
client := tabby.NewClient(tabby.WithBearerToken("your-bearer-token"))
```

## Error Handling

The library provides structured error types that implement the `tabby.Error` interface:

```go
resp, err := client.Completions().Create(ctx, req)
if err != nil {
    switch e := err.(type) {
    case *tabby.APIError:
        fmt.Printf("API Error: %s (Status: %d)\n", e.Message, e.StatusCode)
    case *tabby.ValidationError:
        fmt.Printf("Validation Error: %s (Field: %s)\n", e.Message, e.Field)
    case *tabby.RequestError:
        fmt.Printf("Request Error: %s\n", e.Message)
    default:
        fmt.Printf("Unknown error: %v\n", err)
    }
    return
}
```

## Detailed Documentation

See [docs/api.md](docs/api.md) for detailed API documentation.

## Examples

Check out the [examples](examples) directory for complete, runnable examples:

- [Basic Completions](examples/completions/basic)
- [Streaming Completions](examples/completions/streaming)
- [Chat Completions](examples/chat/basic)
- [Streaming Chat Completions](examples/chat/streaming)
- [Embeddings](examples/embeddings/basic)
- [Model Management](examples/models)

## License

This library is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request