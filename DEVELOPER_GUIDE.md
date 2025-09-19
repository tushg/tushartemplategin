# Tushar Template Gin - Developer Guide

## Table of Contents
1. [Project Structure](#project-structure)
2. [Naming Conventions](#naming-conventions)
3. [Package Organization](#package-organization)
4. [Layer Responsibilities](#layer-responsibilities)
5. [Error Handling Standards](#error-handling-standards)
6. [Model Design Guidelines](#model-design-guidelines)
7. [Repository Pattern](#repository-pattern)
8. [Service Layer Guidelines](#service-layer-guidelines)
9. [API Design Standards](#api-design-standards)
10. [Testing Standards](#testing-standards)
11. [Code Review Checklist](#code-review-checklist)

---

## Project Structure

```
tushartemplategin/
├── cmd/                          # Application entry points
│   ├── server/                   # Main server application
│   └── migrate/                  # Database migration tool
├── internal/                     # Private application code
│   └── domains/                  # Domain-specific modules
│       ├── productregistration/  # Product domain
│       └── health/               # Health check domain
├── pkg/                         # Public library code
│   ├── config/                   # Configuration management
│   ├── constants/                # Application constants
│   ├── database/                 # Database abstractions
│   ├── errors/                   # Custom error handling
│   ├── interfaces/               # Interface definitions
│   ├── logger/                   # Logging utilities
│   └── middleware/               # HTTP middleware
├── scripts/                     # Build and deployment scripts
├── docs/                        # Documentation
├── examples/                    # Usage examples
├── mocks/                       # Generated mocks
└── logs/                        # Log files (gitignored)
```

---

## Naming Conventions

### Files and Folders
- **Files**: Use `snake_case.go` (e.g., `user_service.go`, `product_repository.go`)
  - **Justification**: Snake case is the Go community standard for file names, making code more readable and consistent with Go conventions.
- **Folders**: Use `lowercase` (e.g., `productregistration`, `userauthentication`)
  - **Justification**: Lowercase folder names prevent case sensitivity issues across different operating systems and follow Go package naming conventions.
- **Domain folders**: Use descriptive names (e.g., `productregistration`, `orderprocessing`)
  - **Justification**: Descriptive names make the codebase self-documenting and help developers quickly understand the purpose of each domain.

### Packages
- **Package names**: Use `lowercase`, single word (e.g., `productregistration`, `errors`)
  - **Justification**: Single word package names are easier to import and use, following Go's package naming best practices.
- **Avoid**: Generic names like `common`, `utils`, `helpers`
  - **Justification**: Generic names make code harder to understand and maintain, leading to "god packages" that violate single responsibility principle.
- **Domain packages**: Match folder name (e.g., `package productregistration`)
  - **Justification**: Matching package and folder names prevents confusion and makes imports more intuitive.

### Variables and Functions
- **Variables**: Use `camelCase` for private, `PascalCase` for public
  - **Justification**: Follows Go naming conventions and makes visibility clear at a glance.
- **Functions**: Use `PascalCase` for public functions, `camelCase` for private
  - **Justification**: Consistent with Go's export rules and makes API surface clear.

```go
// ✅ GOOD
var userID int64
var productName string
var isActive bool

func CreateProduct(ctx context.Context, req *CreateProductRequest) (*Product, error)
func GetUserByID(ctx context.Context, id int64) (*User, error)

// ❌ BAD
var userId int64
var ProductName string
var IsActive bool

func createProduct(ctx context.Context, req *CreateProductRequest) (*Product, error)
func GetUserById(ctx context.Context, id int64) (*User, error)
```

### Structs and Interfaces
- **Structs**: Use `PascalCase` for public structs, `camelCase` for private
  - **Justification**: Follows Go naming conventions and makes struct visibility clear.
- **Interfaces**: Use `PascalCase` and descriptive names
  - **Justification**: Clear interface names make code more readable and self-documenting.

```go
// ✅ GOOD
type ProductRegistration struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    SKU         string    `json:"sku"`
    CreatedAt   time.Time `json:"created_at"`
}

type Service interface {
    CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductRegistration, *errors.AppError)
}

// ❌ BAD
type productRegistration struct {
    id          int64
    name        string
    sku         string
    created_at  time.Time
}
```

---

## Package Organization

### Domain Package Structure
- **Justification**: Clear separation of concerns makes code more maintainable and testable. Each file has a single responsibility, making it easier to locate and modify specific functionality.

```
internal/domains/productregistration/
├── interfaces.go      # Domain interfaces
├── models.go          # Domain models and DTOs
├── repository.go      # Repository implementation
├── service.go         # Service implementation
├── routes.go          # HTTP routes and handlers
├── repository_test.go # Repository tests
├── service_test.go    # Service tests
└── README.md          # Domain documentation
```

### Required Files in Each Domain
1. **`interfaces.go`** - All domain interfaces
   - **Justification**: Centralizing interfaces makes dependencies clear and enables easy mocking for testing.
2. **`models.go`** - Domain models, DTOs, and requests/responses
   - **Justification**: Separating models from business logic makes data structures reusable and easier to maintain.
3. **`repository.go`** - Data access layer
   - **Justification**: Isolating data access logic makes it easier to change storage backends and test business logic independently.
4. **`service.go`** - Business logic layer
   - **Justification**: Centralizing business logic makes it easier to test and maintain complex business rules.
5. **`routes.go`** - HTTP handlers and routing
   - **Justification**: Separating HTTP concerns from business logic makes the code more modular and testable.
6. **`*_test.go`** - Unit tests for each layer
   - **Justification**: Co-locating tests with implementation makes it easier to maintain and understand test coverage.

---

## Layer Responsibilities

### 1. Models Layer (`models.go`)
**Responsibility**: Define data structures and validation rules
**Justification**: Centralizing data structures with validation rules ensures consistency across the application and makes data validation automatic and reliable.

```go
// Domain Entity
type ProductRegistration struct {
    ID          int64     `json:"id" db:"id"`
    Name        string    `json:"name" db:"name" validate:"required,min=1,max=255"`
    SKU         string    `json:"sku" db:"sku" validate:"required,min=1,max=50"`
    Price       float64   `json:"price" db:"price" validate:"required,min=0"`
    Stock       int       `json:"stock" db:"stock" validate:"min=0"`
    IsActive    bool      `json:"is_active" db:"is_active"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Request DTOs
type CreateProductRequest struct {
    Name        string  `json:"name" validate:"required,min=1,max=255"`
    Description string  `json:"description" validate:"max=1000"`
    Category    string  `json:"category" validate:"required,min=1,max=100"`
    Price       float64 `json:"price" validate:"required,min=0"`
    SKU         string  `json:"sku" validate:"required,min=1,max=50"`
    Stock       int     `json:"stock" validate:"min=0"`
    IsActive    bool    `json:"is_active"`
}

// Response DTOs
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

### 2. Repository Layer (`repository.go`)
**Responsibility**: Data access and database operations
**Justification**: Isolating data access logic makes it easier to change storage backends, test business logic independently, and maintain database-specific optimizations.

```go
// Repository Interface
type Repository interface {
    Create(ctx context.Context, product *ProductRegistration) (*ProductRegistration, *errors.AppError)
    GetByID(ctx context.Context, id int64) (*ProductRegistration, *errors.AppError)
    Update(ctx context.Context, id int64, product *ProductRegistration) (*ProductRegistration, *errors.AppError)
    Delete(ctx context.Context, id int64) *errors.AppError
    List(ctx context.Context, req *ProductListRequest) ([]*ProductRegistration, int64, *errors.AppError)
    GetBySKU(ctx context.Context, sku string) (*ProductRegistration, *errors.AppError)
    Exists(ctx context.Context, id int64) (bool, *errors.AppError)
    SKUExists(ctx context.Context, sku string, excludeID *int64) (bool, *errors.AppError)
}

// Repository Implementation
type ProductRepository struct {
    db     interfaces.Database
    logger interfaces.Logger
}

func (r *ProductRepository) Create(ctx context.Context, product *ProductRegistration) (*ProductRegistration, *errors.AppError) {
    // Database operations
    // Error handling with custom errors
    // Logging
}
```

### 3. Service Layer (`service.go`)
**Responsibility**: Business logic and orchestration
**Justification**: Centralizing business logic makes it easier to test, maintain, and ensure business rules are consistently applied across the application.

```go
// Service Interface
type Service interface {
    CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductRegistration, *errors.AppError)
    GetProduct(ctx context.Context, id int64) (*ProductRegistration, *errors.AppError)
    UpdateProduct(ctx context.Context, id int64, req *UpdateProductRequest) (*ProductRegistration, *errors.AppError)
    DeleteProduct(ctx context.Context, id int64) *errors.AppError
    ListProducts(ctx context.Context, req *ProductListRequest) (*ProductListResponse, *errors.AppError)
}

// Service Implementation
type ProductService struct {
    repo   Repository
    logger interfaces.Logger
}

func (s *ProductService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductRegistration, *errors.AppError) {
    // Business logic validation
    // Call repository
    // Error handling
    // Logging
}
```

### 4. Routes Layer (`routes.go`)
**Responsibility**: HTTP handling and request/response mapping
**Justification**: Separating HTTP concerns from business logic makes the code more modular, testable, and allows for easy changes to API structure without affecting business logic.

```go
func RegisterRoutes(router *gin.RouterGroup) {
    productGroup := router.Group("/products")
    {
        productGroup.POST("", createProductHandler)
        productGroup.GET("", listProductsHandler)
        productGroup.GET("/:id", getProductHandler)
        productGroup.PUT("/:id", updateProductHandler)
        productGroup.DELETE("/:id", deleteProductHandler)
    }
}

func createProductHandler(c *gin.Context) {
    // Request validation
    // Call service
    // Response formatting
    // Error handling
}
```

---

## Error Handling Standards

### 1. Use Custom Errors in Business Layer
**Justification**: Custom errors provide structured, consistent error responses with proper HTTP status codes and rich context, making debugging easier and providing better user experience.
```go
// ✅ GOOD - Business layer returns custom errors
func (s *ProductService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductRegistration, *errors.AppError) {
    if exists {
        return nil, errors.NewWithDetails(errors.ErrCodeProductSKUExists, "Product SKU already exists", 
            fmt.Sprintf("Product with SKU '%s' already exists", req.SKU), http.StatusConflict).WithField("sku", req.SKU)
    }
    return product, nil
}

// ❌ BAD - Don't return standard errors in business layer
func (s *ProductService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductRegistration, error) {
    if exists {
        return nil, fmt.Errorf("product SKU already exists")
    }
    return product, nil
}
```

### 2. Error Handling Patterns
**Justification**: Consistent error handling patterns across layers ensure predictable behavior, easier debugging, and maintainable code.
```go
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

// Service Layer - Business logic errors
func (s *ProductService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductRegistration, *errors.AppError) {
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

// Handler Layer - HTTP error responses
func createProductHandler(c *gin.Context) {
    product, appErr := productService.CreateProduct(ctx, &req)
    if appErr != nil {
        middleware.HandleAppError(c, appErr)
        return
    }
    c.JSON(http.StatusCreated, ProductResponse{Product: *product})
}
```

---

## Model Design Guidelines

### 1. Struct Tags
**Justification**: Proper struct tags ensure automatic validation, correct JSON serialization, and proper database mapping, reducing boilerplate code and potential errors.
```go
type ProductRegistration struct {
    ID          int64     `json:"id" db:"id" validate:"required"`
    Name        string    `json:"name" db:"name" validate:"required,min=1,max=255"`
    SKU         string    `json:"sku" db:"sku" validate:"required,min=1,max=50"`
    Price       float64   `json:"price" db:"price" validate:"required,min=0"`
    Stock       int       `json:"stock" db:"stock" validate:"min=0"`
    IsActive    bool      `json:"is_active" db:"is_active"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
```

### 2. Validation Rules
**Justification**: Centralized validation rules ensure data integrity and provide consistent error messages, making the API more reliable and user-friendly.
```go
// Required fields
validate:"required"

// String length
validate:"min=1,max=255"

// Numeric ranges
validate:"min=0,max=100"

// Email format
validate:"email"

// Custom validation
validate:"required,min=1,max=50,alphanum"
```

### 3. JSON Naming
- Use `snake_case` for JSON fields
  - **Justification**: Snake case is the standard for JSON APIs and provides better readability for API consumers.
- Use `camelCase` for Go struct fields
  - **Justification**: Follows Go naming conventions and makes the code consistent with Go standards.
- Use `db` tags for database column mapping
  - **Justification**: Explicit database mapping prevents errors and makes the relationship between struct fields and database columns clear.

---

## Constants and Configuration

### 1. Constants (`pkg/constants/`)
**Justification**: Centralizing constants prevents magic numbers, makes configuration changes easier, and ensures consistency across the application.
```go
// http.go
package constants

const (
    DefaultPageSize = 10
    MaxPageSize     = 100
    DefaultTimeout  = 30 * time.Second
)

// error_codes.go
package constants

const (
    ErrorCodeProductNotFound = "PRODUCT_NOT_FOUND"
    ErrorCodeProductSKUExists = "PRODUCT_SKU_EXISTS"
)
```

### 2. Configuration (`pkg/config/`)
**Justification**: Structured configuration management makes environment-specific settings easy to manage and prevents configuration-related errors.
```go
type Config struct {
    Server   ServerConfig   `json:"server"`
    Database DatabaseConfig `json:"database"`
    Logger   LoggerConfig   `json:"logger"`
}

type ServerConfig struct {
    Port         int           `json:"port" validate:"required,min=1,max=65535"`
    ReadTimeout  time.Duration `json:"read_timeout" validate:"required"`
    WriteTimeout time.Duration `json:"write_timeout" validate:"required"`
}
```

---

## API Design Standards

### 1. RESTful Endpoints
**Justification**: RESTful design makes APIs intuitive, predictable, and follows industry standards, making it easier for developers to understand and use the API.
```go
// ✅ GOOD
GET    /api/v1/products           # List products
GET    /api/v1/products/:id       # Get product by ID
POST   /api/v1/products           # Create product
PUT    /api/v1/products/:id       # Update product
DELETE /api/v1/products/:id       # Delete product

// ❌ BAD
GET    /api/v1/getProducts
POST   /api/v1/createProduct
PUT    /api/v1/updateProduct/:id
```

### 2. HTTP Status Codes
**Justification**: Proper HTTP status codes help clients understand the nature of errors and handle them appropriately, improving the overall API experience.
```go
// Success responses
200 OK          # Successful GET, PUT
201 Created     # Successful POST
204 No Content  # Successful DELETE

// Client errors
400 Bad Request         # Invalid request data
401 Unauthorized        # Authentication required
403 Forbidden          # Insufficient permissions
404 Not Found          # Resource not found
409 Conflict           # Business rule violation
422 Unprocessable Entity # Validation error

// Server errors
500 Internal Server Error # System error
503 Service Unavailable   # Service down
```

### 3. Request/Response Format
**Justification**: Consistent request/response formats make the API predictable and easier to use, reducing integration time and potential errors.
```go
// Request
type CreateProductRequest struct {
    Name        string  `json:"name" validate:"required,min=1,max=255"`
    Description string  `json:"description" validate:"max=1000"`
    Category    string  `json:"category" validate:"required,min=1,max=100"`
    Price       float64 `json:"price" validate:"required,min=0"`
    SKU         string  `json:"sku" validate:"required,min=1,max=50"`
    Stock       int     `json:"stock" validate:"min=0"`
    IsActive    bool    `json:"is_active"`
}

// Response
type ProductResponse struct {
    Product ProductRegistration `json:"product"`
}

// Error Response
{
    "error": "PRODUCT_NOT_FOUND",
    "message": "Product not found",
    "details": "Product with ID 123 not found",
    "fields": {
        "product_id": 123
    }
}
```

---

## Testing Standards

### 1. Test File Naming
**Justification**: Consistent test file naming makes it easy to locate tests and ensures they are automatically discovered by Go's testing framework.
```
product_service.go      → product_service_test.go
user_repository.go      → user_repository_test.go
auth_middleware.go      → auth_middleware_test.go
```

### 2. Test Structure
**Justification**: Table-driven tests make it easy to test multiple scenarios, reduce code duplication, and ensure comprehensive test coverage.
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
            wantErr: errors.NewWithDetails(errors.ErrCodeProductSKUExists, "Product SKU already exists", "Product with SKU 'EXISTING-SKU' already exists", http.StatusConflict),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### 3. Mock Generation
**Justification**: Generated mocks ensure consistency, reduce maintenance overhead, and provide type-safe mocking for testing.
```go
//go:generate mockgen -source=interfaces.go -destination=mocks/mock_interfaces.go

type MockProductService struct {
    ctrl     *gomock.Controller
    recorder *MockProductServiceMockRecorder
}

func (m *MockProductService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductRegistration, *errors.AppError) {
    // Mock implementation
}
```

---

## Code Review Checklist

### ✅ Pre-Submit Checklist
- [ ] All files follow naming conventions
- [ ] Package structure is correct
- [ ] Error handling uses custom errors in business layer
- [ ] Models have proper validation tags
- [ ] Repository methods return `*errors.AppError`
- [ ] Service methods return `*errors.AppError`
- [ ] Handlers use `middleware.HandleAppError`
- [ ] All public functions have documentation
- [ ] Unit tests cover all business logic
- [ ] No hardcoded values (use constants)
- [ ] Proper logging with context
- [ ] Database queries use transactions where needed
- [ ] Input validation is implemented
- [ ] Error responses are consistent

### ✅ Code Quality
- [ ] No unused imports
- [ ] No unused variables
- [ ] Functions are focused and single-purpose
- [ ] No magic numbers
- [ ] Proper error wrapping
- [ ] Context is passed correctly
- [ ] No race conditions
- [ ] Memory leaks avoided

### ✅ Security
- [ ] Input validation implemented
- [ ] SQL injection prevention
- [ ] No sensitive data in logs
- [ ] Proper authentication/authorization
- [ ] Rate limiting considered
- [ ] CORS configured correctly

---

## Quick Reference

### Error Handling
```go
// Business layer - return custom errors
func (s *Service) Method() (*Model, *errors.AppError)

// Infrastructure layer - return standard errors
func (r *Repository) Method() error

// Handler layer - use middleware
middleware.HandleAppError(c, appErr)
```

### Model Validation
```go
type Model struct {
    Field string `json:"field" db:"field" validate:"required,min=1,max=255"`
}
```

### Database Operations
```go
// Use transactions for multiple operations
err := r.db.WithTransaction(ctx, func(tx *sql.Tx) error {
    // Multiple database operations
    return nil
})
```

### Logging
```go
s.logger.Info(ctx, "Operation started", interfaces.Fields{
    "user_id": userID,
    "action":  "create_product",
})
```

This guide should be followed by all team members to ensure consistency and reduce code review comments. For questions or clarifications, refer to the existing codebase or ask the team lead.

