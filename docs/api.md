# TabbyAPI Go Client API Documentation

This document provides detailed information on the Go TabbyAPI client library's API. It covers all available services, methods, parameters, and includes example usage.

## Table of Contents

- [Client](#client)
  - [Creating a Client](#creating-a-client)
  - [Client Options](#client-options)
- [Authentication](#authentication)
- [Services](#services)
  - [Completions](#completions)
  - [Chat](#chat)
  - [Embeddings](#embeddings)
  - [Models](#models)
  - [LoRA Adapters](#lora-adapters)
  - [Templates](#templates)
  - [Tokens](#tokens)
  - [Sampling](#sampling)
  - [Health](#health)
  - [Auth](#auth)
- [Error Handling](#error-handling)
- [Streaming Responses](#streaming-responses)

## Client

The `Client` interface is the main entry point for interacting with TabbyAPI. It provides access to all services and configuration options.

### Creating a Client

```go
import "github.com/pixelsquared/go-tabbyapi/tabby"

// Create a client with default options
client := tabby.NewClient()

// Always close the client when done
defer client.Close()
```

### Client Options

The client supports several configuration options:

```go
client := tabby.NewClient(
    // Set the base URL (default is http://localhost:8080)
    tabby.WithBaseURL("http://your-tabby-server:8080"),
    
    // Set the authentication method
    tabby.WithAPIKey("your-api-key"),
    // OR
    tabby.WithAdminKey("your-admin-key"),
    // OR
    tabby.WithBearerToken("your-bearer-token"),
    
    // Set request timeout (default is 30 seconds)
    tabby.WithTimeout(60 * time.Second),
    
    // Set a custom HTTP client
    tabby.WithHTTPClient(&http.Client{
        Timeout: 45 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns:    10,
            IdleConnTimeout: 30 * time.Second,
        },
    }),
    
    // Set a retry policy
    tabby.WithRetryPolicy(tabby.DefaultRetryPolicy()),
    // OR a custom policy
    tabby.WithRetryPolicy(&tabby.SimpleRetryPolicy{
        MaxRetryCount: 5,
        RetryDelayFunc: func(attempts int) time.Duration {
            return time.Duration(attempts) * time.Second
        },
        RetryableFunc: func(resp *http.Response, err error) bool {
            return err != nil || resp.StatusCode >= 500
        },
    }),
)
```

## Authentication

The client supports three authentication methods:

1. **API Key** - For standard API access:
   ```go
   client := tabby.NewClient(tabby.WithAPIKey("your-api-key"))
   ```

2. **Admin Key** - For administrative operations:
   ```go
   client := tabby.NewClient(tabby.WithAdminKey("your-admin-key"))
   ```

3. **Bearer Token** - For OAuth or JWT authentication:
   ```go
   client := tabby.NewClient(tabby.WithBearerToken("your-bearer-token"))
   ```

## Services

### Completions

The `CompletionsService` handles text completion requests.

#### Methods

##### Create

Generates a completion for the provided prompt.

**Parameters:**
- `ctx` - Context for the request
- `req` - A `*CompletionRequest` object with the following fields:
  - `Prompt` - String (required)
  - `MaxTokens` - Maximum number of tokens to generate (optional)
  - `Temperature` - Controls randomness (0.0-2.0, default varies by model)
  - `TopP` - Top-p sampling (0.0-1.0)
  - `TopK` - Top-k sampling
  - `Stream` - Set to false for non-streaming responses
  - `Stop` - Array of strings to stop generation at
  - `Model` - Model identifier (optional)
  - `JSONSchema` - JSON schema for structured output generation (optional)

**Returns:**
- `*CompletionResponse` - The completion response containing generated text
- `error` - Any error that occurred

**Example:**
```go
req := &tabby.CompletionRequest{
    Prompt:      "Write a function that calculates the factorial of a number",
    MaxTokens:   150,
    Temperature: 0.7,
    TopP:        0.9,
    Stop:        []string{"\n\n"},
}

resp, err := client.Completions().Create(ctx, req)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Println(resp.Choices[0].Text)
```

**Example with JSON schema:**
```go
// Define a JSON schema for structured output
personSchema := map[string]interface{}{
    "type": "object",
    "properties": map[string]interface{}{
        "name": map[string]interface{}{
            "type":        "string",
            "description": "The person's full name",
        },
        "age": map[string]interface{}{
            "type":        "integer",
            "description": "The person's age in years",
            "minimum":     0,
        },
        "email": map[string]interface{}{
            "type":        "string",
            "description": "The person's email address",
            "format":      "email",
        },
    },
    "required": []string{"name", "age", "email"},
}

req := &tabby.CompletionRequest{
    Prompt:      "Generate information about a software developer named John",
    MaxTokens:   150,
    Temperature: 0.7,
    JSONSchema:  personSchema,
}

resp, err := client.Completions().Create(ctx, req)
if err != nil {
    log.Fatalf("Error: %v", err)
}

// The response will contain properly formatted JSON that adheres to the schema
fmt.Println(resp.Choices[0].Text)

// You can parse the JSON directly
var person struct {
    Name  string `json:"name"`
    Age   int    `json:"age"`
    Email string `json:"email"`
}
json.Unmarshal([]byte(resp.Choices[0].Text), &person)
```

##### CreateStream

Generates a streaming completion for the provided prompt.

**Parameters:**
- Same as `Create`, but `Stream` should be set to `true`

**Returns:**
- `CompletionStream` - A stream of completion responses
- `error` - Any error that occurred

**Example:**
```go
req := &tabby.CompletionRequest{
    Prompt:      "Explain quantum computing in simple terms",
    MaxTokens:   300,
    Temperature: 0.8,
    TopP:        0.95,
    Stream:      true,
}

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
    
    fmt.Print(response.Choices[0].Text)
}
```

### Chat

The `ChatService` handles chat completion requests.

#### Methods

##### Create

Generates a chat completion for the provided conversation.

**Parameters:**
- `ctx` - Context for the request
- `req` - A `*ChatCompletionRequest` object with the following fields:
  - `Messages` - Array of `ChatMessage` objects (required)
  - `MaxTokens` - Maximum number of tokens to generate (optional)
  - `Temperature` - Controls randomness (0.0-2.0, default varies by model)
  - `TopP` - Top-p sampling (0.0-1.0)
  - `TopK` - Top-k sampling
  - `Stream` - Set to false for non-streaming responses
  - `Stop` - Array of strings to stop generation at
  - `Model` - Model identifier (optional)
  - `JSONSchema` - JSON schema for structured output generation (optional)

**Returns:**
- `*ChatCompletionResponse` - The chat completion response
- `error` - Any error that occurred

**Example:**
```go
req := &tabby.ChatCompletionRequest{
    Messages: []tabby.ChatMessage{
        {
            Role:    tabby.ChatMessageRoleSystem,
            Content: "You are a helpful assistant specialized in programming.",
        },
        {
            Role:    tabby.ChatMessageRoleUser,
            Content: "How do I implement a binary search tree in Go?",
        },
    },
    MaxTokens:   250,
    Temperature: 0.7,
}

resp, err := client.Chat().Create(ctx, req)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Println(resp.Choices[0].Message.Content)
```

**Example with JSON schema:**
```go
// Define a JSON schema for a product review
reviewSchema := map[string]interface{}{
    "type": "object",
    "properties": map[string]interface{}{
        "product_name": map[string]interface{}{
            "type":        "string",
            "description": "The name of the product being reviewed",
        },
        "rating": map[string]interface{}{
            "type":        "integer",
            "description": "The rating from 1 to 5 stars",
            "minimum":     1,
            "maximum":     5,
        },
        "review_content": map[string]interface{}{
            "type":        "string",
            "description": "The detailed review text",
        },
    },
    "required": []string{"product_name", "rating", "review_content"},
}

req := &tabby.ChatCompletionRequest{
    Messages: []tabby.ChatMessage{
        {
            Role:    tabby.ChatMessageRoleUser,
            Content: "Write a review for the new laptop I bought yesterday.",
        },
    },
    JSONSchema: reviewSchema,
}

resp, err := client.Chat().Create(ctx, req)
if err != nil {
    log.Fatalf("Error: %v", err)
}

// Response will contain a properly formatted JSON review
fmt.Println(resp.Choices[0].Message.Content)

// Parse the JSON directly
var review struct {
    ProductName   string `json:"product_name"`
    Rating        int    `json:"rating"`
    ReviewContent string `json:"review_content"`
}
json.Unmarshal([]byte(resp.Choices[0].Message.Content), &review)
```

##### CreateStream

Generates a streaming chat completion.

**Parameters:**
- Same as `Create`, but `Stream` should be set to `true`

**Returns:**
- `ChatCompletionStream` - A stream of chat completion responses
- `error` - Any error that occurred

**Example:**
```go
req := &tabby.ChatCompletionRequest{
    Messages: []tabby.ChatMessage{
        {
            Role:    tabby.ChatMessageRoleSystem,
            Content: "You are a helpful assistant.",
        },
        {
            Role:    tabby.ChatMessageRoleUser,
            Content: "Tell me a story about a time traveler.",
        },
    },
    MaxTokens:   500,
    Temperature: 0.8,
    Stream:      true,
}

stream, err := client.Chat().CreateStream(ctx, req)
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
    
    if response.Choices[0].Delta.Content != "" {
        fmt.Print(response.Choices[0].Delta.Content)
    }
}
```

### Embeddings

The `EmbeddingsService` handles embedding generation.

#### Methods

##### Create

Generates embeddings for the provided input.

**Parameters:**
- `ctx` - Context for the request
- `req` - A `*EmbeddingsRequest` object with the following fields:
  - `Input` - String or array of strings (required)
  - `Model` - Embedding model identifier (optional)
  - `EncodingFormat` - Output format, e.g., "float" or "base64" (optional)

**Returns:**
- `*EmbeddingsResponse` - The embeddings response
- `error` - Any error that occurred

**Example:**
```go
// Single input
req := &tabby.EmbeddingsRequest{
    Input: "The quick brown fox jumps over the lazy dog.",
}

// Multiple inputs
req := &tabby.EmbeddingsRequest{
    Input: []string{
        "The quick brown fox jumps over the lazy dog.",
        "Machine learning models can process natural language.",
    },
}

resp, err := client.Embeddings().Create(ctx, req)
if err != nil {
    log.Fatalf("Error: %v", err)
}

// Process embeddings
for i, embedding := range resp.Data {
    fmt.Printf("Embedding %d:\n", i)
    
    // Convert the embedding to a usable format
    switch e := embedding.Embedding.(type) {
    case []interface{}:
        // Convert to []float64
        floatEmbed := make([]float64, len(e))
        for i, v := range e {
            if f, ok := v.(float64); ok {
                floatEmbed[i] = f
            }
        }
        
        // Use floatEmbed for your application
        // e.g., calculate cosine similarity between embeddings
    case string:
        // Handle base64 encoded embeddings
        // Decode the base64 string to get the raw bytes
    }
}
```

### Models

The `ModelsService` handles model management.

#### Methods

##### List

Returns all available models.

**Parameters:**
- `ctx` - Context for the request

**Returns:**
- `*ModelList` - List of available models
- `error` - Any error that occurred

**Example:**
```go
models, err := client.Models().List(ctx)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Println("Available models:")
for _, model := range models.Data {
    fmt.Printf("- %s (created: %s)\n", model.ID, 
        time.Unix(model.Created, 0).Format(time.RFC3339))
}
```

##### Get

Returns the currently loaded model.

**Parameters:**
- `ctx` - Context for the request

**Returns:**
- `*ModelCard` - Information about the currently loaded model
- `error` - Any error that occurred

**Example:**
```go
model, err := client.Models().Get(ctx)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Printf("Current model: %s\n", model.ID)
if model.Parameters != nil {
    fmt.Printf("Max sequence length: %d\n", model.Parameters.MaxSeqLen)
    fmt.Printf("Cache size: %d\n", model.Parameters.CacheSize)
}
```

##### Load

Loads a model with the specified parameters.

**Parameters:**
- `ctx` - Context for the request
- `req` - A `*ModelLoadRequest` object with the following fields:
  - `ModelName` - Name of the model to load (required)
  - `MaxSeqLen` - Maximum sequence length (optional)
  - `RopeScale` - RoPE scaling factor (optional)
  - `RopeAlpha` - RoPE alpha value or "auto" (optional)
  - `GPUSplit` - GPU memory split for multiple GPUs (optional)
  - `CacheSize` - KV cache size (optional)
  - `CacheMode` - Cache mode (optional)
  - `ChunkSize` - Chunk size (optional)
  - `PromptTemplate` - Prompt template to use (optional)

**Returns:**
- `*ModelLoadResponse` - Response with loading information
- `error` - Any error that occurred

**Example:**
```go
req := &tabby.ModelLoadRequest{
    ModelName: "TinyLlama-1.1B-Chat-v1.0",
    MaxSeqLen: 2048,
    CacheSize: 2000,
}

resp, err := client.Models().Load(ctx, req)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Printf("Model loading status: %s\n", resp.Status)
```

##### LoadStream

Loads a model and returns a stream of loading progress.

**Parameters:**
- Same as `Load`

**Returns:**
- `ModelLoadStream` - A stream of model loading updates
- `error` - Any error that occurred

**Example:**
```go
req := &tabby.ModelLoadRequest{
    ModelName: "TinyLlama-1.1B-Chat-v1.0",
    MaxSeqLen: 2048,
}

stream, err := client.Models().LoadStream(ctx, req)
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
    
    fmt.Printf("Loading status: %s (Module %d/%d)\n", 
        response.Status, response.Module, response.Modules)
}
```

##### Unload

Unloads the currently loaded model.

**Parameters:**
- `ctx` - Context for the request

**Returns:**
- `error` - Any error that occurred

**Example:**
```go
err := client.Models().Unload(ctx)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Println("Model unloaded successfully")
```

##### GetProps

Returns model properties.

**Parameters:**
- `ctx` - Context for the request

**Returns:**
- `*ModelPropsResponse` - Model properties
- `error` - Any error that occurred

**Example:**
```go
props, err := client.Models().GetProps(ctx)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Printf("Total slots: %d\n", props.TotalSlots)
fmt.Printf("Chat template: %s\n", props.ChatTemplate)
```

##### Download

Downloads a model from HuggingFace.

**Parameters:**
- `ctx` - Context for the request
- `req` - A `*DownloadRequest` object with the following fields:
  - `RepoID` - HuggingFace repository ID (required)
  - `RepoType` - Repository type (optional)
  - `FolderName` - Folder name (optional)
  - `Revision` - Revision (optional)
  - `Token` - HuggingFace token (optional)
  - `Include` - Files to include (optional)
  - `Exclude` - Files to exclude (optional)

**Returns:**
- `*DownloadResponse` - Response with download information
- `error` - Any error that occurred

**Example:**
```go
req := &tabby.DownloadRequest{
    RepoID:   "TinyLlama/TinyLlama-1.1B-Chat-v1.0",
    Revision: "main",
}

resp, err := client.Models().Download(ctx, req)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Printf("Model downloaded to: %s\n", resp.DownloadPath)
```

##### Other Model Methods

The client also provides additional methods for model management:

- **ListDraft** - Lists draft models
- **ListEmbedding** - Lists embedding models
- **GetEmbedding** - Gets the current embedding model
- **LoadEmbedding** - Loads an embedding model
- **UnloadEmbedding** - Unloads the current embedding model

### LoRA Adapters

The `LoraService` handles LoRA adapter management.

#### Methods

##### List

Returns all available LoRA adapters.

**Parameters:**
- `ctx` - Context for the request

**Returns:**
- `*LoraList` - List of available LoRA adapters
- `error` - Any error that occurred

**Example:**
```go
loras, err := client.Lora().List(ctx)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Println("Available LoRA adapters:")
for _, lora := range loras.Data {
    fmt.Printf("- %s (scaling: %.2f)\n", lora.ID, lora.Scaling)
}
```

##### GetActive

Returns the currently loaded LoRA adapters.

**Parameters:**
- `ctx` - Context for the request

**Returns:**
- `*LoraList` - List of currently active LoRA adapters
- `error` - Any error that occurred

**Example:**
```go
active, err := client.Lora().GetActive(ctx)
if err != nil {
    log.Fatalf("Error: %v", err)
}

if len(active.Data) == 0 {
    fmt.Println("No LoRA adapters currently active")
} else {
    fmt.Println("Active LoRA adapters:")
    for _, lora := range active.Data {
        fmt.Printf("- %s (scaling: %.2f)\n", lora.ID, lora.Scaling)
    }
}
```

##### Load

Loads LoRA adapters.

**Parameters:**
- `ctx` - Context for the request
- `req` - A `*LoraLoadRequest` object with the following fields:
  - `Loras` - Array of `LoraLoadInfo` objects with:
    - `Name` - Name of the LoRA adapter
    - `Scaling` - Scaling factor
  - `SkipQueue` - Whether to skip the loading queue

**Returns:**
- `*LoraLoadResponse` - Response with loading results
- `error` - Any error that occurred

**Example:**
```go
req := &tabby.LoraLoadRequest{
    Loras: []tabby.LoraLoadInfo{
        {
            Name:    "alpaca-lora",
            Scaling: 0.8,
        },
    },
}

resp, err := client.Lora().Load(ctx, req)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Printf("Successfully loaded: %v\n", resp.Success)
if len(resp.Failure) > 0 {
    fmt.Printf("Failed to load: %v\n", resp.Failure)
}
```

##### Unload

Unloads all currently loaded LoRA adapters.

**Parameters:**
- `ctx` - Context for the request

**Returns:**
- `error` - Any error that occurred

**Example:**
```go
err := client.Lora().Unload(ctx)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Println("All LoRA adapters unloaded")
```

### Templates

The `TemplatesService` handles prompt template management.

#### Methods

##### List

Returns all available templates.

**Parameters:**
- `ctx` - Context for the request

**Returns:**
- `*TemplateList` - List of available templates
- `error` - Any error that occurred

**Example:**
```go
templates, err := client.Templates().List(ctx)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Println("Available templates:")
for _, template := range templates.Data {
    fmt.Printf("- %s\n", template)
}
```

##### Switch

Changes the active template.

**Parameters:**
- `ctx` - Context for the request
- `req` - A `*TemplateSwitchRequest` object with:
  - `PromptTemplateName` - Name of the template to switch to

**Returns:**
- `error` - Any error that occurred

**Example:**
```go
req := &tabby.TemplateSwitchRequest{
    PromptTemplateName: "llama-2-chat",
}

err := client.Templates().Switch(ctx, req)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Println("Template switched successfully")
```

##### Unload

Unloads the currently selected template.

**Parameters:**
- `ctx` - Context for the request

**Returns:**
- `error` - Any error that occurred

**Example:**
```go
err := client.Templates().Unload(ctx)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Println("Template unloaded successfully")
```

### Tokens

The `TokensService` handles token encoding and decoding.

#### Methods

##### Encode

Encodes text into tokens.

**Parameters:**
- `ctx` - Context for the request
- `req` - A `*TokenEncodeRequest` object with:
  - `Text` - String or `[]ChatMessage` to encode
  - `AddBOSToken` - Whether to add beginning-of-sequence token
  - `EncodeSpecialTokens` - Whether to encode special tokens
  - `DecodeSpecialTokens` - Whether to decode special tokens

**Returns:**
- `*TokenEncodeResponse` - The encoding response
- `error` - Any error that occurred

**Example:**
```go
req := &tabby.TokenEncodeRequest{
    Text: "Hello, world!",
}

resp, err := client.Tokens().Encode(ctx, req)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Printf("Tokens: %v\n", resp.Tokens)
fmt.Printf("Token count: %d\n", resp.Length)
```

##### Decode

Decodes tokens into text.

**Parameters:**
- `ctx` - Context for the request
- `req` - A `*TokenDecodeRequest` object with:
  - `Tokens` - Array of token IDs to decode
  - `AddBOSToken` - Whether to add beginning-of-sequence token
  - `EncodeSpecialTokens` - Whether to encode special tokens
  - `DecodeSpecialTokens` - Whether to decode special tokens

**Returns:**
- `*TokenDecodeResponse` - The decoding response
- `error` - Any error that occurred

**Example:**
```go
req := &tabby.TokenDecodeRequest{
    Tokens: []int{15043, 3186, 29991},
}

resp, err := client.Tokens().Decode(ctx, req)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Printf("Decoded text: %s\n", resp.Text)
```

### Sampling

The `SamplingService` handles sampling parameter overrides.

#### Methods

##### ListOverrides

Returns all sampler overrides.

**Parameters:**
- `ctx` - Context for the request

**Returns:**
- `*SamplerOverrideListResponse` - List of sampler overrides
- `error` - Any error that occurred

**Example:**
```go
overrides, err := client.Sampling().ListOverrides(ctx)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Printf("Selected preset: %s\n", overrides.SelectedPreset)
fmt.Println("Available presets:")
for _, preset := range overrides.Presets {
    fmt.Printf("- %s\n", preset)
}
```

##### SwitchOverride

Changes the active sampler override.

**Parameters:**
- `ctx` - Context for the request
- `req` - A `*SamplerOverrideSwitchRequest` object with:
  - `Preset` - Preset name (optional)
  - `Overrides` - Map of override parameters (optional)

**Returns:**
- `error` - Any error that occurred

**Example:**
```go
// Switch to a preset
req := &tabby.SamplerOverrideSwitchRequest{
    Preset: "creative",
}

// Or use custom overrides
req := &tabby.SamplerOverrideSwitchRequest{
    Overrides: map[string]interface{}{
        "temperature": 0.9,
        "top_p": 0.95,
    },
}

err := client.Sampling().SwitchOverride(ctx, req)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Println("Sampling override applied successfully")
```

##### UnloadOverride

Unloads the currently selected override preset.

**Parameters:**
- `ctx` - Context for the request

**Returns:**
- `error` - Any error that occurred

**Example:**
```go
err := client.Sampling().UnloadOverride(ctx)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Println("Sampling override unloaded successfully")
```

### Health

The `HealthService` handles health checks.

#### Methods

##### Check

Returns the current health status.

**Parameters:**
- `ctx` - Context for the request

**Returns:**
- `*HealthCheckResponse` - The health check response
- `error` - Any error that occurred

**Example:**
```go
health, err := client.Health().Check(ctx)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Printf("Health status: %s\n", health.Status)
if len(health.Issues) > 0 {
    fmt.Println("Issues:")
    for _, issue := range health.Issues {
        fmt.Printf("- %s: %s\n", issue.Time, issue.Description)
    }
}
```

### Auth

The `AuthService` handles authentication permissions.

#### Methods

##### GetPermission

Returns the access level for the current authentication.

**Parameters:**
- `ctx` - Context for the request

**Returns:**
- `*AuthPermissionResponse` - The permission response
- `error` - Any error that occurred

**Example:**
```go
perm, err := client.Auth().GetPermission(ctx)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Printf("Current permission level: %s\n", perm.Permission)
```

## Error Handling

The library provides structured error types that implement the `tabby.Error` interface:

```go
resp, err := client.Completions().Create(ctx, req)
if err != nil {
    // Check error type
    switch e := err.(type) {
    case *tabby.APIError:
        fmt.Printf("API Error: %s (Status: %d)\n", e.Message, e.StatusCode)
        if e.RequestID != "" {
            fmt.Printf("Request ID: %s\n", e.RequestID)
        }
    case *tabby.ValidationError:
        fmt.Printf("Validation Error: %s (Field: %s)\n", e.Message, e.Field)
    case *tabby.RequestError:
        fmt.Printf("Request Error: %s\n", e.Message)
        if e.Err != nil {
            fmt.Printf("Underlying error: %v\n", e.Err)
        }
    case *tabby.StreamError:
        fmt.Printf("Stream Error: %s\n", e.Message)
        if e.Err != nil {
            fmt.Printf("Underlying error: %v\n", e.Err)
        }
    default:
        fmt.Printf("Unknown error: %v\n", err)
    }
    return
}
```

The library also provides predefined error variables:

- `tabby.ErrInvalidRequest` - Invalid request parameters
- `tabby.ErrAuthentication` - Authentication failed
- `tabby.ErrPermission` - Permission denied
- `tabby.ErrNotFound` - Resource not found
- `tabby.ErrServerError` - Server error
- `tabby.ErrTimeout` - Request timed out
- `tabby.ErrCanceled` - Request canceled
- `tabby.ErrStreamClosed` - Stream closed

## Streaming Responses

When working with streaming responses, always close the stream when done:

```go
stream, err := client.Completions().CreateStream(ctx, req)
if err != nil {
    log.Fatalf("Error: %v", err)
}
defer stream.Close() // Important: always close the stream

for {
    response, err := stream.Recv()
    
    // Check for end of stream
    if err == io.EOF {
        break
    }
    
    // Check for stream closed error
    if err == tabby.ErrStreamClosed {
        fmt.Println("Stream was closed")
        break
    }
    
    // Handle other errors
    if err != nil {
        log.Fatalf("Stream error: %v", err)
    }
    
    // Process the response
    // ...
}