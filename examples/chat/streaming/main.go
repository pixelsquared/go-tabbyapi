// Package main provides an example of using the TabbyAPI client for streaming chat completions.
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

	// Create a chat completion request with stream option enabled
	req := &tabby.ChatCompletionRequest{
		Messages: []tabby.ChatMessage{
			{
				Role:    tabby.ChatMessageRoleSystem,
				Content: "You are a creative writing assistant that specializes in crafting stories.",
			},
			{
				Role:    tabby.ChatMessageRoleUser,
				Content: "Write a short story about a future where AI and humans collaborate to solve climate change.",
			},
		},
		MaxTokens:   300,
		Temperature: 0.8,
		TopP:        0.95,
		TopK:        50,
		Stream:      true, // Enable streaming response
	}

	fmt.Println("Starting streaming chat completion...")

	// Call the streaming API
	stream, err := client.Chat().CreateStream(ctx, req)
	if err != nil {
		log.Fatalf("Error creating chat completion stream: %v", err)
	}
	// Ensure the stream is closed properly
	defer stream.Close()

	// Process the streaming response
	var fullContent string
	fmt.Println("Assistant: ")

	// Track the current role for rendering purposes (as ChatMessageRole)
	var currentRole tabby.ChatMessageRole

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
			delta := response.Choices[0].Delta

			// Print role change if it happens (usually only on first chunk)
			if delta.Role != "" && delta.Role != currentRole {
				currentRole = delta.Role
			}

			// Add content to our running total if present
			if delta.Content != "" {
				fullContent += delta.Content
				fmt.Print(delta.Content) // Print without newline for continuous output
			}

			// If we received a finish reason, we're done
			if response.Choices[0].FinishReason != "" {
				fmt.Printf("\n\nFinish reason: %s\n", response.Choices[0].FinishReason)
				break
			}
		}
	}

	fmt.Println("\n\nStreaming complete!")
	fmt.Printf("Total generated content length: %d characters\n", len(fullContent))

	// Example of how to use the streamed content in a follow-up request
	fmt.Println("\nYou could continue the conversation with a follow-up request like this:")
	fmt.Printf("User: Now continue the story with a twist ending.\n")
	fmt.Println("Example code for follow-up request:")

	// Using multiple Print statements instead of a raw string literal to avoid newline issues
	fmt.Print("followUpReq := &tabby.ChatCompletionRequest{\n")
	fmt.Print("	Messages: []tabby.ChatMessage{\n")
	fmt.Print("		{\n")
	fmt.Print("			Role:    tabby.ChatMessageRoleSystem,\n")
	fmt.Print("			Content: \"You are a creative writing assistant that specializes in crafting stories.\",\n")
	fmt.Print("		},\n")
	fmt.Print("		{\n")
	fmt.Print("			Role:    tabby.ChatMessageRoleUser,\n")
	fmt.Print("			Content: \"Write a short story about a future where AI and humans collaborate to solve climate change.\",\n")
	fmt.Print("		},\n")
	fmt.Print("		{\n")
	fmt.Print("			Role:    tabby.ChatMessageRoleAssistant,\n")
	fmt.Print("			Content: fullContent,\n")
	fmt.Print("		},\n")
	fmt.Print("		{\n")
	fmt.Print("			Role:    tabby.ChatMessageRoleUser,\n")
	fmt.Print("			Content: \"Now continue the story with a twist ending.\",\n")
	fmt.Print("		},\n")
	fmt.Print("	},\n")
	fmt.Print("	// Additional parameters would be set here\n")
	fmt.Print("}\n")
}

// getEnvOrDefault returns the value of the environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
