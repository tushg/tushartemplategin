# Coding Standards - Quick Reference

## File and Folder Naming

### ✅ Correct Examples
```
internal/domains/productregistration/
├── interfaces.go          # Domain interfaces
├── models.go              # Domain models and DTOs
├── repository.go          # Repository implementation
├── service.go             # Service implementation
├── routes.go              # HTTP routes and handlers
├── repository_test.go     # Repository tests
└── service_test.go        # Service tests

pkg/errors/
├── app_error.go           # Core error types
├── app_error_test.go      # Error tests
└── USAGE_GUIDE.md         # Usage documentation
```

### ❌ Incorrect Examples
```
internal/domains/ProductRegistration/  # Wrong: PascalCase
├── Interfaces.go                      # Wrong: PascalCase
├── Models.go                          # Wrong: PascalCase
├── Repository.go                      # Wrong: PascalCase
└── Service.go                         # Wrong: PascalCase

pkg/Common/                            # Wrong: Generic name
├── Utils.go                           # Wrong: Generic name
└── Helpers.go                         # Wrong: Generic name
```

## Package Naming

### ✅ Correct Examples
```go
package productregistration  // Matches folder name
package errors              // Single word, lowercase
package middleware          // Single word, lowercase
package config              // Single word, lowercase
```

### ❌ Incorrect Examples
```go
package ProductRegistration  // Wrong: PascalCase
package product_registration // Wrong: snake_case
package common               // Wrong: Too generic
package utils                // Wrong: Too generic
package helpers              // Wrong: Too generic
```

## Variable and Function Naming

### ✅ Correct Examples
```go
// Variables
var userID int64
var productName string
var isActive bool
var createdAt time.Time

// Functions
func CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductRegistration, *errors.AppError)
func GetUserByID(ctx context.Context, id int64) (*User, *errors.AppError)
func ValidateRequest(req *CreateProductRequest) *errors.AppError

// Constants
const (
    DefaultPageSize = 10
    MaxPageSize     = 100
    ErrorCodeProductNotFound = "PRODUCT_NOT_FOUND"
)
```

### ❌ Incorrect Examples
```go
// Variables
var userId int64           // Wrong: camelCase
var ProductName string     // Wrong: PascalCase
var IsActive bool          // Wrong: PascalCase
var created_at time.Time   // Wrong: snake_case

// Functions
func createProduct(ctx context.Context, req *CreateProductRequest) (*ProductRegistration, *errors.AppError)  // Wrong: lowercase
func GetUserById(ctx context.Context, id int64) (*User, *errors.AppError)  // Wrong: camelCase
func validateRequest(req *CreateProductRequest) *errors.AppError  // Wrong: lowercase

// Constants
const (
    defaultPageSize = 10   // Wrong: lowercase
    MAX_PAGE_SIZE = 100    // Wrong: ALL_CAPS
    errorCodeProductNotFound = "PRODUCT_NOT_FOUND"  // Wrong: lowercase
)
```

## Struct and Interface Naming

### ✅ Correct Examples
```go
// Structs
type ProductRegistration struct {
    ID          int64     `json:"id" db:"id"`
    Name        string    `json:"name" db:"name"`
    SKU         string    `json:"sku" db:"sku"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Interfaces
type Service interface {
    CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductRegistration, *errors.AppError)
    GetProduct(ctx context.Context, id int64) (*ProductRegistration, *errors.AppError)
}

// Request/Response DTOs
type CreateProductRequest struct {
    Name string `json:"name" validate:"required"`
    SKU  string `json:"sku" validate:"required"`
}

type ProductResponse struct {
    Product ProductRegistration `json:"product"`
}
```

### ❌ Incorrect Examples
```go
// Structs
type productRegistration struct {  // Wrong: lowercase
    id          int64             // Wrong: lowercase
    name        string            // Wrong: lowercase
    sku         string            // Wrong: lowercase
    created_at  time.Time         // Wrong: snake_case
}

// Interfaces
type service interface {  // Wrong: lowercase
    createProduct(ctx context.Context, req *CreateProductRequest) (*ProductRegistration, *errors.AppError)  // Wrong: lowercase
}

// Request/Response DTOs
type createProductRequest struct {  // Wrong: lowercase
    name string `json:"name"`       // Wrong: lowercase
    sku  string `json:"sku"`        // Wrong: lowercase
}
```

## Error Handling Patterns

### ✅ Correct Examples
```go
// Business Layer - Return custom errors
func (s *ProductService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductRegistration, *errors.AppError) {
    if exists {
        return nil, errors.NewWithDetails(errors.ErrCodeProductSKUExists, "Product SKU already exists", 
            fmt.Sprintf("Product with SKU '%s' already exists", req.SKU), http.StatusConflict).WithField("sku", req.SKU)
    }
    return product, nil
}

// Repository Layer - Wrap infrastructure errors
func (r *ProductRepository) Create(ctx context.Context, product *ProductRegistration) (*ProductRegistration, *errors.AppError) {
    if err := r.db.WithTransaction(ctx, func(tx *sql.Tx) error {
        return tx.QueryRowContext(ctx, query, ...).Scan(...)
    }); err != nil {
        return nil, errors.NewWithError(errors.ErrCodeDatabaseQuery, "Failed to create product", 
            http.StatusInternalServerError, err).WithField("operation", "create product")
    }
    return product, nil
}

// Handler Layer - Use middleware
func createProductHandler(c *gin.Context) {
    product, appErr := productService.CreateProduct(ctx, &req)
    if appErr != nil {
        middleware.HandleAppError(c, appErr)
        return
    }
    c.JSON(http.StatusCreated, ProductResponse{Product: *product})
}
```

### ❌ Incorrect Examples
```go
// Business Layer - Don't return standard errors
func (s *ProductService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductRegistration, error) {
    if exists {
        return nil, fmt.Errorf("product SKU already exists")  // Wrong: standard error
    }
    return product, nil
}

// Repository Layer - Don't return standard errors
func (r *ProductRepository) Create(ctx context.Context, product *ProductRegistration) (*ProductRegistration, error) {
    if err := r.db.WithTransaction(ctx, func(tx *sql.Tx) error {
        return tx.QueryRowContext(ctx, query, ...).Scan(...)
    }); err != nil {
        return nil, err  // Wrong: standard error
    }
    return product, nil
}

// Handler Layer - Don't handle errors manually
func createProductHandler(c *gin.Context) {
    product, err := productService.CreateProduct(ctx, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})  // Wrong: manual handling
        return
    }
    c.JSON(http.StatusCreated, ProductResponse{Product: *product})
}
```

## Model Validation

### ✅ Correct Examples
```go
type ProductRegistration struct {
    ID          int64     `json:"id" db:"id" validate:"required"`
    Name        string    `json:"name" db:"name" validate:"required,min=1,max=255"`
    SKU         string    `json:"sku" db:"sku" validate:"required,min=1,max=50,alphanum"`
    Price       float64   `json:"price" db:"price" validate:"required,min=0"`
    Stock       int       `json:"stock" db:"stock" validate:"min=0"`
    IsActive    bool      `json:"is_active" db:"is_active"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateProductRequest struct {
    Name        string  `json:"name" validate:"required,min=1,max=255"`
    Description string  `json:"description" validate:"max=1000"`
    Category    string  `json:"category" validate:"required,min=1,max=100"`
    Price       float64 `json:"price" validate:"required,min=0"`
    SKU         string  `json:"sku" validate:"required,min=1,max=50,alphanum"`
    Stock       int     `json:"stock" validate:"min=0"`
    IsActive    bool    `json:"is_active"`
}
```

### ❌ Incorrect Examples
```go
type ProductRegistration struct {
    id          int64     `json:"id" db:"id"`                    // Wrong: lowercase
    name        string    `json:"name" db:"name"`                // Wrong: lowercase
    sku         string    `json:"sku" db:"sku"`                  // Wrong: lowercase
    price       float64   `json:"price" db:"price"`              // Wrong: lowercase
    stock       int       `json:"stock" db:"stock"`              // Wrong: lowercase
    is_active   bool      `json:"is_active" db:"is_active"`      // Wrong: snake_case
    created_at  time.Time `json:"created_at" db:"created_at"`    // Wrong: snake_case
    updated_at  time.Time `json:"updated_at" db:"updated_at"`    // Wrong: snake_case
}

type CreateProductRequest struct {
    Name        string  `json:"name"`                            // Wrong: No validation
    Description string  `json:"description"`                     // Wrong: No validation
    Category    string  `json:"category"`                        // Wrong: No validation
    Price       float64 `json:"price"`                           // Wrong: No validation
    SKU         string  `json:"sku"`                             // Wrong: No validation
    Stock       int     `json:"stock"`                           // Wrong: No validation
    IsActive    bool    `json:"is_active"`                       // Wrong: No validation
}
```

## API Design

### ✅ Correct Examples
```go
// RESTful endpoints
GET    /api/v1/products           # List products
GET    /api/v1/products/:id       # Get product by ID
POST   /api/v1/products           # Create product
PUT    /api/v1/products/:id       # Update product
DELETE /api/v1/products/:id       # Delete product

// HTTP status codes
200 OK          # Successful GET, PUT
201 Created     # Successful POST
204 No Content  # Successful DELETE
400 Bad Request # Invalid request data
404 Not Found   # Resource not found
409 Conflict    # Business rule violation
500 Internal Server Error # System error

// Response format
type ProductResponse struct {
    Product ProductRegistration `json:"product"`
}

type ProductListResponse struct {
    Products []ProductRegistration `json:"products"`
    Total    int64                 `json:"total"`
    Page     int                   `json:"page"`
    Limit    int                   `json:"limit"`
}
```

### ❌ Incorrect Examples
```go
// Non-RESTful endpoints
GET    /api/v1/getProducts        # Wrong: Verb in URL
POST   /api/v1/createProduct      # Wrong: Verb in URL
PUT    /api/v1/updateProduct/:id  # Wrong: Verb in URL
DELETE /api/v1/deleteProduct/:id  # Wrong: Verb in URL

// Wrong HTTP status codes
200 OK          # Wrong: Should be 201 for POST
200 OK          # Wrong: Should be 204 for DELETE
500 Internal Server Error # Wrong: Should be 400 for validation errors

// Wrong response format
type ProductResponse struct {
    product ProductRegistration `json:"product"`  // Wrong: lowercase
    success bool                `json:"success"`  // Wrong: Unnecessary field
    message string              `json:"message"`  // Wrong: Unnecessary field
}
```

## Testing

### ✅ Correct Examples
```go
func TestProductService_CreateProduct(t *testing.T) {
    tests := []struct {
        name    string
        req     *CreateProductRequest
        want    *ProductRegistration
        wantErr *errors.AppError
    }{
        {
            name: "successful creation",
            req: &CreateProductRequest{
                Name: "Test Product",
                SKU:  "TEST-001",
            },
            want: &ProductRegistration{
                Name: "Test Product",
                SKU:  "TEST-001",
            },
            wantErr: nil,
        },
        {
            name: "duplicate SKU",
            req: &CreateProductRequest{
                Name: "Test Product",
                SKU:  "EXISTING-SKU",
            },
            want:    nil,
            wantErr: errors.NewWithDetails(errors.ErrCodeProductSKUExists, "Product SKU already exists", 
                "Product with SKU 'EXISTING-SKU' already exists", http.StatusConflict),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### ❌ Incorrect Examples
```go
func TestProductService_CreateProduct(t *testing.T) {
    // Wrong: No test cases
    // Wrong: No table-driven tests
    // Wrong: No error cases
    // Wrong: No edge cases
}
```

## Quick Reference

### File Naming
- Files: `snake_case.go`
- Folders: `lowercase`
- Tests: `*_test.go`

### Package Naming
- Single word, lowercase
- Match folder name
- Avoid generic names

### Variable Naming
- Variables: `camelCase`
- Constants: `PascalCase` or `UPPER_CASE`
- Functions: `PascalCase` (public), `camelCase` (private)

### Error Handling
- Business layer: `*errors.AppError`
- Infrastructure layer: `error`
- Handler layer: `middleware.HandleAppError`

### Model Validation
- Use `validate` tags
- Use `json` tags for API
- Use `db` tags for database

### API Design
- RESTful endpoints
- Proper HTTP status codes
- Consistent response format

