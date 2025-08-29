package interfaces

import (
	"time"
)

// DatabaseConfig defines the interface for database configuration
type DatabaseConfig interface {
	GetType() string
	GetPostgres() PostgresConfig
	GetSQLite() SQLiteConfig
	GetMySQL() MySQLConfig
}

// PostgresConfig defines the interface for PostgreSQL configuration
type PostgresConfig interface {
	GetHost() string
	GetPort() int
	GetName() string
	GetUsername() string
	GetPassword() string
	GetSSLMode() string
	GetMaxRetries() int
	GetRetryDelay() time.Duration
	GetTimeout() time.Duration
	GetMaxOpenConns() int
	GetMaxIdleConns() int
	GetConnMaxLifetime() time.Duration
	GetConnMaxIdleTime() time.Duration
}

// SQLiteConfig defines the interface for SQLite configuration
type SQLiteConfig interface {
	GetPath() string
	GetTimeout() time.Duration
}

// MySQLConfig defines the interface for MySQL configuration
type MySQLConfig interface {
	GetHost() string
	GetPort() int
	GetName() string
	GetUsername() string
	GetPassword() string
	GetTimeout() time.Duration
}
