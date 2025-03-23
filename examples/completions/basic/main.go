// Package main provides a basic example of using the TabbyAPI client for text completions.
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

	// Create a new TabbyAPI client with a 30-second timeout
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

	// Define the completion request
	req := &tabby.CompletionRequest{
		// The prompt can be a string or an array of strings
		Prompt:      "Once upon a time, there was a programmer who",
		MaxTokens:   100,           // Maximum number of tokens to generate
		Temperature: 0.7,           // Controls randomness: 0.0 is deterministic, higher values are more random
		TopP:        0.9,           // Top-p sampling: 1.0 is no filtering, lower is more focused
		TopK:        40,            // Top-k sampling: higher values allow more diverse completions
		Stop:        []string{"."}, // Stop generation at these strings (optional)
	}

	// Call the API to generate a completion
	fmt.Println("Generating completion...")
	resp, err := client.Completions().Create(ctx, req)
	if err != nil {
		log.Fatalf("Error generating completion: %v", err)
	}

	// Extract and print the generated text
	if len(resp.Choices) > 0 {
		fmt.Println("\nGenerated Completion:")
		fmt.Println(resp.Choices[0].Text)

		// Print optional finish reason if available
		if resp.Choices[0].FinishReason != "" {
			fmt.Printf("\nFinish reason: %s\n", resp.Choices[0].FinishReason)
		}
	} else {
		fmt.Println("No completion text was generated")
	}

	// Print token usage information if available
	if resp.Usage != nil {
		fmt.Printf("\nToken Usage:\n")
		fmt.Printf("  Prompt tokens: %d\n", resp.Usage.PromptTokens)
		fmt.Printf("  Completion tokens: %d\n", resp.Usage.CompletionTokens)
		fmt.Printf("  Total tokens: %d\n", resp.Usage.TotalTokens)
	}
}

// getEnvOrDefault returns the value of the environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
