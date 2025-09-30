# Message Catalog Service

A comprehensive message catalog service that provides multi-language message management with caching, hot-reloading, and thread-safe operations.

## Features

- **Multi-language Support**: Load messages in different languages
- **Caching**: High-performance in-memory caching with thread safety
- **Hot Reloading**: Dynamic catalog reloading without service restart
- **Error Handling**: Comprehensive error handling with AppError pattern
- **Health Monitoring**: Service health checks and monitoring
- **Generic Design**: Extensible for different message schemas

## Quick Start

### 1. Configuration

Add to your `config.json`:

```json
{
  "message_catalog": {
    "default_language": "en",
    "supported_languages": ["en", "fr"],
    "catalog_path": "./pkg/alert/catalog/message_catalog",
    "reload_interval_seconds": 300,
    "cache_enabled": true
  }
}
```

### 2. Basic Usage

```go
// Initialize service
messageCatalogService := messagecatalog.NewMessageCatalogService(cfg.GetMessageCatalog(), logger)

// Get message by code (default language)
message, err := messageCatalogService.GetMessageByCode(ctx, "ABC001")

// Get message in specific language
message, err := messageCatalogService.GetMessage(ctx, "ABC001", "fr")

// Get messages by category
messages, err := messageCatalogService.GetMessagesByCategory(ctx, "registration", "en")

// Get messages by severity
messages, err := messageCatalogService.GetMessagesBySeverity(ctx, "critical", "en")
```

### 3. Health Monitoring

```go
// Health check
err := messageCatalogService.HealthCheck(ctx)

// Reload catalog
err := messageCatalogService.ReloadCatalog(ctx)

// List available languages
languages, err := messageCatalogService.ListAvailableLanguages(ctx)
```

## Message File Format

Create JSON files for each language in the catalog path:

```json
{
  "MESSAGE_CODE": {
    "messagecode": "MESSAGE_CODE",
    "category": "category_name",
    "severity": "severity_level",
    "message": "Short message",
    "detailed_description": "Detailed message with %s parameters",
    "response_action": "Action to take",
    "metadata": {
      "priority": "high",
      "auto_resolve": "false"
    }
  }
}
```

## API Reference

### Service Interface

```go
type Service interface {
    // Get message by code and language
    GetMessage(ctx context.Context, messageCode, language string) (*Message, error)
    
    // Get message by code using default language
    GetMessageByCode(ctx context.Context, messageCode string) (*Message, error)
    
    // Get messages by category and language
    GetMessagesByCategory(ctx context.Context, category, language string) ([]*Message, error)
    
    // Get messages by severity and language
    GetMessagesBySeverity(ctx context.Context, severity, language string) ([]*Message, error)
    
    // List available languages
    ListAvailableLanguages(ctx context.Context) ([]string, error)
    
    // Reload catalog from files
    ReloadCatalog(ctx context.Context) error
    
    // Health check
    HealthCheck(ctx context.Context) error
}
```

### Message Model

```go
type Message struct {
    MessageCode         string            `json:"messagecode"`
    Category            string            `json:"category"`
    Severity            string            `json:"severity"`
    Message             string            `json:"message"`
    DetailedDescription string            `json:"detailed_description"`
    ResponseAction      string            `json:"response_action"`
    Metadata            map[string]string `json:"metadata,omitempty"`
    CreatedAt           time.Time         `json:"created_at"`
    UpdatedAt           time.Time         `json:"updated_at"`
}
```

## Performance Considerations

### Caching
- Messages are cached in memory for fast access
- Cache is thread-safe with read-write mutex
- Configurable cache settings
- Automatic cache invalidation on reload

### File Loading
- Lazy loading: files are loaded only when needed
- Batch loading: all languages loaded during reload
- Error handling: graceful degradation on file errors
- Hot reloading: catalog can be reloaded without restart

### Memory Management
- Efficient data structures
- Garbage collection friendly
- Configurable cache limits
- Memory leak prevention

## Error Handling

The service uses the AppError pattern for consistent error handling:

- **File Not Found**: `ErrCodeNotFound` when language file doesn't exist
- **JSON Parse Error**: `ErrCodeBadRequest` when JSON is malformed
- **Message Not Found**: `ErrCodeNotFound` when message code doesn't exist
- **Directory Error**: `ErrCodeInternalServer` when directory access fails

## Testing

Run the example to test the service:

```bash
go run ./examples/message_catalog_example.go
```

## Dependencies

- `tushartemplategin/pkg/config` - Configuration management
- `tushartemplategin/pkg/errors` - Error handling
- `tushartemplategin/pkg/interfaces` - Logger interface

## Design Documents

- [Design Document](DESIGN_DOCUMENT.md) - Complete architectural overview
- [Flow Diagram](diagrams/flow_diagram.mmd) - Service flow visualization
- [Class Diagram](diagrams/class_diagram.mmd) - Class relationships
- [Sequence Diagram](diagrams/sequence_diagram.mmd) - Interaction sequences
- [Error Handling](diagrams/error_handling_sequence.mmd) - Error scenarios
- [Performance](diagrams/performance_optimization.mmd) - Performance optimization
