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

// Note: We removed the Handler interface since we're now using closures in routes.go
// This makes the module more flexible and easier to test
// The handlers are now defined as closures in routes.go for better dependency injection
