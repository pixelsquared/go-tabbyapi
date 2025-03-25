# Health Service

The Health service provides functionality for checking the health status of the TabbyAPI server. It allows applications to verify that the server is running properly and to detect any issues that might affect its operation.

## Interface

```go
// HealthService handles health checks to verify the TabbyAPI server status.
type HealthService interface {
	// Check returns the current health status of the TabbyAPI server.
	Check(ctx context.Context) (*HealthCheckResponse, error)
}
```

## Health Types

### HealthCheckResponse

The `HealthCheckResponse` struct provides information about the server health status:

```go
type HealthCheckResponse struct {
	Status string           `json:"status"`           // "ok" or "unhealthy"
	Issues []UnhealthyEvent `json:"issues,omitempty"` // List of issues if status is "unhealthy"
}

type UnhealthyEvent struct {
	Time        string `json:"time"`        // Timestamp of the issue
	Description string `json:"description"` // Description of the issue
}
```

## Examples

### Basic Health Check

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
	
	// Check the health status of the TabbyAPI server
	healthResp, err := client.Health().Check(ctx)
	if err != nil {
		log.Fatalf("Error checking health: %v", err)
	}
	
	fmt.Printf("TabbyAPI server status: %s\n", healthResp.Status)
	
	if healthResp.Status == "ok" {
		fmt.Println("The server is healthy and ready to process requests")
	} else {
		fmt.Println("The server is currently experiencing issues:")
		for _, issue := range healthResp.Issues {
			fmt.Printf("- [%s] %s\n", issue.Time, issue.Description)
		}
	}
}
```

### Health Check with Retry

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
	
	// Wait for the server to be healthy with timeout
	if err := waitForHealthy(client, 5*time.Minute); err != nil {
		log.Fatalf("Server failed to become healthy: %v", err)
	}
	
	fmt.Println("Server is healthy! Proceeding with operations...")
	
	// Continue with other operations...
}

// waitForHealthy polls the health endpoint until the server is healthy or timeout
func waitForHealthy(client tabby.Client, timeout time.Duration) error {
	startTime := time.Now()
	checkInterval := 5 * time.Second
	
	fmt.Println("Waiting for TabbyAPI server to be healthy...")
	
	for {
		// Check if we've exceeded the timeout
		if time.Since(startTime) > timeout {
			return fmt.Errorf("timeout exceeded while waiting for server to become healthy")
		}
		
		// Create a context for this single health check
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		
		// Check health
		healthResp, err := client.Health().Check(ctx)
		
		// Always cancel the context to avoid leaks
		cancel()
		
		if err != nil {
			fmt.Printf("Health check failed: %v. Retrying in %v...\n", err, checkInterval)
		} else if healthResp.Status == "ok" {
			return nil // Server is healthy
		} else {
			fmt.Printf("Server is unhealthy with %d issue(s). Retrying in %v...\n", 
				len(healthResp.Issues), checkInterval)
		}
		
		// Wait before next attempt
		time.Sleep(checkInterval)
	}
}
```

### Monitoring Health Status

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	
	"github.com/pixelsquared/go-tabbyapi/tabby"
)

func main() {
	client := tabby.NewClient(
		tabby.WithBaseURL("http://localhost:8080"),
		tabby.WithAPIKey("your-api-key"),
	)
	defer client.Close()
	
	// Channel to receive termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	// Channel to stop the monitoring goroutine
	stopChan := make(chan struct{})
	var wg sync.WaitGroup
	
	// Start health monitoring
	wg.Add(1)
	go monitorHealth(client, stopChan, &wg)
	
	// Wait for termination signal
	<-sigChan
	fmt.Println("\nReceived termination signal. Shutting down...")
	
	// Stop the monitoring goroutine
	close(stopChan)
	
	// Wait for the goroutine to finish
	wg.Wait()
	fmt.Println("Health monitoring stopped")
}

// monitorHealth periodically checks the health status of the TabbyAPI server
func monitorHealth(client tabby.Client, stopChan <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	
	// Check health every 30 seconds
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	fmt.Println("Starting health monitoring...")
	
	// Perform initial health check
	performHealthCheck(client)
	
	// Continue monitoring until stopped
	for {
		select {
		case <-ticker.C:
			performHealthCheck(client)
		case <-stopChan:
			fmt.Println("Stopping health monitoring...")
			return
		}
	}
}

// performHealthCheck executes a single health check and logs the result
func performHealthCheck(client tabby.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	healthResp, err := client.Health().Check(ctx)
	
	timestamp := time.Now().Format(time.RFC3339)
	
	if err != nil {
		fmt.Printf("[%s] Health check error: %v\n", timestamp, err)
		return
	}
	
	if healthResp.Status == "ok" {
		fmt.Printf("[%s] Health status: OK\n", timestamp)
	} else {
		fmt.Printf("[%s] Health status: UNHEALTHY with %d issue(s)\n", 
			timestamp, len(healthResp.Issues))
		for _, issue := range healthResp.Issues {
			fmt.Printf("  - [%s] %s\n", issue.Time, issue.Description)
		}
	}
}
```

## Error Handling

The Health service operations may fail due to various reasons:

```go
healthResp, err := client.Health().Check(ctx)
if err != nil {
	var apiErr *tabby.APIError
	if errors.As(err, &apiErr) {
		// Handle API-specific errors
		fmt.Printf("API Error: %s (Code: %s)\n", apiErr.Error(), apiErr.Code())
	} else {
		// Handle other types of errors (connectivity, timeout, etc.)
		var reqErr *tabby.RequestError
		if errors.As(err, &reqErr) {
			if errors.Is(reqErr.Unwrap(), context.DeadlineExceeded) {
				fmt.Println("Health check timed out - server might be overloaded")
			} else {
				fmt.Printf("Request error: %v\n", reqErr)
			}
		} else {
			fmt.Printf("Unknown error: %v\n", err)
		}
	}
	
	// Take appropriate action based on the error
	// For example, retry, fail gracefully, or alert administrators
}
```

## Best Practices

1. **Regular Health Checks**: Implement periodic health checks in production applications to detect issues early:
   ```go
   func scheduleHealthChecks(client tabby.Client, interval time.Duration) {
       ticker := time.NewTicker(interval)
       defer ticker.Stop()
       
       for range ticker.C {
           ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
           _, err := client.Health().Check(ctx)
           cancel()
           
           if err != nil {
               log.Printf("Health check failed: %v", err)
               // Implement alerting or recovery logic
           }
       }
   }
   ```

2. **Readiness Checks**: Use health checks before starting critical operations to ensure the server is ready:
   ```go
   func ensureServerReady(client tabby.Client) error {
       ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
       defer cancel()
       
       resp, err := client.Health().Check(ctx)
       if err != nil {
           return fmt.Errorf("server health check failed: %w", err)
       }
       
       if resp.Status != "ok" {
           return fmt.Errorf("server is unhealthy: %d issues reported", len(resp.Issues))
       }
       
       return nil
   }
   ```

3. **Graceful Degradation**: Handle unhealthy server status gracefully in your application:
   ```go
   func performOperationWithHealthCheck(client tabby.Client) error {
       // Check health first
       healthResp, err := client.Health().Check(ctx)
       if err != nil || healthResp.Status != "ok" {
           // Log the issue
           log.Printf("Server health check failed: %v", err)
           
           // Fall back to cached data or alternative service
           return useBackupService()
       }
       
       // Proceed with normal operation
       return performMainOperation(client)
   }
   ```

4. **Timeout Management**: Use appropriate timeouts for health checks:
   ```go
   // Short timeout for routine health checks
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   defer cancel()
   ```

5. **Error Handling**: Implement robust error handling for health checks:
   ```go
   healthResp, err := client.Health().Check(ctx)
   if err != nil {
       // For transient errors, retry with backoff
       if isTransientError(err) {
           return retryWithBackoff(func() error {
               r, e := client.Health().Check(ctx)
               if e == nil && r.Status == "ok" {
                   return nil
               }
               return e
           })
       }
       
       // For persistent errors, take appropriate action
       return handlePersistentError(err)
   }
   ```

6. **Integration with Monitoring**: Integrate health checks with your monitoring system:
   ```go
   func reportHealthToMonitoring(client tabby.Client) {
       healthResp, err := client.Health().Check(ctx)
       
       if err != nil {
           monitoring.ReportMetric("tabby_health_check_error", 1)
           monitoring.ReportEvent("TabbyAPI health check failed", err.Error())
       } else {
           monitoring.ReportMetric("tabby_health_check_success", 1)
           
           if healthResp.Status != "ok" {
               monitoring.ReportMetric("tabby_health_issues", len(healthResp.Issues))
               for _, issue := range healthResp.Issues {
                   monitoring.ReportEvent("TabbyAPI health issue", issue.Description)
               }
           } else {
               monitoring.ReportMetric("tabby_health_issues", 0)
           }
       }
   }