# Alert Service Design Document

## Overview

The Alert Service is a comprehensive system designed to manage and process alerts with multi-language support through a Message Catalog Service. The system follows the existing project's architectural patterns and provides a generic, customizable solution for alert management.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Design Principles](#design-principles)
3. [Component Architecture](#component-architecture)
4. [Data Flow](#data-flow)
5. [Class Diagrams](#class-diagrams)
6. [Sequence Diagrams](#sequence-diagrams)
7. [Configuration](#configuration)
8. [Usage Examples](#usage-examples)
9. [Error Handling](#error-handling)
10. [Performance Considerations](#performance-considerations)

## Architecture Overview

The Alert Service consists of two main components:

### 1. Message Catalog Service
- **Purpose**: Manages multi-language message templates and metadata
- **Features**: 
  - Language-specific message loading
  - Caching for performance
  - Hot-reloading of message files
  - Generic design for different JSON schemas

### 2. Alert Service
- **Purpose**: Processes alerts and integrates with the Message Catalog
- **Features**:
  - Alert creation, retrieval, and management
  - Message formatting with parameters
  - Multi-language support
  - Alert statistics and filtering

## Design Principles

### 1. **Separation of Concerns**
- Message catalog separated from alert processing
- Clear interfaces between components
- Single responsibility for each service

### 2. **Interface Segregation**
- Separate interfaces for catalog and alert operations
- Clean abstraction layers
- Easy mocking for testing

### 3. **Dependency Injection**
- Services injected through constructors
- Loose coupling between components
- Easy testing and maintenance

### 4. **Generic Design**
- Catalog service designed to handle different JSON schemas
- Configurable message loading strategies
- Extensible for future requirements

### 5. **Error Handling**
- Consistent error handling using AppError pattern
- Detailed error context and logging
- Graceful degradation

## Component Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Alert Service Layer                      │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐    ┌─────────────────────────────────┐ │
│  │   AlertService  │    │     MessageCatalogService      │ │
│  │                 │    │                                 │ │
│  │ - CreateAlert   │    │ - GetMessage                    │ │
│  │ - ProcessAlert  │    │ - GetMessagesByCategory         │ │
│  │ - ListAlerts    │    │ - ReloadCatalog                 │ │
│  │ - GetAlertStats │    │ - ListAvailableLanguages       │ │
│  └─────────────────┘    └─────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                    Message Catalog Layer                    │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐    ┌─────────────────────────────────┐ │
│  │ FileMessageLoader│   │         Message Cache           │ │
│  │                 │    │                                 │ │
│  │ - LoadMessages  │    │ - Language-based caching        │ │
│  │ - LoadAllLangs  │    │ - Thread-safe operations        │ │
│  │ - GetAvailable  │    │ - Configurable cache settings   │ │
│  └─────────────────┘    └─────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                    File System Layer                        │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐    ┌─────────────────────────────────┐ │
│  │   en.json       │    │         fr.json                 │ │
│  │                 │    │                                 │ │
│  │ - English msgs  │    │ - French messages               │ │
│  │ - Message codes │    │ - Localized content             │ │
│  │ - Metadata      │    │ - Consistent structure          │ │
│  └─────────────────┘    └─────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## Data Flow

### 1. Alert Processing Flow
```
Client Request → Alert Service → Message Catalog → File System
     ↓              ↓              ↓              ↓
Process Alert → Get Message → Load JSON → Return Message
     ↓              ↓              ↓              ↓
Format Message ← Apply Params ← Parse JSON ← Read File
     ↓
Return Alert Response
```

### 2. Message Loading Flow
```
Service Request → Check Cache → Load from File → Parse JSON
     ↓              ↓              ↓              ↓
Return Message ← Update Cache ← Validate Data ← Extract Fields
```

## Configuration

### Message Catalog Configuration
```json
{
  "catalog": {
    "default_language": "en",
    "supported_languages": ["en", "fr"],
    "catalog_path": "./pkg/alert/catalog/message_catalog",
    "reload_interval_seconds": 300,
    "cache_enabled": true
  }
}
```

### Alert Service Configuration
```json
{
  "alert": {
    "default_severity": "medium",
    "auto_expire_hours": 24,
    "max_alerts_per_category": 1000,
    "cleanup_interval_hours": 1
  }
}
```

## Usage Examples

### Basic Alert Processing
```go
// Process an alert with message code
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

### Multi-language Support
```go
// Process alert in French
processReqFr := &service.ProcessAlertRequest{
    MessageCode: "ABC001",
    Language:    "fr",
    Parameters: map[string]string{
        "s": "LICENCE-PROD-001",
        "d": "2024-12-31",
    },
}

alertResponseFr, err := alertService.ProcessAlert(ctx, processReqFr)
```

## Error Handling

The Alert Service uses the existing AppError pattern for consistent error handling:

- **Validation Errors**: Invalid message codes, missing parameters
- **File System Errors**: Missing language files, JSON parsing errors
- **Service Errors**: Catalog loading failures, alert processing errors
- **Configuration Errors**: Invalid configuration values

## Performance Considerations

### Caching Strategy
- Language-based message caching
- Configurable cache settings
- Thread-safe cache operations
- Automatic cache invalidation

### File Loading
- Lazy loading of language files
- Batch loading for multiple languages
- Error handling for missing files
- Hot-reloading support

### Memory Management
- Efficient data structures
- Garbage collection friendly
- Configurable cache limits
- Memory leak prevention

## Future Enhancements

1. **Database Integration**: Store alerts in database instead of memory
2. **Message Templates**: Support for more complex template engines
3. **Notification Channels**: Email, SMS, webhook integrations
4. **Alert Rules**: Configurable alert rules and conditions
5. **Analytics**: Advanced alert analytics and reporting
6. **API Endpoints**: REST API for external integrations
