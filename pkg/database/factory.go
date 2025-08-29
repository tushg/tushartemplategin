package database

import (
	"context"
	"fmt"
	"strings"

	"tushartemplategin/pkg/database/postgres"
	"tushartemplategin/pkg/interfaces"
)

// DatabaseType represents supported database types
type DatabaseType string

const (
	DatabaseTypePostgreSQL DatabaseType = "postgres"
	DatabaseTypeSQLite     DatabaseType = "sqlite"
	DatabaseTypeMySQL      DatabaseType = "mysql"
)

// String returns the string representation of the database type
func (dt DatabaseType) String() string {
	return string(dt)
}

// IsValid checks if the database type is supported
func (dt DatabaseType) IsValid() bool {
	switch dt {
	case DatabaseTypePostgreSQL, DatabaseTypeSQLite, DatabaseTypeMySQL:
		return true
	default:
		return false
	}
}

// DatabaseFactory creates and manages database instances
type DatabaseFactory struct {
	logger interfaces.Logger
}

// NewDatabaseFactory creates a new database factory instance
func NewDatabaseFactory(log interfaces.Logger) *DatabaseFactory {
	return &DatabaseFactory{
		logger: log,
	}
}

// CreateDatabase creates a database instance based on configuration
func (df *DatabaseFactory) CreateDatabase(cfg interfaces.DatabaseConfig) (interfaces.Database, error) {
	if cfg == nil {
		return nil, fmt.Errorf("database configuration is required")
	}

	// Validate database type
	dbType := DatabaseType(cfg.GetType())
	if !dbType.IsValid() {
		return nil, fmt.Errorf("unsupported database type: %s. Supported types: %s",
			cfg.GetType(), df.getSupportedTypes())
	}

	df.logger.Info(context.Background(), "Creating database instance", interfaces.Fields{
		"type": dbType.String(),
		"host": df.getDatabaseHost(cfg),
	})

	// Create database instance based on type
	var db interfaces.Database
	var err error

	switch dbType {
	case DatabaseTypePostgreSQL:
		db, err = df.createPostgreSQL(cfg)
	case DatabaseTypeSQLite:
		db, err = df.createSQLite(cfg)
	case DatabaseTypeMySQL:
		db, err = df.createMySQL(cfg)
	default:
		return nil, fmt.Errorf("database type %s is not implemented yet", dbType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create %s database: %w", dbType, err)
	}

	df.logger.Info(context.Background(), "Database instance created successfully", interfaces.Fields{
		"type": dbType.String(),
	})

	return db, nil
}

// createPostgreSQL creates a PostgreSQL database instance
func (df *DatabaseFactory) createPostgreSQL(cfg interfaces.DatabaseConfig) (interfaces.Database, error) {
	if cfg.GetPostgres() == nil {
		return nil, fmt.Errorf("PostgreSQL configuration is required")
	}

	// Validate required PostgreSQL configuration
	if err := df.validatePostgreSQLConfig(cfg.GetPostgres()); err != nil {
		return nil, fmt.Errorf("invalid PostgreSQL configuration: %w", err)
	}

	return postgres.NewPostgresDB(cfg.GetPostgres(), df.logger), nil
}

// createSQLite creates a SQLite database instance
func (df *DatabaseFactory) createSQLite(cfg interfaces.DatabaseConfig) (interfaces.Database, error) {
	if cfg.GetSQLite() == nil {
		return nil, fmt.Errorf("SQLite configuration is required")
	}

	// Validate required SQLite configuration
	if err := df.validateSQLiteConfig(cfg.GetSQLite()); err != nil {
		return nil, fmt.Errorf("invalid SQLite configuration: %w", err)
	}

	// TODO: Implement SQLite database when needed
	// return sqlite.NewSQLiteDB(cfg.SQLite, df.logger), nil
	return nil, fmt.Errorf("SQLite database implementation not yet available")
}

// createMySQL creates a MySQL database instance
func (df *DatabaseFactory) createMySQL(cfg interfaces.DatabaseConfig) (interfaces.Database, error) {
	if cfg.GetMySQL() == nil {
		return nil, fmt.Errorf("MySQL configuration is required")
	}

	// TODO: Implement MySQL database when needed
	// return mysql.NewMySQLDB(cfg.MySQL, df.logger), nil
	return nil, fmt.Errorf("MySQL database implementation not yet available")
}

// validatePostgreSQLConfig validates PostgreSQL configuration
func (df *DatabaseFactory) validatePostgreSQLConfig(cfg interfaces.PostgresConfig) error {
	if cfg.GetHost() == "" {
		return fmt.Errorf("host is required")
	}
	if cfg.GetPort() <= 0 || cfg.GetPort() > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	if cfg.GetName() == "" {
		return fmt.Errorf("database name is required")
	}
	if cfg.GetUsername() == "" {
		return fmt.Errorf("username is required")
	}
	if cfg.GetPassword() == "" {
		df.logger.Warn(context.Background(), "PostgreSQL password is empty - this may cause connection issues", interfaces.Fields{})
	}
	return nil
}

// validateSQLiteConfig validates SQLite configuration
func (df *DatabaseFactory) validateSQLiteConfig(cfg interfaces.SQLiteConfig) error {
	if cfg.GetPath() == "" {
		return fmt.Errorf("file path is required")
	}
	return nil
}

// getDatabaseHost returns the database host for logging purposes
func (df *DatabaseFactory) getDatabaseHost(cfg interfaces.DatabaseConfig) string {
	switch DatabaseType(cfg.GetType()) {
	case DatabaseTypePostgreSQL:
		if cfg.GetPostgres() != nil {
			return fmt.Sprintf("%s:%d", cfg.GetPostgres().GetHost(), cfg.GetPostgres().GetPort())
		}
	case DatabaseTypeSQLite:
		if cfg.GetSQLite() != nil {
			return cfg.GetSQLite().GetPath()
		}
	case DatabaseTypeMySQL:
		if cfg.GetMySQL() != nil {
			return fmt.Sprintf("%s:%d", cfg.GetMySQL().GetHost(), cfg.GetMySQL().GetPort())
		}
	}
	return "unknown"
}

// getSupportedTypes returns a comma-separated list of supported database types
func (df *DatabaseFactory) getSupportedTypes() string {
	types := []string{
		DatabaseTypePostgreSQL.String(),
		DatabaseTypeSQLite.String(),
		DatabaseTypeMySQL.String(),
	}
	return strings.Join(types, ", ")
}

// GetDatabaseType returns the database type from configuration
func (df *DatabaseFactory) GetDatabaseType(cfg interfaces.DatabaseConfig) DatabaseType {
	return DatabaseType(cfg.GetType())
}

// IsDatabaseTypeSupported checks if a specific database type is supported
func (df *DatabaseType) IsDatabaseTypeSupported(dbType string) bool {
	return DatabaseType(dbType).IsValid()
}
