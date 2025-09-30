# Alert Service Implementation Guide

## Quick Start

### 1. Directory Structure Setup
```
pkg/alert/
├── catalog/
│   ├── interfaces.go
│   ├── models.go
│   ├── service.go
│   └── message_catalog/
│       ├── en.json
│       └── fr.json
├── service/
│   ├── interfaces.go
│   ├── models.go
│   └── service.go
└── config/
    └── alert_config.json
```

### 2. Basic Usage

```go
// Initialize logger
logger, _ := logger.NewLogger(&logger.Config{
    Level:  "info",
    Format: "console",
    Output: "stdout",
})

// Initialize catalog configuration
catalogConfig := &catalog.CatalogConfig{
    DefaultLanguage:    "en",
    SupportedLanguages: []string{"en", "fr"},
    CatalogPath:        "./pkg/alert/catalog/message_catalog",
    ReloadInterval:     300,
    CacheEnabled:       true,
}

// Initialize services
messageCatalog := catalog.NewMessageCatalogService(catalogConfig, logger)
alertService := service.NewAlertService(messageCatalog, logger)

// Process an alert
processReq := &service.ProcessAlertRequest{
    MessageCode: "ABC001",
    Language:    "en",
    Parameters: map[string]string{
        "s": "PROD-LICENSE-001",
        "d": "2024-12-31",
    },
}

alertResponse, err := alertService.ProcessAlert(ctx, processReq)
```

### 3. Configuration

The Alert Service supports extensive configuration through JSON files:

```json
{
  "catalog": {
    "default_language": "en",
    "supported_languages": ["en", "fr", "es"],
    "catalog_path": "./pkg/alert/catalog/message_catalog",
    "reload_interval_seconds": 300,
    "cache_enabled": true
  },
  "alert": {
    "default_severity": "medium",
    "auto_expire_hours": 24,
    "max_alerts_per_category": 1000,
    "cleanup_interval_hours": 1
  }
}
```

## Advanced Features

### Custom Message Loaders

You can implement custom message loaders for different data sources:

```go
type DatabaseMessageLoader struct {
    db     interfaces.Database
    logger interfaces.Logger
}

func (l *DatabaseMessageLoader) LoadMessages(ctx context.Context, language string) (map[string]*Message, error) {
    // Implementation for database loading
}
```

### Custom Alert Storage

Implement custom alert storage backends:

```go
type DatabaseAlertStorage struct {
    db     interfaces.Database
    logger interfaces.Logger
}

func (s *DatabaseAlertStorage) StoreAlert(ctx context.Context, alert *Alert) error {
    // Implementation for database storage
}
```

## Testing

### Unit Tests

```go
func TestAlertService_ProcessAlert(t *testing.T) {
    // Setup
    mockCatalog := &MockMessageCatalog{}
    logger := &MockLogger{}
    service := NewAlertService(mockCatalog, logger)
    
    // Test
    req := &ProcessAlertRequest{
        MessageCode: "ABC001",
        Language:    "en",
        Parameters:  map[string]string{"s": "test"},
    }
    
    result, err := service.ProcessAlert(context.Background(), req)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "ABC001", result.MessageCode)
}
```

### Integration Tests

```go
func TestAlertService_Integration(t *testing.T) {
    // Setup real services
    config := &CatalogConfig{
        DefaultLanguage: "en",
        CatalogPath:     "./testdata/messages",
        CacheEnabled:    true,
    }
    
    logger, _ := logger.NewLogger(&logger.Config{Level: "debug"})
    catalog := NewMessageCatalogService(config, logger)
    service := NewAlertService(catalog, logger)
    
    // Test end-to-end flow
    req := &ProcessAlertRequest{
        MessageCode: "ABC001",
        Language:    "en",
        Parameters:  map[string]string{"s": "test-license"},
    }
    
    result, err := service.ProcessAlert(context.Background(), req)
    
    assert.NoError(t, err)
    assert.Contains(t, result.FormattedMessage, "test-license")
}
```

## Performance Optimization

### Caching Strategy

- Enable caching for production environments
- Set appropriate cache expiration times
- Monitor cache hit rates
- Implement cache warming strategies

### Memory Management

- Use appropriate data structures
- Implement memory limits
- Monitor memory usage
- Implement garbage collection strategies

### File I/O Optimization

- Use efficient file reading methods
- Implement file watching for hot reloading
- Use appropriate buffer sizes
- Implement connection pooling for database operations

## Monitoring and Observability

### Logging

The Alert Service integrates with the existing logging system:

```go
// Structured logging with context
logger.Info(ctx, "Alert processed successfully", interfaces.Fields{
    "alert_id":     alertResponse.AlertID,
    "message_code": req.MessageCode,
    "language":     req.Language,
    "severity":     alertResponse.Severity,
})
```

### Metrics

Implement custom metrics for monitoring:

```go
// Alert processing metrics
alertProcessingDuration := prometheus.NewHistogramVec(
    prometheus.HistogramOpts{
        Name: "alert_processing_duration_seconds",
        Help: "Time taken to process alerts",
    },
    []string{"message_code", "language", "severity"},
)

// Cache metrics
cacheHitRate := prometheus.NewGaugeVec(
    prometheus.GaugeOpts{
        Name: "message_catalog_cache_hit_rate",
        Help: "Cache hit rate for message catalog",
    },
    []string{"language"},
)
```

## Troubleshooting

### Common Issues

1. **Message not found**: Check message code and language
2. **File loading errors**: Verify file paths and permissions
3. **JSON parsing errors**: Validate JSON syntax
4. **Cache issues**: Check cache configuration and memory usage

### Debug Mode

Enable debug logging for troubleshooting:

```go
config := &logger.Config{
    Level:  "debug",
    Format: "console",
    Output: "stdout",
}
```

### Health Checks

Implement health checks for the Alert Service:

```go
func (s *AlertServiceImpl) HealthCheck(ctx context.Context) error {
    // Check catalog availability
    languages, err := s.catalog.ListAvailableLanguages(ctx)
    if err != nil {
        return fmt.Errorf("catalog health check failed: %w", err)
    }
    
    if len(languages) == 0 {
        return fmt.Errorf("no languages available in catalog")
    }
    
    return nil
}
```

## Security Considerations

### Input Validation

- Validate all input parameters
- Sanitize user-provided data
- Implement rate limiting
- Use parameterized queries for database operations

### Access Control

- Implement proper authentication
- Use authorization for sensitive operations
- Log all access attempts
- Implement audit trails

### Data Protection

- Encrypt sensitive data
- Use secure file permissions
- Implement data retention policies
- Regular security audits

## API Reference

### AlertService Interface

```go
type AlertService interface {
    // CreateAlert creates a new alert
    CreateAlert(ctx context.Context, req *CreateAlertRequest) (*Alert, error)
    
    // GetAlert retrieves an alert by ID
    GetAlert(ctx context.Context, id string) (*Alert, error)
    
    // UpdateAlert updates an existing alert
    UpdateAlert(ctx context.Context, id string, req *UpdateAlertRequest) (*Alert, error)
    
    // DeleteAlert deletes an alert
    DeleteAlert(ctx context.Context, id string) error
    
    // ListAlerts retrieves alerts with filtering and pagination
    ListAlerts(ctx context.Context, req *ListAlertsRequest) (*ListAlertsResponse, error)
    
    // GetAlertByMessageCode retrieves alert information by message code
    GetAlertByMessageCode(ctx context.Context, messageCode, language string) (*AlertResponse, error)
    
    // ProcessAlert processes an alert and returns formatted message
    ProcessAlert(ctx context.Context, req *ProcessAlertRequest) (*AlertResponse, error)
    
    // GetAlertStats returns statistics about alerts
    GetAlertStats(ctx context.Context) (*AlertStats, error)
}
```

### MessageCatalog Interface

```go
type MessageCatalog interface {
    // GetMessage retrieves a message by code and language
    GetMessage(ctx context.Context, messageCode, language string) (*Message, error)
    
    // GetMessageByCode retrieves a message by code using default language
    GetMessageByCode(ctx context.Context, messageCode string) (*Message, error)
    
    // GetMessagesByCategory retrieves all messages for a specific category and language
    GetMessagesByCategory(ctx context.Context, category, language string) ([]*Message, error)
    
    // GetMessagesBySeverity retrieves all messages for a specific severity and language
    GetMessagesBySeverity(ctx context.Context, severity, language string) ([]*Message, error)
    
    // ListAvailableLanguages returns all available languages
    ListAvailableLanguages(ctx context.Context) ([]string, error)
    
    // ReloadCatalog reloads the message catalog from files
    ReloadCatalog(ctx context.Context) error
}
```

## Examples

### Complete Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "tushartemplategin/pkg/alert/catalog"
    "tushartemplategin/pkg/alert/service"
    "tushartemplategin/pkg/logger"
)

func main() {
    // Initialize logger
    loggerInstance, err := logger.NewLogger(&logger.Config{
        Level:  "info",
        Format: "console",
        Output: "stdout",
    })
    if err != nil {
        log.Fatal("Failed to initialize logger:", err)
    }

    // Initialize catalog configuration
    catalogConfig := &catalog.CatalogConfig{
        DefaultLanguage:    "en",
        SupportedLanguages: []string{"en", "fr"},
        CatalogPath:        "./pkg/alert/catalog/message_catalog",
        ReloadInterval:     300,
        CacheEnabled:       true,
    }

    // Initialize message catalog
    messageCatalog := catalog.NewMessageCatalogService(catalogConfig, loggerInstance)

    // Initialize alert service
    alertService := service.NewAlertService(messageCatalog, loggerInstance)

    ctx := context.Background()

    // Example 1: Process an alert with message code
    fmt.Println("=== Example 1: Process Alert by Message Code ===")
    processReq := &service.ProcessAlertRequest{
        MessageCode: "ABC001",
        Language:    "en",
        Parameters: map[string]string{
            "s": "PROD-LICENSE-001",
            "d": "2024-12-31",
        },
        Metadata: map[string]string{
            "user_id": "12345",
            "source":  "license_monitor",
        },
    }

    alertResponse, err := alertService.ProcessAlert(ctx, processReq)
    if err != nil {
        log.Printf("Error processing alert: %v", err)
    } else {
        fmt.Printf("Alert ID: %s\n", alertResponse.AlertID)
        fmt.Printf("Message: %s\n", alertResponse.Message)
        fmt.Printf("Formatted Message: %s\n", alertResponse.FormattedMessage)
        fmt.Printf("Response Action: %s\n", alertResponse.ResponseAction)
        fmt.Printf("Severity: %s\n", alertResponse.Severity)
        fmt.Printf("Category: %s\n", alertResponse.Category)
    }

    // Example 2: Process alert in French
    fmt.Println("\n=== Example 2: Process Alert in French ===")
    processReqFr := &service.ProcessAlertRequest{
        MessageCode: "ABC001",
        Language:    "fr",
        Parameters: map[string]string{
            "s": "LICENCE-PROD-001",
            "d": "2024-12-31",
        },
    }

    alertResponseFr, err := alertService.ProcessAlert(ctx, processReqFr)
    if err != nil {
        log.Printf("Error processing alert: %v", err)
    } else {
        fmt.Printf("Alert ID: %s\n", alertResponseFr.AlertID)
        fmt.Printf("Message: %s\n", alertResponseFr.Message)
        fmt.Printf("Formatted Message: %s\n", alertResponseFr.FormattedMessage)
        fmt.Printf("Response Action: %s\n", alertResponseFr.ResponseAction)
    }

    // Example 3: Create a new alert
    fmt.Println("\n=== Example 3: Create New Alert ===")
    createReq := &service.CreateAlertRequest{
        MessageCode: "ABC004",
        Title:       "Custom Security Alert",
        Description: "Custom security alert description",
        Parameters: map[string]string{
            "s": "user-authentication",
        },
        Metadata: map[string]string{
            "alert_type": "custom",
            "created_by": "admin",
        },
    }

    alert, err := alertService.CreateAlert(ctx, createReq)
    if err != nil {
        log.Printf("Error creating alert: %v", err)
    } else {
        fmt.Printf("Created Alert ID: %s\n", alert.ID)
        fmt.Printf("Alert Title: %s\n", alert.Title)
        fmt.Printf("Alert Severity: %s\n", alert.Severity)
        fmt.Printf("Alert Category: %s\n", alert.Category)
    }

    // Example 4: List alerts
    fmt.Println("\n=== Example 4: List Alerts ===")
    listReq := &service.ListAlertsRequest{
        Page:   1,
        Limit:  10,
        Status: "active",
    }

    listResponse, err := alertService.ListAlerts(ctx, listReq)
    if err != nil {
        log.Printf("Error listing alerts: %v", err)
    } else {
        fmt.Printf("Total Alerts: %d\n", listResponse.Total)
        fmt.Printf("Current Page: %d\n", listResponse.Page)
        fmt.Printf("Alerts in this page: %d\n", len(listResponse.Alerts))
    }

    // Example 5: Get alert statistics
    fmt.Println("\n=== Example 5: Alert Statistics ===")
    stats, err := alertService.GetAlertStats(ctx)
    if err != nil {
        log.Printf("Error getting alert stats: %v", err)
    } else {
        fmt.Printf("Total Alerts: %d\n", stats.TotalAlerts)
        fmt.Printf("Active Alerts: %d\n", stats.ActiveAlerts)
        fmt.Printf("Processed Alerts: %d\n", stats.ProcessedAlerts)
        fmt.Printf("Expired Alerts: %d\n", stats.ExpiredAlerts)
    }
}
```

This implementation guide provides comprehensive documentation for implementing and using the Alert Service system.
