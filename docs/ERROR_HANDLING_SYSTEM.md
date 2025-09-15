# Error Handling System Documentation

## Table of Contents
1. [Overview](#overview)
2. [Error Code Standards](#error-code-standards)
3. [Layer-by-Layer Guidelines](#layer-by-layer-guidelines)
4. [Implementation Examples](#implementation-examples)
5. [Best Practices](#best-practices)
6. [Migration Guide](#migration-guide)

## Overview

This document outlines the comprehensive error handling system implemented using **Approach 3: OpenTelemetry-Compatible Error System**. The system provides structured error handling with error codes, stack traces, component tracking, and distributed tracing support.

### Key Features
- ✅ **Structured Error Objects** with error codes, messages, and metadata
- ✅ **Stack Trace Capture** for debugging and troubleshooting
- ✅ **Component Tracking** to identify error sources
- ✅ **OpenTelemetry Integration** with trace ID and span ID support
- ✅ **External API Error Wrapping** for 3rd party service errors
- ✅ **Retry Logic Support** with retryable error identification
- ✅ **Severity Levels** for error prioritization
- ✅ **Standardized HTTP Responses** with consistent error format

## Error Code Standards

### Naming Convention

Error codes follow the pattern: `{COMPONENT}_{OPERATION}_{TYPE}`

#### Examples:
- `PRODUCT_CREATE_FAILED` - Product creation operation failed
- `USER_VALIDATION_ERROR` - User validation operation error
- `PAYMENT_API_ERROR` - Payment API operation error
- `DB_CONNECTION_FAILED` - Database connection operation failed

### Error Code Categories

#### 1. System Errors (`SYSTEM_*`)
- `SYSTEM_INTERNAL_ERROR` - General system errors
- `SYSTEM_CONFIG_ERROR` - Configuration-related errors
- `SYSTEM_DATABASE_ERROR` - Database system errors
- `SYSTEM_NETWORK_ERROR` - Network-related errors
- `SYSTEM_TIMEOUT_ERROR` - Timeout errors
- `SYSTEM_MEMORY_ERROR` - Memory-related errors
- `SYSTEM_DISK_ERROR` - Disk-related errors

#### 2. Validation Errors (`VALIDATION_*`)
- `VALIDATION_REQUIRED_FIELD` - Required field missing
- `VALIDATION_INVALID_FORMAT` - Invalid data format
- `VALIDATION_INVALID_RANGE` - Value out of valid range
- `VALIDATION_INVALID_TYPE` - Invalid data type
- `VALIDATION_DUPLICATE_VALUE` - Duplicate value detected
- `VALIDATION_CONSTRAINT_VIOLATION` - Constraint violation

#### 3. Authentication & Authorization (`AUTH_*`)
- `AUTH_INVALID_CREDENTIALS` - Invalid login credentials
- `AUTH_TOKEN_EXPIRED` - Authentication token expired
- `AUTH_TOKEN_INVALID` - Invalid authentication token
- `AUTH_INSUFFICIENT_PERMISSIONS` - Insufficient permissions
- `AUTH_ACCOUNT_LOCKED` - Account locked
- `AUTH_ACCOUNT_DISABLED` - Account disabled
- `AUTH_SESSION_EXPIRED` - Session expired

#### 4. Product Service Errors (`PRODUCT_*`)
- `PRODUCT_CREATE_FAILED` - Product creation failed
- `PRODUCT_CREATE_VALIDATION_ERROR` - Product creation validation error
- `PRODUCT_CREATE_DUPLICATE_SKU` - Duplicate SKU during creation
- `PRODUCT_NOT_FOUND` - Product not found
- `PRODUCT_GET_FAILED` - Product retrieval failed
- `PRODUCT_LIST_FAILED` - Product listing failed
- `PRODUCT_UPDATE_FAILED` - Product update failed
- `PRODUCT_UPDATE_VALIDATION_ERROR` - Product update validation error
- `PRODUCT_UPDATE_CONFLICT` - Product update conflict
- `PRODUCT_DELETE_FAILED` - Product deletion failed
- `PRODUCT_DELETE_NOT_FOUND` - Product not found for deletion
- `PRODUCT_DELETE_CONSTRAINT_ERROR` - Product deletion constraint error
- `PRODUCT_STOCK_UPDATE_FAILED` - Stock update failed
- `PRODUCT_STOCK_INSUFFICIENT` - Insufficient stock
- `PRODUCT_STOCK_INVALID_QUANTITY` - Invalid stock quantity

#### 5. Database Errors (`DB_*`)
- `DB_CONNECTION_FAILED` - Database connection failed
- `DB_CONNECTION_TIMEOUT` - Database connection timeout
- `DB_CONNECTION_LOST` - Database connection lost
- `DB_QUERY_FAILED` - Database query failed
- `DB_QUERY_TIMEOUT` - Database query timeout
- `DB_QUERY_SYNTAX_ERROR` - Database query syntax error
- `DB_TRANSACTION_FAILED` - Database transaction failed
- `DB_TRANSACTION_ROLLBACK` - Database transaction rollback
- `DB_TRANSACTION_DEADLOCK` - Database transaction deadlock
- `DB_CONSTRAINT_VIOLATION` - Database constraint violation
- `DB_FOREIGN_KEY_VIOLATION` - Foreign key constraint violation
- `DB_UNIQUE_VIOLATION` - Unique constraint violation
- `DB_NOT_NULL_VIOLATION` - Not null constraint violation

#### 6. External API Errors (`*_API_*`)
- `PAYMENT_API_ERROR` - Payment service API error
- `PAYMENT_API_NETWORK_ERROR` - Payment service network error
- `PAYMENT_API_TIMEOUT_ERROR` - Payment service timeout error
- `PAYMENT_API_AUTH_ERROR` - Payment service authentication error
- `PAYMENT_API_RATE_LIMIT_ERROR` - Payment service rate limit error
- `PAYMENT_API_VALIDATION_ERROR` - Payment service validation error
- `PAYMENT_API_SERVICE_ERROR` - Payment service internal error
- `PAYMENT_API_UNKNOWN_ERROR` - Payment service unknown error
- `INVENTORY_API_ERROR` - Inventory service API error
- `INVENTORY_API_NETWORK_ERROR` - Inventory service network error
- `INVENTORY_API_TIMEOUT_ERROR` - Inventory service timeout error
- `INVENTORY_API_AUTH_ERROR` - Inventory service authentication error
- `INVENTORY_API_RATE_LIMIT_ERROR` - Inventory service rate limit error
- `INVENTORY_API_SERVICE_ERROR` - Inventory service internal error
- `INVENTORY_PRODUCT_NOT_FOUND` - Product not found in inventory
- `INVENTORY_INSUFFICIENT_STOCK` - Insufficient stock in inventory
- `NOTIFICATION_API_ERROR` - Notification service API error
- `NOTIFICATION_API_NETWORK_ERROR` - Notification service network error
- `NOTIFICATION_API_TIMEOUT_ERROR` - Notification service timeout error
- `NOTIFICATION_API_AUTH_ERROR` - Notification service authentication error
- `NOTIFICATION_API_SERVICE_ERROR` - Notification service internal error

#### 7. HTTP Errors (`HTTP_*`)
- `HTTP_BAD_REQUEST` - Bad request (400)
- `HTTP_UNAUTHORIZED` - Unauthorized (401)
- `HTTP_FORBIDDEN` - Forbidden (403)
- `HTTP_NOT_FOUND` - Not found (404)
- `HTTP_METHOD_NOT_ALLOWED` - Method not allowed (405)
- `HTTP_CONFLICT` - Conflict (409)
- `HTTP_UNPROCESSABLE_ENTITY` - Unprocessable entity (422)
- `HTTP_TOO_MANY_REQUESTS` - Too many requests (429)
- `HTTP_INTERNAL_SERVER_ERROR` - Internal server error (500)
- `HTTP_BAD_GATEWAY` - Bad gateway (502)
- `HTTP_SERVICE_UNAVAILABLE` - Service unavailable (503)
- `HTTP_GATEWAY_TIMEOUT` - Gateway timeout (504)

#### 8. Business Logic Errors (`*_*_FAILED`)
- `ORDER_CREATE_FAILED` - Order creation failed
- `ORDER_UPDATE_FAILED` - Order update failed
- `ORDER_CANCEL_FAILED` - Order cancellation failed
- `ORDER_PROCESSING_FAILED` - Order processing failed
- `ORDER_PAYMENT_FAILED` - Order payment failed
- `ORDER_INVENTORY_FAILED` - Order inventory check failed
- `USER_CREATE_FAILED` - User creation failed
- `USER_UPDATE_FAILED` - User update failed
- `USER_DELETE_FAILED` - User deletion failed
- `USER_EMAIL_ALREADY_EXISTS` - User email already exists
- `USER_USERNAME_ALREADY_EXISTS` - Username already exists

### Component Names

Standard component names for error tracking:

- `product-service` - Product management service
- `user-service` - User management service
- `order-service` - Order processing service
- `payment-service` - Payment processing service
- `inventory-service` - Inventory management service
- `notification-service` - Notification service
- `database` - Database operations
- `cache` - Cache operations
- `queue` - Message queue operations
- `router` - HTTP router
- `middleware` - Middleware components
- `validator` - Data validation
- `serializer` - Data serialization
- `logger` - Logging system
- `config` - Configuration management
- `health` - Health check system

## Layer-by-Layer Guidelines

### 1. Router Layer (HTTP Handlers)

**Purpose**: Handle HTTP requests and responses, convert errors to HTTP status codes

**Responsibilities**:
- Parse and validate HTTP requests
- Call service layer methods
- Convert service errors to HTTP responses
- Add request context (trace ID, correlation ID)

**Error Handling Pattern**:
```go
func createProductHandler(c *gin.Context) {
    ctx := c.Request.Context()
    
    // Parse request
    var req CreateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // Use error helper for validation errors
        appErr := errorHelper.HandleValidationError(
            ctx, 
            "request_body", 
            "Invalid request body format", 
            errors.ComponentRouter,
        )
        c.Error(appErr) // Let middleware handle the response
        return
    }
    
    // Call service
    product, err := productService.CreateProduct(ctx, &req)
    if err != nil {
        c.Error(err) // Pass through service error
        return
    }
    
    // Return success response
    c.JSON(http.StatusCreated, ProductResponse{Product: *product})
}
```

**Key Points**:
- Use `c.Error(err)` to pass errors to middleware
- Don't manually create HTTP responses for errors
- Always pass context to service layer
- Use error helper for common validation errors

### 2. Service Layer (Business Logic)

**Purpose**: Implement business logic, orchestrate operations, handle business rules

**Responsibilities**:
- Implement business logic and rules
- Coordinate between different components
- Handle business-level validation
- Call repository and external services
- Transform data between layers

**Error Handling Pattern**:
```go
func (s *ProductService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductRegistration, error) {
    // Business validation
    if req.Price < 0 {
        return nil, s.errorHelper.HandleValidationError(
            ctx,
            "price",
            "Price cannot be negative",
            errors.ComponentProductService,
        )
    }
    
    // Check business rules
    exists, err := s.repo.SKUExists(ctx, req.SKU, nil)
    if err != nil {
        return nil, s.errorHelper.HandleDatabaseError(
            ctx,
            "SKU_EXISTS_CHECK",
            errors.ComponentProductService,
            err,
        )
    }
    
    if exists {
        return nil, s.errorHelper.HandleDuplicateValueError(
            ctx,
            "sku",
            req.SKU,
            errors.ComponentProductService,
        )
    }
    
    // Call external services
    if err := s.inventoryService.CheckStock(ctx, req.SKU, req.Stock); err != nil {
        // External service error is already wrapped
        return nil, err
    }
    
    // Create product
    product, err := s.repo.Create(ctx, product)
    if err != nil {
        return nil, s.errorHelper.HandleDatabaseError(
            ctx,
            "PRODUCT_CREATE",
            errors.ComponentProductService,
            err,
        )
    }
    
    return product, nil
}
```

**Key Points**:
- Use error helper for common error types
- Don't wrap errors that are already AppError
- Pass context to all method calls
- Use appropriate error codes for business logic

### 3. Repository Layer (Data Access)

**Purpose**: Handle data persistence, database operations, data transformation

**Responsibilities**:
- Execute database queries
- Handle database connections and transactions
- Transform data between database and domain models
- Handle database-specific errors

**Error Handling Pattern**:
```go
func (r *ProductRepository) Create(ctx context.Context, product *ProductRegistration) (*ProductRegistration, error) {
    query := `INSERT INTO products (name, description, category, price, sku, stock, is_active, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
    
    err := r.db.QueryRowContext(ctx, query, 
        product.Name, product.Description, product.Category, 
        product.Price, product.SKU, product.Stock, 
        product.IsActive, product.CreatedAt, product.UpdatedAt,
    ).Scan(&product.ID)
    
    if err != nil {
        // Handle specific database errors
        if isConstraintViolation(err) {
            return nil, r.errorHelper.HandleConstraintError(ctx, errors.ComponentDatabase, err)
        }
        
        return nil, r.errorHelper.HandleDatabaseError(
            ctx,
            "PRODUCT_CREATE",
            errors.ComponentDatabase,
            err,
        )
    }
    
    return product, nil
}
```

**Key Points**:
- Use error helper for database-specific errors
- Handle constraint violations appropriately
- Always pass context to database operations
- Use appropriate component names

### 4. External Service Layer

**Purpose**: Handle communication with external APIs and services

**Responsibilities**:
- Make HTTP requests to external services
- Handle external service responses
- Wrap external service errors
- Implement retry logic and circuit breakers

**Error Handling Pattern**:
```go
func (ps *PaymentService) ProcessPayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
    // Make HTTP request
    resp, err := ps.client.Do(httpReq)
    if err != nil {
        // Network error
        return nil, ps.errorHelper.HandlePaymentAPIError(
            ctx,
            0, // No HTTP status for network errors
            "",
            err,
        )
    }
    defer resp.Body.Close()
    
    // Handle different HTTP status codes
    switch resp.StatusCode {
    case http.StatusOK:
        // Success case
        var paymentResp PaymentResponse
        if err := json.NewDecoder(resp.Body).Decode(&paymentResp); err != nil {
            return nil, ps.errorHelper.HandlePaymentAPIError(
                ctx,
                resp.StatusCode,
                "Invalid JSON response",
                err,
            )
        }
        return &paymentResp, nil
        
    case http.StatusBadRequest:
        body, _ := io.ReadAll(resp.Body)
        return nil, ps.errorHelper.HandlePaymentAPIError(
            ctx,
            resp.StatusCode,
            string(body),
            fmt.Errorf("validation error"),
        )
        
    case http.StatusUnauthorized:
        return nil, ps.errorHelper.HandlePaymentAPIError(
            ctx,
            resp.StatusCode,
            "Invalid API key",
            fmt.Errorf("authentication error"),
        )
        
    default:
        body, _ := io.ReadAll(resp.Body)
        return nil, ps.errorHelper.HandlePaymentAPIError(
            ctx,
            resp.StatusCode,
            string(body),
            fmt.Errorf("unknown error"),
        )
    }
}
```

**Key Points**:
- Use specific error helpers for each external service
- Handle different HTTP status codes appropriately
- Preserve external service error details
- Implement proper error wrapping

### 5. Utility Functions

**Purpose**: Provide reusable utility functions across the application

**Responsibilities**:
- Implement common functionality
- Handle utility-specific errors
- Provide error-safe operations

**Error Handling Pattern**:
```go
func ValidateEmail(email string) error {
    if email == "" {
        return errors.NewValidationError(
            errors.VALIDATION_REQUIRED_FIELD,
            "Email is required",
            errors.ComponentValidator,
        )
    }
    
    if !isValidEmailFormat(email) {
        return errors.NewValidationError(
            errors.VALIDATION_INVALID_FORMAT,
            "Invalid email format",
            errors.ComponentValidator,
        )
    }
    
    return nil
}

func ParseJSON(data []byte, v interface{}) error {
    if err := json.Unmarshal(data, v); err != nil {
        return errors.NewInternalError(
            errors.SYSTEM_INTERNAL_ERROR,
            "Failed to parse JSON",
            errors.ComponentSerializer,
        )
    }
    
    return nil
}
```

**Key Points**:
- Use appropriate error codes for utility functions
- Don't wrap errors unnecessarily
- Use appropriate component names
- Provide clear error messages

### 6. Main Function

**Purpose**: Application entry point, initialization, and error handling

**Responsibilities**:
- Initialize application components
- Handle startup errors
- Set up error handling middleware
- Handle graceful shutdown

**Error Handling Pattern**:
```go
func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }
    
    // Initialize logger
    appLogger, err := logger.NewLogger(&cfg.Log)
    if err != nil {
        log.Fatalf("Failed to initialize logger: %v", err)
    }
    
    // Initialize database
    db, err := database.NewDatabaseFactory(appLogger).CreateDatabase(&cfg.Database)
    if err != nil {
        appLogger.Fatal(context.Background(), "Failed to create database", interfaces.Fields{
            "error": err.Error(),
        })
    }
    
    // Set up router with error middleware
    router := gin.New()
    router.Use(errors.ErrorHandlerMiddleware())
    
    // ... rest of initialization
}
```

**Key Points**:
- Use fatal errors for critical startup failures
- Initialize error handling middleware early
- Log errors with proper context
- Handle graceful shutdown

## Implementation Examples

### Example 1: Product Creation with External Service Integration

```go
// Service layer
func (s *ProductService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductRegistration, error) {
    // Business validation
    if req.Price < 0 {
        return nil, s.errorHelper.HandleValidationError(
            ctx, "price", "Price cannot be negative", errors.ComponentProductService,
        )
    }
    
    // Check SKU uniqueness
    exists, err := s.repo.SKUExists(ctx, req.SKU, nil)
    if err != nil {
        return nil, s.errorHelper.HandleDatabaseError(
            ctx, "SKU_EXISTS_CHECK", errors.ComponentProductService, err,
        )
    }
    
    if exists {
        return nil, s.errorHelper.HandleDuplicateValueError(
            ctx, "sku", req.SKU, errors.ComponentProductService,
        )
    }
    
    // Check inventory availability
    if err := s.inventoryService.CheckStock(ctx, req.SKU, req.Stock); err != nil {
        return nil, err // Already wrapped by external service
    }
    
    // Process payment
    paymentResp, err := s.paymentService.ProcessPayment(ctx, &PaymentRequest{
        Amount: req.Price * float64(req.Stock),
        Currency: "USD",
        OrderID: generateOrderID(),
    })
    if err != nil {
        return nil, err // Already wrapped by external service
    }
    
    // Create product
    product := &ProductRegistration{
        Name: req.Name,
        Description: req.Description,
        Category: req.Category,
        Price: req.Price,
        SKU: req.SKU,
        Stock: req.Stock,
        IsActive: req.IsActive,
    }
    
    createdProduct, err := s.repo.Create(ctx, product)
    if err != nil {
        return nil, s.errorHelper.HandleDatabaseError(
            ctx, "PRODUCT_CREATE", errors.ComponentProductService, err,
        )
    }
    
    return createdProduct, nil
}
```

### Example 2: Error Response Format

```json
{
  "error": {
    "code": "PRODUCT_CREATE_DUPLICATE_SKU",
    "message": "Duplicate value for field: sku",
    "component": "product-service",
    "timestamp": "2024-01-15T10:30:45.123456Z",
    "trace_id": "550e8400e29b41d4a716446655440000",
    "severity": "medium",
    "source": "internal",
    "retryable": false
  },
  "request_id": "req_12345678",
  "path": "/api/v1/products",
  "method": "POST",
  "timestamp": "2024-01-15T10:30:45.123456Z"
}
```

### Example 3: External API Error Response

```json
{
  "error": {
    "code": "PAYMENT_API_ERROR",
    "message": "Payment service error",
    "component": "payment-service",
    "service": "stripe-payment-service",
    "timestamp": "2024-01-15T10:30:45.123456Z",
    "trace_id": "550e8400e29b41d4a716446655440000",
    "severity": "high",
    "source": "external",
    "retryable": true,
    "external_api": {
      "service_name": "stripe-payment-service",
      "endpoint": "https://api.stripe.com/v1/payments",
      "http_status_code": 500,
      "response_body": "{\"error\":\"Internal server error\"}",
      "request_id": "req_1234567890"
    }
  },
  "request_id": "req_12345678",
  "path": "/api/v1/products",
  "method": "POST",
  "timestamp": "2024-01-15T10:30:45.123456Z"
}
```

## Best Practices

### 1. Error Code Naming
- Use UPPER_CASE with underscores
- Follow the pattern: `{COMPONENT}_{OPERATION}_{TYPE}`
- Be descriptive and specific
- Avoid generic error codes when possible

### 2. Error Messages
- Use clear, user-friendly messages
- Avoid technical jargon in user-facing messages
- Include relevant context when helpful
- Keep messages concise but informative

### 3. Component Names
- Use consistent naming across the application
- Use lowercase with hyphens
- Be descriptive of the actual component
- Follow the established naming convention

### 4. Error Wrapping
- Don't wrap errors that are already AppError
- Preserve original error context
- Add relevant metadata when wrapping
- Use appropriate error codes for wrapped errors

### 5. Context Usage
- Always pass context to service methods
- Include trace ID and correlation ID
- Use context for cancellation and timeouts
- Don't store context in structs

### 6. Logging
- Log errors with appropriate level
- Include relevant context in log fields
- Don't log the same error multiple times
- Use structured logging format

### 7. Testing
- Test error scenarios thoroughly
- Verify error codes and messages
- Test error wrapping and unwrapping
- Mock external service errors

## Migration Guide

### Step 1: Update Imports
```go
import "tushartemplategin/pkg/errors"
```

### Step 2: Replace fmt.Errorf with AppError
```go
// Before
return nil, fmt.Errorf("failed to create product: %w", err)

// After
return nil, s.errorHelper.HandleDatabaseError(
    ctx, "PRODUCT_CREATE", errors.ComponentProductService, err,
)
```

### Step 3: Update Error Handling in Handlers
```go
// Before
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
        "error": constants.ERROR_500_INTERNAL_SERVER_ERROR,
        "details": err.Error(),
    })
    return
}

// After
if err != nil {
    c.Error(err) // Let middleware handle the response
    return
}
```

### Step 4: Add Error Middleware
```go
router.Use(errors.ErrorHandlerMiddleware())
```

### Step 5: Update Service Layer
- Add error helper to service structs
- Replace fmt.Errorf with appropriate error helper methods
- Use proper error codes and component names

### Step 6: Test Error Scenarios
- Verify error responses are properly formatted
- Test error codes and HTTP status codes
- Verify trace ID and correlation ID are included
- Test external service error handling

## Conclusion

This error handling system provides a comprehensive, production-ready solution for managing errors across all layers of your microservice. It ensures consistent error handling, proper error tracking, and seamless integration with observability tools.

The system is designed to be:
- **Scalable**: Easy to add new error types and components
- **Maintainable**: Clear separation of concerns and consistent patterns
- **Observable**: Rich error data for monitoring and debugging
- **Traceable**: Full request tracing across service boundaries
- **User-Friendly**: Clear error messages and appropriate HTTP status codes
