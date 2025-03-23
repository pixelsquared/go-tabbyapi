// Package main provides an example of using JSON schema with TabbyAPI for structured completions.
package main

import (
	"context"
	"encoding/json"
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

	// Define a JSON schema for a person object
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
			"interests": map[string]interface{}{
				"type":        "array",
				"description": "The person's hobbies and interests",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
		},
		"required": []string{"name", "age", "email"},
	}

	// Define the completion request with JSON schema
	req := &tabby.CompletionRequest{
		Prompt:      "Generate information about a software developer named John who loves golang",
		MaxTokens:   256,
		Temperature: 0.7,
		TopP:        0.9,
		JSONSchema:  personSchema,
	}

	// Call the API to generate a structured completion
	fmt.Println("Generating structured completion...")
	resp, err := client.Completions().Create(ctx, req)
	if err != nil {
		log.Fatalf("Error generating completion: %v", err)
	}

	// Extract and print the generated text
	if len(resp.Choices) > 0 {
		fmt.Println("\nGenerated JSON:")
		fmt.Println(resp.Choices[0].Text)

		// Parse the JSON response to verify it matches our schema
		var person map[string]interface{}
		err := json.Unmarshal([]byte(resp.Choices[0].Text), &person)
		if err != nil {
			fmt.Printf("\nError parsing JSON: %v\n", err)
		} else {
			// Print the parsed object more nicely
			prettyJSON, _ := json.MarshalIndent(person, "", "  ")
			fmt.Println("\nParsed Person Object:")
			fmt.Println(string(prettyJSON))
		}

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
