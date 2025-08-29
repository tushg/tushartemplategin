package health

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"tushartemplategin/mocks"
)

func TestNewHealthRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger(ctrl)

	repo := NewHealthRepository(mockLogger)
	assert.NotNil(t, repo)
}

func TestHealthRepository_GetHealth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger(ctrl)
	repo := NewHealthRepository(mockLogger)

	ctx := context.Background()
	health, err := repo.GetHealth(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, health)
	assert.Equal(t, "healthy", health.Status)
	assert.Equal(t, "tushar-service", health.Service)
	assert.Equal(t, "1.0.0", health.Version)
	assert.WithinDuration(t, time.Now(), health.Timestamp, 2*time.Second)
}

func TestHealthRepository_Health(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger(ctrl)
	repo := NewHealthRepository(mockLogger)

	ctx := context.Background()
	err := repo.Health(ctx)

	assert.NoError(t, err)
}

func TestHealthRepository_GetReadiness(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger(ctrl)
	repo := NewHealthRepository(mockLogger)

	ctx := context.Background()
	readiness, err := repo.GetReadiness(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, readiness)
	assert.Equal(t, "ready", readiness.Status)
	assert.Equal(t, "not_required", readiness.Database)
	assert.Equal(t, "tushar-service", readiness.Service)
	assert.WithinDuration(t, time.Now(), readiness.Timestamp, 2*time.Second)
}

func TestHealthRepository_GetLiveness(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger(ctrl)
	repo := NewHealthRepository(mockLogger)

	ctx := context.Background()
	liveness, err := repo.GetLiveness(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, liveness)
	assert.Equal(t, "alive", liveness.Status)
	assert.Equal(t, "tushar-service", liveness.Service)
	assert.WithinDuration(t, time.Now(), liveness.Timestamp, 2*time.Second)
}

func TestHealthRepository_UpdateHealth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger(ctrl)
	repo := NewHealthRepository(mockLogger)

	ctx := context.Background()
	status := &HealthStatus{
		Status:    "degraded",
		Timestamp: time.Now(),
		Service:   "test-service",
		Version:   "2.0.0",
	}

	// Expect logger to be called with the update
	mockLogger.EXPECT().Info(ctx, "Health status update requested", gomock.Any())

	err := repo.UpdateHealth(ctx, status)
	assert.NoError(t, err)
}

func TestHealthRepository_GetHealthHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocks.NewMockLogger(ctrl)
	repo := NewHealthRepository(mockLogger)

	ctx := context.Background()

	// Expect logger to be called
	mockLogger.EXPECT().Info(ctx, "Health history requested", gomock.Any())

	history, err := repo.GetHealthHistory(ctx, 10)

	assert.NoError(t, err)
	assert.NotNil(t, history)
	assert.Empty(t, history) // Should return empty slice
}
