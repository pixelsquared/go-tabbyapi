// Package main provides an example of using the TabbyAPI client for model management.
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
	adminKey := os.Getenv("TABBY_ADMIN_KEY") // Admin key for model management

	// Create a new TabbyAPI client with admin privileges
	client := tabby.NewClient(
		tabby.WithBaseURL(endpoint),
		tabby.WithAdminKey(adminKey), // Note: Model management requires admin access
		tabby.WithTimeout(60*time.Second),
	)
	// Ensure the client is closed properly
	defer client.Close()

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// List available models
	fmt.Println("Listing available models...")
	modelList, err := client.Models().List(ctx)
	if err != nil {
		log.Fatalf("Error listing models: %v", err)
	}

	// Display the models
	fmt.Printf("\nFound %d models:\n", len(modelList.Data))
	for i, model := range modelList.Data {
		fmt.Printf("\n%d. Model ID: %s\n", i+1, model.ID)
		fmt.Printf("   Owner: %s\n", model.OwnedBy)
		fmt.Printf("   Created: %s\n", time.Unix(model.Created, 0).Format(time.RFC3339))

		// Print parameters if available
		if model.Parameters != nil {
			fmt.Println("   Parameters:")
			if model.Parameters.MaxSeqLen > 0 {
				fmt.Printf("     - Max Sequence Length: %d\n", model.Parameters.MaxSeqLen)
			}
			if model.Parameters.RopeScale > 0 {
				fmt.Printf("     - RoPE Scale: %.2f\n", model.Parameters.RopeScale)
			}
			if model.Parameters.MaxBatchSize > 0 {
				fmt.Printf("     - Max Batch Size: %d\n", model.Parameters.MaxBatchSize)
			}
			if model.Parameters.CacheSize > 0 {
				fmt.Printf("     - Cache Size: %d\n", model.Parameters.CacheSize)
			}
			if model.Parameters.CacheMode != "" {
				fmt.Printf("     - Cache Mode: %s\n", model.Parameters.CacheMode)
			}
			if model.Parameters.PromptTemplate != "" {
				fmt.Printf("     - Prompt Template: %s\n", model.Parameters.PromptTemplate)
			}
		}
	}

	// Get the currently loaded model (if any)
	fmt.Println("\nGetting currently loaded model...")
	currentModel, err := client.Models().Get(ctx)
	if err != nil {
		fmt.Printf("Error getting current model: %v\n", err)
		fmt.Println("No model is currently loaded.")
	} else {
		fmt.Printf("\nCurrently loaded model: %s\n", currentModel.ID)
	}

	// Get model properties if a model is loaded
	if currentModel != nil {
		fmt.Println("\nGetting model properties...")
		props, err := client.Models().GetProps(ctx)
		if err != nil {
			fmt.Printf("Error getting model properties: %v\n", err)
		} else {
			fmt.Printf("\nModel Properties:\n")
			fmt.Printf("  Total Slots: %d\n", props.TotalSlots)
			fmt.Printf("  Chat Template: %s\n", props.ChatTemplate)
			if props.DefaultGenerationSettings != nil {
				fmt.Printf("  Default Context Length: %d\n", props.DefaultGenerationSettings.NCtx)
			}
		}
	}

	// Example of loading a model (commented out to avoid unintentional model loads)
	/*
		fmt.Println("\nLoading a model...")
		loadReq := &tabby.ModelLoadRequest{
			ModelName: "mistralai/Mistral-7B-Instruct-v0.2",  // Example model name
			MaxSeqLen: 4096,                                   // Optional: Context length
			RopeScale: 1.0,                                   // Optional: RoPE scaling factor
			CacheSize: 2000,                                  // Optional: KV cache size in MB
		}

		loadResp, err := client.Models().Load(ctx, loadReq)
		if err != nil {
			log.Fatalf("Error loading model: %v", err)
		}

		fmt.Printf("\nModel loading initiated:\n")
		fmt.Printf("  Model type: %s\n", loadResp.ModelType)
		fmt.Printf("  Module: %d of %d\n", loadResp.Module, loadResp.Modules)
		fmt.Printf("  Status: %s\n", loadResp.Status)
	*/

	// Example of unloading a model (commented out to avoid unintentional unloads)
	/*
		fmt.Println("\nUnloading the current model...")
		err = client.Models().Unload(ctx)
		if err != nil {
			log.Fatalf("Error unloading model: %v", err)
		}
		fmt.Println("Model successfully unloaded")
	*/

	// Listing embedding models example
	fmt.Println("\nListing available embedding models...")
	embeddingModelList, err := client.Models().ListEmbedding(ctx)
	if err != nil {
		fmt.Printf("Error listing embedding models: %v\n", err)
	} else {
		fmt.Printf("\nFound %d embedding models:\n", len(embeddingModelList.Data))
		for i, model := range embeddingModelList.Data {
			fmt.Printf("%d. %s\n", i+1, model.ID)
		}
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
