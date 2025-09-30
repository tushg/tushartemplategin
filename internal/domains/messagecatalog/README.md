# Message Catalog Service

A generic, extensible service that manages message templates across multiple catalogs with multi-language support.

## 🚀 **Quick Start**

### **1. Configuration**

Add message catalog configuration to your `config.json`:

```json
{
  "message_catalog": {
    "default_language": "en-US",
    "supported_languages": ["en-US", "fr-FR", "es-ES"],
    "cache_enabled": true,
    "cache_ttl_seconds": 3600,
    "reload_interval_seconds": 300,
    "catalogs": [
      {
        "name": "alert",
        "path": "./pkg/alert/catalog",
        "enabled": true,
        "structure_file": "messagecatelog.json",
        "language_file_pattern": "messagecatelog-{lang}.json"
      },
      {
        "name": "audit",
        "path": "./pkg/audit/catalog",
        "enabled": true,
        "structure_file": "messagecatelog.json",
        "language_file_pattern": "messagecatelog-{lang}.json"
      }
    ]
  }
}
```

### **2. File Structure**

Create your catalog files:

```
pkg/
├── alert/
│   └── catalog/
│       ├── messagecatelog.json          # Structure definitions
│       ├── messagecatelog-en-US.json    # English translations
│       └── messagecatelog-fr-FR.json    # French translations
└── audit/
    └── catalog/
        ├── messagecatelog.json          # Structure definitions
        ├── messagecatelog-en-US.json    # English translations
        └── messagecatelog-fr-FR.json    # French translations
```

### **3. Structure File Format**

```json
{
  "ABC0001": {
    "message_code": "ABC0001",
    "category": "Registration",
    "severity": "CRITICAL",
    "component": "Reg"
  }
}
```

### **4. Language File Format**

```json
{
  "ABC0001": {
    "message": "Registration is about to expire",
    "detailed_description": "Registration is about to expire on {{.exp_date}}",
    "response_action": "Renew the registration before expiry"
  }
}
```

## 📖 **API Reference**

### **Service Interface**

```go
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
```

### **Message Request**

```go
type MessageRequest struct {
    MessageCode string                 `json:"message_code" validate:"required"`
    Language    string                 `json:"language,omitempty"`
    CatalogName string                 `json:"catalog_name" validate:"required"`
    Parameters  map[string]interface{} `json:"parameters,omitempty"`
}
```

### **Message Response**

```go
type MessageResponse struct {
    MessageCode         string                 `json:"message_code"`
    Category            string                 `json:"category"`
    Severity            string                 `json:"severity"`
    Component           string                 `json:"component"`
    Message             string                 `json:"message"`
    DetailedDescription string                 `json:"detailed_description"`
    ResponseAction      string                 `json:"response_action"`
    Language            string                 `json:"language"`
    CatalogName         string                 `json:"catalog_name"`
    FormattedMessage    string                 `json:"formatted_message"`
    Metadata            map[string]interface{} `json:"metadata,omitempty"`
    Timestamp           time.Time              `json:"timestamp"`
}
```

## 🔧 **Usage Examples**

### **Basic Usage**

```go
// Get a message
req := &messagecatalog.MessageRequest{
    MessageCode: "ABC0001",
    CatalogName: "alert",
    Language:    "en-US",
    Parameters: map[string]interface{}{
        "exp_date": "2024-12-31",
    },
}

message, err := messageCatalogService.GetMessage(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Message: %s\n", message.FormattedMessage)
```

### **Alert Service Integration**

```go
type AlertService struct {
    messageCatalog messagecatalog.Service
    logger         interfaces.Logger
}

func (s *AlertService) ProcessAlert(ctx context.Context, alertCode string, parameters map[string]interface{}) (*AlertResponse, error) {
    req := &messagecatalog.MessageRequest{
        MessageCode: alertCode,
        CatalogName: "alert",
        Language:    "en-US",
        Parameters:  parameters,
    }
    
    message, err := s.messageCatalog.GetMessage(ctx, req)
    if err != nil {
        return nil, err
    }
    
    return &AlertResponse{
        Code:        message.MessageCode,
        Category:    message.Category,
        Severity:    message.Severity,
        Message:     message.FormattedMessage,
        Description: message.DetailedDescription,
        Action:      message.ResponseAction,
    }, nil
}
```

### **Audit Service Integration**

```go
type AuditService struct {
    messageCatalog messagecatalog.Service
    logger         interfaces.Logger
}

func (s *AuditService) LogEvent(ctx context.Context, eventCode string, parameters map[string]interface{}) (*AuditEvent, error) {
    req := &messagecatalog.MessageRequest{
        MessageCode: eventCode,
        CatalogName: "audit",
        Language:    "en-US",
        Parameters:  parameters,
    }
    
    message, err := s.messageCatalog.GetMessage(ctx, req)
    if err != nil {
        return nil, err
    }
    
    return &AuditEvent{
        EventCode:     message.MessageCode,
        EventCategory: message.Category,
        RiskLevel:     message.Severity,
        Description:   message.FormattedMessage,
        Details:       message.DetailedDescription,
        Action:        message.ResponseAction,
    }, nil
}
```

## 🎯 **Features**

- **✅ Multi-Catalog Support**: Handle multiple message catalogs (Alert, Audit, etc.)
- **✅ Multi-Language Support**: Support for multiple languages with fallback
- **✅ Template Parameters**: Dynamic parameter substitution in messages
- **✅ High-Performance Caching**: In-memory caching with configurable TTL
- **✅ Thread-Safe**: Concurrent access support
- **✅ Configuration-Driven**: Easy to add new catalogs
- **✅ Health Checks**: Built-in health monitoring
- **✅ Statistics**: Catalog usage statistics and metrics

## 🔒 **Error Handling**

The service provides comprehensive error handling:

- **Validation Errors**: Invalid request parameters
- **File Not Found**: Missing catalog or language files
- **Parse Errors**: Invalid JSON format
- **Message Not Found**: Message code doesn't exist
- **Cache Errors**: Cache operation failures

## 📊 **Performance**

- **Caching**: High-performance in-memory caching
- **Thread-Safe**: Concurrent access support
- **Memory Efficient**: Optimized data structures
- **Fast Lookup**: O(1) cache lookups

## 🧪 **Testing**

The service includes comprehensive testing:

- **Unit Tests**: Individual component testing
- **Integration Tests**: End-to-end testing
- **Performance Tests**: Load and stress testing
- **Error Tests**: Error condition testing

## 📈 **Monitoring**

Built-in monitoring and metrics:

- **Health Checks**: Service availability
- **Statistics**: Usage metrics
- **Cache Metrics**: Hit/miss ratios
- **Performance Metrics**: Response times

## 🔧 **Configuration Options**

| Option | Description | Default |
|--------|-------------|---------|
| `default_language` | Default language for messages | `en-US` |
| `supported_languages` | List of supported languages | `["en-US", "fr-FR"]` |
| `cache_enabled` | Enable caching | `true` |
| `cache_ttl_seconds` | Cache TTL in seconds | `3600` |
| `reload_interval_seconds` | Reload interval in seconds | `300` |

## 📚 **Documentation**

- [Design Document](DESIGN_DOCUMENT.md) - Comprehensive design documentation
- [Class Diagram](diagrams/class_diagram.mmd) - Class relationships
- [Sequence Diagram](diagrams/sequence_diagram.mmd) - Interaction flow
- [Flow Diagram](diagrams/flow_diagram.mmd) - Process flow
- [Data Flow](diagrams/data_flow_diagram.mmd) - Data flow diagram
- [Cache Architecture](diagrams/cache_architecture.mmd) - Caching strategy
