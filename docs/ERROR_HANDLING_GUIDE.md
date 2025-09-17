# Error Handling System Guide

This guide explains how to use the custom error handling system and middleware in the Tushar Template Gin project.

## Overview

The project now includes a comprehensive error handling system with:
- **Structured Error Types**: Custom `AppError` type with standardized error codes
- **Centralized Error Handling**: Middleware that catches and formats all errors
- **Request Validation**: Middleware for validating JSON requests
- **Consistent Error Responses**: Standardized JSON error responses across all endpoints

## Error Types

### AppError Structure

```go
type AppError struct {
    Code       ErrorCode              `json:"code"`                 // Standardized error code
    Message    string                 `json:"message"`              // Human-readable error message
    Details    string                 `json:"details,omitempty"`    // Additional error details
    HTTPStatus int                    `json:"-"`                    // HTTP status code
    Fields     map[string]interface{} `json:"fields,omitempty"`     // Additional context fields
    Err        error                  `json:"-"`                    // Original error
}
```

### Predefined Error Codes

#### General HTTP Errors
- `INTERNAL_SERVER_ERROR` (500)
- `BAD_REQUEST` (400)
- `UNAUTHORIZED` (401)
- `FORBIDDEN` (403)
- `NOT_FOUND` (404)
- `CONFLICT` (409)
- `UNPROCESSABLE_ENTITY` (422)
- `TOO_MANY_REQUESTS` (429)
- `SERVICE_UNAVAILABLE` (503)

#### Validation Errors
- `VALIDATION_FAILED` (422)
- `INVALID_PARAMETER` (400)
- `MISSING_PARAMETER` (400)

#### Business Logic Errors
- `PRODUCT_NOT_FOUND` (404)
- `PRODUCT_SKU_EXISTS` (409)
- `PRODUCT_CREATE_FAILED` (500)
- `PRODUCT_UPDATE_FAILED` (500)
- `PRODUCT_DELETE_FAILED` (500)
- `INVALID_STOCK` (422)

#### Database Errors
- `DATABASE_CONNECTION_ERROR` (500)
- `DATABASE_QUERY_ERROR` (500)
- `DATABASE_TRANSACTION_ERROR` (500)

## Usage Examples

### 1. Creating Custom Errors

```go
// Simple error
err := errors.NewBadRequest("Invalid request")

// Error with details
err := errors.NewBadRequestWithDetails("Invalid request", "Missing required field 'name'")

// Error with context fields
err := errors.NewProductNotFound(123).WithField("user_id", 456)

// Error wrapping existing error
err := errors.NewInternalServerErrorWithError("Database operation failed", dbErr)

// Business logic error
err := errors.NewProductSKUExists("ABC-123")
```

### 2. Using in Route Handlers

```go
func createProductHandler(c *gin.Context) {
    // Parse request
    var req CreateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        middleware.HandleAppError(c, errors.NewBadRequestWithDetails("Invalid request body", err.Error()))
        return
    }

    // Business logic
    product, err := productService.CreateProduct(ctx, &req)
    if err != nil {
        // Check if it's already an AppError
        if appErr := errors.GetAppError(err); appErr != nil {
            middleware.HandleAppError(c, appErr)
        } else {
            middleware.HandleAppError(c, errors.NewInternalServerErrorWithError("Failed to create product", err))
        }
        return
    }

    c.JSON(http.StatusCreated, product)
}
```

### 3. Using in Service Layer

```go
func (s *ProductService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*Product, error) {
    // Check if SKU exists
    exists, err := s.repo.SKUExists(ctx, req.SKU, nil)
    if err != nil {
        return nil, errors.NewDatabaseError("check SKU existence", err)
    }

    if exists {
        return nil, errors.NewProductSKUExists(req.SKU)
    }

    // Create product
    product, err := s.repo.Create(ctx, product)
    if err != nil {
        return nil, errors.NewDatabaseError("create product", err)
    }

    return product, nil
}
```

### 4. Using in Repository Layer

```go
func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*Product, error) {
    var product Product
    err := r.db.QueryRow(ctx, "SELECT * FROM products WHERE id = $1", id).Scan(&product)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.NewProductNotFound(id)
        }
        return nil, errors.NewDatabaseError("get product by ID", err)
    }
    return &product, nil
}
```

## Middleware

### Error Handling Middleware

The `ErrorHandlerMiddleware` provides:
- **Panic Recovery**: Catches panics and converts them to 500 errors
- **Error Conversion**: Converts standard Go errors to AppError
- **Structured Logging**: Logs errors with context
- **Consistent Responses**: Sends standardized JSON error responses

### Validation Middleware

The `ValidationMiddleware` provides:
- **JSON Validation**: Validates JSON request bodies
- **Request Body Preservation**: Restores request body for handlers
- **Basic Validation**: Performs basic request validation

### Usage in Server Setup

```go
// Set up middleware in order
router.Use(middleware.ErrorHandlerMiddleware(logger))      // First - catches all errors
router.Use(middleware.CorrelationIDMiddleware())           // Second - adds correlation IDs
router.Use(middleware.SecurityHeaders())                   // Third - adds security headers
router.Use(middleware.ValidationMiddleware(logger))        // Fourth - validates requests

// Set up 404 and 405 handlers
router.NoRoute(middleware.NotFoundHandler(logger))
router.NoMethod(middleware.MethodNotAllowedHandler(logger))
```

## Error Response Format

All errors are returned in a consistent JSON format:

```json
{
  "error": "PRODUCT_NOT_FOUND",
  "message": "Product not found",
  "details": "Product with ID 123 not found",
  "fields": {
    "product_id": 123
  }
}
```

## Best Practices

### 1. Error Creation
- Use predefined error constructors when possible
- Add context fields for debugging
- Wrap underlying errors to preserve stack traces

### 2. Error Handling
- Check if error is already an AppError before wrapping
- Use appropriate error codes for different scenarios
- Log errors with sufficient context

### 3. Middleware Order
- Error handling middleware should be first
- Validation middleware should be after error handling
- Security middleware should be early in the chain

### 4. Service Layer
- Return AppError types from service methods
- Don't expose internal implementation details
- Use business logic error codes

### 5. Repository Layer
- Wrap database errors with context
- Use specific error codes for different database operations
- Handle "not found" cases explicitly

## Testing Error Handling

```go
func TestCreateProduct_InvalidSKU(t *testing.T) {
    // Setup
    service := NewProductService(mockRepo, mockLogger)
    
    // Test
    _, err := service.CreateProduct(ctx, &CreateProductRequest{
        SKU: "INVALID-SKU",
    })
    
    // Assert
    assert.Error(t, err)
    appErr := errors.GetAppError(err)
    assert.NotNil(t, appErr)
    assert.Equal(t, errors.ErrCodeValidationFailed, appErr.Code)
}
```

## Migration from Old Error Handling

### Before (Old Way)
```go
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
        "error":   constants.ERROR_500_INTERNAL_SERVER_ERROR,
        "details": err.Error(),
    })
    return
}
```

### After (New Way)
```go
if err != nil {
    if appErr := errors.GetAppError(err); appErr != nil {
        middleware.HandleAppError(c, appErr)
    } else {
        middleware.HandleAppError(c, errors.NewInternalServerErrorWithError("Operation failed", err))
    }
    return
}
```

## Benefits

1. **Consistency**: All errors follow the same format
2. **Debugging**: Rich context and structured logging
3. **Maintainability**: Centralized error handling logic
4. **Type Safety**: Structured error types prevent mistakes
5. **Extensibility**: Easy to add new error types and codes
6. **Testing**: Easy to test error scenarios
7. **Documentation**: Self-documenting error codes and messages
