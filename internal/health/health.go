// Package health provides functionality for checking TabbyAPI's health status.
package health

import (
	"context"
	"fmt"

	"github.com/pixelsquared/go-tabbyapi/internal/rest"
	"github.com/pixelsquared/go-tabbyapi/tabby"
)

// Service implements the HealthService interface for TabbyAPI.
type Service struct {
	client *rest.Client
}

// New creates a new HealthService instance using the provided REST client.
func New(client *rest.Client) *Service {
	return &Service{client: client}
}

// Check returns the current health status of the TabbyAPI service.
func (s *Service) Check(ctx context.Context) (*tabby.HealthCheckResponse, error) {
	var response tabby.HealthCheckResponse
	err := s.client.Get(ctx, "v1/health", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to check health: %w", err)
	}
	return &response, nil
}
