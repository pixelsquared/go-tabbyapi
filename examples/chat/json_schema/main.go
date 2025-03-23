// Package main provides an example of using JSON schema with TabbyAPI for structured chat completions.
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

	// Define a JSON schema for a product review
	reviewSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"product_name": map[string]interface{}{
				"type":        "string",
				"description": "The name of the product being reviewed",
			},
			"rating": map[string]interface{}{
				"type":        "integer",
				"description": "The rating from 1 to 5 stars",
				"minimum":     1,
				"maximum":     5,
			},
			"review_title": map[string]interface{}{
				"type":        "string",
				"description": "A short title for the review",
			},
			"review_content": map[string]interface{}{
				"type":        "string",
				"description": "The detailed review text",
			},
			"pros": map[string]interface{}{
				"type":        "array",
				"description": "List of positive aspects of the product",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
			"cons": map[string]interface{}{
				"type":        "array",
				"description": "List of negative aspects of the product",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
		},
		"required": []string{"product_name", "rating", "review_title", "review_content", "pros", "cons"},
	}

	// Define the chat completion request with JSON schema
	req := &tabby.ChatCompletionRequest{
		Messages: []tabby.ChatMessage{
			{
				Role:    tabby.ChatMessageRoleSystem,
				Content: "You are a product review assistant that helps users write structured reviews.",
			},
			{
				Role:    tabby.ChatMessageRoleUser,
				Content: "I just bought a new laptop and I want to write a review. It's a XYZ UltraBook Pro with 16GB RAM and a 1TB SSD. The battery life is excellent but the keyboard is not very comfortable.",
			},
		},
		MaxTokens:   300,
		Temperature: 0.7,
		TopP:        0.9,
		JSONSchema:  reviewSchema,
	}

	// Call the API to generate a structured chat completion
	fmt.Println("Generating structured chat completion...")
	resp, err := client.Chat().Create(ctx, req)
	if err != nil {
		log.Fatalf("Error generating chat completion: %v", err)
	}

	// Extract and print the generated response
	if len(resp.Choices) > 0 {
		content := ""
		if strContent, ok := resp.Choices[0].Message.Content.(string); ok {
			content = strContent
		}

		fmt.Println("\nGenerated JSON Review:")
		fmt.Println(content)

		// Parse the JSON response to verify it matches our schema
		var review map[string]interface{}
		err := json.Unmarshal([]byte(content), &review)
		if err != nil {
			fmt.Printf("\nError parsing JSON: %v\n", err)
		} else {
			// Print the parsed object more nicely
			prettyJSON, _ := json.MarshalIndent(review, "", "  ")
			fmt.Println("\nParsed Review Object:")
			fmt.Println(string(prettyJSON))
		}

		// Print optional finish reason if available
		if resp.Choices[0].FinishReason != "" {
			fmt.Printf("\nFinish reason: %s\n", resp.Choices[0].FinishReason)
		}
	} else {
		fmt.Println("No chat completion was generated")
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
