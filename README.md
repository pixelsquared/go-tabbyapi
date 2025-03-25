# go-tabbyapi

[![Go Reference](https://pkg.go.dev/badge/github.com/pixelsquared/go-tabbyapi.svg)](https://pkg.go.dev/github.com/pixelsquared/go-tabbyapi)
[![Go Report Card](https://goreportcard.com/badge/github.com/pixelsquared/go-tabbyapi)](https://goreportcard.com/report/github.com/pixelsquared/go-tabbyapi)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A comprehensive Go client library for interacting with [TabbyAPI](https://github.com/TabbyML/tabby), an open-source self-hosted AI coding assistant.

## Features

- **Complete API Coverage**: Support for all TabbyAPI endpoints including completions, chat, embeddings, model management, and more
- **Streaming Support**: Real-time token streaming for both completions and chat
- **Type Safety**: Strongly typed requests and responses for better reliability
- **Authentication**: Multiple authentication methods (API key, admin key, bearer token)
- **Configurable**: Customizable HTTP client, timeout settings, and retry policies
- **Error Handling**: Structured error types with detailed information
- **Context Support**: All API calls accept a context for cancellation and deadlines

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
        tabby.WithAPIKey("your-api-key"), // Optional
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

## Usage Examples

### Chat Completion

```go
// Create a chat completion
resp, err := client.Chat().Create(ctx, &tabby.ChatCompletionRequest{
    Messages: []tabby.ChatMessage{
        {
            Role:    tabby.ChatMessageRoleSystem,
            Content: "You are a helpful assistant specialized in programming.",
        },
        {
            Role:    tabby.ChatMessageRoleUser,
            Content: "Write a function to calculate the factorial of a number in Python.",
        },
    },
    MaxTokens:   150,
    Temperature: 0.7,
})

if err != nil {
    log.Fatalf("Error generating chat completion: %v", err)
}

fmt.Println(resp.Choices[0].Message.Content)
```

### Streaming

```go
// Create a streaming completion
stream, err := client.Completions().CreateStream(ctx, &tabby.CompletionRequest{
    Prompt:      "Write a function that sorts an array in Go:",
    MaxTokens:   200,
    Temperature: 0.7,
    Stream:      true,
})

if err != nil {
    log.Fatalf("Error creating stream: %v", err)
}
defer stream.Close()

// Process the streaming responses
for {
    resp, err := stream.Recv()
    if err == io.EOF {
        break  // Stream completed
    }
    if err != nil {
        log.Fatalf("Error receiving from stream: %v", err)
    }
    
    // Print each token as it's generated
    fmt.Print(resp.Choices[0].Text)
}
```

### Generating Embeddings

```go
// Generate embeddings for text
resp, err := client.Embeddings().Create(ctx, &tabby.EmbeddingsRequest{
    Input: []string{
        "The quick brown fox jumps over the lazy dog",
        "Hello world",
    },
})

if err != nil {
    log.Fatalf("Error generating embeddings: %v", err)
}

// Process the embeddings
for i, embedding := range resp.Data {
    if values, ok := embedding.Embedding.([]float32); ok {
        fmt.Printf("Embedding %d: dimension=%d\n", i+1, len(values))
    }
}
```

### Model Management

```go
// List available models
models, err := client.Models().List(ctx)
if err != nil {
    log.Fatalf("Error listing models: %v", err)
}

fmt.Printf("Available models: %d\n", len(models.Data))
for _, model := range models.Data {
    fmt.Printf("- %s\n", model.ID)
}

// Load a model
resp, err := client.Models().Load(ctx, &tabby.ModelLoadRequest{
    ModelName: "mistral-7b-v0.1",
    MaxSeqLen: 4096,
})

if err != nil {
    log.Fatalf("Error loading model: %v", err)
}

fmt.Printf("Model loaded: Status=%s\n", resp.Status)
```

## Documentation

For detailed documentation, refer to:

- [Getting Started Guide](docs/getting-started.md)
- [Configuration Guide](docs/configuration.md)
- [Error Handling Guide](docs/error-handling.md)
- [Examples](docs/examples.md)
- [API Reference](docs/services/README.md)
  - [Completions Service](docs/services/completions.md)
  - [Chat Service](docs/services/chat.md)
  - [Embeddings Service](docs/services/embeddings.md)
  - [Models Service](docs/services/models.md)
  - [Lora Service](docs/services/lora.md)
  - [Templates Service](docs/services/templates.md)
  - [Tokens Service](docs/services/tokens.md)
  - [Sampling Service](docs/services/sampling.md)
  - [Health Service](docs/services/health.md)
  - [Auth Service](docs/services/auth.md)

## Available Services

| Service | Description |
|---------|-------------|
| `Completions()` | Generate text completions |
| `Chat()` | Chat-based interactions |
| `Embeddings()` | Generate vector embeddings |
| `Models()` | Model management |
| `Lora()` | LoRA adapter management |
| `Templates()` | Manage prompt templates |
| `Tokens()` | Token encoding and decoding |
| `Sampling()` | Sampling parameter management |
| `Health()` | Health checks |
| `Auth()` | Authentication permissions |

## Running the Examples

The library includes several examples for different use cases. To run them:

```bash
# Set environment variables
export TABBY_API_ENDPOINT="http://localhost:8080"
export TABBY_API_KEY="your-api-key"  # Optional

# Run an example
go run examples/completions/basic/main.go
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [TabbyML](https://github.com/TabbyML) for creating TabbyAPI