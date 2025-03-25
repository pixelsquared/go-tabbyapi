# Completions Service

The Completions service allows you to generate text completions based on provided prompts. It supports both synchronous and streaming responses, making it suitable for various applications ranging from code completion to content generation.

## Interface

```go
// CompletionsService handles text completion requests for generating text
// based on provided prompts, with support for both synchronous and streaming responses.
type CompletionsService interface {
	// Create generates a completion for the provided request.
	// This method returns a complete, non-streaming response.
	Create(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)

	// CreateStream generates a streaming completion where tokens are returned
	// incrementally as they're generated.
	// The returned CompletionStream must be closed when no longer needed to release resources.
	CreateStream(ctx context.Context, req *CompletionRequest) (CompletionStream, error)
}
```

## CompletionRequest

The `CompletionRequest` struct defines the parameters for a completion request:

```go
type CompletionRequest struct {
	Prompt      string      `json:"prompt"`           // The text prompt to generate completions for
	MaxTokens   int         `json:"max_tokens,omitempty"`    // Maximum tokens to generate
	Temperature float64     `json:"temperature,omitempty"`   // Controls randomness (0.0-2.0)
	TopP        float64     `json:"top_p,omitempty"`         // Alternative to temperature, nucleus sampling
	TopK        int         `json:"top_k,omitempty"`         // Limit to K most likely tokens
	Stream      bool        `json:"stream,omitempty"`        // Enable streaming responses
	Stop        []string    `json:"stop,omitempty"`          // Sequences where generation should stop
	Model       string      `json:"model,omitempty"`         // Model to use (if not default)
	JSONSchema  interface{} `json:"json_schema,omitempty"`   // JSON schema for structured output
}
```

### Request Parameters

| Parameter   | Type        | Description                                         | Default |
|-------------|-------------|-----------------------------------------------------|---------|
| Prompt      | string      | The text prompt to complete                         | (required) |
| MaxTokens   | int         | Maximum number of tokens to generate                | (model dependent) |
| Temperature | float64     | Controls randomness (higher = more random)          | 1.0 |
| TopP        | float64     | Nucleus sampling parameter (consider tokens with top_p probability mass) | 1.0 |
| TopK        | int         | Only sample from top K most likely tokens           | 0 (disabled) |
| Stream      | bool        | Enable streaming response (token-by-token)          | false |
| Stop        | []string    | Stop sequences to end generation when encountered   | [] |
| Model       | string      | Model ID to use (if multiple available)             | (currently loaded model) |
| JSONSchema  | interface{} | Schema for structured JSON output                   | nil |

## CompletionResponse

The response to a completion request contains the generated text and metadata:

```go
type CompletionResponse struct {
	ID      string                 `json:"id"`      // Unique identifier for the completion
	Object  string                 `json:"object"`  // Type of object (always "completion")
	Created int64                  `json:"created"` // Unix timestamp of creation
	Model   string                 `json:"model"`   // Model used for completion
	Choices []CompletionRespChoice `json:"choices"` // Generated completions
	Usage   *UsageStats            `json:"usage,omitempty"`  // Token usage statistics
}

type CompletionRespChoice struct {
	Text         string              `json:"text"`          // Generated text
	Index        int                 `json:"index"`         // Choice index
	FinishReason string              `json:"finish_reason,omitempty"` // Why generation ended
	LogProbs     *CompletionLogProbs `json:"logprobs,omitempty"`     // Log probabilities (if requested)
}

type UsageStats struct {
	PromptTokens     int `json:"prompt_tokens"`     // Tokens in the prompt
	CompletionTokens int `json:"completion_tokens"` // Tokens in the completion
	TotalTokens      int `json:"total_tokens"`      // Total tokens used
}
```

## Streaming Completions

For streaming completions, use `CreateStream` which returns chunks of the response as they're generated:

```go
type CompletionStream = Stream[*CompletionStreamResponse]

type CompletionStreamResponse struct {
	ID      string                   `json:"id"`
	Object  string                   `json:"object"`
	Created int64                    `json:"created"`
	Model   string                   `json:"model"`
	Choices []CompletionStreamChoice `json:"choices"`
}

type CompletionStreamChoice struct {
	Text         string `json:"text"`          // Text chunk
	Index        int    `json:"index"`         // Choice index
	FinishReason string `json:"finish_reason,omitempty"` // Present only in the final chunk
}
```

## Examples

### Basic Completion

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
	
	req := &tabby.CompletionRequest{
		Prompt:      "func fibonacci(n int) int {",
		MaxTokens:   100,
		Temperature: 0.7,
	}
	
	resp, err := client.Completions().Create(ctx, req)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	
	fmt.Printf("Generated Code:\n%s\n", resp.Choices[0].Text)
	
	if resp.Usage != nil {
		fmt.Printf("\nToken Usage:\n")
		fmt.Printf("  Prompt tokens: %d\n", resp.Usage.PromptTokens)
		fmt.Printf("  Completion tokens: %d\n", resp.Usage.CompletionTokens)
		fmt.Printf("  Total tokens: %d\n", resp.Usage.TotalTokens)
	}
}
```

### Streaming Completion

```go
package main

import (
	"context"
	"fmt"
	"io"
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
	
	req := &tabby.CompletionRequest{
		Prompt:      "Write a function to calculate the factorial of a number:",
		MaxTokens:   150,
		Temperature: 0.7,
		Stream:      true,
	}
	
	stream, err := client.Completions().CreateStream(ctx, req)
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}
	defer stream.Close()
	
	fmt.Println("Streaming response:")
	
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break // Stream completed
		}
		if err != nil {
			log.Fatalf("Error receiving from stream: %v", err)
		}
		
		// Print each chunk as it arrives
		if len(resp.Choices) > 0 {
			fmt.Print(resp.Choices[0].Text)
		}
	}
}
```

### Structured Output with JSON Schema

```go
package main

import (
	"context"
	"encoding/json"
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
	
	// Define JSON schema for structured output
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type": "string",
				"description": "Name of the person",
			},
			"age": map[string]interface{}{
				"type": "integer",
				"description": "Age of the person",
			},
			"skills": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "List of skills",
			},
		},
		"required": []string{"name", "age", "skills"},
	}
	
	req := &tabby.CompletionRequest{
		Prompt:      "Generate information about a software developer",
		MaxTokens:   150,
		Temperature: 0.7,
		JSONSchema:  schema,
	}
	
	resp, err := client.Completions().Create(ctx, req)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	
	fmt.Printf("Generated JSON:\n%s\n", resp.Choices[0].Text)
	
	// Parse the JSON response
	var person map[string]interface{}
	if err := json.Unmarshal([]byte(resp.Choices[0].Text), &person); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	
	fmt.Printf("\nParsed Data:\n")
	fmt.Printf("  Name: %s\n", person["name"])
	fmt.Printf("  Age: %v\n", person["age"])
	fmt.Printf("  Skills: %v\n", person["skills"])
}
```

## Error Handling

The service methods can return several types of errors:

- `*tabby.APIError`: Errors returned by the TabbyAPI server
- `*tabby.RequestError`: Errors related to HTTP requests
- `*tabby.ValidationError`: Errors due to invalid request parameters
- `*tabby.StreamError`: Errors during streaming operations

Example error handling:

```go
resp, err := client.Completions().Create(ctx, req)
if err != nil {
	// Use type assertions to handle different error types
	var apiErr *tabby.APIError
	if errors.As(err, &apiErr) {
		fmt.Printf("API Error (code %s): %s\n", apiErr.Code(), apiErr.Error())
		
		// Check for specific error codes
		if apiErr.Code() == "rate_limit_exceeded" {
			// Handle rate limiting
		}
	}
	
	var reqErr *tabby.RequestError
	if errors.As(err, &reqErr) {
		fmt.Printf("Request Error: %s\n", reqErr.Error())
		
		// Check for specific causes
		if errors.Is(reqErr.Unwrap(), context.DeadlineExceeded) {
			// Handle timeout
		}
	}
	
	log.Fatalf("Error: %v", err)
}
```

## Best Practices

1. **Set Appropriate Timeouts**: Use context with timeouts for all requests, especially for larger completions.

2. **Resource Management**: Always close streams when done:
   ```go
   stream, err := client.Completions().CreateStream(ctx, req)
   if err != nil {
       return err
   }
   defer stream.Close()
   ```

3. **Error Handling**: Implement robust error handling to manage various failure modes.

4. **Token Management**: Be mindful of token usage, especially for large contexts.

5. **Temperature Setting**:
   - Lower temperature (0.1-0.5) for deterministic/factual responses
   - Higher temperature (0.7-1.0) for creative content