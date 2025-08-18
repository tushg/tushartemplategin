package database

import (
	"context"
	"fmt"
	"strings"

	"tushartemplategin/pkg/config"
	"tushartemplategin/pkg/database/postgres"
	"tushartemplategin/pkg/logger"
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
	logger logger.Logger
}

// NewDatabaseFactory creates a new database factory instance
func NewDatabaseFactory(log logger.Logger) *DatabaseFactory {
	return &DatabaseFactory{
		logger: log,
	}
}

// CreateDatabase creates a database instance based on configuration
func (df *DatabaseFactory) CreateDatabase(cfg *config.DatabaseConfig) (Database, error) {
	if cfg == nil {
		return nil, fmt.Errorf("database configuration is required")
	}

	// Validate database type
	dbType := DatabaseType(cfg.Type)
	if !dbType.IsValid() {
		return nil, fmt.Errorf("unsupported database type: %s. Supported types: %s",
			cfg.Type, df.getSupportedTypes())
	}

	df.logger.Info(context.Background(), "Creating database instance", logger.Fields{
		"type": dbType.String(),
		"host": df.getDatabaseHost(cfg),
	})

	// Create database instance based on type
	var db Database
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

	df.logger.Info(context.Background(), "Database instance created successfully", logger.Fields{
		"type": dbType.String(),
	})

	return db, nil
}

// createPostgreSQL creates a PostgreSQL database instance
func (df *DatabaseFactory) createPostgreSQL(cfg *config.DatabaseConfig) (Database, error) {
	if cfg.Postgres == nil {
		return nil, fmt.Errorf("PostgreSQL configuration is required")
	}

	// Validate required PostgreSQL configuration
	if err := df.validatePostgreSQLConfig(cfg.Postgres); err != nil {
		return nil, fmt.Errorf("invalid PostgreSQL configuration: %w", err)
	}

	return postgres.NewPostgresDB(cfg.Postgres, df.logger), nil
}

// createSQLite creates a SQLite database instance
func (df *DatabaseFactory) createSQLite(cfg *config.DatabaseConfig) (Database, error) {
	if cfg.SQLite == nil {
		return nil, fmt.Errorf("SQLite configuration is required")
	}

	// Validate required SQLite configuration
	if err := df.validateSQLiteConfig(cfg.SQLite); err != nil {
		return nil, fmt.Errorf("invalid SQLite configuration: %w", err)
	}

	// TODO: Implement SQLite database when needed
	// return sqlite.NewSQLiteDB(cfg.SQLite, df.logger), nil
	return nil, fmt.Errorf("SQLite database implementation not yet available")
}

// createMySQL creates a MySQL database instance
func (df *DatabaseFactory) createMySQL(cfg *config.DatabaseConfig) (Database, error) {
	if cfg.MySQL == nil {
		return nil, fmt.Errorf("MySQL configuration is required")
	}

	// TODO: Implement MySQL database when needed
	// return mysql.NewMySQLDB(cfg.MySQL, df.logger), nil
	return nil, fmt.Errorf("MySQL database implementation not yet available")
}

// validatePostgreSQLConfig validates PostgreSQL configuration
func (df *DatabaseFactory) validatePostgreSQLConfig(cfg *config.PostgresConfig) error {
	if cfg.Host == "" {
		return fmt.Errorf("host is required")
	}
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	if cfg.Name == "" {
		return fmt.Errorf("database name is required")
	}
	if cfg.Username == "" {
		return fmt.Errorf("username is required")
	}
	if cfg.Password == "" {
		df.logger.Warn(context.Background(), "PostgreSQL password is empty - this may cause connection issues", logger.Fields{})
	}
	return nil
}

// validateSQLiteConfig validates SQLite configuration
func (df *DatabaseFactory) validateSQLiteConfig(cfg *config.SQLiteConfig) error {
	if cfg.FilePath == "" {
		return fmt.Errorf("file path is required")
	}
	return nil
}

// getDatabaseHost returns the database host for logging purposes
func (df *DatabaseFactory) getDatabaseHost(cfg *config.DatabaseConfig) string {
	switch DatabaseType(cfg.Type) {
	case DatabaseTypePostgreSQL:
		if cfg.Postgres != nil {
			return fmt.Sprintf("%s:%d", cfg.Postgres.Host, cfg.Postgres.Port)
		}
	case DatabaseTypeSQLite:
		if cfg.SQLite != nil {
			return cfg.SQLite.FilePath
		}
	case DatabaseTypeMySQL:
		if cfg.MySQL != nil {
			return fmt.Sprintf("%s:%d", cfg.MySQL.Host, cfg.MySQL.Port)
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
func (df *DatabaseFactory) GetDatabaseType(cfg *config.DatabaseConfig) DatabaseType {
	return DatabaseType(cfg.Type)
}

// IsDatabaseTypeSupported checks if a specific database type is supported
func (df *DatabaseType) IsDatabaseTypeSupported(dbType string) bool {
	return DatabaseType(dbType).IsValid()
}
