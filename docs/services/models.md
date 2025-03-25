# Models Service

The Models service provides functionality for managing AI models in TabbyAPI. It enables operations such as listing available models, loading and unloading models, retrieving model information, and managing embedding models.

## Interface

```go
// ModelsService handles model management operations including listing, loading,
// unloading, and querying information about models.
type ModelsService interface {
	// List returns all available models.
	List(ctx context.Context) (*ModelList, error)

	// Get returns the currently loaded model.
	Get(ctx context.Context) (*ModelCard, error)

	// Load loads a model with the specified parameters.
	Load(ctx context.Context, req *ModelLoadRequest) (*ModelLoadResponse, error)

	// LoadStream loads a model and returns a stream of loading progress.
	LoadStream(ctx context.Context, req *ModelLoadRequest) (ModelLoadStream, error)

	// Unload unloads the currently loaded model.
	Unload(ctx context.Context) error

	// GetProps returns model properties.
	GetProps(ctx context.Context) (*ModelPropsResponse, error)

	// Download downloads a model from HuggingFace.
	Download(ctx context.Context, req *DownloadRequest) (*DownloadResponse, error)

	// ListDraft returns all available draft models.
	ListDraft(ctx context.Context) (*ModelList, error)

	// ListEmbedding returns all available embedding models.
	ListEmbedding(ctx context.Context) (*ModelList, error)

	// GetEmbedding returns the currently loaded embedding model.
	GetEmbedding(ctx context.Context) (*ModelCard, error)

	// LoadEmbedding loads an embedding model.
	LoadEmbedding(ctx context.Context, req *EmbeddingModelLoadRequest) (*ModelLoadResponse, error)

	// UnloadEmbedding unloads the current embedding model.
	UnloadEmbedding(ctx context.Context) error
}
```

## Model Types and Information

### ModelCard

The `ModelCard` struct provides information about a model:

```go
type ModelCard struct {
	ID         string               `json:"id"`         // Unique identifier for the model
	Object     string               `json:"object"`     // Type of object
	Created    int64                `json:"created"`    // Unix timestamp of creation
	OwnedBy    string               `json:"owned_by"`   // Owner of the model
	Parameters *ModelCardParameters `json:"parameters,omitempty"` // Model parameters
}

type ModelCardParameters struct {
	MaxSeqLen      int     `json:"max_seq_len,omitempty"`      // Maximum sequence length
	RopeScale      float64 `json:"rope_scale,omitempty"`       // RoPE scaling factor
	RopeAlpha      float64 `json:"rope_alpha,omitempty"`       // RoPE alpha parameter
	MaxBatchSize   int     `json:"max_batch_size,omitempty"`   // Maximum batch size
	CacheSize      int     `json:"cache_size,omitempty"`       // KV cache size
	CacheMode      string  `json:"cache_mode,omitempty"`       // Cache mode
	ChunkSize      int     `json:"chunk_size,omitempty"`       // Chunk size
	PromptTemplate string  `json:"prompt_template,omitempty"`  // Prompt template
	UseVision      bool    `json:"use_vision,omitempty"`       // Whether the model supports vision
}
```

### ModelList

The `ModelList` struct contains a list of available models:

```go
type ModelList struct {
	Object string      `json:"object"` // Type of object (always "list")
	Data   []ModelCard `json:"data"`   // Array of model cards
}
```

## Managing Models

### Loading Models

To load a model, use the `Load` method with a `ModelLoadRequest`:

```go
type ModelLoadRequest struct {
	ModelName      string      `json:"model_name"`              // Name of the model to load
	MaxSeqLen      int         `json:"max_seq_len,omitempty"`   // Maximum sequence length
	RopeScale      float64     `json:"rope_scale,omitempty"`    // RoPE scaling factor
	RopeAlpha      interface{} `json:"rope_alpha,omitempty"`    // RoPE alpha (float64 or "auto")
	GPUSplit       []float64   `json:"gpu_split,omitempty"`     // GPU split ratio
	CacheSize      int         `json:"cache_size,omitempty"`    // KV cache size
	CacheMode      string      `json:"cache_mode,omitempty"`    // Cache mode
	ChunkSize      int         `json:"chunk_size,omitempty"`    // Chunk size
	PromptTemplate string      `json:"prompt_template,omitempty"` // Prompt template
}
```

Loading a model returns a `ModelLoadResponse`:

```go
type ModelLoadResponse struct {
	ModelType string `json:"model_type"` // Type of the model
	Module    int    `json:"module"`     // Current module being loaded
	Modules   int    `json:"modules"`    // Total number of modules
	Status    string `json:"status"`     // Loading status
}
```

### Streaming Model Loading

For large models, use `LoadStream` to get loading progress updates:

```go
type ModelLoadStream = Stream[*ModelLoadResponse]
```

### Model Properties

To get the properties of the currently loaded model, use `GetProps`:

```go
type ModelPropsResponse struct {
	TotalSlots                int                             `json:"total_slots"`
	ChatTemplate              string                          `json:"chat_template"`
	DefaultGenerationSettings *ModelDefaultGenerationSettings `json:"default_generation_settings"`
}

type ModelDefaultGenerationSettings struct {
	NCtx int `json:"n_ctx"` // Context size
}
```

## Managing Embedding Models

Embedding models are managed separately from regular models:

```go
type EmbeddingModelLoadRequest struct {
	EmbeddingModelName string `json:"embedding_model_name"`     // Name of the embedding model
	EmbeddingsDevice   string `json:"embeddings_device,omitempty"` // Device to load the model on
}
```

## Downloading Models

Models can be downloaded from HuggingFace:

```go
type DownloadRequest struct {
	RepoID     string   `json:"repo_id"`               // HuggingFace repository ID
	RepoType   string   `json:"repo_type,omitempty"`   // Repository type
	FolderName string   `json:"folder_name,omitempty"` // Folder to save the model
	Revision   string   `json:"revision,omitempty"`    // Git revision
	Token      string   `json:"token,omitempty"`       // HuggingFace token for private repos
	Include    []string `json:"include,omitempty"`     // Files to include
	Exclude    []string `json:"exclude,omitempty"`     // Files to exclude
}

type DownloadResponse struct {
	DownloadPath string `json:"download_path"` // Path where the model was downloaded
}
```

## Examples

### Listing Available Models

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
	client := tabby.NewClient(
		tabby.WithBaseURL("http://localhost:8080"),
		tabby.WithAdminKey("your-admin-key"), // Note: Admin key required for model operations
	)
	defer client.Close()
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// List all available models
	models, err := client.Models().List(ctx)
	if err != nil {
		log.Fatalf("Error listing models: %v", err)
	}
	
	fmt.Printf("Found %d available models:\n", len(models.Data))
	for i, model := range models.Data {
		fmt.Printf("%d. %s (owned by: %s)\n", i+1, model.ID, model.OwnedBy)
		if model.Parameters != nil {
			fmt.Printf("   - Max Sequence Length: %d\n", model.Parameters.MaxSeqLen)
			if model.Parameters.UseVision {
				fmt.Printf("   - Supports vision\n")
			}
		}
	}
}
```

### Getting Current Model

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
	client := tabby.NewClient(
		tabby.WithBaseURL("http://localhost:8080"),
		tabby.WithAdminKey("your-admin-key"),
	)
	defer client.Close()
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Get the currently loaded model
	current, err := client.Models().Get(ctx)
	if err != nil {
		log.Fatalf("Error getting current model: %v", err)
	}
	
	fmt.Printf("Currently loaded model: %s\n", current.ID)
	
	// Get model properties
	props, err := client.Models().GetProps(ctx)
	if err != nil {
		log.Fatalf("Error getting model properties: %v", err)
	}
	
	fmt.Printf("Model properties:\n")
	fmt.Printf("- Total slots: %d\n", props.TotalSlots)
	fmt.Printf("- Chat template: %s\n", props.ChatTemplate)
	if props.DefaultGenerationSettings != nil {
		fmt.Printf("- Default context size: %d\n", props.DefaultGenerationSettings.NCtx)
	}
}
```

### Loading a Model

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
	client := tabby.NewClient(
		tabby.WithBaseURL("http://localhost:8080"),
		tabby.WithAdminKey("your-admin-key"),
	)
	defer client.Close()
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute) // Longer timeout for model loading
	defer cancel()
	
	// Load a model with custom parameters
	loadReq := &tabby.ModelLoadRequest{
		ModelName:      "mistral-7b-v0.1",
		MaxSeqLen:      4096,
		RopeScale:      1.0,
		RopeAlpha:      "auto", // Auto-determine based on model
		CacheSize:      512,     // Cache size in MB
		PromptTemplate: "mistral", // Use Mistral prompt template
	}
	
	fmt.Printf("Loading model %s...\n", loadReq.ModelName)
	resp, err := client.Models().Load(ctx, loadReq)
	if err != nil {
		log.Fatalf("Error loading model: %v", err)
	}
	
	fmt.Printf("Model loaded successfully\n")
	fmt.Printf("- Model type: %s\n", resp.ModelType)
	fmt.Printf("- Status: %s\n", resp.Status)
}
```

### Streaming Model Loading Progress

```go
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"
	
	"github.com/pixelsquared/go-tabbyapi/tabby"
)

func main() {
	client := tabby.NewClient(
		tabby.WithBaseURL("http://localhost:8080"),
		tabby.WithAdminKey("your-admin-key"),
	)
	defer client.Close()
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute) // Longer timeout for model loading
	defer cancel()
	
	// Load a model with streaming progress
	loadReq := &tabby.ModelLoadRequest{
		ModelName: "llama2-70b",
		MaxSeqLen: 4096,
	}
	
	fmt.Printf("Loading model %s with streaming progress...\n", loadReq.ModelName)
	stream, err := client.Models().LoadStream(ctx, loadReq)
	if err != nil {
		log.Fatalf("Error creating load stream: %v", err)
	}
	defer stream.Close()
	
	// Process loading progress updates
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break // Loading completed
		}
		if err != nil {
			log.Fatalf("Error receiving from stream: %v", err)
		}
		
		// Update progress
		fmt.Printf("Loading progress: Module %d of %d, Status: %s\n", resp.Module, resp.Modules, resp.Status)
	}
	
	fmt.Println("Model loaded successfully")
}
```

### Downloading a Model from HuggingFace

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
	client := tabby.NewClient(
		tabby.WithBaseURL("http://localhost:8080"),
		tabby.WithAdminKey("your-admin-key"),
		tabby.WithTimeout(30*time.Minute), // Longer timeout for downloads
	)
	defer client.Close()
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	
	// Download a model from HuggingFace
	downloadReq := &tabby.DownloadRequest{
		RepoID:   "microsoft/phi-2",
		Revision: "main",
		// Token: "hf_...",  // For private repos
		Include: []string{"*.safetensors", "*.json", "tokenizer.model"},
	}
	
	fmt.Printf("Downloading model from %s...\n", downloadReq.RepoID)
	resp, err := client.Models().Download(ctx, downloadReq)
	if err != nil {
		log.Fatalf("Error downloading model: %v", err)
	}
	
	fmt.Printf("Model downloaded successfully to: %s\n", resp.DownloadPath)
}
```

### Managing Embedding Models

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
	client := tabby.NewClient(
		tabby.WithBaseURL("http://localhost:8080"),
		tabby.WithAdminKey("your-admin-key"),
	)
	defer client.Close()
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	
	// List available embedding models
	embeddingModels, err := client.Models().ListEmbedding(ctx)
	if err != nil {
		log.Fatalf("Error listing embedding models: %v", err)
	}
	
	fmt.Printf("Available embedding models:\n")
	for i, model := range embeddingModels.Data {
		fmt.Printf("%d. %s\n", i+1, model.ID)
	}
	
	// Load an embedding model
	loadReq := &tabby.EmbeddingModelLoadRequest{
		EmbeddingModelName: "BAAI/bge-small-en-v1.5",
		EmbeddingsDevice:   "cuda:0", // Load on first GPU
	}
	
	fmt.Printf("Loading embedding model %s...\n", loadReq.EmbeddingModelName)
	resp, err := client.Models().LoadEmbedding(ctx, loadReq)
	if err != nil {
		log.Fatalf("Error loading embedding model: %v", err)
	}
	fmt.Printf("Embedding model loaded successfully\n")
	
	// Get the current embedding model
	current, err := client.Models().GetEmbedding(ctx)
	if err != nil {
		log.Fatalf("Error getting current embedding model: %v", err)
	}
	
	fmt.Printf("Current embedding model: %s\n", current.ID)
}
```

## Error Handling

The Models service operations may fail due to various reasons:

```go
resp, err := client.Models().Load(ctx, loadReq)
if err != nil {
	var apiErr *tabby.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.Code() {
		case "permission_error":
			log.Fatalf("You need admin permissions to load models. Check your API key")
		case "invalid_request":
			log.Fatalf("Invalid request parameters: %s", apiErr.Error())
		case "not_found":
			log.Fatalf("Model not found: %s", loadReq.ModelName)
		default:
			log.Fatalf("API Error: %s", apiErr.Error())
		}
	} else {
		// Handle other types of errors like network issues
		log.Fatalf("Error loading model: %v", err)
	}
}
```

## Best Practices

1. **Admin Authentication**: Most model management operations require admin privileges:
   ```go
   client := tabby.NewClient(
       tabby.WithAdminKey("your-admin-key"),
   )
   ```

2. **Longer Timeouts**: Model loading and downloading can take significant time:
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
   ```

3. **Resource Management**: Always close streams when done:
   ```go
   stream, err := client.Models().LoadStream(ctx, req)
   if err != nil {
       return err
   }
   defer stream.Close()
   ```

4. **Model Selection**:
   - For completion tasks, select models with strong language generation capabilities
   - For embedding tasks, use models specifically designed for embeddings
   - Consider hardware limitations when selecting model size

5. **Parameter Tuning**:
   - `MaxSeqLen`: Higher values enable longer context but require more memory
   - `RopeScale`/`RopeAlpha`: Adjusts the model's effective context length
   - `CacheSize`: Larger cache improves performance but requires more memory
   - `ChunkSize`: Adjust for optimal throughput based on hardware

6. **Error Recovery**: Implement robust error handling, especially for long-running operations:
   ```go
   if err != nil {
       // Log the error
       log.Printf("Error loading model: %v", err)
       
       // Attempt to clean up if needed
       _ = client.Models().Unload(ctx)
       
       // Potentially retry with different parameters
       loadReq.MaxSeqLen = 2048 // Reduce parameters
       resp, err = client.Models().Load(ctx, loadReq)
   }