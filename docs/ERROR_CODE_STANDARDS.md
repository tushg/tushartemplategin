# Error Code Standards and Guidelines

## Table of Contents
1. [Overview](#overview)
2. [Naming Conventions](#naming-conventions)
3. [Error Code Structure](#error-code-structure)
4. [Category Definitions](#category-definitions)
5. [Component Standards](#component-standards)
6. [HTTP Status Code Mapping](#http-status-code-mapping)
7. [Error Message Guidelines](#error-message-guidelines)
8. [Examples and Patterns](#examples-and-patterns)

## Overview

This document defines the comprehensive error code standards for the microservice application. Error codes provide a structured way to identify, categorize, and handle different types of errors across the system.

## Naming Conventions

### General Rules
- **Format**: `UPPER_CASE_WITH_UNDERSCORES`
- **Pattern**: `{COMPONENT}_{OPERATION}_{TYPE}`
- **Length**: Maximum 50 characters
- **Characters**: Only letters, numbers, and underscores
- **Uniqueness**: Each error code must be unique across the system

### Examples
```go
// ✅ Good
PRODUCT_CREATE_FAILED
USER_VALIDATION_ERROR
PAYMENT_API_TIMEOUT_ERROR
DB_CONNECTION_LOST

// ❌ Bad
product_create_failed          // Wrong case
USER_CREATE_FAILED_ERROR      // Redundant "ERROR"
PRODUCT_CREATE_FAILED_ERROR   // Redundant "ERROR"
PRODUCT_CREATE_FAILED_ERROR_1 // Numbers not allowed
```

## Error Code Structure

### 1. Component Prefix
Identifies the system component or service where the error occurred.

**Format**: `{COMPONENT}_`

**Examples**:
- `PRODUCT_` - Product management service
- `USER_` - User management service
- `ORDER_` - Order processing service
- `PAYMENT_` - Payment processing service
- `INVENTORY_` - Inventory management service
- `NOTIFICATION_` - Notification service
- `DB_` - Database operations
- `HTTP_` - HTTP-related errors
- `SYSTEM_` - System-level errors
- `VALIDATION_` - Validation errors
- `AUTH_` - Authentication/authorization errors

### 2. Operation Suffix
Describes the specific operation that failed.

**Format**: `{OPERATION}_`

**Common Operations**:
- `CREATE` - Creation operations
- `READ` / `GET` - Retrieval operations
- `UPDATE` - Update operations
- `DELETE` - Deletion operations
- `LIST` - Listing operations
- `VALIDATE` - Validation operations
- `AUTHENTICATE` - Authentication operations
- `AUTHORIZE` - Authorization operations
- `CONNECT` - Connection operations
- `QUERY` - Query operations
- `TRANSACTION` - Transaction operations

### 3. Type Suffix
Describes the type or nature of the error.

**Format**: `{TYPE}`

**Common Types**:
- `FAILED` - General failure
- `ERROR` - General error
- `TIMEOUT` - Timeout error
- `NOT_FOUND` - Resource not found
- `DUPLICATE` - Duplicate value
- `INVALID` - Invalid value
- `UNAUTHORIZED` - Unauthorized access
- `FORBIDDEN` - Forbidden access
- `CONFLICT` - Conflict error
- `CONSTRAINT_VIOLATION` - Constraint violation
- `NETWORK_ERROR` - Network-related error
- `SERVICE_ERROR` - External service error
- `VALIDATION_ERROR` - Validation error
- `AUTH_ERROR` - Authentication error
- `RATE_LIMIT_ERROR` - Rate limit exceeded

## Category Definitions

### 1. System Errors (`SYSTEM_*`)

**Purpose**: System-level errors that affect the entire application

**Pattern**: `SYSTEM_{OPERATION}_{TYPE}`

**Examples**:
```go
SYSTEM_INTERNAL_ERROR     // General system error
SYSTEM_CONFIG_ERROR       // Configuration error
SYSTEM_DATABASE_ERROR     // Database system error
SYSTEM_NETWORK_ERROR      // Network system error
SYSTEM_TIMEOUT_ERROR      // System timeout error
SYSTEM_MEMORY_ERROR       // Memory system error
SYSTEM_DISK_ERROR         // Disk system error
```

**Usage**:
- Critical system failures
- Configuration problems
- Infrastructure issues
- Resource exhaustion

### 2. Validation Errors (`VALIDATION_*`)

**Purpose**: Data validation and input validation errors

**Pattern**: `VALIDATION_{TYPE}`

**Examples**:
```go
VALIDATION_REQUIRED_FIELD        // Required field missing
VALIDATION_INVALID_FORMAT        // Invalid data format
VALIDATION_INVALID_RANGE         // Value out of valid range
VALIDATION_INVALID_TYPE          // Invalid data type
VALIDATION_DUPLICATE_VALUE       // Duplicate value detected
VALIDATION_CONSTRAINT_VIOLATION  // Constraint violation
```

**Usage**:
- Input validation failures
- Data format errors
- Business rule violations
- Constraint violations

### 3. Authentication & Authorization (`AUTH_*`)

**Purpose**: Authentication and authorization related errors

**Pattern**: `AUTH_{TYPE}`

**Examples**:
```go
AUTH_INVALID_CREDENTIALS        // Invalid login credentials
AUTH_TOKEN_EXPIRED              // Authentication token expired
AUTH_TOKEN_INVALID              // Invalid authentication token
AUTH_INSUFFICIENT_PERMISSIONS   // Insufficient permissions
AUTH_ACCOUNT_LOCKED             // Account locked
AUTH_ACCOUNT_DISABLED           // Account disabled
AUTH_SESSION_EXPIRED            // Session expired
```

**Usage**:
- Login failures
- Token validation errors
- Permission checks
- Account status issues

### 4. Product Service Errors (`PRODUCT_*`)

**Purpose**: Product management specific errors

**Pattern**: `PRODUCT_{OPERATION}_{TYPE}`

**Examples**:
```go
// Creation
PRODUCT_CREATE_FAILED           // Product creation failed
PRODUCT_CREATE_VALIDATION_ERROR // Product creation validation error
PRODUCT_CREATE_DUPLICATE_SKU    // Duplicate SKU during creation

// Retrieval
PRODUCT_NOT_FOUND               // Product not found
PRODUCT_GET_FAILED              // Product retrieval failed
PRODUCT_LIST_FAILED             // Product listing failed

// Update
PRODUCT_UPDATE_FAILED           // Product update failed
PRODUCT_UPDATE_VALIDATION_ERROR // Product update validation error
PRODUCT_UPDATE_CONFLICT         // Product update conflict

// Deletion
PRODUCT_DELETE_FAILED           // Product deletion failed
PRODUCT_DELETE_NOT_FOUND        // Product not found for deletion
PRODUCT_DELETE_CONSTRAINT_ERROR // Product deletion constraint error

// Stock Management
PRODUCT_STOCK_UPDATE_FAILED     // Stock update failed
PRODUCT_STOCK_INSUFFICIENT      // Insufficient stock
PRODUCT_STOCK_INVALID_QUANTITY  // Invalid stock quantity
```

**Usage**:
- Product CRUD operations
- Business logic validation
- Stock management
- SKU validation

### 5. Database Errors (`DB_*`)

**Purpose**: Database operation specific errors

**Pattern**: `DB_{OPERATION}_{TYPE}`

**Examples**:
```go
// Connection
DB_CONNECTION_FAILED     // Database connection failed
DB_CONNECTION_TIMEOUT    // Database connection timeout
DB_CONNECTION_LOST       // Database connection lost

// Query
DB_QUERY_FAILED          // Database query failed
DB_QUERY_TIMEOUT         // Database query timeout
DB_QUERY_SYNTAX_ERROR    // Database query syntax error

// Transaction
DB_TRANSACTION_FAILED    // Database transaction failed
DB_TRANSACTION_ROLLBACK  // Database transaction rollback
DB_TRANSACTION_DEADLOCK  // Database transaction deadlock

// Constraints
DB_CONSTRAINT_VIOLATION  // Database constraint violation
DB_FOREIGN_KEY_VIOLATION // Foreign key constraint violation
DB_UNIQUE_VIOLATION      // Unique constraint violation
DB_NOT_NULL_VIOLATION    // Not null constraint violation
```

**Usage**:
- Database connection issues
- Query execution failures
- Transaction problems
- Constraint violations

### 6. External API Errors (`*_API_*`)

**Purpose**: External service integration errors

**Pattern**: `{SERVICE}_API_{TYPE}`

**Examples**:
```go
// Payment Service
PAYMENT_API_ERROR              // Payment service API error
PAYMENT_API_NETWORK_ERROR      // Payment service network error
PAYMENT_API_TIMEOUT_ERROR      // Payment service timeout error
PAYMENT_API_AUTH_ERROR         // Payment service authentication error
PAYMENT_API_RATE_LIMIT_ERROR   // Payment service rate limit error
PAYMENT_API_VALIDATION_ERROR   // Payment service validation error
PAYMENT_API_SERVICE_ERROR      // Payment service internal error
PAYMENT_API_UNKNOWN_ERROR      // Payment service unknown error

// Inventory Service
INVENTORY_API_ERROR            // Inventory service API error
INVENTORY_API_NETWORK_ERROR    // Inventory service network error
INVENTORY_API_TIMEOUT_ERROR    // Inventory service timeout error
INVENTORY_API_AUTH_ERROR       // Inventory service authentication error
INVENTORY_API_RATE_LIMIT_ERROR // Inventory service rate limit error
INVENTORY_API_SERVICE_ERROR    // Inventory service internal error
INVENTORY_PRODUCT_NOT_FOUND    // Product not found in inventory
INVENTORY_INSUFFICIENT_STOCK   // Insufficient stock in inventory

// Notification Service
NOTIFICATION_API_ERROR         // Notification service API error
NOTIFICATION_API_NETWORK_ERROR // Notification service network error
NOTIFICATION_API_TIMEOUT_ERROR // Notification service timeout error
NOTIFICATION_API_AUTH_ERROR    // Notification service authentication error
NOTIFICATION_API_SERVICE_ERROR // Notification service internal error
```

**Usage**:
- External service communication
- API integration failures
- Service-specific business errors
- Network and timeout issues

### 7. HTTP Errors (`HTTP_*`)

**Purpose**: HTTP protocol specific errors

**Pattern**: `HTTP_{TYPE}`

**Examples**:
```go
// Client Errors (4xx)
HTTP_BAD_REQUEST           // Bad request (400)
HTTP_UNAUTHORIZED          // Unauthorized (401)
HTTP_FORBIDDEN             // Forbidden (403)
HTTP_NOT_FOUND             // Not found (404)
HTTP_METHOD_NOT_ALLOWED    // Method not allowed (405)
HTTP_CONFLICT              // Conflict (409)
HTTP_UNPROCESSABLE_ENTITY  // Unprocessable entity (422)
HTTP_TOO_MANY_REQUESTS     // Too many requests (429)

// Server Errors (5xx)
HTTP_INTERNAL_SERVER_ERROR // Internal server error (500)
HTTP_BAD_GATEWAY           // Bad gateway (502)
HTTP_SERVICE_UNAVAILABLE   // Service unavailable (503)
HTTP_GATEWAY_TIMEOUT       // Gateway timeout (504)
```

**Usage**:
- HTTP protocol errors
- Status code mapping
- Request/response issues
- Gateway errors

### 8. Business Logic Errors (`*_*_FAILED`)

**Purpose**: Business process specific errors

**Pattern**: `{PROCESS}_{OPERATION}_FAILED`

**Examples**:
```go
// Order Processing
ORDER_CREATE_FAILED           // Order creation failed
ORDER_UPDATE_FAILED           // Order update failed
ORDER_CANCEL_FAILED           // Order cancellation failed
ORDER_PROCESSING_FAILED       // Order processing failed
ORDER_PAYMENT_FAILED          // Order payment failed
ORDER_INVENTORY_FAILED        // Order inventory check failed

// User Management
USER_CREATE_FAILED            // User creation failed
USER_UPDATE_FAILED            // User update failed
USER_DELETE_FAILED            // User deletion failed
USER_EMAIL_ALREADY_EXISTS     // User email already exists
USER_USERNAME_ALREADY_EXISTS  // Username already exists
```

**Usage**:
- Business process failures
- Workflow errors
- Process-specific validations
- Business rule violations

## Component Standards

### Component Naming Convention
- **Format**: `lowercase-with-hyphens`
- **Length**: Maximum 30 characters
- **Characters**: Only lowercase letters, numbers, and hyphens
- **Uniqueness**: Each component name must be unique

### Standard Components
```go
const (
    ComponentProductService    = "product-service"
    ComponentUserService       = "user-service"
    ComponentOrderService      = "order-service"
    ComponentPaymentService    = "payment-service"
    ComponentInventoryService  = "inventory-service"
    ComponentNotificationService = "notification-service"
    ComponentDatabase          = "database"
    ComponentCache             = "cache"
    ComponentQueue             = "queue"
    ComponentRouter            = "router"
    ComponentMiddleware        = "middleware"
    ComponentValidator         = "validator"
    ComponentSerializer        = "serializer"
    ComponentLogger            = "logger"
    ComponentConfig            = "config"
    ComponentHealth            = "health"
)
```

### Component Usage Guidelines
1. **Service Components**: Use for business logic services
2. **Infrastructure Components**: Use for database, cache, queue
3. **HTTP Components**: Use for router, middleware
4. **Utility Components**: Use for validator, serializer, logger
5. **System Components**: Use for config, health

## HTTP Status Code Mapping

### Error Code to HTTP Status Code Mapping

| Error Code Pattern | HTTP Status | Description |
|-------------------|-------------|-------------|
| `VALIDATION_*` | 400 | Bad Request |
| `AUTH_INVALID_CREDENTIALS` | 401 | Unauthorized |
| `AUTH_TOKEN_EXPIRED` | 401 | Unauthorized |
| `AUTH_TOKEN_INVALID` | 401 | Unauthorized |
| `AUTH_INSUFFICIENT_PERMISSIONS` | 403 | Forbidden |
| `AUTH_ACCOUNT_LOCKED` | 403 | Forbidden |
| `AUTH_ACCOUNT_DISABLED` | 403 | Forbidden |
| `*_NOT_FOUND` | 404 | Not Found |
| `*_DUPLICATE_*` | 409 | Conflict |
| `*_CONFLICT` | 409 | Conflict |
| `HTTP_TOO_MANY_REQUESTS` | 429 | Too Many Requests |
| `*_CONSTRAINT_VIOLATION` | 422 | Unprocessable Entity |
| `SYSTEM_*` | 500 | Internal Server Error |
| `*_API_ERROR` | 502 | Bad Gateway |
| `*_API_TIMEOUT_ERROR` | 504 | Gateway Timeout |
| `*_API_SERVICE_ERROR` | 503 | Service Unavailable |

### Default Mapping
- **Unknown Error Codes**: 500 Internal Server Error
- **External API Errors**: 502 Bad Gateway
- **System Errors**: 500 Internal Server Error
- **Validation Errors**: 400 Bad Request

## Error Message Guidelines

### Message Structure
1. **Primary Message**: Clear, user-friendly description
2. **Details**: Technical details for debugging
3. **Context**: Relevant context information

### Message Examples

#### Good Messages
```go
// Clear and specific
"Product with SKU 'ABC123' already exists"
"Invalid email format: 'user@'"
"Database connection timeout after 30 seconds"
"Payment service returned error: Invalid card number"

// User-friendly
"Please provide a valid email address"
"Product not found"
"Access denied: Insufficient permissions"
"Service temporarily unavailable"
```

#### Bad Messages
```go
// Too technical
"SQLSTATE 23505: duplicate key value violates unique constraint"
"HTTP 500 Internal Server Error"
"Error: null pointer exception"

// Too vague
"Error occurred"
"Something went wrong"
"Invalid input"
"Failed"
```

### Message Guidelines
1. **Be Specific**: Include relevant details
2. **Be Clear**: Use simple, understandable language
3. **Be Helpful**: Provide actionable information
4. **Be Consistent**: Use consistent terminology
5. **Be Concise**: Keep messages brief but informative

## Examples and Patterns

### 1. Service Layer Error Creation
```go
// Validation error
return nil, s.errorHelper.HandleValidationError(
    ctx, "price", "Price cannot be negative", errors.ComponentProductService,
)

// Database error
return nil, s.errorHelper.HandleDatabaseError(
    ctx, "PRODUCT_CREATE", errors.ComponentProductService, err,
)

// Business logic error
return nil, s.errorHelper.HandleBusinessLogicError(
    ctx, errors.PRODUCT_CREATE_DUPLICATE_SKU, 
    "Product with SKU already exists", errors.ComponentProductService,
)
```

### 2. External API Error Handling
```go
// Payment service error
return nil, s.errorHelper.HandlePaymentAPIError(
    ctx, httpStatus, responseBody, err,
)

// Inventory service error
return nil, s.errorHelper.HandleInventoryAPIError(
    ctx, httpStatus, responseBody, err,
)
```

### 3. Error Code Usage in Tests
```go
func TestProductCreation(t *testing.T) {
    // Test validation error
    _, err := service.CreateProduct(ctx, invalidRequest)
    assert.Error(t, err)
    assert.Equal(t, errors.VALIDATION_REQUIRED_FIELD, errors.GetAppError(err).Code)
    
    // Test duplicate SKU error
    _, err = service.CreateProduct(ctx, duplicateSKURequest)
    assert.Error(t, err)
    assert.Equal(t, errors.PRODUCT_CREATE_DUPLICATE_SKU, errors.GetAppError(err).Code)
}
```

### 4. Error Response Format
```json
{
  "error": {
    "code": "PRODUCT_CREATE_DUPLICATE_SKU",
    "message": "Product with SKU 'ABC123' already exists",
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

## Conclusion

This error code standard provides a comprehensive framework for error handling across the microservice application. By following these guidelines, you ensure:

- **Consistency**: Uniform error handling across all components
- **Maintainability**: Easy to understand and modify error codes
- **Debugging**: Clear error identification and troubleshooting
- **Monitoring**: Effective error tracking and alerting
- **User Experience**: Clear and helpful error messages

Remember to:
1. Follow the naming conventions strictly
2. Use appropriate error codes for each scenario
3. Provide clear and helpful error messages
4. Map error codes to appropriate HTTP status codes
5. Test error scenarios thoroughly
6. Document any new error codes you create
