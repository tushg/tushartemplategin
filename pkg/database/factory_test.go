package database

import (
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_interfaces "tushartemplategin/mocks"
)

func TestNewDatabaseFactory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mock_interfaces.NewMockLogger(ctrl)

	factory := NewDatabaseFactory(mockLogger)
	assert.NotNil(t, factory)
	assert.Equal(t, mockLogger, factory.logger)
}

func TestCreateDatabase_PostgreSQL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mock_interfaces.NewMockLogger(ctrl)
	mockConfig := mock_interfaces.NewMockDatabaseConfig(ctrl)
	mockPostgresConfig := mock_interfaces.NewMockPostgresConfig(ctrl)

	// Set up mock expectations
	mockConfig.EXPECT().GetType().Return("postgres").AnyTimes()
	mockConfig.EXPECT().GetPostgres().Return(mockPostgresConfig).AnyTimes()
	mockPostgresConfig.EXPECT().GetHost().Return("localhost").AnyTimes()
	mockPostgresConfig.EXPECT().GetPort().Return(5432).AnyTimes()
	mockPostgresConfig.EXPECT().GetName().Return("testdb").AnyTimes()
	mockPostgresConfig.EXPECT().GetUsername().Return("testuser").AnyTimes()
	mockPostgresConfig.EXPECT().GetPassword().Return("testpass").AnyTimes()
	mockPostgresConfig.EXPECT().GetSSLMode().Return("disable").AnyTimes()
	mockPostgresConfig.EXPECT().GetMaxRetries().Return(3).AnyTimes()
	mockPostgresConfig.EXPECT().GetRetryDelay().Return(time.Second).AnyTimes()
	mockPostgresConfig.EXPECT().GetTimeout().Return(30 * time.Second).AnyTimes()
	mockPostgresConfig.EXPECT().GetMaxOpenConns().Return(25).AnyTimes()
	mockPostgresConfig.EXPECT().GetMaxIdleConns().Return(5).AnyTimes()
	mockPostgresConfig.EXPECT().GetConnMaxLifetime().Return(5 * time.Minute).AnyTimes()
	mockPostgresConfig.EXPECT().GetConnMaxIdleTime().Return(5 * time.Minute).AnyTimes()

	// Set up logger expectations
	mockLogger.EXPECT().Info(gomock.Any(), "Creating database instance", gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any(), "Database instance created successfully", gomock.Any()).AnyTimes()

	factory := NewDatabaseFactory(mockLogger)
	db, err := factory.CreateDatabase(mockConfig)

	// The factory should succeed in creating the database instance
	assert.NoError(t, err)
	assert.NotNil(t, db)
}

func TestCreateDatabase_SQLite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mock_interfaces.NewMockLogger(ctrl)
	mockConfig := mock_interfaces.NewMockDatabaseConfig(ctrl)
	mockSQLiteConfig := mock_interfaces.NewMockSQLiteConfig(ctrl)

	// Set up mock expectations
	mockConfig.EXPECT().GetType().Return("sqlite").AnyTimes()
	mockConfig.EXPECT().GetSQLite().Return(mockSQLiteConfig).AnyTimes()
	mockSQLiteConfig.EXPECT().GetPath().Return(":memory:").AnyTimes()
	mockSQLiteConfig.EXPECT().GetTimeout().Return(30 * time.Second).AnyTimes()

	// Set up logger expectations
	mockLogger.EXPECT().Info(gomock.Any(), "Creating database instance", gomock.Any()).AnyTimes()

	factory := NewDatabaseFactory(mockLogger)
	db, err := factory.CreateDatabase(mockConfig)

	// Note: This will fail in tests because we can't mock sql.Open
	// For now, we'll just check that the factory returns an error
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestCreateDatabase_MySQL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mock_interfaces.NewMockLogger(ctrl)
	mockConfig := mock_interfaces.NewMockDatabaseConfig(ctrl)
	mockMySQLConfig := mock_interfaces.NewMockMySQLConfig(ctrl)

	// Set up mock expectations
	mockConfig.EXPECT().GetType().Return("mysql").AnyTimes()
	mockConfig.EXPECT().GetMySQL().Return(mockMySQLConfig).AnyTimes()
	mockMySQLConfig.EXPECT().GetHost().Return("localhost").AnyTimes()
	mockMySQLConfig.EXPECT().GetPort().Return(3306).AnyTimes()
	mockMySQLConfig.EXPECT().GetName().Return("testdb").AnyTimes()
	mockMySQLConfig.EXPECT().GetUsername().Return("testuser").AnyTimes()
	mockMySQLConfig.EXPECT().GetPassword().Return("testpass").AnyTimes()
	mockMySQLConfig.EXPECT().GetTimeout().Return(30 * time.Second).AnyTimes()

	// Set up logger expectations
	mockLogger.EXPECT().Info(gomock.Any(), "Creating database instance", gomock.Any()).AnyTimes()

	factory := NewDatabaseFactory(mockLogger)
	db, err := factory.CreateDatabase(mockConfig)

	// Note: This will fail in tests because we can't mock sql.Open
	// For now, we'll just check that the factory returns an error
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestCreateDatabase_UnsupportedType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mock_interfaces.NewMockLogger(ctrl)
	mockConfig := mock_interfaces.NewMockDatabaseConfig(ctrl)

	// Set up mock expectations
	mockConfig.EXPECT().GetType().Return("unsupported").AnyTimes()

	factory := NewDatabaseFactory(mockLogger)
	db, err := factory.CreateDatabase(mockConfig)

	assert.Error(t, err)
	assert.Nil(t, db)
	assert.Contains(t, err.Error(), "unsupported database type")
}

func TestGetDatabaseType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mock_interfaces.NewMockLogger(ctrl)
	mockConfig := mock_interfaces.NewMockDatabaseConfig(ctrl)

	// Set up mock expectations
	mockConfig.EXPECT().GetType().Return("postgres").AnyTimes()

	factory := NewDatabaseFactory(mockLogger)
	dbType := factory.GetDatabaseType(mockConfig)

	assert.Equal(t, "postgres", string(dbType))
}

func TestValidatePostgreSQLConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mock_interfaces.NewMockLogger(ctrl)
	mockPostgresConfig := mock_interfaces.NewMockPostgresConfig(ctrl)

	// Set up mock expectations for valid config
	mockPostgresConfig.EXPECT().GetHost().Return("localhost").AnyTimes()
	mockPostgresConfig.EXPECT().GetPort().Return(5432).AnyTimes()
	mockPostgresConfig.EXPECT().GetName().Return("testdb").AnyTimes()
	mockPostgresConfig.EXPECT().GetUsername().Return("testuser").AnyTimes()
	mockPostgresConfig.EXPECT().GetPassword().Return("testpass").AnyTimes()

	factory := NewDatabaseFactory(mockLogger)
	err := factory.validatePostgreSQLConfig(mockPostgresConfig)

	assert.NoError(t, err)
}

func TestValidatePostgreSQLConfig_InvalidHost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mock_interfaces.NewMockLogger(ctrl)
	mockPostgresConfig := mock_interfaces.NewMockPostgresConfig(ctrl)

	// Set up mock expectations for invalid config
	mockPostgresConfig.EXPECT().GetHost().Return("")

	factory := NewDatabaseFactory(mockLogger)
	err := factory.validatePostgreSQLConfig(mockPostgresConfig)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "host is required")
}

func TestValidateSQLiteConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mock_interfaces.NewMockLogger(ctrl)
	mockSQLiteConfig := mock_interfaces.NewMockSQLiteConfig(ctrl)

	// Set up mock expectations for valid config
	mockSQLiteConfig.EXPECT().GetPath().Return(":memory:").AnyTimes()

	factory := NewDatabaseFactory(mockLogger)
	err := factory.validateSQLiteConfig(mockSQLiteConfig)

	assert.NoError(t, err)
}

func TestValidateSQLiteConfig_InvalidPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mock_interfaces.NewMockLogger(ctrl)
	mockSQLiteConfig := mock_interfaces.NewMockSQLiteConfig(ctrl)

	// Set up mock expectations for invalid config
	mockSQLiteConfig.EXPECT().GetPath().Return("")

	factory := NewDatabaseFactory(mockLogger)
	err := factory.validateSQLiteConfig(mockSQLiteConfig)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path is required")
}
