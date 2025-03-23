# Go TabbyAPI Client Library Architecture

This document outlines the architecture for a Go client library for TabbyAPI, focusing on complete feature coverage with idiomatic Go implementations.

## 1. Design Principles

The library will adhere to the following principles:

1. **Complete API Coverage**: Implement all TabbyAPI endpoints and features
2. **Idiomatic Go**: Follow Go conventions and best practices
3. **Type Safety**: Leverage Go's type system for compile-time safety
4. **Context Awareness**: Support context-based cancellation and timeouts
5. **Clear Error Handling**: Provide detailed and structured error information
6. **Minimal Dependencies**: Rely primarily on the standard library
7. **Efficient Streaming**: First-class support for server-sent events (SSE)
8. **Comprehensive Documentation**: Clear godoc and examples

## 2. Package Structure

```
go-tabbyapi/
├── tabby/             # Core package
│   ├── client.go      # Main client and interfaces
│   ├── options.go     # Client configuration
│   ├── types.go       # Shared types
│   └── errors.go      # Error definitions
├── internal/          # Internal implementation details
│   ├── rest/          # REST client implementation
│   ├── stream/        # Streaming implementation
│   └── auth/          # Authentication utilities
├── models/            # Types and services for models
├── completions/       # Types and services for completions
├── chat/              # Types and services for chat completions
├── embeddings/        # Types and services for embeddings
├── tokens/            # Types and services for tokenization
├── templates/         # Types and services for templates
├── lora/              # Types and services for LoRA adapters
├── sampling/          # Types and services for sampling parameters
├── examples/          # Usage examples
│   ├── chat/          # Chat completion examples
│   ├── completions/   # Text completion examples
│   ├── embeddings/    # Embedding examples
│   ├── models/        # Model management examples
│   └── streaming/     # Streaming examples
└── docs/              # Documentation
    └── api.md         # API documentation
```

## 3. Core Interfaces

### 3.1 Client Interface

```go
// Client provides access to the TabbyAPI
type Client interface {
    // API Services
    Completions() CompletionsService
    Chat() ChatService
    Embeddings() EmbeddingsService
    Models() ModelsService
    Lora() LoraService
    Templates() TemplatesService
    Tokens() TokensService
    Sampling() SamplingService
    Health() HealthService
    Auth() AuthService
    
    // Close releases resources used by the client
    Close() error
    
    // Client options
    WithBaseURL(url string) Client
    WithHTTPClient(client *http.Client) Client
    WithAPIKey(key string) Client
    WithAdminKey(key string) Client
    WithBearerToken(token string) Client
    WithTimeout(timeout time.Duration) Client
    WithRetryPolicy(policy RetryPolicy) Client
}
```

### 3.2 Service Interfaces

Each API area will have a dedicated service interface:

```go
// CompletionsService handles text completion requests
type CompletionsService interface {
    // Create generates a completion for the provided request
    Create(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)
    
    // CreateStream generates a streaming completion
    CreateStream(ctx context.Context, req *CompletionRequest) (CompletionStream, error)
}

// ChatService handles chat completion requests
type ChatService interface {
    // Create generates a chat completion for the provided request
    Create(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error)
    
    // CreateStream generates a streaming chat completion
    CreateStream(ctx context.Context, req *ChatCompletionRequest) (ChatCompletionStream, error)
}

// ModelsService handles model management
type ModelsService interface {
    // List returns all available models
    List(ctx context.Context) (*ModelList, error)
    
    // Get returns the currently loaded model
    Get(ctx context.Context) (*ModelCard, error)
    
    // Load loads a model with the specified parameters
    Load(ctx context.Context, req *ModelLoadRequest) (*ModelLoadResponse, error)
    
    // LoadStream loads a model and returns a stream of loading progress
    LoadStream(ctx context.Context, req *ModelLoadRequest) (ModelLoadStream, error)
    
    // Unload unloads the currently loaded model
    Unload(ctx context.Context) error
    
    // GetProps returns model properties
    GetProps(ctx context.Context) (*ModelPropsResponse, error)
    
    // Download downloads a model from HuggingFace
    Download(ctx context.Context, req *DownloadRequest) (*DownloadResponse, error)
    
    // ListDraft returns all available draft models
    ListDraft(ctx context.Context) (*ModelList, error)
    
    // ListEmbedding returns all available embedding models
    ListEmbedding(ctx context.Context) (*ModelList, error)
    
    // GetEmbedding returns the currently loaded embedding model
    GetEmbedding(ctx context.Context) (*ModelCard, error)
    
    // LoadEmbedding loads an embedding model
    LoadEmbedding(ctx context.Context, req *EmbeddingModelLoadRequest) (*ModelLoadResponse, error)
    
    // UnloadEmbedding unloads the current embedding model
    UnloadEmbedding(ctx context.Context) error
}

// EmbeddingsService handles embedding generation
type EmbeddingsService interface {
    // Create generates embeddings for the provided input
    Create(ctx context.Context, req *EmbeddingsRequest) (*EmbeddingsResponse, error)
}

// LoraService handles LoRA adapter management
type LoraService interface {
    // List returns all available LoRA adapters
    List(ctx context.Context) (*LoraList, error)
    
    // GetActive returns the currently loaded LoRA adapters
    GetActive(ctx context.Context) (*LoraList, error)
    
    // Load loads LoRA adapters
    Load(ctx context.Context, req *LoraLoadRequest) (*LoraLoadResponse, error)
    
    // Unload unloads all currently loaded LoRA adapters
    Unload(ctx context.Context) error
}

// TokensService handles token encoding and decoding
type TokensService interface {
    // Encode encodes text into tokens
    Encode(ctx context.Context, req *TokenEncodeRequest) (*TokenEncodeResponse, error)
    
    // Decode decodes tokens into text
    Decode(ctx context.Context, req *TokenDecodeRequest) (*TokenDecodeResponse, error)
}

// TemplatesService handles prompt template management
type TemplatesService interface {
    // List returns all available templates
    List(ctx context.Context) (*TemplateList, error)
    
    // Switch changes the active template
    Switch(ctx context.Context, req *TemplateSwitchRequest) error
    
    // Unload unloads the currently selected template
    Unload(ctx context.Context) error
}

// SamplingService handles sampling parameter overrides
type SamplingService interface {
    // ListOverrides returns all sampler overrides
    ListOverrides(ctx context.Context) (*SamplerOverrideListResponse, error)
    
    // SwitchOverride changes the active sampler override
    SwitchOverride(ctx context.Context, req *SamplerOverrideSwitchRequest) error
    
    // UnloadOverride unloads the currently selected override preset
    UnloadOverride(ctx context.Context) error
}

// HealthService handles health checks
type HealthService interface {
    // Check returns the current health status
    Check(ctx context.Context) (*HealthCheckResponse, error)
}

// AuthService handles authentication permissions
type AuthService interface {
    // GetPermission returns the access level for the current authentication
    GetPermission(ctx context.Context) (*AuthPermissionResponse, error)
}
```

### 3.3 Stream Interfaces

```go
// Stream is a generic interface for SSE streams
type Stream[T any] interface {
    // Recv returns the next item from the stream
    Recv() (T, error)
    
    // Close releases resources associated with the stream
    Close() error
}

// CompletionStream is a stream of completion responses
type CompletionStream = Stream[*CompletionStreamResponse]

// ChatCompletionStream is a stream of chat completion responses
type ChatCompletionStream = Stream[*ChatCompletionStreamResponse]

// ModelLoadStream is a stream of model loading progress updates
type ModelLoadStream = Stream[*ModelLoadResponse]
```

## 4. Type Definitions

The library will include comprehensive type definitions matching the TabbyAPI schemas:

```go
// CompletionRequest matches the TabbyAPI completion request schema
type CompletionRequest struct {
    Prompt           interface{} `json:"prompt"`             // String or array of strings
    MaxTokens        int         `json:"max_tokens,omitempty"`
    Temperature      float64     `json:"temperature,omitempty"`
    TopP             float64     `json:"top_p,omitempty"`
    TopK             int         `json:"top_k,omitempty"`
    Stream           bool        `json:"stream,omitempty"`
    Stop             interface{} `json:"stop,omitempty"`     // String or array of strings
    Model            string      `json:"model,omitempty"`
    // Additional parameters from the API specification
}

// ChatCompletionRequest matches the TabbyAPI chat completion request schema
type ChatCompletionRequest struct {
    Messages         []ChatMessage `json:"messages"`
    MaxTokens        int           `json:"max_tokens,omitempty"`
    Temperature      float64       `json:"temperature,omitempty"`
    TopP             float64       `json:"top_p,omitempty"`
    TopK             int           `json:"top_k,omitempty"`
    Stream           bool          `json:"stream,omitempty"`
    Stop             interface{}   `json:"stop,omitempty"`
    Model            string        `json:"model,omitempty"`
    // Additional parameters from the API specification
}

// All additional types will be similarly defined to match the API specification

// Helper types for improved ergonomics
type ChatMessageRole string

const (
    ChatMessageRoleUser      ChatMessageRole = "user"
    ChatMessageRoleAssistant ChatMessageRole = "assistant"
    ChatMessageRoleSystem    ChatMessageRole = "system"
    ChatMessageRoleTool      ChatMessageRole = "tool"
)

// ChatMessage represents a message in a chat completion request
type ChatMessage struct {
    Role    ChatMessageRole `json:"role"`
    Content interface{}     `json:"content"`      // String or array of ChatMessageContent
    // Additional fields from the API specification
}
```

## 5. Authentication

The library will support all authentication methods in TabbyAPI:

```go
// Authenticator provides authentication for API requests
type Authenticator interface {
    // Apply adds authentication to the provided request
    Apply(req *http.Request)
}

// Three concrete authenticator implementations:

// APIKeyAuthenticator uses the X-API-Key header
type APIKeyAuthenticator struct {
    Key string
}

// AdminKeyAuthenticator uses the X-Admin-Key header
type AdminKeyAuthenticator struct {
    Key string
}

// BearerTokenAuthenticator uses the Authorization header
type BearerTokenAuthenticator struct {
    Token string
}
```

## 6. Error Handling

The library will provide structured error types:

```go
// Error is the interface implemented by all errors in the library
type Error interface {
    error
    Code() string
    HTTPStatusCode() int
}

// APIError represents an error returned by the TabbyAPI
type APIError struct {
    StatusCode int
    Message    string
    Details    interface{}
    RequestID  string
}

// ValidationError represents a validation error
type ValidationError struct {
    Field   string
    Message string
    Type    string
}

// Predefined error types
var (
    ErrInvalidRequest = errors.New("invalid request parameters")
    ErrAuthentication = errors.New("authentication failed")
    ErrPermission     = errors.New("permission denied")
    ErrNotFound       = errors.New("resource not found")
    ErrServerError    = errors.New("server error")
    ErrTimeout        = errors.New("request timed out")
    ErrCanceled       = errors.New("request canceled")
    ErrStreamClosed   = errors.New("stream closed")
)
```

## 7. Configuration and Options

The library will use a functional options pattern:

```go
// Option configures a Client
type Option func(*clientImpl)

// WithBaseURL sets the base URL for the API
func WithBaseURL(url string) Option {
    return func(c *clientImpl) {
        c.baseURL = url
    }
}

// WithHTTPClient sets the HTTP client
func WithHTTPClient(client *http.Client) Option {
    return func(c *clientImpl) {
        c.httpClient = client
    }
}

// WithAPIKey sets the API key for authentication
func WithAPIKey(key string) Option {
    return func(c *clientImpl) {
        c.auth = &APIKeyAuthenticator{Key: key}
    }
}

// Additional options...
```

## 8. Client Implementation

The core client implementation:

```go
// NewClient creates a new TabbyAPI client
func NewClient(options ...Option) Client {
    c := &clientImpl{
        baseURL:    "http://localhost:8080",
        httpClient: &http.Client{Timeout: 30 * time.Second},
    }
    
    for _, opt := range options {
        opt(c)
    }
    
    return c
}

// clientImpl is the concrete implementation of the Client interface
type clientImpl struct {
    baseURL    string
    httpClient *http.Client
    auth       Authenticator
    retryPolicy RetryPolicy
}

// Implementations of service methods...
```

## 9. Streaming Implementation

```go
// streamImpl implements the Stream interface
type streamImpl[T any] struct {
    ctx      context.Context
    cancel   context.CancelFunc
    response *http.Response
    reader   *bufio.Reader
    closed   bool
    mu       sync.Mutex
}

// NewStream creates a new stream from an HTTP response
func NewStream[T any](ctx context.Context, resp *http.Response) Stream[T] {
    ctx, cancel := context.WithCancel(ctx)
    return &streamImpl[T]{
        ctx:      ctx,
        cancel:   cancel,
        response: resp,
        reader:   bufio.NewReader(resp.Body),
    }
}

// Recv reads the next item from the stream
func (s *streamImpl[T]) Recv() (T, error) {
    // Implementation details...
}

// Close closes the stream
func (s *streamImpl[T]) Close() error {
    // Implementation details...
}
```

## 10. Endpoint to Method Mapping

To ensure complete coverage, here's a mapping from TabbyAPI endpoints to client methods:

| TabbyAPI Endpoint | HTTP Method | Client Method |
|-------------------|-------------|---------------|
| `/v1/completions` | POST | `client.Completions().Create()` |
| `/v1/chat/completions` | POST | `client.Chat().Create()` |
| `/v1/embeddings` | POST | `client.Embeddings().Create()` |
| `/health` | GET | `client.Health().Check()` |
| `/.well-known/serviceinfo` | GET | `client.ServiceInfo().Get()` |
| `/v1/model/list` | GET | `client.Models().List()` |
| `/v1/models` | GET | `client.Models().List()` |
| `/v1/model` | GET | `client.Models().Get()` |
| `/props` | GET | `client.Models().GetProps()` |
| `/v1/model/draft/list` | GET | `client.Models().ListDraft()` |
| `/v1/model/load` | POST | `client.Models().Load()` |
| `/v1/model/unload` | POST | `client.Models().Unload()` |
| `/v1/download` | POST | `client.Models().Download()` |
| `/v1/lora/list` | GET | `client.Lora().List()` |
| `/v1/loras` | GET | `client.Lora().List()` |
| `/v1/lora` | GET | `client.Lora().GetActive()` |
| `/v1/lora/load` | POST | `client.Lora().Load()` |
| `/v1/lora/unload` | POST | `client.Lora().Unload()` |
| `/v1/model/embedding/list` | GET | `client.Models().ListEmbedding()` |
| `/v1/model/embedding` | GET | `client.Models().GetEmbedding()` |
| `/v1/model/embedding/load` | POST | `client.Models().LoadEmbedding()` |
| `/v1/model/embedding/unload` | POST | `client.Models().UnloadEmbedding()` |
| `/v1/token/encode` | POST | `client.Tokens().Encode()` |
| `/v1/token/decode` | POST | `client.Tokens().Decode()` |
| `/v1/auth/permission` | GET | `client.Auth().GetPermission()` |
| `/v1/template/list` | GET | `client.Templates().List()` |
| `/v1/templates` | GET | `client.Templates().List()` |
| `/v1/template/switch` | POST | `client.Templates().Switch()` |
| `/v1/template/unload` | POST | `client.Templates().Unload()` |
| `/v1/sampling/override/list` | GET | `client.Sampling().ListOverrides()` |
| `/v1/sampling/overrides` | GET | `client.Sampling().ListOverrides()` |
| `/v1/sampling/override/switch` | POST | `client.Sampling().SwitchOverride()` |
| `/v1/sampling/override/unload` | POST | `client.Sampling().UnloadOverride()` |

## 11. Example Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/yourusername/go-tabbyapi/tabby"
    "github.com/yourusername/go-tabbyapi/chat"
)

func main() {
    // Create a client with API key authentication
    client := tabby.NewClient(
        tabby.WithBaseURL("http://localhost:8080"),
        tabby.WithAPIKey(os.Getenv("TABBY_API_KEY")),
        tabby.WithTimeout(30 * time.Second),
    )
    
    // Create a context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Check the health of the service
    health, err := client.Health().Check(ctx)
    if err != nil {
        log.Fatalf("Health check failed: %v", err)
    }
    fmt.Printf("Service health: %s\n", health.Status)
    
    // Create a chat completion request
    req := &chat.CompletionRequest{
        Messages: []chat.Message{
            {
                Role:    chat.ChatMessageRoleUser,
                Content: "Write a short poem about Go programming.",
            },
        },
        MaxTokens:   150,
        Temperature: 0.7,
    }
    
    // Send the request
    resp, err := client.Chat().Create(ctx, req)
    if err != nil {
        log.Fatalf("Chat completion failed: %v", err)
    }
    
    // Print the response
    fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
    fmt.Printf("Usage: %d tokens\n", resp.Usage.TotalTokens)
}
```

## 12. Streaming Example

```go
package main

import (
    "context"
    "fmt"
    "io"
    "log"
    "os"
    "time"

    "github.com/yourusername/go-tabbyapi/tabby"
    "github.com/yourusername/go-tabbyapi/chat"
)

func main() {
    // Create a client with API key authentication
    client := tabby.NewClient(
        tabby.WithBaseURL("http://localhost:8080"),
        tabby.WithAPIKey(os.Getenv("TABBY_API_KEY")),
    )
    
    // Create a context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Create a streaming chat completion request
    req := &chat.CompletionRequest{
        Messages: []chat.Message{
            {
                Role:    chat.ChatMessageRoleUser,
                Content: "Write a short poem about Go programming.",
            },
        },
        MaxTokens:   150,
        Temperature: 0.7,
        Stream:      true,
    }
    
    // Send the streaming request
    stream, err := client.Chat().CreateStream(ctx, req)
    if err != nil {
        log.Fatalf("Chat completion failed: %v", err)
    }
    defer stream.Close()
    
    // Process the stream
    for {
        resp, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatalf("Stream error: %v", err)
        }
        
        // Print the partial response
        if len(resp.Choices) > 0 && resp.Choices[0].Delta.Content != "" {
            fmt.Print(resp.Choices[0].Delta.Content)
        }
    }
    fmt.Println()
}
```

## 13. Architecture Diagram

```mermaid
classDiagram
    class Client {
        +Completions() CompletionsService
        +Chat() ChatService
        +Embeddings() EmbeddingsService
        +Models() ModelsService
        +Lora() LoraService
        +Templates() TemplatesService
        +Tokens() TokensService
        +Sampling() SamplingService
        +Health() HealthService
        +Auth() AuthService
        +Close() error
        +WithBaseURL(string) Client
        +WithHTTPClient(*http.Client) Client
        +WithAPIKey(string) Client
        +WithAdminKey(string) Client
        +WithBearerToken(string) Client
        +WithTimeout(time.Duration) Client
        +WithRetryPolicy(RetryPolicy) Client
    }
    
    class CompletionsService {
        +Create(ctx, *CompletionRequest) (*CompletionResponse, error)
        +CreateStream(ctx, *CompletionRequest) (CompletionStream, error)
    }
    
    class ChatService {
        +Create(ctx, *ChatCompletionRequest) (*ChatCompletionResponse, error)
        +CreateStream(ctx, *ChatCompletionRequest) (ChatCompletionStream, error)
    }
    
    class ModelsService {
        +List(ctx) (*ModelList, error)
        +Get(ctx) (*ModelCard, error)
        +Load(ctx, *ModelLoadRequest) (*ModelLoadResponse, error)
        +LoadStream(ctx, *ModelLoadRequest) (ModelLoadStream, error)
        +Unload(ctx) error
        +GetProps(ctx) (*ModelPropsResponse, error)
        +Download(ctx, *DownloadRequest) (*DownloadResponse, error)
        +ListDraft(ctx) (*ModelList, error)
        +ListEmbedding(ctx) (*ModelList, error)
        +GetEmbedding(ctx) (*ModelCard, error)
        +LoadEmbedding(ctx, *EmbeddingModelLoadRequest) (*ModelLoadResponse, error)
        +UnloadEmbedding(ctx) error
    }
    
    class Stream~T~ {
        +Recv() (T, error)
        +Close() error
    }
    
    class Authenticator {
        +Apply(req *http.Request)
    }
    
    class APIKeyAuthenticator {
        +Apply(req *http.Request)
    }
    
    class AdminKeyAuthenticator {
        +Apply(req *http.Request)
    }
    
    class BearerTokenAuthenticator {
        +Apply(req *http.Request)
    }
    
    class Error {
        +Error() string
        +Code() string
        +HTTPStatusCode() int
    }
    
    Client --> CompletionsService
    Client --> ChatService
    Client --> ModelsService
    Client --> "Other Services..."
    Client --> Authenticator
    
    Authenticator <|-- APIKeyAuthenticator
    Authenticator <|-- AdminKeyAuthenticator
    Authenticator <|-- BearerTokenAuthenticator
    
    CompletionsService --> Stream
    ChatService --> Stream
    ModelsService --> Stream
```

## 14. Testing Strategy

The library will include comprehensive tests:

1. **Unit Tests**: For core functionality and individual components
2. **Integration Tests**: For interactions between components
3. **Mock Tests**: Using `httptest` to simulate API responses
4. **End-to-End Tests**: Real API calls (configurable for CI/CD)
5. **Stream Tests**: Special focus on streaming functionality
6. **Example Tests**: Executable examples as tests

## 15. Documentation Strategy

The library will include:

1. **Godoc**: Complete documentation for all exported items
2. **Examples**: Well-documented examples for common use cases
3. **README**: Quick start guide and basic usage
4. **Wiki**: Extended documentation and best practices
5. **Diagrams**: Visual representation of the architecture

## 16. Versioning and Compatibility

The library will follow semantic versioning (SemVer):

- **Major Version**: Breaking changes to the API
- **Minor Version**: New features in a backward-compatible manner
- **Patch Version**: Backward-compatible bug fixes

## Conclusion

This architecture document outlines a comprehensive Go client library for TabbyAPI with complete feature coverage and idiomatic Go implementations. The design prioritizes clear interfaces, efficient streaming, robust error handling, and complete documentation, ensuring that Go developers can easily and effectively work with TabbyAPI.