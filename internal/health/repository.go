package health

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"tushartemplategin/pkg/database"
	"tushartemplategin/pkg/logger"
)

// HealthRepository implements the Repository interface for health data
type HealthRepository struct {
	db     database.Database
	txMgr  *database.TxManager
	logger logger.Logger
}

// NewHealthRepository creates a new health repository
func NewHealthRepository(db database.Database, log logger.Logger) Repository {
	return &HealthRepository{
		db:     db,
		txMgr:  database.NewTxManager(db, log),
		logger: log,
	}
}

// GetHealth returns the overall health status of the service
func (r *HealthRepository) GetHealth(ctx context.Context) (*HealthStatus, error) {
	query := `SELECT status, timestamp, service, version FROM health_status WHERE id = $1 AND deleted_at IS NULL`

	var status HealthStatus
	err := r.db.Driver().QueryRowContext(ctx, query, "service_health").Scan(
		&status.Status,
		&status.Timestamp,
		&status.Service,
		&status.Version,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Return default status if no record exists
			return &HealthStatus{
				Status:    "healthy",
				Timestamp: time.Now(),
				Service:   "tushar-service",
				Version:   "1.0.0",
			}, nil
		}
		r.logger.Error(ctx, "Failed to get health status from database", logger.Fields{"error": err.Error()})
		// Return default status if database error
		return &HealthStatus{
			Status:    "healthy",
			Timestamp: time.Now(),
			Service:   "tushar-service",
			Version:   "1.0.0",
		}, nil
	}

	return &status, nil
}

// GetReadiness returns the readiness status for Kubernetes readiness probes
func (r *HealthRepository) GetReadiness(ctx context.Context) (*ReadinessStatus, error) {
	// Check database connection health
	dbHealth := "connected"
	if err := r.db.Health(ctx); err != nil {
		dbHealth = "disconnected"
		r.logger.Error(ctx, "Database health check failed", logger.Fields{"error": err.Error()})
	}

	return &ReadinessStatus{
		Status:    "ready",
		Timestamp: time.Now(),
		Database:  dbHealth,
		Service:   "tushar-service",
	}, nil
}

// GetLiveness returns the liveness status for Kubernetes liveness probes
func (r *HealthRepository) GetLiveness(ctx context.Context) (*LivenessStatus, error) {
	return &LivenessStatus{
		Status:    "alive",
		Timestamp: time.Now(),
		Service:   "gin-service",
	}, nil
}

// UpdateHealth updates the health status in database with transaction support
func (r *HealthRepository) UpdateHealth(ctx context.Context, status *HealthStatus) error {
	return r.txMgr.WithTransaction(ctx, func(tx *sql.Tx) error {
		query := `
			INSERT INTO health_status (id, status, timestamp, service, version, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $6)
			ON CONFLICT (id) 
			DO UPDATE SET 
				status = EXCLUDED.status,
				timestamp = EXCLUDED.timestamp,
				service = EXCLUDED.service,
				version = EXCLUDED.version,
				updated_at = EXCLUDED.updated_at
		`

		_, err := tx.ExecContext(ctx, query,
			"service_health",
			status.Status,
			status.Timestamp,
			status.Service,
			status.Version,
			time.Now(),
		)

		if err != nil {
			return fmt.Errorf("failed to update health status: %w", err)
		}

		return nil
	})
}

// GetHealthHistory retrieves health status history
func (r *HealthRepository) GetHealthHistory(ctx context.Context, limit int) ([]*HealthStatus, error) {
	query := `
		SELECT status, timestamp, service, version 
		FROM health_status 
		WHERE deleted_at IS NULL 
		ORDER BY timestamp DESC 
		LIMIT $1
	`

	rows, err := r.db.Driver().QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query health history: %w", err)
	}
	defer rows.Close()

	var history []*HealthStatus
	for rows.Next() {
		var status HealthStatus
		if err := rows.Scan(&status.Status, &status.Timestamp, &status.Service, &status.Version); err != nil {
			return nil, fmt.Errorf("failed to scan health status: %w", err)
		}
		history = append(history, &status)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating health history: %w", err)
	}

	return history, nil
}
