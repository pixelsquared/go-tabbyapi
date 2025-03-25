# Chat Service

The Chat service enables multi-turn conversations with AI models using different message roles. It provides support for both synchronous and streaming responses, allowing for interactive dialogues with the model.

## Interface

```go
// ChatService handles chat completion requests for multi-turn conversations
// with support for different message roles (system, user, assistant).
type ChatService interface {
	// Create generates a chat completion for the provided request.
	// This method returns a complete, non-streaming response.
	Create(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error)

	// CreateStream generates a streaming chat completion where tokens are returned
	// incrementally as they're generated.
	// The returned ChatCompletionStream must be closed when no longer needed to release resources.
	CreateStream(ctx context.Context, req *ChatCompletionRequest) (ChatCompletionStream, error)
}
```

## Chat Message Roles

The Chat service uses different message roles to structure conversations:

```go
type ChatMessageRole string

const (
	// ChatMessageRoleUser represents a message from the user
	ChatMessageRoleUser ChatMessageRole = "user"

	// ChatMessageRoleAssistant represents a message from the assistant
	ChatMessageRoleAssistant ChatMessageRole = "assistant"

	// ChatMessageRoleSystem represents a system message, typically used for instructions
	ChatMessageRoleSystem ChatMessageRole = "system"

	// ChatMessageRoleTool represents a tool message, typically used for function calls
	ChatMessageRoleTool ChatMessageRole = "tool"
)
```

Each role serves a specific purpose:
- **System**: Sets behavior instructions, context, or constraints for the assistant
- **User**: Represents user queries or inputs
- **Assistant**: Contains responses from the assistant
- **Tool**: Contains outputs from tool or function calls

## ChatCompletionRequest

The `ChatCompletionRequest` struct defines the parameters for a chat request:

```go
type ChatCompletionRequest struct {
	Messages    []ChatMessage `json:"messages"`        // Conversation history
	MaxTokens   int           `json:"max_tokens,omitempty"`  // Maximum tokens to generate
	Temperature float64       `json:"temperature,omitempty"` // Controls randomness
	TopP        float64       `json:"top_p,omitempty"`       // Alternative to temperature
	TopK        int           `json:"top_k,omitempty"`       // Limit to K most likely tokens
	Stream      bool          `json:"stream,omitempty"`      // Enable streaming responses
	Stop        []string      `json:"stop,omitempty"`        // Sequences where generation should stop
	Model       string        `json:"model,omitempty"`       // Model to use (if not default)
	JSONSchema  interface{}   `json:"json_schema,omitempty"` // JSON schema for structured output
}

type ChatMessage struct {
	Role    ChatMessageRole `json:"role"`    // Role of the message sender
	Content interface{}     `json:"content"` // String or array of ChatMessageContent
}
```

### Request Parameters

| Parameter   | Type            | Description                                         | Default |
|-------------|-----------------|-----------------------------------------------------|---------|
| Messages    | []ChatMessage   | Array of messages representing the conversation     | (required) |
| MaxTokens   | int             | Maximum number of tokens to generate                | (model dependent) |
| Temperature | float64         | Controls randomness (higher = more random)          | 1.0 |
| TopP        | float64         | Nucleus sampling parameter                          | 1.0 |
| TopK        | int             | Only sample from top K most likely tokens           | 0 (disabled) |
| Stream      | bool            | Enable streaming response (token-by-token)          | false |
| Stop        | []string        | Stop sequences to end generation when encountered   | [] |
| Model       | string          | Model ID to use (if multiple available)             | (currently loaded model) |
| JSONSchema  | interface{}     | Schema for structured JSON output                   | nil |

## Multimodal Content

The ChatMessage content can be either a string or an array of content parts with different types:

```go
// ChatMessageContent represents a part of a message content
type ChatMessageContent struct {
	Type     string        `json:"type"`                  // Content type ("text" or "image_url")
	Text     string        `json:"text,omitempty"`        // Text content
	ImageURL *ChatImageURL `json:"image_url,omitempty"`   // Image URL content
}

// ChatImageURL represents an image URL in a chat message
type ChatImageURL struct {
	URL string `json:"url"`  // The URL of the image
}
```

This allows for creating multimodal messages that contain both text and images.

## ChatCompletionResponse

The response to a chat completion request contains the model's reply:

```go
type ChatCompletionResponse struct {
	ID      string                     `json:"id"`      // Unique identifier for the completion
	Object  string                     `json:"object"`  // Type of object (always "chat.completion")
	Created int64                      `json:"created"` // Unix timestamp of creation
	Model   string                     `json:"model"`   // Model used for completion
	Choices []ChatCompletionRespChoice `json:"choices"` // Generated chat choices
	Usage   *UsageStats                `json:"usage,omitempty"`  // Token usage statistics
}

type ChatCompletionRespChoice struct {
	Index        int         `json:"index"`        // Choice index
	Message      ChatMessage `json:"message"`      // The assistant's response message
	FinishReason string      `json:"finish_reason,omitempty"` // Why generation ended
}

type UsageStats struct {
	PromptTokens     int `json:"prompt_tokens"`     // Tokens in the prompt
	CompletionTokens int `json:"completion_tokens"` // Tokens in the completion
	TotalTokens      int `json:"total_tokens"`      // Total tokens used
}
```

## Streaming Chat Completions

For streaming chat completions, use `CreateStream` which returns delta updates incrementally:

```go
type ChatCompletionStream = Stream[*ChatCompletionStreamResponse]

type ChatCompletionStreamResponse struct {
	ID      string                       `json:"id"`
	Object  string                       `json:"object"`
	Created int64                        `json:"created"`
	Model   string                       `json:"model"`
	Choices []ChatCompletionStreamChoice `json:"choices"`
}

type ChatCompletionStreamChoice struct {
	Index        int    `json:"index"`         // Choice index
	Delta        *Delta `json:"delta"`         // Partial content update
	FinishReason string `json:"finish_reason,omitempty"` // Present only in the final chunk
}

type Delta struct {
	Role    ChatMessageRole `json:"role,omitempty"`    // Present only in the first chunk
	Content string          `json:"content,omitempty"` // Partial content
}
```

## Examples

### Basic Chat Completion

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
	
	req := &tabby.ChatCompletionRequest{
		Messages: []tabby.ChatMessage{
			{
				Role:    tabby.ChatMessageRoleSystem,
				Content: "You are a helpful assistant specialized in programming and technology.",
			},
			{
				Role:    tabby.ChatMessageRoleUser,
				Content: "What are the benefits of using dependency injection in software design?",
			},
		},
		MaxTokens:   150,
		Temperature: 0.7,
	}
	
	resp, err := client.Chat().Create(ctx, req)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	
	if len(resp.Choices) > 0 {
		fmt.Println("Assistant's Response:")
		fmt.Println(resp.Choices[0].Message.Content)
	}
	
	if resp.Usage != nil {
		fmt.Printf("\nToken Usage:\n")
		fmt.Printf("  Prompt tokens: %d\n", resp.Usage.PromptTokens)
		fmt.Printf("  Completion tokens: %d\n", resp.Usage.CompletionTokens)
		fmt.Printf("  Total tokens: %d\n", resp.Usage.TotalTokens)
	}
}
```

### Multi-turn Conversation

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
	
	// Initial conversation
	messages := []tabby.ChatMessage{
		{
			Role:    tabby.ChatMessageRoleSystem,
			Content: "You are a helpful assistant specialized in programming and technology.",
		},
		{
			Role:    tabby.ChatMessageRoleUser,
			Content: "What is the difference between REST and GraphQL?",
		},
	}
	
	req := &tabby.ChatCompletionRequest{
		Messages:    messages,
		MaxTokens:   150,
		Temperature: 0.7,
	}
	
	// First response
	resp, err := client.Chat().Create(ctx, req)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	
	if len(resp.Choices) > 0 {
		fmt.Println("Initial Response:")
		fmt.Println(resp.Choices[0].Message.Content)
		
		// Add the assistant's response to the conversation
		messages = append(messages, tabby.ChatMessage{
			Role:    tabby.ChatMessageRoleAssistant,
			Content: resp.Choices[0].Message.Content,
		})
		
		// Add a follow-up question
		messages = append(messages, tabby.ChatMessage{
			Role:    tabby.ChatMessageRoleUser,
			Content: "Which one is better for mobile applications and why?",
		})
		
		// Create a new request with the updated conversation
		req.Messages = messages
		
		// Get the follow-up response
		resp, err = client.Chat().Create(ctx, req)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		
		if len(resp.Choices) > 0 {
			fmt.Println("\nFollow-up Response:")
			fmt.Println(resp.Choices[0].Message.Content)
		}
	}
}
```

### Streaming Chat Completion

```go
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
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
	
	req := &tabby.ChatCompletionRequest{
		Messages: []tabby.ChatMessage{
			{
				Role:    tabby.ChatMessageRoleSystem,
				Content: "You are a helpful assistant specialized in programming and technology.",
			},
			{
				Role:    tabby.ChatMessageRoleUser,
				Content: "Explain the concept of containers in Docker and how they differ from virtual machines.",
			},
		},
		MaxTokens:   300,
		Temperature: 0.7,
		Stream:      true,
	}
	
	stream, err := client.Chat().CreateStream(ctx, req)
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}
	defer stream.Close()
	
	fmt.Println("Streaming response:")
	
	var fullResponse strings.Builder
	
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break // Stream completed
		}
		if err != nil {
			log.Fatalf("Error receiving from stream: %v", err)
		}
		
		// Process each chunk
		if len(resp.Choices) > 0 {
			chunk := resp.Choices[0].Delta.Content
			fmt.Print(chunk)
			fullResponse.WriteString(chunk)
		}
	}
	
	fmt.Println("\n\nFull response received and assembled.")
}
```

### Multimodal Chat (With Image)

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
	
	// Create a multimodal message with text and image
	content := []tabby.ChatMessageContent{
		{
			Type: "text",
			Text: "What's in this image? Describe it in detail.",
		},
		{
			Type: "image_url",
			ImageURL: &tabby.ChatImageURL{
				URL: "https://example.com/image.jpg",
			},
		},
	}
	
	req := &tabby.ChatCompletionRequest{
		Messages: []tabby.ChatMessage{
			{
				Role:    tabby.ChatMessageRoleSystem,
				Content: "You are a helpful assistant that can analyze images.",
			},
			{
				Role:    tabby.ChatMessageRoleUser,
				Content: content,
			},
		},
		MaxTokens:   200,
		Temperature: 0.7,
	}
	
	resp, err := client.Chat().Create(ctx, req)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	
	if len(resp.Choices) > 0 {
		fmt.Println("Assistant's Response:")
		fmt.Println(resp.Choices[0].Message.Content)
	}
}
```

## Error Handling

The Chat service can return the same types of errors as other services:

```go
resp, err := client.Chat().Create(ctx, req)
if err != nil {
	var apiErr *tabby.APIError
	if errors.As(err, &apiErr) {
		fmt.Printf("API Error (code %s): %s\n", apiErr.Code(), apiErr.Error())
		
		switch apiErr.Code() {
		case "invalid_request":
			fmt.Println("Check your request parameters")
		case "authentication_error":
			fmt.Println("Check your API key")
		case "permission_error":
			fmt.Println("Your API key doesn't have sufficient permissions")
		}
	}
	
	log.Fatalf("Error: %v", err)
}
```

## Best Practices

1. **Conversation Context**: Keep the conversation history for multi-turn chats, but be mindful of token limits.

2. **System Instructions**: Use system messages to set behavior, constraints, or provide context.

3. **Streaming for Long Responses**: Use streaming for longer responses to improve user experience.

4. **Resource Management**: Always close streams when done:
   ```go
   stream, err := client.Chat().CreateStream(ctx, req)
   if err != nil {
       return err
   }
   defer stream.Close()
   ```

5. **Temperature Settings**:
   - Lower temperature (0.1-0.5) for factual, consistent responses
   - Higher temperature (0.7-1.0) for more creative, varied responses

6. **Message Role Usage**:
   - System: Set overall behavior and constraints (typically one message at the start)
   - User: Provide inputs, questions, or requests
   - Assistant: Include previous responses in the conversation history
   - Tool: Include outputs from tool or function calls when relevant

7. **Token Management**:
   - Be aware of model context window limits
   - Consider truncating or summarizing very long conversation histories