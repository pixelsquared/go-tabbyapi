# TabbyAPI Services

This section provides detailed documentation for each of the services available in the go-tabbyapi client library.

## Available Services

- [Completions Service](completions.md): Generate text completions based on prompts
- [Chat Service](chat.md): Handle chat-based interactions with multi-role messages
- [Embeddings Service](embeddings.md): Generate vector embeddings from text inputs
- [Models Service](models.md): Manage AI models (list, load, unload, etc.)
- [Lora Service](lora.md): Manage LoRA adapters for fine-tuning models
- [Templates Service](templates.md): Manage prompt templates for different model types
- [Tokens Service](tokens.md): Handle token encoding and decoding operations
- [Sampling Service](sampling.md): Manage sampling parameters for text generation
- [Health Service](health.md): Check TabbyAPI server health status
- [Auth Service](auth.md): Manage authentication permissions

## Accessing Services

Each service is accessible through the main client interface:

```go
// Create a client
client := tabby.NewClient(
    tabby.WithBaseURL("http://localhost:8080"),
    tabby.WithAPIKey("your-api-key"),
)

// Access a specific service
completionsService := client.Completions()
chatService := client.Chat()
modelsService := client.Models()
// ... and so on
```

## Common Patterns

Most service methods follow these common patterns:

1. **Resource Creation**:
   ```go
   resource, err := service.Create(ctx, request)
   ```

2. **Resource Listing**:
   ```go
   list, err := service.List(ctx)
   ```

3. **Resource Retrieval**:
   ```go
   resource, err := service.Get(ctx)
   ```

4. **Resource Management**:
   ```go
   err := service.Load(ctx, request)
   err := service.Unload(ctx)
   ```

5. **Streaming Operations**:
   ```go
   stream, err := service.CreateStream(ctx, request)
   defer stream.Close()
   
   for {
       chunk, err := stream.Recv()
       if err == io.EOF {
           break
       }
       if err != nil {
           return err
       }
       // Process chunk
   }
   ```

## Context Usage

All service methods accept a `context.Context` as their first parameter, which can be used for:

- **Timeouts**: Set deadlines for API requests
- **Cancellation**: Cancel in-progress requests
- **Request scoping**: Associate values with the request lifecycle

Example with timeout:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

resp, err := client.Completions().Create(ctx, request)
```

## Error Handling

Service methods return strongly typed errors that implement the `tabby.Error` interface:

```go
resp, err := client.Completions().Create(ctx, request)
if err != nil {
    var apiErr *tabby.APIError
    if errors.As(err, &apiErr) {
        // Handle API-specific errors
        log.Printf("API Error (Code: %s, Status: %d): %s", 
            apiErr.Code(), apiErr.HTTPStatusCode(), apiErr.Error())
    } else {
        // Handle other errors
        log.Printf("Error: %v", err)
    }
}