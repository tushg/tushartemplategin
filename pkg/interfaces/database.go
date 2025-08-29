package interfaces

import (
	"context"
	"database/sql"
	"time"
)

// Database defines the main interface for database operations
type Database interface {
	// Connection management
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Health(ctx context.Context) error

	// Transaction management
	BeginTx(ctx context.Context) (*sql.Tx, error)
	WithTransaction(ctx context.Context, fn func(*sql.Tx) error) error

	// Close resources
	Close() error
}

// DBInterface defines the interface for database operations
// This allows us to mock the database for testing
type DBInterface interface {
	PingContext(ctx context.Context) error
	Close() error
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Stats() sql.DBStats
	SetMaxOpenConns(n int)
	SetMaxIdleConns(n int)
	SetConnMaxLifetime(d time.Duration)
	SetConnMaxIdleTime(d time.Duration)
}

// TxInterface defines the interface for transaction operations
type TxInterface interface {
	Commit() error
	Rollback() error
}

// Repository defines the base interface for all repositories
// This is a minimal interface that can be extended by domain-specific repositories
type Repository interface {
	// Health check for the repository
	Health(ctx context.Context) error
}
