package health

import (
	"context"
	"time"

	"tushartemplategin/pkg/interfaces"
)

// HealthRepository implements the Repository interface for health data
type HealthRepository struct {
	logger interfaces.Logger
}

// NewHealthRepository creates a new health repository
func NewHealthRepository(log interfaces.Logger) Repository {
	return &HealthRepository{
		logger: log,
	}
}

// GetHealth returns the overall health status of the service
func (r *HealthRepository) GetHealth(ctx context.Context) (*HealthStatus, error) {
	return &HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   "tushar-service",
		Version:   "1.0.0",
	}, nil
}

// Health implements the Repository interface
func (r *HealthRepository) Health(ctx context.Context) error {
	// Simple health check - always return healthy
	return nil
}

// GetReadiness returns the readiness status for Kubernetes readiness probes
func (r *HealthRepository) GetReadiness(ctx context.Context) (*ReadinessStatus, error) {
	return &ReadinessStatus{
		Status:    "ready",
		Timestamp: time.Now(),
		Database:  "not_required", // No database dependency
		Service:   "tushar-service",
	}, nil
}

// GetLiveness returns the liveness status for Kubernetes liveness probes
func (r *HealthRepository) GetLiveness(ctx context.Context) (*LivenessStatus, error) {
	return &LivenessStatus{
		Status:    "alive",
		Timestamp: time.Now(),
		Service:   "tushar-service",
	}, nil
}

// UpdateHealth updates the health status (in-memory only)
func (r *HealthRepository) UpdateHealth(ctx context.Context, status *HealthStatus) error {
	// Log the update but don't persist to database
	r.logger.Info(ctx, "Health status update requested", interfaces.Fields{
		"status":  status.Status,
		"service": status.Service,
		"version": status.Version,
	})
	return nil
}

// GetHealthHistory retrieves health status history (empty for in-memory implementation)
func (r *HealthRepository) GetHealthHistory(ctx context.Context, limit int) ([]*HealthStatus, error) {
	// Return empty history since we're not persisting to database
	r.logger.Info(ctx, "Health history requested", interfaces.Fields{"limit": limit})
	return []*HealthStatus{}, nil
}
