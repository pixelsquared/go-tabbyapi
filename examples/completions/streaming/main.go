// Package main provides an example of using the TabbyAPI client for streaming text completions.
package main

import (
	"context"
	"fmt"
	"io"
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
		tabby.WithTimeout(60*time.Second), // Longer timeout for streaming
	)
	// Ensure the client is closed properly
	defer client.Close()

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Define the completion request with stream option enabled
	req := &tabby.CompletionRequest{
		Prompt:      "Write a short story about artificial intelligence in the year 2050.",
		MaxTokens:   300,
		Temperature: 0.8,
		TopP:        0.95,
		TopK:        50,
		Stream:      true, // Enable streaming response
	}

	fmt.Println("Starting streaming completion...")

	// Call the streaming API
	stream, err := client.Completions().CreateStream(ctx, req)
	if err != nil {
		log.Fatalf("Error creating completion stream: %v", err)
	}
	// Ensure the stream is closed properly
	defer stream.Close()

	// Process the streaming response
	var fullText string
	for {
		// Receive the next chunk from the stream
		response, err := stream.Recv()

		// Check for end of stream or errors
		if err != nil {
			if err == io.EOF {
				// End of stream
				break
			}
			// Check if it's a stream closed error from the client
			if err == tabby.ErrStreamClosed {
				fmt.Println("\nStream was closed")
				break
			}
			// Handle other errors
			log.Fatalf("Error receiving from stream: %v", err)
		}

		// Process the response chunk
		if len(response.Choices) > 0 {
			chunk := response.Choices[0].Text
			fullText += chunk

			// Print the chunk (without newline to allow continuous output)
			fmt.Print(chunk)

			// If we received a finish reason, we're done
			if response.Choices[0].FinishReason != "" {
				fmt.Printf("\n\nFinish reason: %s\n", response.Choices[0].FinishReason)
				break
			}
		}
	}

	fmt.Println("\n\nStreaming complete!")
	fmt.Printf("Total generated text length: %d characters\n", len(fullText))
}

// getEnvOrDefault returns the value of the environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
