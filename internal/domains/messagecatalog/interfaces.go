package messagecatalog

import "context"

// Service defines the interface for message catalog operations
type Service interface {
	// Message operations
	GetMessage(ctx context.Context, req *MessageRequest) (*MessageResponse, error)
	GetMessageByCode(ctx context.Context, messageCode, catalogName, language string) (*MessageResponse, error)
	GetMessagesByCategory(ctx context.Context, category, catalogName, language string) ([]*MessageResponse, error)
	GetMessagesBySeverity(ctx context.Context, severity, catalogName, language string) ([]*MessageResponse, error)

	// Catalog management
	ReloadCatalog(ctx context.Context, catalogName string) error
	ReloadAllCatalogs(ctx context.Context) error
	HealthCheck(ctx context.Context) error

	// Catalog information
	ListAvailableCatalogs(ctx context.Context) ([]string, error)
	ListAvailableLanguages(ctx context.Context, catalogName string) ([]string, error)
	GetCatalogInfo(ctx context.Context, catalogName string) (*CatalogInfo, error)
	GetCatalogStats(ctx context.Context) (*CatalogStats, error)
}

// CatalogLoader defines interface for loading catalog data
type CatalogLoader interface {
	LoadCatalogStructure(ctx context.Context, catalogName string) (map[string]interface{}, error)
	LoadLanguageFile(ctx context.Context, catalogName, language string) (map[string]interface{}, error)
	ListAvailableLanguages(ctx context.Context, catalogName string) ([]string, error)
}

// CacheManager defines interface for caching operations
type CacheManager interface {
	Get(ctx context.Context, key string) (interface{}, bool)
	Set(ctx context.Context, key string, value interface{}, ttl int) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
	GetStats(ctx context.Context) map[string]interface{}
}
