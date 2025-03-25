# Sampling Service

The Sampling service provides functionality for managing sampling parameters that control text generation behavior. It allows setting and overriding default sampling parameters globally or using predefined presets.

## Interface

```go
// SamplingService handles sampling parameter overrides and presets
// to control text generation behavior across different requests.
type SamplingService interface {
	// ListOverrides returns all sampler overrides and presets.
	ListOverrides(ctx context.Context) (*SamplerOverrideListResponse, error)

	// SwitchOverride changes the active sampler override.
	SwitchOverride(ctx context.Context, req *SamplerOverrideSwitchRequest) error

	// UnloadOverride unloads the currently selected override preset.
	UnloadOverride(ctx context.Context) error
}
```

## Sampling Types

### SamplerOverrideListResponse

The `SamplerOverrideListResponse` struct provides information about available and current sampling parameters:

```go
type SamplerOverrideListResponse struct {
	SelectedPreset string                 `json:"selected_preset,omitempty"` // Currently selected preset
	Overrides      map[string]interface{} `json:"overrides"`                 // Current overrides
	Presets        []string               `json:"presets"`                   // Available presets
}
```

### SamplerOverrideSwitchRequest

To change the active sampling parameters, use the `SwitchOverride` method with a `SamplerOverrideSwitchRequest`:

```go
type SamplerOverrideSwitchRequest struct {
	Preset    string                 `json:"preset,omitempty"`    // Preset name to use
	Overrides map[string]interface{} `json:"overrides,omitempty"` // Custom override values
}
```

You can either specify a named preset or provide custom override values. Common override parameters include:

- `temperature`: Controls randomness (higher values = more random)
- `top_p`: Controls nucleus sampling (consider tokens with top_p probability mass)
- `top_k`: Limit to top K most likely tokens
- `repetition_penalty`: Penalty for repeating tokens
- `presence_penalty`: Penalty for tokens already in the text
- `frequency_penalty`: Penalty for frequent tokens

## Examples

### Listing Available Sampling Presets and Overrides

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
		tabby.WithAdminKey("your-admin-key"), // Admin key required for sampling operations
	)
	defer client.Close()
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// List all available sampling presets and current overrides
	samplingInfo, err := client.Sampling().ListOverrides(ctx)
	if err != nil {
		log.Fatalf("Error listing sampling overrides: %v", err)
	}
	
	// Display available presets
	fmt.Printf("Available sampling presets (%d):\n", len(samplingInfo.Presets))
	for i, preset := range samplingInfo.Presets {
		fmt.Printf("%d. %s\n", i+1, preset)
	}
	
	// Display current preset (if any)
	if samplingInfo.SelectedPreset != "" {
		fmt.Printf("\nCurrently selected preset: %s\n", samplingInfo.SelectedPreset)
	} else {
		fmt.Println("\nNo preset is currently selected")
	}
	
	// Display current overrides (if any)
	if len(samplingInfo.Overrides) > 0 {
		fmt.Println("\nCurrent overrides:")
		for param, value := range samplingInfo.Overrides {
			fmt.Printf("- %s: %v\n", param, value)
		}
	} else {
		fmt.Println("\nNo overrides are currently set")
	}
}
```

### Switching to a Preset

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
	
	// First, get the list of available presets
	samplingInfo, err := client.Sampling().ListOverrides(ctx)
	if err != nil {
		log.Fatalf("Error listing sampling presets: %v", err)
	}
	
	if len(samplingInfo.Presets) == 0 {
		log.Fatalf("No sampling presets available")
	}
	
	// Choose a preset (for example, the first one)
	presetName := samplingInfo.Presets[0]
	
	// Switch to the selected preset
	switchReq := &tabby.SamplerOverrideSwitchRequest{
		Preset: presetName,
	}
	
	fmt.Printf("Switching to preset: %s\n", presetName)
	err = client.Sampling().SwitchOverride(ctx, switchReq)
	if err != nil {
		log.Fatalf("Error switching sampling preset: %v", err)
	}
	
	fmt.Printf("Successfully switched to preset: %s\n", presetName)
	
	// Verify the switch
	updatedInfo, err := client.Sampling().ListOverrides(ctx)
	if err != nil {
		log.Fatalf("Error getting updated sampling info: %v", err)
	}
	
	fmt.Printf("Current preset: %s\n", updatedInfo.SelectedPreset)
	fmt.Println("Applied overrides:")
	for param, value := range updatedInfo.Overrides {
		fmt.Printf("- %s: %v\n", param, value)
	}
}
```

### Setting Custom Overrides

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
	
	// Define custom sampling parameters
	customOverrides := map[string]interface{}{
		"temperature":        0.5,
		"top_p":              0.9,
		"top_k":              50,
		"repetition_penalty": 1.1,
	}
	
	// Set custom overrides
	switchReq := &tabby.SamplerOverrideSwitchRequest{
		Overrides: customOverrides,
	}
	
	fmt.Println("Setting custom sampling overrides...")
	err := client.Sampling().SwitchOverride(ctx, switchReq)
	if err != nil {
		log.Fatalf("Error setting sampling overrides: %v", err)
	}
	
	fmt.Println("Successfully set custom sampling overrides")
	
	// Verify the overrides
	updatedInfo, err := client.Sampling().ListOverrides(ctx)
	if err != nil {
		log.Fatalf("Error getting updated sampling info: %v", err)
	}
	
	fmt.Println("Applied overrides:")
	for param, value := range updatedInfo.Overrides {
		fmt.Printf("- %s: %v\n", param, value)
	}
}
```

### Unloading Sampling Overrides

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
	
	// Check current overrides before unloading
	beforeInfo, err := client.Sampling().ListOverrides(ctx)
	if err != nil {
		log.Fatalf("Error getting sampling info: %v", err)
	}
	
	if beforeInfo.SelectedPreset == "" && len(beforeInfo.Overrides) == 0 {
		fmt.Println("No sampling overrides are currently active")
		return
	}
	
	// Unload the current sampling overrides
	fmt.Println("Unloading current sampling overrides...")
	err = client.Sampling().UnloadOverride(ctx)
	if err != nil {
		log.Fatalf("Error unloading sampling overrides: %v", err)
	}
	
	fmt.Println("Successfully unloaded sampling overrides")
	
	// Verify the unload
	afterInfo, err := client.Sampling().ListOverrides(ctx)
	if err != nil {
		log.Fatalf("Error getting updated sampling info: %v", err)
	}
	
	if afterInfo.SelectedPreset == "" && len(afterInfo.Overrides) == 0 {
		fmt.Println("Confirmed: No sampling overrides are active")
	} else {
		fmt.Println("Warning: Sampling overrides are still active")
	}
}
```

### Testing Different Sampling Parameters

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
	
	// Define a standard prompt to test with different sampling parameters
	prompt := "Once upon a time in a distant galaxy,"
	
	// Define different sampling parameter sets to test
	testParams := []map[string]interface{}{
		{
			"temperature": 0.1, // Very deterministic
			"top_p":       0.9,
		},
		{
			"temperature": 0.7, // Balanced
			"top_p":       0.9,
		},
		{
			"temperature": 1.4, // Very creative/random
			"top_p":       0.95,
		},
	}
	
	// Try each parameter set
	for i, params := range testParams {
		// Set the sampling overrides
		switchReq := &tabby.SamplerOverrideSwitchRequest{
			Overrides: params,
		}
		
		fmt.Printf("\n--- Testing sampling parameters %d ---\n", i+1)
		fmt.Printf("Parameters: %v\n", params)
		
		err := client.Sampling().SwitchOverride(ctx, switchReq)
		if err != nil {
			fmt.Printf("Error setting sampling parameters: %v\n", err)
			continue
		}
		
		// Generate text with these parameters
		completionReq := &tabby.CompletionRequest{
			Prompt:    prompt,
			MaxTokens: 50,
		}
		
		completionResp, err := client.Completions().Create(ctx, completionReq)
		if err != nil {
			fmt.Printf("Error generating text: %v\n", err)
			continue
		}
		
		if len(completionResp.Choices) > 0 {
			fmt.Printf("\nGenerated text with temperature=%.1f:\n", params["temperature"])
			fmt.Printf("%s%s\n", prompt, completionResp.Choices[0].Text)
		}
		
		// Reset overrides before trying the next set
		_ = client.Sampling().UnloadOverride(ctx)
		
		// Add a small delay between tests
		time.Sleep(2 * time.Second)
	}
	
	// Final cleanup
	err := client.Sampling().UnloadOverride(ctx)
	if err != nil {
		log.Printf("Warning: Error unloading final sampling overrides: %v", err)
	}
}
```

## Error Handling

The Sampling service operations may fail due to various reasons:

```go
err := client.Sampling().SwitchOverride(ctx, switchReq)
if err != nil {
	var apiErr *tabby.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.Code() {
		case "permission_error":
			log.Fatalf("You need admin permissions to manage sampling parameters")
		case "not_found":
			log.Fatalf("Preset not found: %s", switchReq.Preset)
		case "invalid_request":
			log.Fatalf("Invalid request parameters: %s", apiErr.Error())
		default:
			log.Fatalf("API Error: %s", apiErr.Error())
		}
	} else {
		log.Fatalf("Error switching sampling parameters: %v", err)
	}
}
```

## Best Practices

1. **Admin Authentication**: Sampling operations require admin permissions:
   ```go
   client := tabby.NewClient(
       tabby.WithAdminKey("your-admin-key"),
   )
   ```

2. **Parameter Tuning**: Different tasks benefit from different sampling parameters:

   | Parameter     | Low Values (0.1-0.3)          | Medium Values (0.5-0.7)     | High Values (0.9-1.5)      |
   |---------------|-------------------------------|-----------------------------|-----------------------------|
   | Temperature   | More deterministic, factual   | Balanced                    | More random, creative       |
   | Top P         | Focus on most likely tokens   | Balanced diversity          | Consider unlikely tokens    |
   | Top K         | Limited vocabulary            | Moderate vocabulary         | Wide vocabulary             |
   | Rep. Penalty  | Allow repetition              | Some repetition control     | Strong repetition control   |

3. **Presets vs. Custom**: Use presets for well-tested parameter combinations, or custom overrides for specific requirements.

4. **Global vs. Request-Specific**: Remember that sampling overrides are global and affect all requests. For request-specific parameters, set them directly in the request:
   ```go
   // These override the global sampling parameters for this specific request
   req := &tabby.CompletionRequest{
       Prompt:      "Hello",
       Temperature: 0.8,
       TopP:        0.9,
   }
   ```

5. **Experiment and Compare**: Test different sampling parameters to find the best balance for your specific use case:
   ```go
   for temp := 0.1; temp <= 1.5; temp += 0.2 {
       // Test with different temperature values
   }
   ```

6. **Reset After Testing**: Always reset to default parameters when done experimenting:
   ```go
   defer func() {
       ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
       defer cancel()
       _ = client.Sampling().UnloadOverride(ctx)
   }()
   ```

7. **Error Recovery**: Implement robust error handling for sampling operations:
   ```go
   if err != nil {
       // Log the error
       log.Printf("Error with sampling operation: %v", err)
       
       // Attempt to reset to default state
       _ = client.Sampling().UnloadOverride(ctx)
   }
   ```

8. **Model-Specific Parameters**: Different models may respond differently to the same sampling parameters. Adjust based on the specific model you're using.