# Tokens Service

The Tokens service provides functionality for token-level operations, including encoding text into token IDs and decoding token IDs back into text. This service is particularly useful for understanding how models process text, calculating token usage, and manipulating text at the token level.

## Interface

```go
// TokensService handles tokenization operations including encoding text to token IDs
// and decoding token IDs back to text.
type TokensService interface {
	// Encode converts text into token IDs.
	Encode(ctx context.Context, req *TokenEncodeRequest) (*TokenEncodeResponse, error)

	// Decode converts token IDs back into text.
	Decode(ctx context.Context, req *TokenDecodeRequest) (*TokenDecodeResponse, error)
}
```

## Token Types

### TokenEncodeRequest

The `TokenEncodeRequest` struct defines the parameters for an encoding request:

```go
type TokenEncodeRequest struct {
	Text                interface{} `json:"text"` // String or []ChatMessage
	AddBOSToken         bool        `json:"add_bos_token,omitempty"`         // Add beginning-of-sequence token
	EncodeSpecialTokens bool        `json:"encode_special_tokens,omitempty"` // Include special tokens in encoding
	DecodeSpecialTokens bool        `json:"decode_special_tokens,omitempty"` // Include special tokens in decoding
}
```

The `Text` field can be either:
- A string for simple text encoding
- An array of `ChatMessage` objects for conversation encoding

### TokenEncodeResponse

The response to a token encoding request:

```go
type TokenEncodeResponse struct {
	Tokens []int `json:"tokens"` // Array of token IDs
	Length int   `json:"length"` // Number of tokens
}
```

### TokenDecodeRequest

The `TokenDecodeRequest` struct defines the parameters for a decoding request:

```go
type TokenDecodeRequest struct {
	Tokens              []int `json:"tokens"`                            // Array of token IDs to decode
	AddBOSToken         bool  `json:"add_bos_token,omitempty"`          // Add beginning-of-sequence token
	EncodeSpecialTokens bool  `json:"encode_special_tokens,omitempty"`  // Include special tokens in encoding
	DecodeSpecialTokens bool  `json:"decode_special_tokens,omitempty"`  // Include special tokens in decoding
}
```

### TokenDecodeResponse

The response to a token decoding request:

```go
type TokenDecodeResponse struct {
	Text string `json:"text"` // Decoded text
}
```

## Examples

### Encoding Text to Tokens

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
	
	// Create a token encode request for simple text
	req := &tabby.TokenEncodeRequest{
		Text:        "Hello, world! How are large language models tokenized?",
		AddBOSToken: true, // Add beginning-of-sequence token
	}
	
	resp, err := client.Tokens().Encode(ctx, req)
	if err != nil {
		log.Fatalf("Error encoding tokens: %v", err)
	}
	
	fmt.Printf("Encoded %d tokens: %v\n", resp.Length, resp.Tokens)
}
```

### Encoding Chat Messages

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
	
	// Create chat messages to encode
	messages := []tabby.ChatMessage{
		{
			Role:    tabby.ChatMessageRoleSystem,
			Content: "You are a helpful assistant.",
		},
		{
			Role:    tabby.ChatMessageRoleUser,
			Content: "Can you explain how tokenization works?",
		},
		{
			Role:    tabby.ChatMessageRoleAssistant,
			Content: "Tokenization is the process of converting text into numerical token IDs.",
		},
	}
	
	// Encode the chat messages
	req := &tabby.TokenEncodeRequest{
		Text:                messages,
		AddBOSToken:         true,
		EncodeSpecialTokens: true,
	}
	
	resp, err := client.Tokens().Encode(ctx, req)
	if err != nil {
		log.Fatalf("Error encoding chat messages: %v", err)
	}
	
	fmt.Printf("Encoded chat conversation into %d tokens\n", resp.Length)
	fmt.Printf("Token IDs: %v\n", resp.Tokens)
}
```

### Decoding Tokens to Text

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
	
	// Sample token IDs to decode
	tokenIds := []int{15043, 3186, 29889, 29871, 29896, 29947, 29947, 29991, 29947}
	
	// Create a token decode request
	req := &tabby.TokenDecodeRequest{
		Tokens:              tokenIds,
		DecodeSpecialTokens: true,
	}
	
	resp, err := client.Tokens().Decode(ctx, req)
	if err != nil {
		log.Fatalf("Error decoding tokens: %v", err)
	}
	
	fmt.Printf("Decoded text: %s\n", resp.Text)
}
```

### Encode-Decode Roundtrip

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
	
	// Original text to encode and then decode
	originalText := "This is a test of encode-decode roundtrip. Will the text be preserved exactly?"
	
	// Step 1: Encode the text to tokens
	encodeReq := &tabby.TokenEncodeRequest{
		Text: originalText,
	}
	
	encodeResp, err := client.Tokens().Encode(ctx, encodeReq)
	if err != nil {
		log.Fatalf("Error encoding tokens: %v", err)
	}
	
	fmt.Printf("Original text: %s\n", originalText)
	fmt.Printf("Encoded into %d tokens: %v\n", encodeResp.Length, encodeResp.Tokens)
	
	// Step 2: Decode the tokens back to text
	decodeReq := &tabby.TokenDecodeRequest{
		Tokens: encodeResp.Tokens,
	}
	
	decodeResp, err := client.Tokens().Decode(ctx, decodeReq)
	if err != nil {
		log.Fatalf("Error decoding tokens: %v", err)
	}
	
	fmt.Printf("Decoded text: %s\n", decodeResp.Text)
	
	// Step 3: Compare original and decoded text
	if originalText == decodeResp.Text {
		fmt.Println("✅ Roundtrip successful! Text was preserved exactly.")
	} else {
		fmt.Println("❌ Roundtrip changed the text.")
		fmt.Println("   This is normal for some tokenizers, which may normalize whitespace or special characters.")
	}
}
```

### Token Usage Calculator

```go
package main

import (
	"context"
	"fmt"
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
	
	// Function to count tokens in a text
	countTokens := func(text string) (int, error) {
		req := &tabby.TokenEncodeRequest{Text: text}
		resp, err := client.Tokens().Encode(ctx, req)
		if err != nil {
			return 0, err
		}
		return resp.Length, nil
	}
	
	// Example texts to analyze
	texts := []string{
		"This is a short sentence.",
		"This is a longer sentence with more words and should therefore use more tokens when encoded by the language model.",
		strings.Repeat("This is repeated text. ", 10),
		"Technical terms like 'tokenization', 'embeddings', and 'transformer architecture' might be encoded differently.",
		"Code snippets: ```python\ndef hello_world():\n    print('Hello, world!')\n```",
	}
	
	fmt.Println("Token Usage Analysis:")
	fmt.Println("======================")
	
	for i, text := range texts {
		count, err := countTokens(text)
		if err != nil {
			log.Printf("Error counting tokens for text %d: %v", i+1, err)
			continue
		}
		
		// Calculate the token-to-character ratio
		ratio := float64(count) / float64(len(text))
		
		fmt.Printf("\nText %d:\n", i+1)
		fmt.Printf("- Characters: %d\n", len(text))
		fmt.Printf("- Tokens: %d\n", count)
		fmt.Printf("- Tokens/Char Ratio: %.3f\n", ratio)
		fmt.Printf("- First 50 chars: %s\n", text[:min(50, len(text))])
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
```

## Error Handling

The Tokens service operations may fail due to various reasons:

```go
resp, err := client.Tokens().Encode(ctx, req)
if err != nil {
	var apiErr *tabby.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.Code() {
		case "invalid_request":
			log.Fatalf("Invalid request parameters: %s", apiErr.Error())
		case "not_found":
			log.Fatalf("No model is loaded, or the model doesn't have a tokenizer")
		default:
			log.Fatalf("API Error: %s", apiErr.Error())
		}
	} else {
		log.Fatalf("Error encoding tokens: %v", err)
	}
}
```

## Best Practices

1. **Model Dependency**: Tokenization depends on the currently loaded model, so ensure a model is loaded before using the Tokens service.

2. **Special Tokens**: Be mindful of special tokens handling:
   - Use `AddBOSToken: true` when that's consistent with how you'll use the text in completions/chat
   - Special tokens like BOS (beginning of sequence), EOS (end of sequence), etc. can affect the token count

3. **Performance**: For large batches of texts, send them as a single request rather than multiple requests:
   ```go
   // Efficient (single request)
   req := &tabby.TokenEncodeRequest{
       Text: []string{"text1", "text2", "text3"},
   }
   
   // Less efficient (multiple requests)
   for _, text := range texts {
       req := &tabby.TokenEncodeRequest{Text: text}
       // ...
   }
   ```

4. **Token Counting**: Use the Tokens service to calculate exactly how many tokens a prompt will use before sending generation requests:
   ```go
   func calculateTokenUsage(prompt string) (int, error) {
       req := &tabby.TokenEncodeRequest{Text: prompt}
       resp, err := client.Tokens().Encode(ctx, req)
       if err != nil {
           return 0, err
       }
       return resp.Length, nil
   }
   
   // Use this to stay within token limits
   maxPromptTokens := 4096
   promptTokens, _ := calculateTokenUsage(prompt)
   maxGeneratedTokens := maxPromptTokens - promptTokens
   
   if maxGeneratedTokens <= 0 {
       log.Fatal("Prompt is too long")
   }
   
   // Adjust request parameters
   completionReq := &tabby.CompletionRequest{
       Prompt:    prompt,
       MaxTokens: maxGeneratedTokens,
   }
   ```

5. **Token Inspection**: Use the Tokens service to inspect and debug how models tokenize specific text:
   ```go
   func inspectTokenization(text string) error {
       // Encode the text
       encodeResp, err := client.Tokens().Encode(ctx, &tabby.TokenEncodeRequest{Text: text})
       if err != nil {
           return err
       }
       
       // For each token, decode it back to see the token boundaries
       for i, token := range encodeResp.Tokens {
           decodeResp, err := client.Tokens().Decode(ctx, &tabby.TokenDecodeRequest{
               Tokens: []int{token},
           })
           if err != nil {
               return err
           }
           fmt.Printf("Token %d: ID=%d, Text='%s'\n", i, token, decodeResp.Text)
       }
       return nil
   }
   ```

6. **Chat Message Encoding**: When encoding chat messages, the tokens include any template formatting the model uses:
   ```go
   // This includes tokens for the template format, not just the raw text
   messages := []tabby.ChatMessage{
       {Role: tabby.ChatMessageRoleSystem, Content: "You are an assistant."},
       {Role: tabby.ChatMessageRoleUser, Content: "Hello"},
   }
   
   req := &tabby.TokenEncodeRequest{Text: messages}
   ```

7. **Error Recovery**: Implement robust error handling, especially for tokenization which depends on model state:
   ```go
   resp, err := client.Tokens().Encode(ctx, req)
   if err != nil {
       log.Printf("Error encoding tokens: %v", err)
       
       // Check if a model is loaded
       _, modelErr := client.Models().Get(ctx)
       if modelErr != nil {
           log.Printf("No model loaded, attempting to load default model")
           // Attempt to load a model
       }
   }