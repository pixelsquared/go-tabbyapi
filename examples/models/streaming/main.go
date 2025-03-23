// Package main provides an example of using the TabbyAPI client to stream model loading progress.
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
	adminKey := os.Getenv("TABBY_ADMIN_KEY") // Admin key required for model management

	// Create a new TabbyAPI client with admin privileges
	client := tabby.NewClient(
		tabby.WithBaseURL(endpoint),
		tabby.WithAdminKey(adminKey),       // Note: Model management requires admin access
		tabby.WithTimeout(300*time.Second), // Longer timeout for model loading
	)
	// Ensure the client is closed properly
	defer client.Close()

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	// Define the model to load
	// Change this to a model name that's available in your server's model directory
	modelName := getEnvOrDefault("TABBY_MODEL_NAME", "mistralai/Mistral-7B-Instruct-v0.2")

	fmt.Printf("Preparing to load model: %s\n", modelName)
	fmt.Println("Press Ctrl+C to cancel at any time")

	// Create a model load request
	loadReq := &tabby.ModelLoadRequest{
		ModelName: modelName,
		MaxSeqLen: 4096, // Optional: Context length
		RopeScale: 1.0,  // Optional: RoPE scaling factor
		CacheSize: 2000, // Optional: KV cache size in MB
	}

	// Initialize the streaming model load
	fmt.Println("Starting model loading process with streaming progress...")
	stream, err := client.Models().LoadStream(ctx, loadReq)
	if err != nil {
		log.Fatalf("Error initializing model load stream: %v", err)
	}
	// Ensure the stream is closed properly
	defer stream.Close()

	// Track progress
	var lastModule int
	var totalModules int
	startTime := time.Now()

	// Process the stream
	for {
		// Receive the next update from the stream
		response, err := stream.Recv()

		// Check for end of stream or errors
		if err != nil {
			if err == io.EOF {
				// Normal end of stream
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

		// Store total modules count when we first get it
		if totalModules == 0 && response.Modules > 0 {
			totalModules = response.Modules
		}

		// Update progress
		currentModule := response.Module

		// Only print if there's been a change in module number
		if currentModule != lastModule {
			elapsedTime := time.Since(startTime)
			progress := float64(currentModule) / float64(totalModules) * 100.0

			fmt.Printf("Loading module %d of %d (%.1f%%) - Model type: %s - Status: %s - Elapsed: %s\n",
				currentModule, totalModules, progress, response.ModelType, response.Status, elapsedTime.Round(time.Second))

			lastModule = currentModule
		}

		// If we get a status that indicates completion, break the loop
		if response.Status == "loaded" && currentModule >= totalModules {
			break
		}
	}

	// Final confirmation
	totalTime := time.Since(startTime).Round(time.Second)
	fmt.Printf("\nModel loading completed in %s\n", totalTime)

	// Verify the model is loaded by getting the current model
	currentModel, err := client.Models().Get(ctx)
	if err != nil {
		fmt.Printf("Error verifying loaded model: %v\n", err)
	} else if currentModel != nil {
		fmt.Printf("\nSuccessfully loaded model: %s\n", currentModel.ID)

		// Get model properties
		props, err := client.Models().GetProps(ctx)
		if err != nil {
			fmt.Printf("Error getting model properties: %v\n", err)
		} else {
			fmt.Printf("Model context length: %d tokens\n", props.DefaultGenerationSettings.NCtx)
		}
	}

	// Example of testing the model with a simple completion
	fmt.Println("\nTesting the model with a simple completion...")
	testCompletion(client, ctx)
}

// testCompletion tests the loaded model with a simple completion request
func testCompletion(client tabby.Client, ctx context.Context) {
	// Create a simple completion request
	req := &tabby.CompletionRequest{
		Prompt:      "Hello, world!",
		MaxTokens:   20,
		Temperature: 0.7,
	}

	// Call the API
	resp, err := client.Completions().Create(ctx, req)
	if err != nil {
		fmt.Printf("Error testing model with completion: %v\n", err)
		return
	}

	// Print the result
	if len(resp.Choices) > 0 {
		fmt.Println("\nModel test response:")
		fmt.Printf("Prompt: \"Hello, world!\"\n")
		fmt.Printf("Response: \"%s\"\n", resp.Choices[0].Text)
	} else {
		fmt.Println("No completion text was generated")
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
