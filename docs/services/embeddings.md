# Embeddings Service

The Embeddings service allows you to generate vector embeddings from text inputs. These embeddings represent text as numerical vectors that capture semantic meaning, making them useful for semantic search, clustering, classification, and other machine learning tasks.

## Interface

```go
// EmbeddingsService handles embedding generation for text inputs,
// converting text into numerical vector representations.
type EmbeddingsService interface {
	// Create generates embeddings for the provided input.
	// This method converts text into vector embeddings that can be used for
	// semantic search, clustering, classification, and other machine learning tasks.
	//
	// The EmbeddingsRequest can contain either a single string or an array of
	// strings for batch processing. The response contains a corresponding
	// number of embedding vectors.
	//
	// The embedding model must be loaded via ModelsService.LoadEmbedding
	// before using this method, unless a default embedding model is configured.
	Create(ctx context.Context, req *EmbeddingsRequest) (*EmbeddingsResponse, error)
}
```

## EmbeddingsRequest

The `EmbeddingsRequest` struct defines the parameters for an embeddings request:

```go
type EmbeddingsRequest struct {
	Input          interface{} `json:"input"`                    // String or array of strings
	Model          string      `json:"model,omitempty"`          // Model to use (if not default)
	EncodingFormat string      `json:"encoding_format,omitempty"` // Format for returned embeddings
}
```

### Request Parameters

| Parameter      | Type        | Description                                        | Default |
|----------------|-------------|----------------------------------------------------|---------|
| Input          | interface{} | Text to embed (string or []string)                 | (required) |
| Model          | string      | Embedding model ID to use                          | (current embedding model) |
| EncodingFormat | string      | Format of the embedding vectors ("float" or "base64") | "float" |

## EmbeddingsResponse

The response to an embeddings request contains the generated embedding vectors:

```go
type EmbeddingsResponse struct {
	Object string            `json:"object"`  // Type of object (always "list")
	Data   []EmbeddingObject `json:"data"`    // Array of embedding objects
	Model  string            `json:"model"`   // Model used for embeddings
	Usage  UsageInfo         `json:"usage"`   // Token usage information
}

type EmbeddingObject struct {
	Object    string      `json:"object"`    // Type of object (always "embedding")
	Embedding interface{} `json:"embedding"` // Array of floats or base64 string
	Index     int         `json:"index"`     // Index in the input array
}

type UsageInfo struct {
	PromptTokens int `json:"prompt_tokens"` // Tokens in the input
	TotalTokens  int `json:"total_tokens"`  // Total tokens used
}
```

The `Embedding` field contains the vector representation of the input text, either as:
- A slice of float values `[]float32` when `EncodingFormat` is "float" (default)
- A base64-encoded string when `EncodingFormat` is "base64"

## Examples

### Basic Embedding Generation

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
		tabby.WithAPIKey("your-api-key"),
	)
	defer client.Close()
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Create an embedding request for a single string
	req := &tabby.EmbeddingsRequest{
		Input: "The quick brown fox jumps over the lazy dog.",
	}
	
	resp, err := client.Embeddings().Create(ctx, req)
	if err != nil {
		log.Fatalf("Error generating embeddings: %v", err)
	}
	
	// Print information about the response
	fmt.Printf("Model used: %s\n", resp.Model)
	fmt.Printf("Number of embeddings: %d\n", len(resp.Data))
	
	if len(resp.Data) > 0 {
		// Get the first embedding
		embedding := resp.Data[0].Embedding
		
		// Type assertion to get the embedding vector
		if values, ok := embedding.([]float32); ok {
			fmt.Printf("Embedding dimension: %d\n", len(values))
			fmt.Printf("First 5 values: %v\n", values[:5])
		}
	}
	
	// Print token usage
	fmt.Printf("Token usage: %d tokens\n", resp.Usage.TotalTokens)
}
```

### Batch Embedding Generation

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
		tabby.WithAPIKey("your-api-key"),
	)
	defer client.Close()
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Create an embedding request for multiple strings
	req := &tabby.EmbeddingsRequest{
		Input: []string{
			"The quick brown fox jumps over the lazy dog.",
			"Machine learning is a field of artificial intelligence.",
			"Vector embeddings capture semantic meaning in text.",
		},
	}
	
	resp, err := client.Embeddings().Create(ctx, req)
	if err != nil {
		log.Fatalf("Error generating batch embeddings: %v", err)
	}
	
	fmt.Printf("Generated %d embeddings\n", len(resp.Data))
	
	// Process each embedding
	for i, item := range resp.Data {
		// Type assertion to get the embedding vector
		if values, ok := item.Embedding.([]float32); ok {
			fmt.Printf("Embedding %d: dimension=%d, index=%d\n", 
				i+1, len(values), item.Index)
		}
	}
}
```

### Using Base64 Encoding Format

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
		tabby.WithAPIKey("your-api-key"),
	)
	defer client.Close()
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Request base64-encoded embeddings to reduce JSON size for large models
	req := &tabby.EmbeddingsRequest{
		Input:          "The quick brown fox jumps over the lazy dog.",
		EncodingFormat: "base64",
	}
	
	resp, err := client.Embeddings().Create(ctx, req)
	if err != nil {
		log.Fatalf("Error generating embeddings: %v", err)
	}
	
	// Process the base64-encoded embedding
	if len(resp.Data) > 0 {
		if base64Str, ok := resp.Data[0].Embedding.(string); ok {
			fmt.Printf("Base64-encoded embedding: %s... (truncated)\n", base64Str[:50])
			fmt.Printf("Base64 string length: %d\n", len(base64Str))
		}
	}
}
```

### Semantic Search Example

```go
package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"sort"
	"time"
	
	"github.com/pixelsquared/go-tabbyapi/tabby"
)

func main() {
	client := tabby.NewClient(
		tabby.WithBaseURL("http://localhost:8080"),
		tabby.WithAPIKey("your-api-key"),
	)
	defer client.Close()
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Define a corpus of documents
	documents := []string{
		"Machine learning is a field of inquiry devoted to understanding and building methods that 'learn'.",
		"Deep learning is part of a broader family of machine learning methods based on artificial neural networks.",
		"Natural language processing is a subfield of linguistics, computer science, and AI concerned with interactions between computers and human language.",
		"Computer vision is an interdisciplinary field that deals with how computers can gain high-level understanding from digital images or videos.",
		"Reinforcement learning is an area of machine learning concerned with how intelligent agents ought to take actions in an environment.",
	}
	
	// Generate embeddings for all documents
	req := &tabby.EmbeddingsRequest{
		Input: documents,
	}
	
	docEmbs, err := client.Embeddings().Create(ctx, req)
	if err != nil {
		log.Fatalf("Error generating document embeddings: %v", err)
	}
	
	// Generate embedding for the query
	query := "How do computers understand images?"
	queryReq := &tabby.EmbeddingsRequest{
		Input: query,
	}
	
	queryEmb, err := client.Embeddings().Create(ctx, queryReq)
	if err != nil {
		log.Fatalf("Error generating query embedding: %v", err)
	}
	
	// Perform semantic search using cosine similarity
	type SearchResult struct {
		DocumentIndex int
		Document      string
		Similarity    float64
	}
	
	var results []SearchResult
	
	// Extract query embedding vector
	var queryVector []float32
	if v, ok := queryEmb.Data[0].Embedding.([]float32); ok {
		queryVector = v
	} else {
		log.Fatalf("Could not get query embedding vector")
	}
	
	// Calculate similarity for each document
	for i, docEmb := range docEmbs.Data {
		// Extract document embedding vector
		var docVector []float32
		if v, ok := docEmb.Embedding.([]float32); ok {
			docVector = v
		} else {
			continue
		}
		
		// Calculate cosine similarity
		similarity := cosineSimilarity(queryVector, docVector)
		
		results = append(results, SearchResult{
			DocumentIndex: i,
			Document:      documents[i],
			Similarity:    similarity,
		})
	}
	
	// Sort results by similarity (highest first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})
	
	// Print top results
	fmt.Printf("Search query: %s\n\n", query)
	fmt.Println("Top results:")
	for i, result := range results {
		fmt.Printf("%d. [Score: %.4f] %s\n", i+1, result.Similarity, result.Document)
	}
}

// cosineSimilarity calculates the cosine similarity between two vectors
func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}
	
	var dotProduct float64
	var normA float64
	var normB float64
	
	for i := 0; i < len(a); i++ {
		dotProduct += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}
	
	if normA == 0 || normB == 0 {
		return 0
	}
	
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
```

## Using with Models Service

Before generating embeddings, you typically need to load an embedding model:

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
		tabby.WithAdminKey("your-admin-key"), // Admin key for model operations
	)
	defer client.Close()
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	
	// Step 1: List available embedding models
	models, err := client.Models().ListEmbedding(ctx)
	if err != nil {
		log.Fatalf("Error listing embedding models: %v", err)
	}
	
	if len(models.Data) == 0 {
		log.Fatalf("No embedding models available")
	}
	
	fmt.Println("Available embedding models:")
	for i, model := range models.Data {
		fmt.Printf("%d. %s\n", i+1, model.ID)
	}
	
	// Step 2: Load an embedding model
	loadReq := &tabby.EmbeddingModelLoadRequest{
		EmbeddingModelName: models.Data[0].ID, // Use the first available model
		EmbeddingsDevice:   "cuda:0",         // Use GPU if available
	}
	
	_, err = client.Models().LoadEmbedding(ctx, loadReq)
	if err != nil {
		log.Fatalf("Error loading embedding model: %v", err)
	}
	
	fmt.Printf("Loaded embedding model: %s\n", loadReq.EmbeddingModelName)
	
	// Step 3: Generate embeddings
	embReq := &tabby.EmbeddingsRequest{
		Input: "Now we can generate embeddings with the loaded model.",
	}
	
	resp, err := client.Embeddings().Create(ctx, embReq)
	if err != nil {
		log.Fatalf("Error generating embeddings: %v", err)
	}
	
	fmt.Printf("Successfully generated embeddings using model: %s\n", resp.Model)
	if len(resp.Data) > 0 {
		if values, ok := resp.Data[0].Embedding.([]float32); ok {
			fmt.Printf("Embedding dimension: %d\n", len(values))
		}
	}
}
```

## Error Handling

The Embeddings service can return several types of errors:

```go
resp, err := client.Embeddings().Create(ctx, req)
if err != nil {
	var apiErr *tabby.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.Code() {
		case "not_found":
			fmt.Println("No embedding model is loaded. Load an embedding model first using Models().LoadEmbedding()")
		case "invalid_request":
			fmt.Println("Invalid request parameters. Check your input format.")
		default:
			fmt.Printf("API Error: %s\n", apiErr.Error())
		}
	} else {
		fmt.Printf("Error: %v\n", err)
	}
	return
}
```

## Best Practices

1. **Load Embedding Models**: Before using the Embeddings service, ensure an embedding model is loaded using `Models().LoadEmbedding()`.

2. **Batch Processing**: When embedding multiple texts, pass them as a single request rather than making multiple requests:
   ```go
   // Efficient
   req := &tabby.EmbeddingsRequest{
       Input: []string{"text1", "text2", "text3"},
   }
   
   // Less efficient
   // Making three separate requests for each text
   ```

3. **Input Length**: Be mindful of token limits for the embedding model. Very long texts may be truncated.

4. **Encoding Format**: Use "base64" encoding format for very large embedding models to reduce JSON payload size:
   ```go
   req := &tabby.EmbeddingsRequest{
       Input: "text",
       EncodingFormat: "base64",
   }
   ```

5. **Caching**: For frequently used embeddings, consider caching them in your application to avoid regenerating them.

6. **Model Selection**: Choose embedding models based on your specific needs:
   - Smaller models (e.g., "BAAI/bge-small-en") for faster processing
   - Larger models for potentially higher quality embeddings
   - Domain-specific models for specialized applications

7. **Normalization**: For some similarity operations, normalizing embeddings can improve results:
   ```go
   func normalizeVector(v []float32) []float32 {
       result := make([]float32, len(v))
       var sum float32
       
       for _, val := range v {
           sum += val * val
       }
       
       magnitude := float32(math.Sqrt(float64(sum)))
       if magnitude > 0 {
           for i, val := range v {
               result[i] = val / magnitude
           }
       }
       
       return result
   }