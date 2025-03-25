# Auth Service

The Auth service provides functionality for checking authentication permissions and access levels. It allows applications to verify what operations the current authentication credentials can perform.

## Interface

```go
// AuthService handles authentication permissions and access levels.
type AuthService interface {
	// GetPermission returns the access level for the current authentication method.
	GetPermission(ctx context.Context) (*AuthPermissionResponse, error)
}
```

## Auth Types

### AuthPermissionResponse

The `AuthPermissionResponse` struct contains the permission level:

```go
type AuthPermissionResponse struct {
	Permission string `json:"permission"` // Permission level ("none", "read", "write", "admin")
}
```

Common permission levels include:
- `"none"`: No access
- `"read"`: Read-only access to generate completions, embeddings, etc.
- `"write"`: Read and write access, including managing conversations
- `"admin"`: Full access, including administrative operations like model and LoRA management

## Examples

### Checking Current Permissions

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
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Check permissions for the current authentication credentials
	authResp, err := client.Auth().GetPermission(ctx)
	if err != nil {
		log.Fatalf("Error checking permissions: %v", err)
	}
	
	fmt.Printf("Current permission level: %s\n", authResp.Permission)
	
	// Take action based on permission level
	switch authResp.Permission {
	case "none":
		fmt.Println("Warning: No access permissions. Most operations will fail.")
	case "read":
		fmt.Println("Read-only access. You can generate completions and embeddings, but can't perform admin operations.")
	case "write":
		fmt.Println("Write access. You can generate completions, embeddings, and manage conversations.")
	case "admin":
		fmt.Println("Admin access. You have full access to all TabbyAPI operations.")
	default:
		fmt.Printf("Unknown permission level: %s\n", authResp.Permission)
	}
}
```

### Permission-Based Feature Enablement

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
	
	// Check permissions and enable features accordingly
	authResp, err := client.Auth().GetPermission(ctx)
	if err != nil {
		log.Fatalf("Error checking permissions: %v", err)
	}
	
	fmt.Printf("Permission level: %s\n", authResp.Permission)
	
	// Use the client to perform operations based on permission level
	if authResp.Permission == "none" {
		fmt.Println("No valid permissions. Please check your authentication credentials.")
		return
	}
	
	// All permission levels can perform basic completions
	fmt.Println("Basic completion feature: Enabled")
	
	// Example: Try to list models (admin only)
	if authResp.Permission == "admin" {
		fmt.Println("Admin features: Enabled")
		fmt.Println("- Model management")
		fmt.Println("- LoRA adapter management")
		fmt.Println("- Template management")
		fmt.Println("- Sampling parameter management")
		
		// Demonstrate an admin operation
		models, err := client.Models().List(ctx)
		if err != nil {
			fmt.Printf("Error listing models: %v\n", err)
		} else {
			fmt.Printf("Available models: %d\n", len(models.Data))
		}
	} else {
		fmt.Println("Admin features: Disabled (requires admin permission)")
	}
}
```

### Permission Check with Different Auth Methods

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
	// Define API endpoint
	endpoint := "http://localhost:8080"
	
	// Test different authentication methods
	testAuth := func(name string, client tabby.Client) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		
		fmt.Printf("\nTesting authentication method: %s\n", name)
		
		authResp, err := client.Auth().GetPermission(ctx)
		if err != nil {
			fmt.Printf("Error checking permissions: %v\n", err)
			return
		}
		
		fmt.Printf("Permission level: %s\n", authResp.Permission)
	}
	
	// Test with API key
	apiKeyClient := tabby.NewClient(
		tabby.WithBaseURL(endpoint),
		tabby.WithAPIKey("your-api-key"),
	)
	defer apiKeyClient.Close()
	testAuth("API Key", apiKeyClient)
	
	// Test with admin key
	adminKeyClient := tabby.NewClient(
		tabby.WithBaseURL(endpoint),
		tabby.WithAdminKey("your-admin-key"),
	)
	defer adminKeyClient.Close()
	testAuth("Admin Key", adminKeyClient)
	
	// Test with bearer token
	bearerClient := tabby.NewClient(
		tabby.WithBaseURL(endpoint),
		tabby.WithBearerToken("your-bearer-token"),
	)
	defer bearerClient.Close()
	testAuth("Bearer Token", bearerClient)
	
	// Test with no authentication
	noAuthClient := tabby.NewClient(
		tabby.WithBaseURL(endpoint),
	)
	defer noAuthClient.Close()
	testAuth("No Authentication", noAuthClient)
}
```

### Graceful Error Handling for Permission Issues

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	
	"github.com/pixelsquared/go-tabbyapi/tabby"
)

func main() {
	client := tabby.NewClient(
		tabby.WithBaseURL("http://localhost:8080"),
		tabby.WithAPIKey("your-api-key"), // Using regular API key (may not have admin permissions)
	)
	defer client.Close()
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Check permissions first
	authResp, err := client.Auth().GetPermission(ctx)
	if err != nil {
		log.Fatalf("Error checking permissions: %v", err)
	}
	
	fmt.Printf("Current permission level: %s\n", authResp.Permission)
	
	// Attempt to perform an admin operation (loading a model)
	if authResp.Permission == "admin" {
		fmt.Println("You have admin permissions. Proceeding with model operations...")
		
		// Admin operation example: listing models
		models, err := client.Models().List(ctx)
		if err != nil {
			log.Fatalf("Error listing models: %v", err)
		}
		
		fmt.Printf("Found %d models\n", len(models.Data))
	} else {
		fmt.Println("You don't have admin permissions. Attempting model operations anyway...")
		
		// Try the operation and handle the permission error gracefully
		_, err := client.Models().List(ctx)
		
		var apiErr *tabby.APIError
		if errors.As(err, &apiErr) && apiErr.Code() == "permission_error" {
			fmt.Println("Permission denied as expected. Please use an admin key for this operation.")
			fmt.Println("Continuing with regular operations...")
			
			// Fall back to operations that work with current permission level
			fmt.Println("Attempting a regular completion instead:")
			
			// This should work with any permission level that's not "none"
			completionReq := &tabby.CompletionRequest{
				Prompt:    "Hello, world!",
				MaxTokens: 20,
			}
			
			completionResp, err := client.Completions().Create(ctx, completionReq)
			if err != nil {
				fmt.Printf("Error creating completion: %v\n", err)
			} else if len(completionResp.Choices) > 0 {
				fmt.Printf("Completion: %s\n", completionResp.Choices[0].Text)
			}
		} else if err != nil {
			fmt.Printf("Unexpected error: %v\n", err)
		}
	}
}
```

## Error Handling

The Auth service operations may fail due to various reasons:

```go
authResp, err := client.Auth().GetPermission(ctx)
if err != nil {
	var apiErr *tabby.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.Code() {
		case "authentication_error":
			fmt.Println("Authentication failed. Check your API key or credentials.")
		default:
			fmt.Printf("API Error: %s (Code: %s)\n", apiErr.Error(), apiErr.Code())
		}
	} else {
		// Handle other types of errors (connectivity, timeout, etc.)
		var reqErr *tabby.RequestError
		if errors.As(err, &reqErr) {
			fmt.Printf("Request error: %v\n", reqErr)
		} else {
			fmt.Printf("Unknown error: %v\n", err)
		}
	}
}
```

## Best Practices

1. **Check Permissions Early**: Verify permissions at the start of your application to ensure appropriate feature enablement:
   ```go
   func checkPermissions(client tabby.Client) (string, error) {
       ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
       defer cancel()
       
       resp, err := client.Auth().GetPermission(ctx)
       if err != nil {
           return "", err
       }
       
       return resp.Permission, nil
   }
   ```

2. **Feature Gates**: Use permission levels to enable or disable features:
   ```go
   func enableAdminFeatures(ui *UserInterface, permission string) {
       if permission == "admin" {
           ui.EnableModelManagement()
           ui.EnableLoraManagement()
           ui.EnableTemplateManagement()
       } else {
           ui.DisableAdminFeatures()
       }
   }
   ```

3. **Graceful Degradation**: Handle permission issues gracefully:
   ```go
   func performOperation(client tabby.Client) error {
       permission, err := checkPermissions(client)
       if err != nil {
           return err
       }
       
       if permission != "admin" && requiresAdmin(operation) {
           // Fallback to non-admin operation
           return performNonAdminFallback()
       }
       
       // Proceed with original operation
       return performOriginalOperation()
   }
   ```

4. **Use Appropriate Authentication Methods**:
   - For most applications:
     ```go
     client := tabby.NewClient(
         tabby.WithAPIKey("your-api-key"),
     )
     ```
   - For administrative applications:
     ```go
     client := tabby.NewClient(
         tabby.WithAdminKey("your-admin-key"),
     )
     ```

5. **Security Best Practices**:
   - Store keys securely (environment variables, secret management services)
   - Don't hardcode keys in source code
   - Use the least privilege principle (don't use admin keys for regular operations)
   - Implement access control within your application

6. **Permission Caching**: Cache permission results to avoid frequent queries:
   ```go
   var (
       cachedPermission string
       permissionExpiry time.Time
       permissionMutex  sync.Mutex
   )
   
   func getPermissionWithCache(client tabby.Client) (string, error) {
       permissionMutex.Lock()
       defer permissionMutex.Unlock()
       
       // Return cached permission if not expired
       if time.Now().Before(permissionExpiry) {
           return cachedPermission, nil
       }
       
       // Otherwise fetch fresh permission
       ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
       defer cancel()
       
       resp, err := client.Auth().GetPermission(ctx)
       if err != nil {
           return "", err
       }
       
       // Update cache
       cachedPermission = resp.Permission
       permissionExpiry = time.Now().Add(5 * time.Minute)
       
       return cachedPermission, nil
   }