# Examples Guide

This guide provides an overview of the examples included in the go-tabbyapi library. Each example demonstrates a specific use case or feature to help you get started quickly.

## Core Examples

The examples are organized by service and functionality:

- Chat Examples
- Completions Examples
- Embeddings Examples
- Models Examples

## Examples Directory Structure

The examples are located in the `/examples` directory with the following structure:

```
examples/
├── README.md
├── chat/
│   ├── README.md
│   ├── basic/
│   │   └── main.go
│   ├── json_schema/
│   │   └── main.go
│   └── streaming/
│       └── main.go
├── completions/
│   ├── README.md
│   ├── basic/
│   │   └── main.go
│   ├── json_schema/
│   │   └── main.go
│   └── streaming/
│       └── main.go
├── embeddings/
│   ├── README.md
│   └── basic/
│       └── main.go
└── models/
    ├── README.md
    ├── list/
    │   └── main.go
    └── streaming/
        └── main.go
```

## Chat Examples

### Basic Chat

Demonstrates a simple chat conversation with the model, including continuing the conversation with follow-up questions.

**Location:** [examples/chat/basic/main.go](../examples/chat/basic/main.go)

**Key Features:**
- Creating a chat completion with system and user messages
- Processing the assistant's response
- Continuing the conversation by adding the response to the history
- Adding a follow-up question and getting another response

**Usage:**
```bash
go run examples/chat/basic/main.go
```

### JSON Schema Chat

Shows how to generate structured JSON responses from the chat model using a JSON schema.

**Location:** [examples/chat/json_schema/main.go](../examples/chat/json_schema/main.go)

**Key Features:**
- Defining a JSON schema to structure the model's output
- Creating a chat completion with the schema
- Parsing and using the structured JSON response

**Usage:**
```bash
go run examples/chat/json_schema/main.go
```

### Streaming Chat

Demonstrates how to use streaming for chat completions to get tokens as they're generated.

**Location:** [examples/chat/streaming/main.go](../examples/chat/streaming/main.go)

**Key Features:**
- Creating a streaming chat completion
- Processing tokens incrementally as they arrive
- Handling stream completion and errors

**Usage:**
```bash
go run examples/chat/streaming/main.go
```

## Completions Examples

### Basic Completions

Shows how to generate text completions from a prompt.

**Location:** [examples/completions/basic/main.go](../examples/completions/basic/main.go)

**Key Features:**
- Creating a simple text completion
- Configuring generation parameters (temperature, max tokens, etc.)
- Processing the generated text

**Usage:**
```bash
go run examples/completions/basic/main.go
```

### JSON Schema Completions

Demonstrates generating structured JSON output from completions using a JSON schema.

**Location:** [examples/completions/json_schema/main.go](../examples/completions/json_schema/main.go)

**Key Features:**
- Defining a JSON schema to structure model output
- Generating structured text that conforms to the schema
- Parsing and using the structured JSON response

**Usage:**
```bash
go run examples/completions/json_schema/main.go
```

### Streaming Completions

Shows how to use streaming to get completion tokens incrementally as they're generated.

**Location:** [examples/completions/streaming/main.go](../examples/completions/streaming/main.go)

**Key Features:**
- Creating a streaming completion request
- Processing tokens as they arrive in real-time
- Building the complete response incrementally

**Usage:**
```bash
go run examples/completions/streaming/main.go
```

## Embeddings Examples

### Basic Embeddings

Demonstrates how to generate vector embeddings from text inputs.

**Location:** [examples/embeddings/basic/main.go](../examples/embeddings/basic/main.go)

**Key Features:**
- Generating embeddings for single or multiple texts
- Processing and using the embedding vectors
- Understanding embedding dimensions and usage

**Usage:**
```bash
go run examples/embeddings/basic/main.go
```

## Models Examples

### List Models

Shows how to list and get information about available models.

**Location:** [examples/models/list/main.go](../examples/models/list/main.go)

**Key Features:**
- Listing all available models
- Getting information about specific models
- Retrieving model properties and parameters

**Usage:**
```bash
go run examples/models/list/main.go
```

### Streaming Model Loading

Demonstrates how to load a model with streaming progress updates.

**Location:** [examples/models/streaming/main.go](../examples/models/streaming/main.go)

**Key Features:**
- Loading a model with custom parameters
- Streaming progress updates during loading
- Handling loading completion and errors

**Usage:**
```bash
go run examples/models/streaming/main.go
```

## Running Examples

To run any of these examples, make sure you have:

1. A running TabbyAPI server (local or remote)
2. Go 1.18 or higher installed
3. The go-tabbyapi library installed

Then set the required environment variables:

```bash
export TABBY_API_ENDPOINT="http://localhost:8080"  # Your TabbyAPI server URL
export TABBY_API_KEY="your-api-key"                # Optional, if authentication is required
```

And run the example:

```bash
go run examples/path/to/example/main.go
```

## Creating Your Own Examples

When creating your own applications, you can use these examples as starting points. Here's a template to help you get started:

```go
package main

import (
    "context"
    "fmt"
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
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // YOUR CODE HERE
    // Use the client to interact with TabbyAPI
    
    fmt.Println("Example completed successfully!")
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

## Additional Resources

For more comprehensive information on using the TabbyAPI client library, refer to:

- [Getting Started Guide](getting-started.md)
- [Configuration Guide](configuration.md)
- [Error Handling Guide](error-handling.md)
- [API Reference](services/README.md)