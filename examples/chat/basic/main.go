// Package main provides a basic example of using the TabbyAPI client for chat completions.
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

	// Create a chat completion request with a conversation history
	req := &tabby.ChatCompletionRequest{
		// Define the conversation messages
		Messages: []tabby.ChatMessage{
			{
				Role:    tabby.ChatMessageRoleSystem,
				Content: "You are a helpful assistant specialized in programming and technology.",
			},
			{
				Role:    tabby.ChatMessageRoleUser,
				Content: "What are the main benefits of using Go for backend development?",
			},
		},
		MaxTokens:   150, // Maximum number of tokens to generate
		Temperature: 0.7, // Controls randomness
		TopP:        0.9, // Top-p sampling
		TopK:        40,  // Top-k sampling
	}

	// Call the API to generate a chat completion
	fmt.Println("Generating chat completion...")
	resp, err := client.Chat().Create(ctx, req)
	if err != nil {
		log.Fatalf("Error generating chat completion: %v", err)
	}

	// Extract and print the generated response
	if len(resp.Choices) > 0 {
		fmt.Println("\nAssistant's Response:")
		fmt.Println(resp.Choices[0].Message.Content)

		// Print optional finish reason if available
		if resp.Choices[0].FinishReason != "" {
			fmt.Printf("\nFinish reason: %s\n", resp.Choices[0].FinishReason)
		}
	} else {
		fmt.Println("No response was generated")
	}

	// Print token usage information if available
	if resp.Usage != nil {
		fmt.Printf("\nToken Usage:\n")
		fmt.Printf("  Prompt tokens: %d\n", resp.Usage.PromptTokens)
		fmt.Printf("  Completion tokens: %d\n", resp.Usage.CompletionTokens)
		fmt.Printf("  Total tokens: %d\n", resp.Usage.TotalTokens)
	}

	// Example of continuing the conversation
	fmt.Println("\n--- Continuing the conversation ---")

	// Add the assistant's response to the conversation history
	if len(resp.Choices) > 0 {
		req.Messages = append(req.Messages, tabby.ChatMessage{
			Role:    tabby.ChatMessageRoleAssistant,
			Content: resp.Choices[0].Message.Content,
		})
	}

	// Add a follow-up question from the user
	req.Messages = append(req.Messages, tabby.ChatMessage{
		Role:    tabby.ChatMessageRoleUser,
		Content: "What about its concurrency model? How does that work?",
	})

	// Generate another response
	fmt.Println("Generating follow-up response...")
	resp2, err := client.Chat().Create(ctx, req)
	if err != nil {
		log.Fatalf("Error generating follow-up response: %v", err)
	}

	// Print the follow-up response
	if len(resp2.Choices) > 0 {
		fmt.Println("\nAssistant's Follow-up Response:")
		fmt.Println(resp2.Choices[0].Message.Content)
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
