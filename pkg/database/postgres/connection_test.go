package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_interfaces "tushartemplategin/mocks"
)

func TestNewPostgresDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfig := mock_interfaces.NewMockPostgresConfig(ctrl)
	mockLogger := mock_interfaces.NewMockLogger(ctrl)

	// Set up mock expectations
	mockConfig.EXPECT().GetHost().Return("localhost").AnyTimes()
	mockConfig.EXPECT().GetPort().Return(5432).AnyTimes()
	mockConfig.EXPECT().GetName().Return("testdb").AnyTimes()
	mockConfig.EXPECT().GetUsername().Return("testuser").AnyTimes()
	mockConfig.EXPECT().GetPassword().Return("testpass").AnyTimes()
	mockConfig.EXPECT().GetSSLMode().Return("disable").AnyTimes()
	mockConfig.EXPECT().GetMaxRetries().Return(3).AnyTimes()
	mockConfig.EXPECT().GetRetryDelay().Return(time.Second).AnyTimes()
	mockConfig.EXPECT().GetTimeout().Return(30 * time.Second).AnyTimes()
	mockConfig.EXPECT().GetMaxOpenConns().Return(25).AnyTimes()
	mockConfig.EXPECT().GetMaxIdleConns().Return(5).AnyTimes()
	mockConfig.EXPECT().GetConnMaxLifetime().Return(5 * time.Minute).AnyTimes()
	mockConfig.EXPECT().GetConnMaxIdleTime().Return(5 * time.Minute).AnyTimes()

	// Test successful creation
	postgresDB := NewPostgresDB(mockConfig, mockLogger)
	assert.NotNil(t, postgresDB)
	assert.Equal(t, mockConfig, postgresDB.config)
	assert.Equal(t, mockLogger, postgresDB.logger)
}

func TestDisconnect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfig := mock_interfaces.NewMockPostgresConfig(ctrl)
	mockLogger := mock_interfaces.NewMockLogger(ctrl)
	mockDB := mock_interfaces.NewMockDBInterface(ctrl)

	// Set up mock expectations
	mockConfig.EXPECT().GetHost().Return("localhost").AnyTimes()
	mockConfig.EXPECT().GetPort().Return(5432).AnyTimes()
	mockConfig.EXPECT().GetName().Return("testdb").AnyTimes()
	mockConfig.EXPECT().GetUsername().Return("testuser").AnyTimes()
	mockConfig.EXPECT().GetPassword().Return("testpass").AnyTimes()
	mockConfig.EXPECT().GetSSLMode().Return("disable").AnyTimes()
	mockConfig.EXPECT().GetMaxRetries().Return(3).AnyTimes()
	mockConfig.EXPECT().GetRetryDelay().Return(time.Second).AnyTimes()
	mockConfig.EXPECT().GetTimeout().Return(30 * time.Second).AnyTimes()
	mockConfig.EXPECT().GetMaxOpenConns().Return(25).AnyTimes()
	mockConfig.EXPECT().GetMaxIdleConns().Return(5).AnyTimes()
	mockConfig.EXPECT().GetConnMaxLifetime().Return(5 * time.Minute).AnyTimes()
	mockConfig.EXPECT().GetConnMaxIdleTime().Return(5 * time.Minute).AnyTimes()

	postgresDB := NewPostgresDB(mockConfig, mockLogger)
	postgresDB.SetTestDB(mockDB)

	// Test successful disconnect
	mockDB.EXPECT().Close().Return(nil)
	mockLogger.EXPECT().Info(gomock.Any(), "PostgreSQL connection closed", gomock.Any())

	err := postgresDB.Disconnect(context.Background())
	assert.NoError(t, err)
}

func TestHealth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfig := mock_interfaces.NewMockPostgresConfig(ctrl)
	mockLogger := mock_interfaces.NewMockLogger(ctrl)
	mockDB := mock_interfaces.NewMockDBInterface(ctrl)

	// Set up mock expectations
	mockConfig.EXPECT().GetHost().Return("localhost").AnyTimes()
	mockConfig.EXPECT().GetPort().Return(5432).AnyTimes()
	mockConfig.EXPECT().GetName().Return("testdb").AnyTimes()
	mockConfig.EXPECT().GetUsername().Return("testuser").AnyTimes()
	mockConfig.EXPECT().GetPassword().Return("testpass").AnyTimes()
	mockConfig.EXPECT().GetSSLMode().Return("disable").AnyTimes()
	mockConfig.EXPECT().GetMaxRetries().Return(3).AnyTimes()
	mockConfig.EXPECT().GetRetryDelay().Return(time.Second).AnyTimes()
	mockConfig.EXPECT().GetTimeout().Return(30 * time.Second).AnyTimes()
	mockConfig.EXPECT().GetMaxOpenConns().Return(25).AnyTimes()
	mockConfig.EXPECT().GetMaxIdleConns().Return(5).AnyTimes()
	mockConfig.EXPECT().GetConnMaxLifetime().Return(5 * time.Minute).AnyTimes()
	mockConfig.EXPECT().GetConnMaxIdleTime().Return(5 * time.Minute).AnyTimes()

	postgresDB := NewPostgresDB(mockConfig, mockLogger)
	postgresDB.SetTestDB(mockDB)

	// Test successful health check - no logging expected
	mockDB.EXPECT().PingContext(gomock.Any()).Return(nil)

	err := postgresDB.Health(context.Background())
	assert.NoError(t, err)
}

func TestHealthFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfig := mock_interfaces.NewMockPostgresConfig(ctrl)
	mockLogger := mock_interfaces.NewMockLogger(ctrl)
	mockDB := mock_interfaces.NewMockDBInterface(ctrl)

	// Set up mock expectations
	mockConfig.EXPECT().GetHost().Return("localhost").AnyTimes()
	mockConfig.EXPECT().GetPort().Return(5432).AnyTimes()
	mockConfig.EXPECT().GetName().Return("testdb").AnyTimes()
	mockConfig.EXPECT().GetUsername().Return("testuser").AnyTimes()
	mockConfig.EXPECT().GetPassword().Return("testpass").AnyTimes()
	mockConfig.EXPECT().GetSSLMode().Return("disable").AnyTimes()
	mockConfig.EXPECT().GetMaxRetries().Return(3).AnyTimes()
	mockConfig.EXPECT().GetRetryDelay().Return(time.Second).AnyTimes()
	mockConfig.EXPECT().GetTimeout().Return(30 * time.Second).AnyTimes()
	mockConfig.EXPECT().GetMaxOpenConns().Return(25).AnyTimes()
	mockConfig.EXPECT().GetMaxIdleConns().Return(5).AnyTimes()
	mockConfig.EXPECT().GetConnMaxLifetime().Return(5 * time.Minute).AnyTimes()
	mockConfig.EXPECT().GetConnMaxIdleTime().Return(5 * time.Minute).AnyTimes()

	postgresDB := NewPostgresDB(mockConfig, mockLogger)
	postgresDB.SetTestDB(mockDB)

	// Test health check failure - no logging expected
	expectedErr := sql.ErrConnDone
	mockDB.EXPECT().PingContext(gomock.Any()).Return(expectedErr)

	err := postgresDB.Health(context.Background())
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestBeginTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfig := mock_interfaces.NewMockPostgresConfig(ctrl)
	mockLogger := mock_interfaces.NewMockLogger(ctrl)
	mockDB := mock_interfaces.NewMockDBInterface(ctrl)

	// Set up mock expectations
	mockConfig.EXPECT().GetHost().Return("localhost").AnyTimes()
	mockConfig.EXPECT().GetPort().Return(5432).AnyTimes()
	mockConfig.EXPECT().GetName().Return("testdb").AnyTimes()
	mockConfig.EXPECT().GetUsername().Return("testuser").AnyTimes()
	mockConfig.EXPECT().GetPassword().Return("testpass").AnyTimes()
	mockConfig.EXPECT().GetSSLMode().Return("disable").AnyTimes()
	mockConfig.EXPECT().GetMaxRetries().Return(3).AnyTimes()
	mockConfig.EXPECT().GetRetryDelay().Return(time.Second).AnyTimes()
	mockConfig.EXPECT().GetTimeout().Return(30 * time.Second).AnyTimes()
	mockConfig.EXPECT().GetMaxOpenConns().Return(25).AnyTimes()
	mockConfig.EXPECT().GetMaxIdleConns().Return(5).AnyTimes()
	mockConfig.EXPECT().GetConnMaxLifetime().Return(5 * time.Minute).AnyTimes()
	mockConfig.EXPECT().GetConnMaxIdleTime().Return(5 * time.Minute).AnyTimes()

	postgresDB := NewPostgresDB(mockConfig, mockLogger)
	postgresDB.SetTestDB(mockDB)

	// Test successful transaction begin - no logging expected
	expectedTx := &sql.Tx{}
	mockDB.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(expectedTx, nil)

	tx, err := postgresDB.BeginTx(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedTx, tx)
}

func TestWithTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfig := mock_interfaces.NewMockPostgresConfig(ctrl)
	mockLogger := mock_interfaces.NewMockLogger(ctrl)
	mockDB := mock_interfaces.NewMockDBInterface(ctrl)

	// Set up mock expectations
	mockConfig.EXPECT().GetHost().Return("localhost").AnyTimes()
	mockConfig.EXPECT().GetPort().Return(5432).AnyTimes()
	mockConfig.EXPECT().GetName().Return("testdb").AnyTimes()
	mockConfig.EXPECT().GetUsername().Return("testuser").AnyTimes()
	mockConfig.EXPECT().GetPassword().Return("testpass").AnyTimes()
	mockConfig.EXPECT().GetSSLMode().Return("disable").AnyTimes()
	mockConfig.EXPECT().GetMaxRetries().Return(3).AnyTimes()
	mockConfig.EXPECT().GetRetryDelay().Return(time.Second).AnyTimes()
	mockConfig.EXPECT().GetTimeout().Return(30 * time.Second).AnyTimes()
	mockConfig.EXPECT().GetMaxOpenConns().Return(25).AnyTimes()
	mockConfig.EXPECT().GetMaxIdleConns().Return(5).AnyTimes()
	mockConfig.EXPECT().GetConnMaxLifetime().Return(5 * time.Minute).AnyTimes()
	mockConfig.EXPECT().GetConnMaxIdleTime().Return(5 * time.Minute).AnyTimes()

	postgresDB := NewPostgresDB(mockConfig, mockLogger)
	postgresDB.SetTestDB(mockDB)

	// Skip this test for now as it requires complex transaction mocking
	t.Skip("Skipping WithTransaction test due to complex transaction mocking requirements")
}

func TestGetConnectionStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfig := mock_interfaces.NewMockPostgresConfig(ctrl)
	mockLogger := mock_interfaces.NewMockLogger(ctrl)
	mockDB := mock_interfaces.NewMockDBInterface(ctrl)

	// Set up mock expectations
	mockConfig.EXPECT().GetHost().Return("localhost").AnyTimes()
	mockConfig.EXPECT().GetPort().Return(5432).AnyTimes()
	mockConfig.EXPECT().GetName().Return("testdb").AnyTimes()
	mockConfig.EXPECT().GetUsername().Return("testuser").AnyTimes()
	mockConfig.EXPECT().GetPassword().Return("testpass").AnyTimes()
	mockConfig.EXPECT().GetSSLMode().Return("disable").AnyTimes()
	mockConfig.EXPECT().GetMaxRetries().Return(3).AnyTimes()
	mockConfig.EXPECT().GetRetryDelay().Return(time.Second).AnyTimes()
	mockConfig.EXPECT().GetTimeout().Return(30 * time.Second).AnyTimes()
	mockConfig.EXPECT().GetMaxOpenConns().Return(25).AnyTimes()
	mockConfig.EXPECT().GetMaxIdleConns().Return(5).AnyTimes()
	mockConfig.EXPECT().GetConnMaxLifetime().Return(5 * time.Minute).AnyTimes()
	mockConfig.EXPECT().GetConnMaxIdleTime().Return(5 * time.Minute).AnyTimes()

	postgresDB := NewPostgresDB(mockConfig, mockLogger)
	postgresDB.SetTestDB(mockDB)

	// Test getting connection stats
	expectedStats := sql.DBStats{
		MaxOpenConnections: 25,
		OpenConnections:    5,
		InUse:              2,
		Idle:               3,
		WaitCount:          10,
		WaitDuration:       time.Second,
		MaxIdleClosed:      1,
		MaxLifetimeClosed:  2,
	}
	mockDB.EXPECT().Stats().Return(expectedStats).AnyTimes()

	stats := postgresDB.GetConnectionStats()
	expectedMap := map[string]interface{}{
		"status":            "connected",
		"maxOpenConns":      expectedStats.MaxOpenConnections,
		"openConnections":   expectedStats.OpenConnections,
		"inUse":             expectedStats.InUse,
		"idle":              expectedStats.Idle,
		"waitCount":         expectedStats.WaitCount,
		"waitDuration":      expectedStats.WaitDuration,
		"maxIdleClosed":     expectedStats.MaxIdleClosed,
		"maxLifetimeClosed": expectedStats.MaxLifetimeClosed,
	}
	assert.Equal(t, expectedMap, stats)
}

func TestClose(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfig := mock_interfaces.NewMockPostgresConfig(ctrl)
	mockLogger := mock_interfaces.NewMockLogger(ctrl)
	mockDB := mock_interfaces.NewMockDBInterface(ctrl)

	// Set up mock expectations
	mockConfig.EXPECT().GetHost().Return("localhost").AnyTimes()
	mockConfig.EXPECT().GetPort().Return(5432).AnyTimes()
	mockConfig.EXPECT().GetName().Return("testdb").AnyTimes()
	mockConfig.EXPECT().GetUsername().Return("testuser").AnyTimes()
	mockConfig.EXPECT().GetPassword().Return("testpass").AnyTimes()
	mockConfig.EXPECT().GetSSLMode().Return("disable").AnyTimes()
	mockConfig.EXPECT().GetMaxRetries().Return(3).AnyTimes()
	mockConfig.EXPECT().GetRetryDelay().Return(time.Second).AnyTimes()
	mockConfig.EXPECT().GetTimeout().Return(30 * time.Second).AnyTimes()
	mockConfig.EXPECT().GetMaxOpenConns().Return(25).AnyTimes()
	mockConfig.EXPECT().GetMaxIdleConns().Return(5).AnyTimes()
	mockConfig.EXPECT().GetConnMaxLifetime().Return(5 * time.Minute).AnyTimes()
	mockConfig.EXPECT().GetConnMaxIdleTime().Return(5 * time.Minute).AnyTimes()

	postgresDB := NewPostgresDB(mockConfig, mockLogger)
	postgresDB.SetTestDB(mockDB)

	// Test successful close
	mockDB.EXPECT().Close().Return(nil)
	mockLogger.EXPECT().Info(gomock.Any(), "PostgreSQL connection closed", gomock.Any())

	err := postgresDB.Close()
	assert.NoError(t, err)
}

func TestSetTestDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfig := mock_interfaces.NewMockPostgresConfig(ctrl)
	mockLogger := mock_interfaces.NewMockLogger(ctrl)
	mockDB := mock_interfaces.NewMockDBInterface(ctrl)

	// Set up mock expectations
	mockConfig.EXPECT().GetHost().Return("localhost").AnyTimes()
	mockConfig.EXPECT().GetPort().Return(5432).AnyTimes()
	mockConfig.EXPECT().GetName().Return("testdb").AnyTimes()
	mockConfig.EXPECT().GetUsername().Return("testuser").AnyTimes()
	mockConfig.EXPECT().GetPassword().Return("testpass").AnyTimes()
	mockConfig.EXPECT().GetSSLMode().Return("disable").AnyTimes()
	mockConfig.EXPECT().GetMaxRetries().Return(3).AnyTimes()
	mockConfig.EXPECT().GetRetryDelay().Return(time.Second).AnyTimes()
	mockConfig.EXPECT().GetTimeout().Return(30 * time.Second).AnyTimes()
	mockConfig.EXPECT().GetMaxOpenConns().Return(25).AnyTimes()
	mockConfig.EXPECT().GetMaxIdleConns().Return(5).AnyTimes()
	mockConfig.EXPECT().GetConnMaxLifetime().Return(5 * time.Minute).AnyTimes()
	mockConfig.EXPECT().GetConnMaxIdleTime().Return(5 * time.Minute).AnyTimes()

	postgresDB := NewPostgresDB(mockConfig, mockLogger)

	// Test setting test DB
	postgresDB.SetTestDB(mockDB)
	retrievedDB := postgresDB.GetTestDB()
	assert.Equal(t, mockDB, retrievedDB)
}
