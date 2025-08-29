package health

import (
	"context"

	"tushartemplategin/pkg/interfaces"
)

// HealthService implements the Service interface for health business logic
type HealthService struct {
	repo   Repository
	logger interfaces.Logger
}

// NewHealthService creates a new health service
func NewHealthService(repo Repository, log interfaces.Logger) Service {
	return &HealthService{
		repo:   repo,
		logger: log,
	}
}

// GetHealth returns the overall health status
func (s *HealthService) GetHealth(ctx context.Context) (*HealthStatus, error) {
	s.logger.Info(ctx, "Getting health status", interfaces.Fields{})
	return s.repo.GetHealth(ctx)
}

// GetReadiness returns the readiness status for Kubernetes readiness probes
func (s *HealthService) GetReadiness(ctx context.Context) (*ReadinessStatus, error) {
	s.logger.Info(ctx, "Getting readiness status", interfaces.Fields{})
	return s.repo.GetReadiness(ctx)
}

// GetLiveness returns the liveness status for Kubernetes liveness probes
func (s *HealthService) GetLiveness(ctx context.Context) (*LivenessStatus, error) {
	s.logger.Info(ctx, "Getting liveness status", interfaces.Fields{})
	return s.repo.GetLiveness(ctx)
}
