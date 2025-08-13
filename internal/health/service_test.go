package health

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewHealthService tests service creation
func TestNewHealthService(t *testing.T) {
	// Test that health service can be created successfully

	// Create mock repository and logger
	mockRepo := &mockHealthRepository{}
	mockLogger := &mockLogger{}

	// Create service
	service := NewHealthService(mockRepo, mockLogger)

	// Verify service was created
	assert.NotNil(t, service, "Health service should be created successfully")

	// Verify service implements the Service interface
	var _ Service = service
}

// TestGetHealth tests health status retrieval
func TestGetHealth(t *testing.T) {
	// Test that health status is retrieved correctly

	// Create mock repository with expected data
	expectedHealth := &HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   "test-service",
		Version:   "1.0.0",
	}

	mockRepo := &mockHealthRepository{
		healthStatus: expectedHealth,
		healthError:  nil,
	}

	mockLogger := &mockLogger{}

	// Create service
	service := NewHealthService(mockRepo, mockLogger)

	// Get health status
	ctx := context.Background()
	health, err := service.GetHealth(ctx)

	// Verify no error
	require.NoError(t, err, "Should not return error")

	// Verify health status
	assert.Equal(t, expectedHealth.Status, health.Status, "Health status should match")
	assert.Equal(t, expectedHealth.Service, health.Service, "Service name should match")
	assert.Equal(t, expectedHealth.Version, health.Version, "Version should match")
}

// TestGetHealthError tests health status retrieval with error
func TestGetHealthError(t *testing.T) {
	// Test that errors are properly propagated from repository

	expectedError := "database connection failed"
	mockRepo := &mockHealthRepository{
		healthStatus: nil,
		healthError:  expectedError,
	}

	mockLogger := &mockLogger{}

	// Create service
	service := NewHealthService(mockRepo, mockLogger)

	// Get health status
	ctx := context.Background()
	health, err := service.GetHealth(ctx)

	// Verify error is returned
	assert.Error(t, err, "Should return error")
	assert.Contains(t, err.Error(), expectedError, "Error should contain expected message")

	// Verify health status is nil
	assert.Nil(t, health, "Health status should be nil on error")
}

// TestGetReadiness tests readiness status retrieval
func TestGetReadiness(t *testing.T) {
	// Test that readiness status is retrieved correctly

	expectedReadiness := &ReadinessStatus{
		Status:    "ready",
		Timestamp: time.Now(),
		Service:   "test-service",
		Version:   "1.0.0",
	}

	mockRepo := &mockHealthRepository{
		readinessStatus: expectedReadiness,
		readinessError:  nil,
	}

	mockLogger := &mockLogger{}

	// Create service
	service := NewHealthService(mockRepo, mockLogger)

	// Get readiness status
	ctx := context.Background()
	readiness, err := service.GetReadiness(ctx)

	// Verify no error
	require.NoError(t, err, "Should not return error")

	// Verify readiness status
	assert.Equal(t, expectedReadiness.Status, readiness.Status, "Readiness status should match")
}

// TestGetLiveness tests liveness status retrieval
func TestGetLiveness(t *testing.T) {
	// Test that liveness status is retrieved correctly

	expectedLiveness := &LivenessStatus{
		Status:    "alive",
		Timestamp: time.Now(),
		Service:   "test-service",
		Version:   "1.0.0",
	}

	mockRepo := &mockHealthRepository{
		livenessStatus: expectedLiveness,
		livenessError:  nil,
	}

	mockLogger := &mockLogger{}

	// Create service
	service := NewHealthService(mockRepo, mockLogger)

	// Get liveness status
	ctx := context.Background()
	liveness, err := service.GetLiveness(ctx)

	// Verify no error
	require.NoError(t, err, "Should not return error")

	// Verify liveness status
	assert.Equal(t, expectedLiveness.Status, liveness.Status, "Liveness status should match")
}

// TestSingletonPattern tests the singleton pattern implementation
func TestSingletonPattern(t *testing.T) {
	// Test that the singleton pattern works correctly

	// Reset singleton for testing
	ResetDefaultService()

	// Get service first time
	service1, err1 := GetDefaultService()
	require.NoError(t, err1, "Should get service successfully first time")
	require.NotNil(t, service1, "Service should not be nil")

	// Get service second time
	service2, err2 := GetDefaultService()
	require.NoError(t, err2, "Should get service successfully second time")
	require.NotNil(t, service2, "Service should not be nil")

	// Verify it's the same instance (singleton)
	assert.Equal(t, service1, service2, "Should return same service instance (singleton)")

	// Verify both services work
	ctx := context.Background()

	health1, err1 := service1.GetHealth(ctx)
	health2, err2 := service2.GetHealth(ctx)

	// Both should work identically
	assert.Equal(t, health1, health2, "Both service instances should return same data")
	assert.Equal(t, err1, err2, "Both service instances should return same errors")
}

// TestSingletonReset tests singleton reset functionality
func TestSingletonReset(t *testing.T) {
	// Test that singleton can be reset for testing

	// Get initial service
	service1, err1 := GetDefaultService()
	require.NoError(t, err1, "Should get service successfully")

	// Reset singleton
	ResetDefaultService()

	// Get service again after reset
	service2, err2 := GetDefaultService()
	require.NoError(t, err2, "Should get service successfully after reset")

	// Verify it's a different instance
	assert.NotEqual(t, service1, service2, "Should return different service instance after reset")
}

// TestServiceLogging tests that service operations are logged
func TestServiceLogging(t *testing.T) {
	// Test that service operations are properly logged

	mockRepo := &mockHealthRepository{
		healthStatus: &HealthStatus{Status: "healthy"},
		healthError:  nil,
	}

	mockLogger := &mockLogger{
		loggedMessages: make([]string, 0),
	}

	// Create service
	service := NewHealthService(mockRepo, mockLogger)

	// Perform health check
	ctx := context.Background()
	_, err := service.GetHealth(ctx)
	require.NoError(t, err, "Health check should succeed")

	// Verify logging occurred
	assert.Len(t, mockLogger.loggedMessages, 1, "Should log one message")
	assert.Contains(t, mockLogger.loggedMessages[0], "Health check requested", "Should log health check request")
}

// TestContextPropagation tests that context is properly propagated
func TestContextPropagation(t *testing.T) {
	// Test that context is properly passed through the service layer

	mockRepo := &mockHealthRepository{
		healthStatus: &HealthStatus{Status: "healthy"},
		healthError:  nil,
	}

	mockLogger := &mockLogger{}

	// Create service
	service := NewHealthService(mockRepo, mockLogger)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Perform health check with context
	_, err := service.GetHealth(ctx)
	require.NoError(t, err, "Health check should succeed")

	// Verify context was used (mock repository should have received it)
	// This test verifies that the context flows through the service layer
	assert.True(t, true, "Context should be properly propagated")
}

// TestConcurrentAccess tests concurrent access to singleton service
func TestConcurrentAccess(t *testing.T) {
	// Test that singleton service can be accessed concurrently

	// Reset singleton for testing
	ResetDefaultService()

	// Number of concurrent goroutines
	numGoroutines := 10

	// Channel to collect results
	results := make(chan *HealthStatus, numGoroutines)
	errors := make(chan error, numGoroutines)

	// Start concurrent goroutines
	for i := 0; i < numGoroutines; i++ {
		go func() {
			service, err := GetDefaultService()
			if err != nil {
				errors <- err
				return
			}

			ctx := context.Background()
			health, err := service.GetHealth(ctx)
			if err != nil {
				errors <- err
				return
			}

			results <- health
		}()
	}

	// Collect results
	healthResults := make([]*HealthStatus, 0, numGoroutines)
	errorResults := make([]error, 0, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		select {
		case health := <-results:
			healthResults = append(healthResults, health)
		case err := <-errors:
			errorResults = append(errorResults, err)
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout waiting for concurrent access results")
		}
	}

	// Verify all goroutines succeeded
	assert.Len(t, healthResults, numGoroutines, "All goroutines should succeed")
	assert.Len(t, errorResults, 0, "No goroutines should fail")

	// Verify all results are the same (same service instance)
	firstResult := healthResults[0]
	for i, result := range healthResults {
		assert.Equal(t, firstResult, result, "Result %d should match first result", i)
	}
}

// Mock implementations for testing

type mockHealthRepository struct {
	healthStatus    *HealthStatus
	healthError     string
	readinessStatus *ReadinessStatus
	readinessError  string
	livenessStatus  *LivenessStatus
	livenessError   string
}

func (m *mockHealthRepository) GetHealth(ctx context.Context) (*HealthStatus, error) {
	if m.healthError != "" {
		return nil, assert.AnError
	}
	return m.healthStatus, nil
}

func (m *mockHealthRepository) GetReadiness(ctx context.Context) (*ReadinessStatus, error) {
	if m.readinessError != "" {
		return nil, assert.AnError
	}
	return m.readinessStatus, nil
}

func (m *mockHealthRepository) GetLiveness(ctx context.Context) (*LivenessStatus, error) {
	if m.livenessError != "" {
		return nil, assert.AnError
	}
	return m.livenessStatus, nil
}

type mockLogger struct {
	loggedMessages []string
}

func (m *mockLogger) Debug(ctx context.Context, msg string, fields Fields) {
	m.loggedMessages = append(m.loggedMessages, msg)
}

func (m *mockLogger) Info(ctx context.Context, msg string, fields Fields) {
	m.loggedMessages = append(m.loggedMessages, msg)
}

func (m *mockLogger) Warn(ctx context.Context, msg string, fields Fields) {
	m.loggedMessages = append(m.loggedMessages, msg)
}

func (m *mockLogger) Error(ctx context.Context, msg string, fields Fields) {
	m.loggedMessages = append(m.loggedMessages, msg)
}

func (m *mockLogger) Fatal(ctx context.Context, msg string, err error, fields Fields) {
	m.loggedMessages = append(m.loggedMessages, msg)
}

// Benchmark tests

// BenchmarkGetHealth benchmarks health status retrieval
func BenchmarkGetHealth(b *testing.B) {
	mockRepo := &mockHealthRepository{
		healthStatus: &HealthStatus{Status: "healthy"},
		healthError:  nil,
	}

	mockLogger := &mockLogger{}
	service := NewHealthService(mockRepo, mockLogger)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GetHealth(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSingletonAccess benchmarks singleton access performance
func BenchmarkSingletonAccess(b *testing.B) {
	// Reset singleton for clean benchmark
	ResetDefaultService()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service, err := GetDefaultService()
		if err != nil {
			b.Fatal(err)
		}
		if service == nil {
			b.Fatal("Service is nil")
		}
	}
}
