package health

import "context"

// Repository defines the interface for health data access
type Repository interface {
	GetHealth(ctx context.Context) (*HealthStatus, error)       // Get overall health status
	GetReadiness(ctx context.Context) (*ReadinessStatus, error) // Get readiness status
	GetLiveness(ctx context.Context) (*LivenessStatus, error)   // Get liveness status
}

// Service defines the interface for health business logic
type Service interface {
	GetHealth(ctx context.Context) (*HealthStatus, error)       // Get overall health status
	GetReadiness(ctx context.Context) (*ReadinessStatus, error) // Get readiness status
	GetLiveness(ctx context.Context) (*LivenessStatus, error)   // Get liveness status
}
