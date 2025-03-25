# Lora Service

The Lora service provides functionality for managing Low-Rank Adaptation (LoRA) adapters. LoRA is a technique for fine-tuning large language models with significantly fewer parameters, making it more efficient than full fine-tuning.

## Interface

```go
// LoraService handles Low-Rank Adaptation (LoRA) adapter management.
// LoRA adapters allow fine-tuning language models with significantly fewer parameters.
type LoraService interface {
	// List returns all available LoRA adapters.
	List(ctx context.Context) (*LoraList, error)

	// GetActive returns the currently loaded LoRA adapters.
	GetActive(ctx context.Context) (*LoraList, error)

	// Load loads specified LoRA adapters.
	Load(ctx context.Context, req *LoraLoadRequest) (*LoraLoadResponse, error)

	// Unload unloads all currently loaded LoRA adapters.
	Unload(ctx context.Context) error
}
```

## LoRA Types

### LoraCard

The `LoraCard` struct provides information about a LoRA adapter:

```go
type LoraCard struct {
	ID      string  `json:"id"`      // Unique identifier for the adapter
	Object  string  `json:"object"`  // Type of object (always "lora")
	Created int64   `json:"created"` // Unix timestamp of creation
	OwnedBy string  `json:"owned_by"` // Owner of the adapter
	Scaling float64 `json:"scaling,omitempty"` // Current scaling factor
}
```

### LoraList

The `LoraList` struct contains a list of LoRA adapters:

```go
type LoraList struct {
	Object string     `json:"object"` // Type of object (always "list")
	Data   []LoraCard `json:"data"`   // Array of LoRA adapter cards
}
```

## Loading LoRA Adapters

### LoraLoadRequest

To load LoRA adapters, use the `Load` method with a `LoraLoadRequest`:

```go
type LoraLoadRequest struct {
	Loras     []LoraLoadInfo `json:"loras"`     // Array of LoRA adapters to load
	SkipQueue bool           `json:"skip_queue,omitempty"` // Skip waiting in queue
}

type LoraLoadInfo struct {
	Name    string  `json:"name"`    // Name of the LoRA adapter
	Scaling float64 `json:"scaling,omitempty"` // Scaling factor (usually between 0 and 1)
}
```

The `Scaling` parameter controls the strength of the adaptation. Higher values make the fine-tuning more pronounced, while lower values retain more of the base model's behavior.

### LoraLoadResponse

Loading a LoRA adapter returns a `LoraLoadResponse`:

```go
type LoraLoadResponse struct {
	Success []string `json:"success"` // Array of successfully loaded adapters
	Failure []string `json:"failure"` // Array of adapters that failed to load
}
```

## Examples

### Listing Available LoRA Adapters

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
		tabby.WithAdminKey("your-admin-key"), // Admin key required for LoRA operations
	)
	defer client.Close()
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// List all available LoRA adapters
	loras, err := client.Lora().List(ctx)
	if err != nil {
		log.Fatalf("Error listing LoRA adapters: %v", err)
	}
	
	fmt.Printf("Found %d available LoRA adapters:\n", len(loras.Data))
	for i, lora := range loras.Data {
		fmt.Printf("%d. %s (owned by: %s)\n", i+1, lora.ID, lora.OwnedBy)
		fmt.Printf("   Created: %s\n", time.Unix(lora.Created, 0).Format(time.RFC3339))
	}
}
```

### Getting Currently Active LoRA Adapters

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
		tabby.WithAdminKey("your-admin-key"),
	)
	defer client.Close()
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Get currently active LoRA adapters
	active, err := client.Lora().GetActive(ctx)
	if err != nil {
		log.Fatalf("Error getting active LoRA adapters: %v", err)
	}
	
	if len(active.Data) == 0 {
		fmt.Println("No LoRA adapters are currently active")
	} else {
		fmt.Printf("Currently active LoRA adapters (%d):\n", len(active.Data))
		for i, lora := range active.Data {
			fmt.Printf("%d. %s (scaling: %.2f)\n", i+1, lora.ID, lora.Scaling)
		}
	}
}
```

### Loading LoRA Adapters

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
		tabby.WithAdminKey("your-admin-key"),
	)
	defer client.Close()
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Load LoRA adapters with different scaling factors
	loadReq := &tabby.LoraLoadRequest{
		Loras: []tabby.LoraLoadInfo{
			{
				Name:    "code-lora", // A coding-focused LoRA
				Scaling: 0.8,        // Strong influence
			},
			{
				Name:    "writing-style-lora", // A writing style LoRA
				Scaling: 0.5,                 // Moderate influence
			},
		},
		SkipQueue: false, // Wait if other operations are in progress
	}
	
	fmt.Println("Loading LoRA adapters...")
	resp, err := client.Lora().Load(ctx, loadReq)
	if err != nil {
		log.Fatalf("Error loading LoRA adapters: %v", err)
	}
	
	// Check which adapters were loaded successfully
	if len(resp.Success) > 0 {
		fmt.Printf("Successfully loaded adapters: %v\n", resp.Success)
	}
	
	// Check if any adapters failed to load
	if len(resp.Failure) > 0 {
		fmt.Printf("Failed to load adapters: %v\n", resp.Failure)
	}
	
	// Verify the loaded adapters
	active, err := client.Lora().GetActive(ctx)
	if err != nil {
		log.Fatalf("Error getting active LoRA adapters: %v", err)
	}
	
	fmt.Printf("Currently active LoRA adapters (%d):\n", len(active.Data))
	for i, lora := range active.Data {
		fmt.Printf("%d. %s (scaling: %.2f)\n", i+1, lora.ID, lora.Scaling)
	}
}
```

### Unloading LoRA Adapters

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
		tabby.WithAdminKey("your-admin-key"),
	)
	defer client.Close()
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Check currently active LoRA adapters before unloading
	activeBefore, err := client.Lora().GetActive(ctx)
	if err != nil {
		log.Fatalf("Error getting active LoRA adapters: %v", err)
	}
	
	if len(activeBefore.Data) == 0 {
		fmt.Println("No LoRA adapters are currently active")
		return
	}
	
	fmt.Printf("Currently active LoRA adapters (%d):\n", len(activeBefore.Data))
	for i, lora := range activeBefore.Data {
		fmt.Printf("%d. %s\n", i+1, lora.ID)
	}
	
	// Unload all active LoRA adapters
	fmt.Println("Unloading all LoRA adapters...")
	err = client.Lora().Unload(ctx)
	if err != nil {
		log.Fatalf("Error unloading LoRA adapters: %v", err)
	}
	
	fmt.Println("Successfully unloaded all LoRA adapters")
	
	// Verify that no adapters are active
	activeAfter, err := client.Lora().GetActive(ctx)
	if err != nil {
		log.Fatalf("Error getting active LoRA adapters: %v", err)
	}
	
	if len(activeAfter.Data) == 0 {
		fmt.Println("Confirmed: No LoRA adapters are active")
	} else {
		fmt.Printf("Warning: %d LoRA adapters are still active\n", len(activeAfter.Data))
	}
}
```

### Using LoRA with Text Generation

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
		tabby.WithAdminKey("your-admin-key"),
	)
	defer client.Close()
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	
	// Step 1: Load a LoRA adapter
	loadReq := &tabby.LoraLoadRequest{
		Loras: []tabby.LoraLoadInfo{
			{
				Name:    "coding-assistant-lora",
				Scaling: 0.8,
			},
		},
	}
	
	fmt.Println("Loading LoRA adapter...")
	loadResp, err := client.Lora().Load(ctx, loadReq)
	if err != nil {
		log.Fatalf("Error loading LoRA adapter: %v", err)
	}
	
	if len(loadResp.Success) == 0 {
		log.Fatalf("Failed to load LoRA adapter")
	}
	
	fmt.Printf("Successfully loaded LoRA adapter: %s\n", loadResp.Success[0])
	
	// Step 2: Use the model with the loaded LoRA for text generation
	fmt.Println("Generating text with LoRA-adapted model...")
	
	chatReq := &tabby.ChatCompletionRequest{
		Messages: []tabby.ChatMessage{
			{
				Role:    tabby.ChatMessageRoleSystem,
				Content: "You are a coding assistant specialized in Python.",
			},
			{
				Role:    tabby.ChatMessageRoleUser,
				Content: "Write a function to find the nth Fibonacci number using recursion with memoization.",
			},
		},
		MaxTokens:   500,
		Temperature: 0.7,
	}
	
	chatResp, err := client.Chat().Create(ctx, chatReq)
	if err != nil {
		log.Fatalf("Error generating text: %v", err)
	}
	
	if len(chatResp.Choices) > 0 {
		fmt.Println("\nGenerated text with LoRA-adapted model:")
		fmt.Println(chatResp.Choices[0].Message.Content)
	}
	
	// Step 3: Unload the LoRA adapter when done
	fmt.Println("\nUnloading LoRA adapter...")
	err = client.Lora().Unload(ctx)
	if err != nil {
		log.Fatalf("Error unloading LoRA adapter: %v", err)
	}
	
	fmt.Println("LoRA adapter unloaded successfully")
}
```

## Error Handling

The Lora service operations may fail due to various reasons:

```go
resp, err := client.Lora().Load(ctx, loadReq)
if err != nil {
	var apiErr *tabby.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.Code() {
		case "permission_error":
			log.Fatalf("You need admin permissions to manage LoRA adapters")
		case "not_found":
			log.Fatalf("One or more of the requested LoRA adapters were not found")
		case "invalid_request":
			log.Fatalf("Invalid request parameters: %s", apiErr.Error())
		default:
			log.Fatalf("API Error: %s", apiErr.Error())
		}
	} else {
		log.Fatalf("Error loading LoRA adapters: %v", err)
	}
}
```

## Best Practices

1. **Admin Authentication**: LoRA operations require admin permissions:
   ```go
   client := tabby.NewClient(
       tabby.WithAdminKey("your-admin-key"),
   )
   ```

2. **Scaling Factors**: Use appropriate scaling factors based on your needs:
   - Higher values (0.7-1.0) for stronger influence of the LoRA adaptation
   - Lower values (0.1-0.5) for subtle influence while preserving base model behavior
   - Multiple adapters can be combined with different scaling factors

3. **Resource Management**: Unload LoRA adapters when they're no longer needed to free up resources:
   ```go
   defer func() {
       ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
       defer cancel()
       _ = client.Lora().Unload(ctx)
   }()
   ```

4. **Error Handling**: Implement robust error handling, especially when loading multiple adapters:
   ```go
   if len(resp.Failure) > 0 {
       fmt.Printf("Warning: Failed to load LoRA adapters: %v\n", resp.Failure)
       // Decide whether to continue or abort based on which adapters failed
   }
   ```

5. **LoRA Compatibility**: Ensure that LoRA adapters are compatible with the loaded base model.

6. **Performance Considerations**:
   - Loading multiple LoRA adapters increases memory usage and may impact performance
   - Be cautious with the number of simultaneous adapters used
   
7. **Storage Management**: Periodically check available LoRA adapters and remove unused ones to save storage space.

8. **Testing LoRA Effects**: Test the effect of different scaling factors on output quality:
   ```go
   // Compare outputs with different scaling factors
   for _, scaling := range []float64{0.2, 0.5, 0.8} {
       req := &tabby.LoraLoadRequest{
           Loras: []tabby.LoraLoadInfo{
               {
                   Name:    "your-lora",
                   Scaling: scaling,
               },
           },
       }
       // Load and test with each scaling factor
   }