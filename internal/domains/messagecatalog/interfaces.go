package messagecatalog

import (
	"context"
)

// Service defines the interface for message catalog operations
type Service interface {
	// Message Operations
	GetMessage(ctx context.Context, messageCode, language string) (*Message, error)
	GetMessageByCode(ctx context.Context, messageCode string) (*Message, error)
	GetMessagesByCategory(ctx context.Context, category, language string) ([]*Message, error)
	GetMessagesBySeverity(ctx context.Context, severity, language string) ([]*Message, error)
	ListAvailableLanguages(ctx context.Context) ([]string, error)

	// Catalog Management
	ReloadCatalog(ctx context.Context) error
	HealthCheck(ctx context.Context) error
}
