# Quick Reference Guide

## File Naming
- **Files**: `snake_case.go` (e.g., `user_service.go`)
- **Folders**: `lowercase` (e.g., `productregistration`)
- **Tests**: `*_test.go` (e.g., `user_service_test.go`)

## Package Naming
- **Package**: `lowercase`, single word (e.g., `productregistration`)
- **Avoid**: `common`, `utils`, `helpers`

## Variable Naming
- **Variables**: `camelCase` (e.g., `userID`, `productName`)
- **Functions**: `PascalCase` public, `camelCase` private
- **Constants**: `PascalCase` or `UPPER_CASE`

## Error Handling
```go
// Business layer - return custom errors
func (s *Service) Method() (*Model, *errors.AppError)

// Infrastructure layer - return standard errors  
func (r *Repository) Method() error

// Handler layer - use middleware
middleware.HandleAppError(c, appErr)
```

## Model Validation
```go
type Model struct {
    Field string `json:"field" db:"field" validate:"required,min=1,max=255"`
}
```

## API Endpoints
```go
GET    /api/v1/products           # List
GET    /api/v1/products/:id       # Get by ID
POST   /api/v1/products           # Create
PUT    /api/v1/products/:id       # Update
DELETE /api/v1/products/:id       # Delete
```

## HTTP Status Codes
- `200` OK - Successful GET, PUT
- `201` Created - Successful POST
- `204` No Content - Successful DELETE
- `400` Bad Request - Invalid data
- `404` Not Found - Resource not found
- `409` Conflict - Business rule violation
- `500` Internal Server Error - System error

## Error Creation
```go
// Simple error
errors.New(code, message, httpStatus)

// With details
errors.NewWithDetails(code, message, details, httpStatus)

// With wrapped error
errors.NewWithError(code, message, httpStatus, err)

// With fields
errors.New(...).WithField(key, value)
```

## Logging
```go
s.logger.Error(ctx, "Operation failed", interfaces.Fields{
    "error":     err.Error(),
    "operation": "operation_name",
    "context":   "additional_context",
})
```

## Database Operations
```go
// Use transactions for multiple operations
err := r.db.WithTransaction(ctx, func(tx *sql.Tx) error {
    // Multiple database operations
    return nil
})
```

## Testing
```go
func TestService_Method(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr *errors.AppError
    }{
        {
            name:    "success case",
            input:   InputType{Field: "value"},
            want:    OutputType{Field: "value"},
            wantErr: nil,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Code Review Checklist
- [ ] File naming follows `snake_case.go`
- [ ] Package naming follows `lowercase`
- [ ] Error handling uses custom errors in business layer
- [ ] Models have proper validation tags
- [ ] API endpoints are RESTful
- [ ] HTTP status codes are appropriate
- [ ] Logging includes context
- [ ] Tests cover all scenarios
- [ ] No hardcoded values
- [ ] Proper error wrapping
