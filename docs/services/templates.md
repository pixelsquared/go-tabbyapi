# Templates Service

The Templates service provides functionality for managing prompt templates in TabbyAPI. Prompt templates define how to format prompts for specific model architectures, which can significantly impact the quality and consistency of model outputs.

## Interface

```go
// TemplatesService handles prompt template management for different model types.
// Templates define how to format prompts for specific model architectures.
type TemplatesService interface {
	// List returns all available prompt templates.
	List(ctx context.Context) (*TemplateList, error)

	// Switch changes the active prompt template.
	Switch(ctx context.Context, req *TemplateSwitchRequest) error

	// Unload unloads the currently selected template.
	Unload(ctx context.Context) error
}
```

## Template Types

### TemplateList

The `TemplateList` struct contains a list of available templates:

```go
type TemplateList struct {
	Object string   `json:"object"` // Type of object (always "list")
	Data   []string `json:"data"`   // Array of template names
}
```

### TemplateSwitchRequest

To switch to a different template, use the `Switch` method with a `TemplateSwitchRequest`:

```go
type TemplateSwitchRequest struct {
	PromptTemplateName string `json:"prompt_template_name"` // Name of the template to activate
}
```

## Examples

### Listing Available Templates

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
		tabby.WithAdminKey("your-admin-key"), // Admin key required for template operations
	)
	defer client.Close()
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// List all available prompt templates
	templates, err := client.Templates().List(ctx)
	if err != nil {
		log.Fatalf("Error listing templates: %v", err)
	}
	
	fmt.Printf("Available prompt templates (%d):\n", len(templates.Data))
	for i, template := range templates.Data {
		fmt.Printf("%d. %s\n", i+1, template)
	}
}
```

### Switching Templates

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
	
	// First, get the list of available templates
	templates, err := client.Templates().List(ctx)
	if err != nil {
		log.Fatalf("Error listing templates: %v", err)
	}
	
	if len(templates.Data) == 0 {
		log.Fatalf("No templates available")
	}
	
	// Choose a template to use (for example, the first one)
	templateName := templates.Data[0]
	
	// Switch to the selected template
	switchReq := &tabby.TemplateSwitchRequest{
		PromptTemplateName: templateName,
	}
	
	fmt.Printf("Switching to template: %s\n", templateName)
	err = client.Templates().Switch(ctx, switchReq)
	if err != nil {
		log.Fatalf("Error switching template: %v", err)
	}
	
	fmt.Printf("Successfully switched to template: %s\n", templateName)
}
```

### Unloading a Template

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
	
	// Unload the currently active template
	fmt.Println("Unloading the current template...")
	err := client.Templates().Unload(ctx)
	if err != nil {
		log.Fatalf("Error unloading template: %v", err)
	}
	
	fmt.Println("Template unloaded successfully")
}
```

### Using Templates with Generation

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
	
	// List available templates
	templates, err := client.Templates().List(ctx)
	if err != nil {
		log.Fatalf("Error listing templates: %v", err)
	}
	
	// Find a specific template
	var templateName string
	for _, t := range templates.Data {
		if t == "llama" || t == "llama2" {
			templateName = t
			break
		}
	}
	
	if templateName == "" {
		log.Fatalf("No suitable template found")
	}
	
	// Switch to the selected template
	switchReq := &tabby.TemplateSwitchRequest{
		PromptTemplateName: templateName,
	}
	
	fmt.Printf("Switching to template: %s\n", templateName)
	err = client.Templates().Switch(ctx, switchReq)
	if err != nil {
		log.Fatalf("Error switching template: %v", err)
	}
	
	// Now generate text using the selected template
	fmt.Println("Generating text with the selected template...")
	
	chatReq := &tabby.ChatCompletionRequest{
		Messages: []tabby.ChatMessage{
			{
				Role:    tabby.ChatMessageRoleSystem,
				Content: "You are a helpful assistant.",
			},
			{
				Role:    tabby.ChatMessageRoleUser,
				Content: "What's the capital of France?",
			},
		},
		MaxTokens:   100,
		Temperature: 0.7,
	}
	
	chatResp, err := client.Chat().Create(ctx, chatReq)
	if err != nil {
		log.Fatalf("Error generating text: %v", err)
	}
	
	if len(chatResp.Choices) > 0 {
		fmt.Println("\nGenerated text with template:", templateName)
		fmt.Println(chatResp.Choices[0].Message.Content)
	}
	
	// Optionally unload the template when done
	err = client.Templates().Unload(ctx)
	if err != nil {
		log.Printf("Warning: Error unloading template: %v", err)
	}
}
```

### Testing Different Templates

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
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	
	// List available templates
	templates, err := client.Templates().List(ctx)
	if err != nil {
		log.Fatalf("Error listing templates: %v", err)
	}
	
	// Define a standard prompt to test with each template
	chatReq := &tabby.ChatCompletionRequest{
		Messages: []tabby.ChatMessage{
			{
				Role:    tabby.ChatMessageRoleSystem,
				Content: "You are a helpful assistant.",
			},
			{
				Role:    tabby.ChatMessageRoleUser,
				Content: "Write a short poem about artificial intelligence.",
			},
		},
		MaxTokens:   150,
		Temperature: 0.7,
	}
	
	// Try each template (limit to max 3 for this example)
	maxTemplates := 3
	if len(templates.Data) < maxTemplates {
		maxTemplates = len(templates.Data)
	}
	
	for i := 0; i < maxTemplates; i++ {
		templateName := templates.Data[i]
		
		// Skip non-chat templates
		if templateName == "completion" || templateName == "raw" {
			continue
		}
		
		// Switch to this template
		switchReq := &tabby.TemplateSwitchRequest{
			PromptTemplateName: templateName,
		}
		
		fmt.Printf("\n--- Testing template: %s ---\n", templateName)
		err = client.Templates().Switch(ctx, switchReq)
		if err != nil {
			fmt.Printf("Error switching to template %s: %v\n", templateName, err)
			continue
		}
		
		// Generate text with this template
		chatResp, err := client.Chat().Create(ctx, chatReq)
		if err != nil {
			fmt.Printf("Error generating text with template %s: %v\n", templateName, err)
			continue
		}
		
		if len(chatResp.Choices) > 0 {
			fmt.Printf("\nResponse with template '%s':\n", templateName)
			fmt.Println(chatResp.Choices[0].Message.Content)
		}
		
		// Unload this template before trying the next one
		_ = client.Templates().Unload(ctx)
		
		// Add a small delay between template switches
		time.Sleep(2 * time.Second)
	}
}
```

## Error Handling

The Templates service operations may fail due to various reasons:

```go
err := client.Templates().Switch(ctx, switchReq)
if err != nil {
	var apiErr *tabby.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.Code() {
		case "permission_error":
			log.Fatalf("You need admin permissions to manage templates")
		case "not_found":
			log.Fatalf("Template not found: %s", switchReq.PromptTemplateName)
		case "invalid_request":
			log.Fatalf("Invalid request parameters: %s", apiErr.Error())
		default:
			log.Fatalf("API Error: %s", apiErr.Error())
		}
	} else {
		log.Fatalf("Error switching template: %v", err)
	}
}
```

## Best Practices

1. **Admin Authentication**: Template operations require admin permissions:
   ```go
   client := tabby.NewClient(
       tabby.WithAdminKey("your-admin-key"),
   )
   ```

2. **Template Selection**: Choose the appropriate template for your model architecture:
   - Use model-specific templates (e.g., "llama", "mistral", "vicuna") for best results
   - Check the model's documentation for recommended templates

3. **Clear State**: Unload templates when switching between different operations:
   ```go
   defer func() {
       ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
       defer cancel()
       _ = client.Templates().Unload(ctx)
   }()
   ```

4. **Error Recovery**: Implement robust error handling for template operations:
   ```go
   if err != nil {
       // Log the error
       log.Printf("Error with template operation: %v", err)
       
       // Attempt to reset to a default state
       _ = client.Templates().Unload(ctx)
   }
   ```

5. **Template Testing**: Compare model outputs with different templates to find the best one for your specific use case.

6. **Template-Model Compatibility**: Ensure that the selected template is compatible with the loaded model to avoid unexpected behaviors or poor-quality outputs.

7. **Consistency**: For production applications, use the same template consistently to ensure predictable model behavior, or test thoroughly if changing templates.