package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"tushartemplategin/pkg/interfaces"
)

// PostgresDB implements the Database interface for PostgreSQL
type PostgresDB struct {
	config interfaces.PostgresConfig
	db     interfaces.DBInterface
	logger interfaces.Logger
}

// NewPostgresDB creates a new PostgreSQL database instance
func NewPostgresDB(cfg interfaces.PostgresConfig, log interfaces.Logger) *PostgresDB {
	return &PostgresDB{
		config: cfg,
		logger: log,
	}
}

// SetTestDB sets the database connection for testing purposes
func (p *PostgresDB) SetTestDB(db interfaces.DBInterface) {
	p.db = db
}

// GetTestDB returns the database interface for testing
func (p *PostgresDB) GetTestDB() interfaces.DBInterface {
	return p.db
}

// Connect establishes a connection to PostgreSQL with retry logic
func (p *PostgresDB) Connect(ctx context.Context) error {
	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		p.config.GetHost(), p.config.GetPort(), p.config.GetName(), p.config.GetUsername(), p.config.GetPassword(), p.config.GetSSLMode())

	var db *sql.DB
	var err error

	// Retry connection with exponential backoff
	for attempt := 0; attempt <= p.config.GetMaxRetries(); attempt++ {
		if attempt > 0 {
			p.logger.Info(ctx, "Retrying database connection", interfaces.Fields{
				"attempt":    attempt,
				"maxRetries": p.config.GetMaxRetries(),
			})

			// Wait before retry (exponential backoff)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(p.config.GetRetryDelay() * time.Duration(attempt)):
			}
		}

		db, err = sql.Open("postgres", dsn)
		if err != nil {
			p.logger.Error(ctx, "Failed to open database connection", interfaces.Fields{
				"attempt": attempt,
				"error":   err.Error(),
			})
			continue
		}

		// Configure connection pool for production
		db.SetMaxOpenConns(p.config.GetMaxOpenConns())
		db.SetMaxIdleConns(p.config.GetMaxIdleConns())
		db.SetConnMaxLifetime(p.config.GetConnMaxLifetime())
		db.SetConnMaxIdleTime(p.config.GetConnMaxIdleTime())

		// Test connection with timeout
		connCtx, cancel := context.WithTimeout(ctx, p.config.GetTimeout())
		if err := db.PingContext(connCtx); err != nil {
			cancel()
			p.logger.Error(ctx, "Failed to ping database", interfaces.Fields{
				"attempt": attempt,
				"error":   err.Error(),
			})
			db.Close()
			continue
		}
		cancel()

		// Connection successful
		p.db = db
		p.logger.Info(ctx, "PostgreSQL connection established", interfaces.Fields{
			"host":         p.config.GetHost(),
			"port":         p.config.GetPort(),
			"db":           p.config.GetName(),
			"maxOpenConns": p.config.GetMaxOpenConns(),
			"maxIdleConns": p.config.GetMaxIdleConns(),
			"attempts":     attempt + 1,
		})
		return nil
	}

	return fmt.Errorf("failed to connect to database after %d attempts: %w", p.config.GetMaxRetries()+1, err)
}

// Disconnect closes the database connection
func (p *PostgresDB) Disconnect(ctx context.Context) error {
	if p.db != nil {
		if err := p.db.Close(); err != nil {
			p.logger.Error(ctx, "Failed to close database connection", interfaces.Fields{"error": err.Error()})
			return err
		}
		p.logger.Info(ctx, "PostgreSQL connection closed", interfaces.Fields{})
	}
	return nil
}

// Health checks database connectivity
func (p *PostgresDB) Health(ctx context.Context) error {
	if p.db == nil {
		return fmt.Errorf("database not connected")
	}
	return p.db.PingContext(ctx)
}

// BeginTx starts a new transaction
func (p *PostgresDB) BeginTx(ctx context.Context) (*sql.Tx, error) {
	if p.db == nil {
		return nil, fmt.Errorf("database not connected")
	}

	ctx, cancel := context.WithTimeout(ctx, p.config.GetTimeout())
	defer cancel()

	return p.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
}

// WithTransaction executes a function within a transaction
func (p *PostgresDB) WithTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := p.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure transaction is rolled back on error or panic
	defer func() {
		if panicVal := recover(); panicVal != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				p.logger.Error(ctx, "Failed to rollback transaction on panic", interfaces.Fields{
					"panic": panicVal,
					"error": rbErr.Error(),
				})
			}
			panic(panicVal) // re-throw panic after rollback
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed and rollback failed: %w (rollback error: %v)", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// GetConnectionStats returns connection pool statistics for monitoring
func (p *PostgresDB) GetConnectionStats() map[string]interface{} {
	if p.db == nil {
		return map[string]interface{}{
			"status": "disconnected",
		}
	}

	return map[string]interface{}{
		"status":            "connected",
		"maxOpenConns":      p.db.Stats().MaxOpenConnections,
		"openConnections":   p.db.Stats().OpenConnections,
		"inUse":             p.db.Stats().InUse,
		"idle":              p.db.Stats().Idle,
		"waitCount":         p.db.Stats().WaitCount,
		"waitDuration":      p.db.Stats().WaitDuration,
		"maxIdleClosed":     p.db.Stats().MaxIdleClosed,
		"maxLifetimeClosed": p.db.Stats().MaxLifetimeClosed,
	}
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	return p.Disconnect(context.Background())
}
