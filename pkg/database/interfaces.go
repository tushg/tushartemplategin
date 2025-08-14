package database

import (
	"context"
	"database/sql"
)

// Database defines the interface for database operations
type Database interface {
	// Connection management
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Health(ctx context.Context) error

	// Transaction management
	BeginTx(ctx context.Context) (*sql.Tx, error)
	WithTransaction(ctx context.Context, fn func(*sql.Tx) error) error

	// Raw database access
	Driver() *sql.DB

	// Close resources
	Close() error
}

// Repository defines the base interface for all repositories
type Repository interface {
	// Health check for the repository
	Health(ctx context.Context) error
}
