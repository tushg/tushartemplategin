# Message Catalog Service Design Document

## Overview

The Message Catalog Service is a comprehensive system designed to manage multi-language message templates and metadata. It provides a generic, customizable solution for message management that can be easily extended for different use cases and message formats.

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

The Message Catalog Service consists of the following main components:

### 1. Service Layer
- **Purpose**: Provides business logic and orchestration
- **Features**: 
  - Message retrieval by code and language
  - Category and severity-based filtering
  - Multi-language support
  - Health monitoring

### 2. File System Layer
- **Purpose**: Manages message data storage
- **Features**:
  - JSON file-based storage
  - Language-specific files
  - Hot-reloading support
  - Error handling for missing files

### 3. Caching Layer
- **Purpose**: Performance optimization
- **Features**:
  - In-memory caching
  - Thread-safe operations
  - Configurable cache settings
  - Automatic cache invalidation

## Design Principles

### 1. **Separation of Concerns**
- Service logic separated from file I/O
- Clear interfaces between components
- Single responsibility for each component

### 2. **Interface Segregation**
- Clean service interfaces
- Easy mocking for testing
- Minimal dependencies

### 3. **Dependency Injection**
- Services injected through constructors
- Loose coupling between components
- Easy testing and maintenance

### 4. **Generic Design**
- Configurable message loading strategies
- Extensible for future requirements
- Support for different JSON schemas

### 5. **Error Handling**
- Consistent error handling using AppError pattern
- Detailed error context and logging
- Graceful degradation

## Component Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Message Catalog Service                  │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────────┐ │
│  │              MessageCatalogService                     │ │
│  │                                                         │ │
│  │ - GetMessage(code, lang)                               │ │
│  │ - GetMessageByCode(code)                               │ │
│  │ - GetMessagesByCategory(cat, lang)                     │ │
│  │ - GetMessagesBySeverity(sev, lang)                     │ │
│  │ - ListAvailableLanguages()                             │ │
│  │ - ReloadCatalog()                                      │ │
│  │ - HealthCheck()                                        │ │
│  └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                    Caching Layer                           │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                  Message Cache                         │ │
│  │                                                         │ │
│  │ - Language-based caching                               │ │
│  │ - Thread-safe operations                               │ │
│  │ - Configurable cache settings                          │ │
│  │ - Automatic cache invalidation                         │ │
│  └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                    File System Layer                       │
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

### 1. Message Retrieval Flow
```
Client Request → MessageCatalogService → Check Cache → Load from File
     ↓              ↓                    ↓              ↓
Return Message ← Parse JSON ← Read File ← File System
```

### 2. Cache Management Flow
```
Service Request → Check Cache → Cache Hit/Miss → Update Cache
     ↓              ↓              ↓              ↓
Return Message ← Return Cached ← Load from File ← Parse JSON
```

### 3. Catalog Reload Flow
```
Reload Request → Scan Directory → Load All Files → Update Cache
     ↓              ↓              ↓              ↓
Confirm Reload ← Parse All JSON ← Read All Files ← File System
```

## Configuration

### Message Catalog Configuration
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

## Usage Examples

### Basic Message Retrieval
```go
// Get message by code (default language)
message, err := messageCatalogService.GetMessageByCode(ctx, "ABC001")

// Get message in specific language
message, err := messageCatalogService.GetMessage(ctx, "ABC001", "fr")

// Get messages by category
messages, err := messageCatalogService.GetMessagesByCategory(ctx, "registration", "en")

// Get messages by severity
messages, err := messageCatalogService.GetMessagesBySeverity(ctx, "critical", "en")
```

### Health Monitoring
```go
// Health check
err := messageCatalogService.HealthCheck(ctx)

// Reload catalog
err := messageCatalogService.ReloadCatalog(ctx)

// List available languages
languages, err := messageCatalogService.ListAvailableLanguages(ctx)
```

## Error Handling

The Message Catalog Service uses the existing AppError pattern for consistent error handling:

- **File Not Found**: Missing language files
- **JSON Parse Errors**: Invalid JSON syntax
- **Message Not Found**: Invalid message codes
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

1. **Database Integration**: Store messages in database
2. **Message Templates**: Support for more complex template engines
3. **Version Control**: Message versioning and history
4. **Analytics**: Message usage analytics and reporting
5. **API Endpoints**: REST API for external integrations
6. **Message Validation**: Schema validation for message files
