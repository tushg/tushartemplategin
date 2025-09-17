# Error Handling Implementation Guide

## Quick Start

### 1. Basic Error Creation

```go
// Simple error
err := errors.NewBadRequest("Invalid request")

// Error with details
err := errors.NewBadRequestWithDetails("Invalid request", "Missing required field 'name'")

// Error with context
err := errors.NewProductNotFound(123).
    WithField("user_id", 456).
    WithDetails("Product was deleted by admin")
```

### 2. Error Handling in Handlers

```go
func createProductHandler(c *gin.Context) {
    var req CreateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        middleware.HandleAppError(c, errors.NewBadRequestWithDetails("Invalid request body", err.Error()))
        return
    }
    
    product, err := productService.CreateProduct(c.Request.Context(), &req)
    if err != nil {
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

### 3. Error Handling in Services

```go
func (s *ProductService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*Product, error) {
    // Business validation
    if req.Price < 0 {
        return nil, errors.NewValidationErrorWithDetails("Invalid price", "Price must be positive").
            WithField("price", req.Price)
    }
    
    // Check SKU uniqueness
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
        return nil, err // Already wrapped as AppError
    }
    
    return product, nil
}
```

### 4. Error Handling in Repositories

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

## Adding New Error Types

### Step 1: Define Error Code

```go
// In pkg/errors/app_error.go
const (
    // Add new error code
    ErrCodeUserNotFound ErrorCode = "USER_NOT_FOUND"
    ErrCodeInsufficientFunds ErrorCode = "INSUFFICIENT_FUNDS"
)
```

### Step 2: Create Constructor Function

```go
// In pkg/errors/app_error.go
func NewUserNotFound(id interface{}) *AppError {
    return NewNotFoundWithDetails("User not found", fmt.Sprintf("User with ID %v not found", id)).
        WithField("user_id", id)
}

func NewInsufficientFunds(required, available float64) *AppError {
    return NewBadRequestWithDetails("Insufficient funds", 
        fmt.Sprintf("Required: $%.2f, Available: $%.2f", required, available)).
        WithField("required_amount", required).
        WithField("available_amount", available)
}
```

### Step 3: Use in Your Code

```go
// In service layer
func (s *OrderService) ProcessOrder(ctx context.Context, orderID string, userID int64) error {
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        if errors.GetAppError(err) != nil && errors.GetAppError(err).Code == errors.ErrCodeNotFound {
            return errors.NewUserNotFound(userID)
        }
        return errors.NewDatabaseError("get user", err)
    }
    
    if user.Balance < order.TotalAmount {
        return errors.NewInsufficientFunds(order.TotalAmount, user.Balance)
    }
    
    return nil
}
```

## Error Response Format

All errors return a consistent JSON format:

```json
{
  "error": "PRODUCT_NOT_FOUND",
  "message": "Product not found",
  "details": "Product with ID 123 not found",
  "fields": {
    "product_id": 123,
    "user_id": 456
  }
}
```

## Testing Errors

### Unit Tests

```go
func TestNewProductNotFound(t *testing.T) {
    err := errors.NewProductNotFound(123)
    
    assert.Equal(t, errors.ErrCodeProductNotFound, err.Code)
    assert.Equal(t, "Product not found", err.Message)
    assert.Equal(t, http.StatusNotFound, err.HTTPStatus)
    assert.Equal(t, 123, err.Fields["product_id"])
}
```

### Integration Tests

```go
func TestGetProduct_NotFound(t *testing.T) {
    server := setupTestServer()
    
    resp := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/api/v1/products/999", nil)
    server.ServeHTTP(resp, req)
    
    assert.Equal(t, http.StatusNotFound, resp.Code)
    
    var errorResp map[string]interface{}
    json.Unmarshal(resp.Body.Bytes(), &errorResp)
    assert.Equal(t, "PRODUCT_NOT_FOUND", errorResp["error"])
}
```

## Common Patterns

### 1. Database Error Handling

```go
// Pattern: Wrap database errors with context
func (r *Repository) Query(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error) {
    rows, err := r.db.Query(ctx, sql, args...)
    if err != nil {
        return nil, errors.NewDatabaseError("execute query", err)
    }
    return rows, nil
}
```

### 2. Business Logic Validation

```go
// Pattern: Validate business rules and return specific errors
func (s *Service) ProcessOrder(ctx context.Context, order *Order) error {
    if order.Amount <= 0 {
        return errors.NewValidationErrorWithDetails("Invalid order amount", "Amount must be positive").
            WithField("amount", order.Amount)
    }
    
    if order.CustomerID == 0 {
        return errors.NewValidationErrorWithDetails("Customer ID required", "Customer ID cannot be zero").
            WithField("customer_id", order.CustomerID)
    }
    
    return nil
}
```

### 3. Error Chain Preservation

```go
// Pattern: Preserve original error while adding context
func (s *Service) ComplexOperation(ctx context.Context) error {
    err := s.step1(ctx)
    if err != nil {
        return errors.NewInternalServerErrorWithError("Step 1 failed", err).
            WithField("operation", "complex_operation").
            WithField("step", "step1")
    }
    
    err = s.step2(ctx)
    if err != nil {
        return errors.NewInternalServerErrorWithError("Step 2 failed", err).
            WithField("operation", "complex_operation").
            WithField("step", "step2")
    }
    
    return nil
}
```

## Best Practices Checklist

- [ ] Use specific error codes for different scenarios
- [ ] Include relevant context in error fields
- [ ] Preserve error chains for debugging
- [ ] Log errors with sufficient context
- [ ] Test error scenarios thoroughly
- [ ] Use appropriate HTTP status codes
- [ ] Provide user-friendly error messages
- [ ] Handle panics gracefully
- [ ] Validate input at the boundary
- [ ] Don't expose internal implementation details

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

This implementation guide provides everything you need to start using the error handling system effectively in your application!
