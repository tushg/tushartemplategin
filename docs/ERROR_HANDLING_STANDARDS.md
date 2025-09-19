# Error Handling Standards

## Table of Contents
1. [Error Handling Strategy](#error-handling-strategy)
2. [Layer-Specific Error Handling](#layer-specific-error-handling)
3. [Error Creation Patterns](#error-creation-patterns)
4. [Error Response Format](#error-response-format)
5. [Error Logging Standards](#error-logging-standards)
6. [Common Error Scenarios](#common-error-scenarios)
7. [Error Handling Checklist](#error-handling-checklist)

---

## Error Handling Strategy

### Why Custom Errors?
**Justification**: Custom errors provide structured, consistent error responses across all microservices, making debugging easier and providing better user experience with meaningful error messages and proper HTTP status codes.

### Hybrid Approach
**Justification**: Using custom errors in business layers and standard errors in infrastructure layers follows Go idioms while providing rich error information where it matters most - in the business logic.

---

## Layer-Specific Error Handling

### 1. Infrastructure Layer (Database, External APIs)
**Pattern**: Return standard `error`
**Justification**: Infrastructure errors are often low-level and should be wrapped by business layers to provide context.

```go
// ✅ CORRECT - Infrastructure layer
type Database interface {
    Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
    WithTransaction(ctx context.Context, fn func(*sql.Tx) error) error
}

func (r *ProductRepository) Create(ctx context.Context, product *ProductRegistration) (*ProductRegistration, *errors.AppError) {
    if err := r.db.WithTransaction(ctx, func(tx *sql.Tx) error {
        return tx.QueryRowContext(ctx, query, ...).Scan(...)
    }); err != nil {
        // Wrap infrastructure error with business context
        return nil, errors.NewWithError(errors.ErrCodeDatabaseQuery, "Failed to create product", 
            http.StatusInternalServerError, err).WithField("operation", "create product")
    }
    return product, nil
}
```

### 2. Business Layer (Service, Repository)
**Pattern**: Return `*errors.AppError`
**Justification**: Business logic errors need rich context, proper HTTP status codes, and structured information for API responses.

```go
// ✅ CORRECT - Business layer
type Service interface {
    CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductRegistration, *errors.AppError)
    GetProduct(ctx context.Context, id int64) (*ProductRegistration, *errors.AppError)
}

func (s *ProductService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductRegistration, *errors.AppError) {
    // Business validation
    exists, appErr := s.repo.SKUExists(ctx, req.SKU, nil)
    if appErr != nil {
        return nil, appErr
    }
    
    if exists {
        return nil, errors.NewWithDetails(errors.ErrCodeProductSKUExists, "Product SKU already exists", 
            fmt.Sprintf("Product with SKU '%s' already exists", req.SKU), http.StatusConflict).WithField("sku", req.SKU)
    }
    
    return s.repo.Create(ctx, product)
}
```

### 3. Presentation Layer (Handlers, Middleware)
**Pattern**: Use `middleware.HandleAppError`
**Justification**: Centralized error handling ensures consistent error responses and proper logging across all endpoints.

```go
// ✅ CORRECT - Handler layer
func createProductHandler(c *gin.Context) {
    var req CreateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        middleware.HandleAppError(c, errors.NewWithDetails(errors.ErrCodeBadRequest, "Invalid request data", 
            err.Error(), http.StatusBadRequest))
        return
    }
    
    product, appErr := productService.CreateProduct(ctx, &req)
    if appErr != nil {
        middleware.HandleAppError(c, appErr)
        return
    }
    
    c.JSON(http.StatusCreated, ProductResponse{Product: *product})
}
```

---

## Error Creation Patterns

### 1. Basic Error Creation
**Pattern**: `errors.New(code, message, httpStatus)`
**Justification**: Simple errors without additional context or details.

```go
// ✅ CORRECT
return nil, errors.New(errors.ErrCodeNotFound, "Product not found", http.StatusNotFound)

// ❌ WRONG - Don't use for errors that need context
return nil, errors.New(errors.ErrCodeProductNotFound, "Product not found", http.StatusNotFound)
```

### 2. Error with Details
**Pattern**: `errors.NewWithDetails(code, message, details, httpStatus)`
**Justification**: When you need to provide additional context or explanation to the error.

```go
// ✅ CORRECT
return nil, errors.NewWithDetails(errors.ErrCodeProductSKUExists, "Product SKU already exists", 
    fmt.Sprintf("Product with SKU '%s' already exists", req.SKU), http.StatusConflict)

// ❌ WRONG - Don't use for simple errors
return nil, errors.NewWithDetails(errors.ErrCodeNotFound, "Product not found", "Product not found", http.StatusNotFound)
```

### 3. Error with Wrapped Error
**Pattern**: `errors.NewWithError(code, message, httpStatus, err)`
**Justification**: When wrapping infrastructure errors or other Go errors to provide business context.

```go
// ✅ CORRECT
return nil, errors.NewWithError(errors.ErrCodeDatabaseQuery, "Failed to create product", 
    http.StatusInternalServerError, err).WithField("operation", "create product")

// ❌ WRONG - Don't use for business logic errors
return nil, errors.NewWithError(errors.ErrCodeProductNotFound, "Product not found", 
    http.StatusNotFound, nil)
```

### 4. Error with Fields
**Pattern**: `errors.New(...).WithField(key, value)`
**Justification**: When you need to provide structured data that can be used for debugging or API responses.

```go
// ✅ CORRECT
return nil, errors.NewWithDetails(errors.ErrCodeProductSKUExists, "Product SKU already exists", 
    fmt.Sprintf("Product with SKU '%s' already exists", req.SKU), http.StatusConflict).WithField("sku", req.SKU)

// ❌ WRONG - Don't use for simple errors
return nil, errors.New(errors.ErrCodeNotFound, "Product not found", http.StatusNotFound).WithField("id", id)
```

---

## Error Response Format

### Standard Error Response
**Justification**: Consistent error format across all microservices makes debugging easier and provides better developer experience.

```json
{
    "error": "PRODUCT_NOT_FOUND",
    "message": "Product not found",
    "details": "Product with ID 123 not found",
    "fields": {
        "product_id": 123,
        "operation": "get_product"
    }
}
```

### HTTP Status Code Mapping
**Justification**: Proper HTTP status codes help clients understand the nature of the error and handle it appropriately.

```go
// Business Logic Errors
ErrCodeProductNotFound     → 404 Not Found
ErrCodeProductSKUExists    → 409 Conflict
ErrCodeInvalidStock        → 422 Unprocessable Entity

// System Errors
ErrCodeInternalServer      → 500 Internal Server Error
ErrCodeDatabaseConnection  → 503 Service Unavailable
ErrCodeServiceUnavailable  → 503 Service Unavailable

// Client Errors
ErrCodeBadRequest          → 400 Bad Request
ErrCodeUnauthorized        → 401 Unauthorized
ErrCodeForbidden           → 403 Forbidden
```

---

## Error Logging Standards

### Log Level Guidelines
**Justification**: Proper log levels help with monitoring and debugging in production environments.

```go
// ✅ CORRECT - Log levels
s.logger.Error(ctx, "Database operation failed", interfaces.Fields{
    "error":     err.Error(),
    "operation": "create_product",
    "user_id":   userID,
})

s.logger.Warn(ctx, "Business rule violation", interfaces.Fields{
    "error":     appErr.Error(),
    "operation": "create_product",
    "sku":       req.SKU,
})

s.logger.Info(ctx, "Operation completed successfully", interfaces.Fields{
    "operation": "create_product",
    "product_id": product.ID,
})
```

### Error Context
**Justification**: Rich context in error logs makes debugging faster and more effective.

```go
// ✅ CORRECT - Rich context
s.logger.Error(ctx, "Failed to create product", interfaces.Fields{
    "error":        err.Error(),
    "operation":    "create_product",
    "user_id":      userID,
    "request_data": map[string]interface{}{
        "name": req.Name,
        "sku":  req.SKU,
    },
    "timestamp":    time.Now().UTC(),
})

// ❌ WRONG - Minimal context
s.logger.Error(ctx, "Failed to create product", interfaces.Fields{
    "error": err.Error(),
})
```

---

## Common Error Scenarios

### 1. Validation Errors
**Pattern**: Use `ErrCodeBadRequest` with validation details
**Justification**: Validation errors are client errors that should return 400 status.

```go
// ✅ CORRECT
if req.Name == "" {
    return nil, errors.NewWithDetails(errors.ErrCodeBadRequest, "Validation failed", 
        "Product name is required", http.StatusBadRequest).WithField("field", "name")
}
```

### 2. Business Rule Violations
**Pattern**: Use specific business error codes
**Justification**: Business rule violations need specific error codes for proper handling.

```go
// ✅ CORRECT
if exists {
    return nil, errors.NewWithDetails(errors.ErrCodeProductSKUExists, "Product SKU already exists", 
        fmt.Sprintf("Product with SKU '%s' already exists", req.SKU), http.StatusConflict).WithField("sku", req.SKU)
}
```

### 3. Resource Not Found
**Pattern**: Use `ErrCodeNotFound` with resource details
**Justification**: Not found errors need to specify what resource was not found.

```go
// ✅ CORRECT
if product == nil {
    return nil, errors.NewWithDetails(errors.ErrCodeNotFound, "Product not found", 
        fmt.Sprintf("Product with ID %d not found", id), http.StatusNotFound).WithField("product_id", id)
}
```

### 4. Database Errors
**Pattern**: Wrap with business context
**Justification**: Database errors need business context to be meaningful.

```go
// ✅ CORRECT
if err := r.db.WithTransaction(ctx, func(tx *sql.Tx) error {
    return tx.QueryRowContext(ctx, query, ...).Scan(...)
}); err != nil {
    return nil, errors.NewWithError(errors.ErrCodeDatabaseQuery, "Failed to create product", 
        http.StatusInternalServerError, err).WithField("operation", "create product")
}
```

---

## Error Handling Checklist

### ✅ Pre-Submit Checklist
- [ ] **Error Types**: Business layer returns `*errors.AppError`, infrastructure returns `error`
- [ ] **Error Codes**: Use appropriate error codes for different scenarios
- [ ] **HTTP Status**: Map error codes to correct HTTP status codes
- [ ] **Error Context**: Include relevant context in error messages
- [ ] **Error Fields**: Add structured data for debugging
- [ ] **Error Logging**: Log errors with appropriate level and context
- [ ] **Error Wrapping**: Wrap infrastructure errors with business context
- [ ] **Error Handling**: Use `middleware.HandleAppError` in handlers
- [ ] **Error Testing**: Test error scenarios in unit tests
- [ ] **Error Documentation**: Document error codes and their meanings

### ✅ Error Code Guidelines
- [ ] **Business Errors**: Use specific codes (e.g., `ErrCodeProductNotFound`)
- [ ] **System Errors**: Use generic codes (e.g., `ErrCodeInternalServer`)
- [ ] **Client Errors**: Use appropriate HTTP error codes
- [ ] **Error Messages**: Use clear, user-friendly messages
- [ ] **Error Details**: Provide additional context when needed
- [ ] **Error Fields**: Include structured data for debugging

### ✅ Error Response Guidelines
- [ ] **Consistent Format**: Use standard error response format
- [ ] **Proper Status Codes**: Map error codes to HTTP status codes
- [ ] **Error Context**: Include relevant context in responses
- [ ] **Error Fields**: Add structured data for API consumers
- [ ] **Error Logging**: Log errors with appropriate level and context

---

## Quick Reference

### Error Creation
```go
// Simple error
errors.New(code, message, httpStatus)

// Error with details
errors.NewWithDetails(code, message, details, httpStatus)

// Error with wrapped error
errors.NewWithError(code, message, httpStatus, err)

// Error with fields
errors.New(...).WithField(key, value)
```

### Error Handling
```go
// Business layer - return custom errors
func (s *Service) Method() (*Model, *errors.AppError)

// Infrastructure layer - return standard errors
func (r *Repository) Method() error

// Handler layer - use middleware
middleware.HandleAppError(c, appErr)
```

### Error Logging
```go
// Error level
s.logger.Error(ctx, "Operation failed", interfaces.Fields{
    "error":     err.Error(),
    "operation": "operation_name",
    "context":   "additional_context",
})

// Warning level
s.logger.Warn(ctx, "Business rule violation", interfaces.Fields{
    "error":     appErr.Error(),
    "operation": "operation_name",
    "context":   "additional_context",
})
```

This error handling standard ensures consistent, maintainable, and debuggable error handling across all microservices in the project.
