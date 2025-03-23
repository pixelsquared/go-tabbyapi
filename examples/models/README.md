# TabbyAPI Go Client - Models Examples

This directory contains examples demonstrating how to use the Go client library for model management with TabbyAPI.

## Examples

### Model Listing and Information

Located in the `list` directory, this example shows how to:
- List all available models on the TabbyAPI server
- Get information about the currently loaded model
- Display model parameters and properties
- List embedding models

```bash
cd list
go run main.go
```

### Model Loading with Streaming Progress

Located in the `streaming` directory, this example demonstrates:
- How to initiate model loading with specific parameters
- Stream real-time progress updates during model loading
- Track loading progress and timing information
- Verify successful model loading and test the model

```bash
cd streaming
go run main.go
```

## Administrative Privileges

Note that model management typically requires administrative privileges. These examples use the admin key for authentication:

```go
client := tabby.NewClient(
    tabby.WithBaseURL(endpoint),
    tabby.WithAdminKey(adminKey),
    tabby.WithTimeout(60*time.Second),
)
```

## Environment Variables

These examples use the following environment variables:

- `TABBY_API_ENDPOINT` - The URL of your TabbyAPI server (defaults to `http://localhost:8080`)
- `TABBY_ADMIN_KEY` - Your admin key for authentication (required for model management)
- `TABBY_MODEL_NAME` - The name of the model to load (for the streaming example)

## Model Management Operations

The examples demonstrate these common operations:

- **List**: Retrieve all available models
- **Get**: Get the currently loaded model
- **Load**: Load a model with specific parameters
- **LoadStream**: Load a model with streaming progress updates
- **Unload**: Unload the currently loaded model
- **GetProps**: Get properties of the loaded model
- **ListEmbedding**: List available embedding models

## Model Loading Parameters

When loading a model, you can specify various parameters:

- `ModelName`: Name of the model to load
- `MaxSeqLen`: Maximum sequence length (context window)
- `RopeScale`: RoPE scaling factor for extending context
- `CacheSize`: Size of the KV cache in MB
- `ChunkSize`: Chunk size for efficient processing
- `PromptTemplate`: Custom prompt template to use

## Safety Considerations

Model loading operations can impact server performance and resources. Exercise caution when:

1. Loading large models on systems with limited resources
2. Unloading models that might be in use by other clients
3. Setting parameters that might affect stability

The examples include some operations (like model unloading) commented out to prevent accidental execution.

## Further Customization

You can modify these examples to:
- Implement model rotation strategies
- Create model management automation scripts
- Build monitoring tools for model usage
- Develop workflows that use different models for different tasks