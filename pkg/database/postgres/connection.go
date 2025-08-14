package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"tushartemplategin/pkg/config"
	"tushartemplategin/pkg/logger"
)

// PostgresDB implements the Database interface for PostgreSQL
type PostgresDB struct {
	config *config.PostgresConfig
	db     *sql.DB
	logger logger.Logger
}

// NewPostgresDB creates a new PostgreSQL database instance
func NewPostgresDB(cfg *config.PostgresConfig, log logger.Logger) *PostgresDB {
	return &PostgresDB{
		config: cfg,
		logger: log,
	}
}

// Connect establishes a connection to PostgreSQL with retry logic
func (p *PostgresDB) Connect(ctx context.Context) error {
	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		p.config.Host, p.config.Port, p.config.Name, p.config.Username, p.config.Password, p.config.SSLMode)

	var db *sql.DB
	var err error

	// Retry connection with exponential backoff
	for attempt := 0; attempt <= p.config.MaxRetries; attempt++ {
		if attempt > 0 {
			p.logger.Info(ctx, "Retrying database connection", logger.Fields{
				"attempt":    attempt,
				"maxRetries": p.config.MaxRetries,
			})

			// Wait before retry (exponential backoff)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(p.config.RetryDelay * time.Duration(attempt)):
			}
		}

		db, err = sql.Open("postgres", dsn)
		if err != nil {
			p.logger.Error(ctx, "Failed to open database connection", logger.Fields{
				"attempt": attempt,
				"error":   err.Error(),
			})
			continue
		}

		// Configure connection pool for production
		db.SetMaxOpenConns(p.config.MaxOpenConns)
		db.SetMaxIdleConns(p.config.MaxIdleConns)
		db.SetConnMaxLifetime(p.config.ConnMaxLifetime)
		db.SetConnMaxIdleTime(p.config.ConnMaxIdleTime)

		// Test connection with timeout
		connCtx, cancel := context.WithTimeout(ctx, p.config.Timeout)
		if err := db.PingContext(connCtx); err != nil {
			cancel()
			p.logger.Error(ctx, "Failed to ping database", logger.Fields{
				"attempt": attempt,
				"error":   err.Error(),
			})
			db.Close()
			continue
		}
		cancel()

		// Connection successful
		p.db = db
		p.logger.Info(ctx, "PostgreSQL connection established", logger.Fields{
			"host":         p.config.Host,
			"port":         p.config.Port,
			"db":           p.config.Name,
			"maxOpenConns": p.config.MaxOpenConns,
			"maxIdleConns": p.config.MaxIdleConns,
			"attempts":     attempt + 1,
		})
		return nil
	}

	return fmt.Errorf("failed to connect to database after %d attempts: %w", p.config.MaxRetries+1, err)
}

// Disconnect closes the database connection
func (p *PostgresDB) Disconnect(ctx context.Context) error {
	if p.db != nil {
		if err := p.db.Close(); err != nil {
			p.logger.Error(ctx, "Failed to close database connection", logger.Fields{"error": err.Error()})
			return err
		}
		p.logger.Info(ctx, "PostgreSQL connection closed", logger.Fields{})
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

	ctx, cancel := context.WithTimeout(ctx, p.config.Timeout)
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
				p.logger.Error(ctx, "Failed to rollback transaction on panic", logger.Fields{
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

// Driver returns the underlying sql.DB instance
func (p *PostgresDB) Driver() *sql.DB {
	return p.db
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
