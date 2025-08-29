package health

import (
	"context"
)

// Service defines the interface for health business logic
type Service interface {
	GetHealth(ctx context.Context) (*HealthStatus, error)       // Get overall health status
	GetReadiness(ctx context.Context) (*ReadinessStatus, error) // Get readiness status
	GetLiveness(ctx context.Context) (*LivenessStatus, error)   // Get liveness status
}

// Repository defines the interface for health data access
type Repository interface {
	// Health check for the repository
	Health(ctx context.Context) error
	GetHealth(ctx context.Context) (*HealthStatus, error)
	GetReadiness(ctx context.Context) (*ReadinessStatus, error)
	GetLiveness(ctx context.Context) (*LivenessStatus, error)
	UpdateHealth(ctx context.Context, status *HealthStatus) error
	GetHealthHistory(ctx context.Context, limit int) ([]*HealthStatus, error)
}
