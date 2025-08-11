package health

import (
	"context"
	"tushartemplategin/pkg/logger"
)

// service implements the Service interface
type service struct {
	repo   Repository    // Health repository for data access
	logger logger.Logger // Logger for recording operations
}

// NewHealthService creates a new health service instance
func NewHealthService(repo Repository, logger logger.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

// GetHealth retrieves the overall health status and logs the operation
func (s *service) GetHealth(ctx context.Context) (*HealthStatus, error) {
	s.logger.Info(ctx, "Health check requested", logger.Fields{})
	return s.repo.GetHealth(ctx)
}

// GetReadiness retrieves the readiness status and logs the operation
func (s *service) GetReadiness(ctx context.Context) (*ReadinessStatus, error) {
	s.logger.Info(ctx, "Readiness check requested", logger.Fields{})
	return s.repo.GetReadiness(ctx)
}

// GetLiveness retrieves the liveness status and logs the operation
func (s *service) GetLiveness(ctx context.Context) (*LivenessStatus, error) {
	s.logger.Info(ctx, "Liveness check requested", logger.Fields{})
	return s.repo.GetLiveness(ctx)
}
