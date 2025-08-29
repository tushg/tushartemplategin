package interfaces

import (
	"context"
)

// Logger defines the interface for logging operations
type Logger interface {
	Debug(ctx context.Context, msg string, fields Fields)
	Info(ctx context.Context, msg string, fields Fields)
	Warn(ctx context.Context, msg string, fields Fields)
	Error(ctx context.Context, msg string, fields Fields)
	Fatal(ctx context.Context, msg string, err error, fields Fields)
}

// Fields represents key-value pairs for structured logging
type Fields map[string]interface{}
