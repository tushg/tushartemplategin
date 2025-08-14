package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"tushartemplategin/pkg/logger"
)

// TxManager handles database transactions with proper error handling
type TxManager struct {
	db     Database
	logger logger.Logger
}

// NewTxManager creates a new transaction manager
func NewTxManager(db Database, log logger.Logger) *TxManager {
	return &TxManager{
		db:     db,
		logger: log,
	}
}

// WithTransaction executes a function within a transaction with automatic rollback on error
func (tm *TxManager) WithTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	return tm.db.WithTransaction(ctx, fn)
}

// WithReadOnlyTransaction executes a function within a read-only transaction
func (tm *TxManager) WithReadOnlyTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := tm.db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin read-only transaction: %w", err)
	}

	// Ensure transaction is rolled back
	defer func() {
		if p := recover(); p != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				tm.logger.Error(ctx, "Failed to rollback read-only transaction on panic", logger.Fields{
					"panic": p,
					"error": rbErr.Error(),
				})
			}
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("read-only transaction failed and rollback failed: %w (rollback error: %v)", err, rbErr)
		}
		return err
	}

	// Always rollback read-only transactions
	if err := tx.Rollback(); err != nil {
		return fmt.Errorf("failed to rollback read-only transaction: %w", err)
	}

	return nil
}

// WithTimeout executes a function with a timeout context
func (tm *TxManager) WithTimeout(ctx context.Context, timeout time.Duration, fn func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return fn(ctx)
}
