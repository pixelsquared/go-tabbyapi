// Package main provides an example of using the TabbyAPI client for generating embeddings.
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
	// Get API endpoint and key from environment variables
	endpoint := getEnvOrDefault("TABBY_API_ENDPOINT", "http://localhost:8080")
	apiKey := os.Getenv("TABBY_API_KEY")

	// Create a new TabbyAPI client
	client := tabby.NewClient(
		tabby.WithBaseURL(endpoint),
		tabby.WithAPIKey(apiKey),
		tabby.WithTimeout(30*time.Second),
	)
	// Ensure the client is closed properly
	defer client.Close()

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create an embeddings request
	// You can provide a single string or an array of strings for batch processing
	req := &tabby.EmbeddingsRequest{
		// Single input
		Input: "The quick brown fox jumps over the lazy dog.",

		// For batch processing, you can use:
		// Input: []string{
		//     "The quick brown fox jumps over the lazy dog.",
		//     "Machine learning models can process natural language.",
		//     "Embeddings are useful for semantic search and clustering.",
		// },

		// You can optionally specify a model, though this will typically
		// use the default embedding model loaded in the TabbyAPI server
		// Model: "your-embedding-model-name",

		// Optional: specify the encoding format (defaults to float)
		// EncodingFormat: "base64",
	}

	// Call the API to generate embeddings
	fmt.Println("Generating embeddings...")
	resp, err := client.Embeddings().Create(ctx, req)
	if err != nil {
		log.Fatalf("Error generating embeddings: %v", err)
	}

	// Process and display the embedding results
	fmt.Printf("\nGenerated %d embeddings\n", len(resp.Data))

	// For each embedding in the response
	for i, embedding := range resp.Data {
		// The embedding can be either a []float64 or a base64 string
		// depending on the encoding_format parameter
		fmt.Printf("\nEmbedding %d:\n", i)

		// Get the embedding type and first few dimensions
		switch e := embedding.Embedding.(type) {
		case []interface{}:
			// It's a float array (typical case)
			fmt.Printf("  Type: Float array\n")
			fmt.Printf("  Dimensions: %d\n", len(e))

			// Print the first few dimensions for demonstration
			fmt.Printf("  First 5 dimensions: [")
			for j := 0; j < min(5, len(e)); j++ {
				if j > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%.6f", e[j])
			}
			fmt.Println("]")

			// Example of how you would convert to []float64 for actual use
			fmt.Println("  (Convert to []float64 for actual use in your application)")

		case string:
			// It's a base64 string (if encoding_format was set to "base64")
			fmt.Printf("  Type: Base64 string\n")
			fmt.Printf("  Length: %d characters\n", len(e))
			fmt.Printf("  Preview: %s...\n", e[:min(30, len(e))])

		default:
			fmt.Printf("  Unexpected embedding type: %T\n", embedding.Embedding)
		}
	}

	// Print token usage information
	fmt.Printf("\nToken Usage:\n")
	fmt.Printf("  Prompt tokens: %d\n", resp.Usage.PromptTokens)
	fmt.Printf("  Total tokens: %d\n", resp.Usage.TotalTokens)

	// Example applications of embeddings
	fmt.Println("\nExample applications of embeddings:")
	fmt.Println("1. Semantic search: Compare document embeddings to query embeddings")
	fmt.Println("2. Clustering: Group similar documents based on embedding proximity")
	fmt.Println("3. Classification: Train classifiers on embeddings")
	fmt.Println("4. Recommendation systems: Suggest content with similar embeddings")
}

// getEnvOrDefault returns the value of the environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
