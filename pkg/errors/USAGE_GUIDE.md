# Minimal Error Handling Usage Guide

## Core Methods (Only 5!)

### 1. `New(code, message, httpStatus)` - Basic Error
```go
// Simple errors
err := errors.New(errors.ErrCodeBadRequest, "Invalid request", http.StatusBadRequest)
err := errors.New(errors.ErrCodeNotFound, "Resource not found", http.StatusNotFound)
err := errors.New(errors.ErrCodeInternalServer, "Internal server error", http.StatusInternalServerError)
```

### 2. `NewWithDetails(code, message, details, httpStatus)` - Error with Details
```go
// Errors with additional context
err := errors.NewWithDetails(errors.ErrCodeBadRequest, "Invalid request", "Missing required field 'name'", http.StatusBadRequest)
err := errors.NewWithDetails(errors.ErrCodeProductNotFound, "Product not found", "Product with ID 123 not found", http.StatusNotFound)
```

### 3. `NewWithError(code, message, httpStatus, err)` - Wrap Existing Errors
```go
// Wrap database errors
err := errors.NewWithError(errors.ErrCodeDatabaseQuery, "Database operation failed", http.StatusInternalServerError, dbErr)
err := errors.NewWithError(errors.ErrCodeInternalServer, "Failed to create product", http.StatusInternalServerError, serviceErr)
```

### 4. `WithField(key, value)` - Add Context (Chainable)
```go
// Add context to any error
err := errors.New(errors.ErrCodeProductNotFound, "Product not found", http.StatusNotFound).WithField("product_id", 123)
err := errors.NewWithDetails(errors.ErrCodeProductSKUExists, "SKU exists", "Product with SKU 'ABC-123' already exists", http.StatusConflict).WithField("sku", "ABC-123")
```

### 5. `IsAppError(err)` & `GetAppError(err)` - Helper Functions
```go
// Check and extract AppError
if errors.IsAppError(err) {
    appErr := errors.GetAppError(err)
    // Use appErr...
}
```

## Usage Across All Layers

### Handler Layer
```go
if err := c.ShouldBindJSON(&req); err != nil {
    middleware.HandleAppError(c, errors.NewWithDetails(errors.ErrCodeBadRequest, "Invalid request body", err.Error(), http.StatusBadRequest))
    return
}

// Service errors
product, err := productService.CreateProduct(ctx, &req)
if err != nil {
    if appErr := errors.GetAppError(err); appErr != nil {
        middleware.HandleAppError(c, appErr)
    } else {
        middleware.HandleAppError(c, errors.NewWithError(errors.ErrCodeInternalServer, "Failed to create product", http.StatusInternalServerError, err))
    }
    return
}
```

### Service Layer
```go
// Business logic errors
if exists {
    return nil, errors.NewWithDetails(errors.ErrCodeProductSKUExists, "Product SKU already exists", fmt.Sprintf("Product with SKU '%s' already exists", sku), http.StatusConflict).WithField("sku", sku)
}

if product == nil {
    return nil, errors.NewWithDetails(errors.ErrCodeProductNotFound, "Product not found", fmt.Sprintf("Product with ID %d not found", id), http.StatusNotFound).WithField("product_id", id)
}

// Wrap repository errors
if err != nil {
    return nil, errors.NewWithError(errors.ErrCodeDatabaseQuery, "Database operation failed", http.StatusInternalServerError, err).WithField("operation", "check SKU existence")
}
```

### Repository Layer
```go
// Wrap database errors
if err != nil {
    return nil, errors.NewWithError(errors.ErrCodeDatabaseQuery, "Failed to create product", http.StatusInternalServerError, err)
}

// Handle not found
if err == sql.ErrNoRows {
    return nil, errors.NewWithDetails(errors.ErrCodeProductNotFound, "Product not found", fmt.Sprintf("Product with ID %d not found", id), http.StatusNotFound).WithField("product_id", id)
}
```

## Common Patterns

### Pattern 1: Simple Error
```go
return errors.New(errors.ErrCodeBadRequest, "Invalid request", http.StatusBadRequest)
```

### Pattern 2: Error with Details
```go
return errors.NewWithDetails(errors.ErrCodeNotFound, "Resource not found", "The requested resource was not found", http.StatusNotFound)
```

### Pattern 3: Error with Context
```go
return errors.New(errors.ErrCodeProductNotFound, "Product not found", http.StatusNotFound).WithField("product_id", id)
```

### Pattern 4: Wrap Existing Error
```go
return errors.NewWithError(errors.ErrCodeDatabaseQuery, "Database operation failed", http.StatusInternalServerError, err)
```

### Pattern 5: Error with Details + Context
```go
return errors.NewWithDetails(errors.ErrCodeProductSKUExists, "SKU exists", fmt.Sprintf("Product with SKU '%s' already exists", sku), http.StatusConflict).WithField("sku", sku)
```

## Benefits

✅ **Only 5 methods** - Easy to remember and use  
✅ **Handles all cases** - Can create any error type with any context  
✅ **Consistent** - Same pattern across all layers  
✅ **Flexible** - Can chain WithField for additional context  
✅ **No confusion** - Clear when to use which method  

## Quick Reference

| Use Case | Method | Example |
|----------|--------|---------|
| Simple error | `New()` | `errors.New(ErrCodeBadRequest, "Invalid request", 400)` |
| Error with details | `NewWithDetails()` | `errors.NewWithDetails(ErrCodeNotFound, "Not found", "Details", 404)` |
| Wrap existing error | `NewWithError()` | `errors.NewWithError(ErrCodeDatabaseQuery, "DB failed", 500, err)` |
| Add context | `WithField()` | `err.WithField("user_id", 123)` |
| Check error type | `IsAppError()` | `errors.IsAppError(err)` |
