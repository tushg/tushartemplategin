# Adding a New Domain - Step-by-Step Guide

This guide provides a comprehensive walkthrough for adding a new domain with CRUD REST endpoints to the Tushar Template Gin project. We'll use the ProductRegistration domain as an example.

## Table of Contents
1. [Overview](#overview)
2. [Project Structure](#project-structure)
3. [Step-by-Step Implementation](#step-by-step-implementation)
4. [Database Setup](#database-setup)
5. [Testing](#testing)
6. [Best Practices](#best-practices)

## Overview

The project follows a clean architecture pattern with the following layers:
- **Models**: Data structures and request/response DTOs
- **Interfaces**: Service and repository contracts
- **Repository**: Data access layer (database operations)
- **Service**: Business logic layer
- **Routes**: HTTP handlers and endpoint definitions

## Project Structure

```
internal/
├── domains/
│   └── [domain-name]/
│       ├── models.go          # Data models and DTOs
│       ├── interfaces.go      # Service and repository interfaces
│       ├── repository.go      # Data access implementation
│       ├── service.go         # Business logic implementation
│       └── routes.go          # HTTP handlers and routing
├── health/                    # Example domain (health checks)
└── ...
```

## Step-by-Step Implementation

### Step 1: Create Domain Directory Structure

Create a new directory for your domain under `internal/domains/`:

```bash
mkdir -p internal/domains/[domain-name]
```

For our example:
```bash
mkdir -p internal/domains/productregistration
```

### Step 2: Define Models and DTOs (`models.go`)

Create the `models.go` file to define your data structures:

```go
package productregistration

import (
	"time"
)

// Main entity structure
type ProductRegistration struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" binding:"required"`
	Description string    `json:"description" db:"description"`
	Category    string    `json:"category" db:"category" binding:"required"`
	Price       float64   `json:"price" db:"price" binding:"required,min=0"`
	SKU         string    `json:"sku" db:"sku" binding:"required"`
	Stock       int       `json:"stock" db:"stock" binding:"required,min=0"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Request DTOs
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=255"`
	Description string  `json:"description" binding:"max=1000"`
	Category    string  `json:"category" binding:"required,min=1,max=100"`
	Price       float64 `json:"price" binding:"required,min=0"`
	SKU         string  `json:"sku" binding:"required,min=1,max=50"`
	Stock       int     `json:"stock" binding:"required,min=0"`
	IsActive    bool    `json:"is_active"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty" binding:"omitempty,min=1,max=255"`
	Description *string  `json:"description,omitempty" binding:"omitempty,max=1000"`
	Category    *string  `json:"category,omitempty" binding:"omitempty,min=1,max=100"`
	Price       *float64 `json:"price,omitempty" binding:"omitempty,min=0"`
	SKU         *string  `json:"sku,omitempty" binding:"omitempty,min=1,max=50"`
	Stock       *int     `json:"stock,omitempty" binding:"omitempty,min=0"`
	IsActive    *bool    `json:"is_active,omitempty"`
}

// Response DTOs
type ProductListResponse struct {
	Products []ProductRegistration `json:"products"`
	Total    int64                 `json:"total"`
	Page     int                   `json:"page"`
	Limit    int                   `json:"limit"`
}

type ProductListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	Limit    int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Category string `form:"category" binding:"omitempty"`
	IsActive *bool  `form:"is_active" binding:"omitempty"`
	Search   string `form:"search" binding:"omitempty"`
}

type ProductResponse struct {
	Product ProductRegistration `json:"product"`
}
```

**Key Points:**
- Use struct tags for JSON serialization (`json:"field_name"`)
- Use struct tags for database mapping (`db:"field_name"`)
- Use Gin binding tags for validation (`binding:"required,min=1"`)
- Use pointers in update requests to distinguish between zero values and omitted fields

### Step 3: Define Interfaces (`interfaces.go`)

Create the `interfaces.go` file to define your service and repository contracts:

```go
package productregistration

import (
	"context"
)

// Service defines the interface for business logic
type Service interface {
	// CRUD operations
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductRegistration, error)
	GetProduct(ctx context.Context, id int64) (*ProductRegistration, error)
	UpdateProduct(ctx context.Context, id int64, req *UpdateProductRequest) (*ProductRegistration, error)
	DeleteProduct(ctx context.Context, id int64) error
	ListProducts(ctx context.Context, req *ProductListRequest) (*ProductListResponse, error)
	
	// Additional business operations
	GetProductBySKU(ctx context.Context, sku string) (*ProductRegistration, error)
	UpdateStock(ctx context.Context, id int64, stock int) error
}

// Repository defines the interface for data access
type Repository interface {
	// CRUD operations
	Create(ctx context.Context, product *ProductRegistration) (*ProductRegistration, error)
	GetByID(ctx context.Context, id int64) (*ProductRegistration, error)
	Update(ctx context.Context, id int64, product *ProductRegistration) (*ProductRegistration, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, req *ProductListRequest) ([]*ProductRegistration, int64, error)
	
	// Additional data operations
	GetBySKU(ctx context.Context, sku string) (*ProductRegistration, error)
	UpdateStock(ctx context.Context, id int64, stock int) error
	Exists(ctx context.Context, id int64) (bool, error)
	SKUExists(ctx context.Context, sku string, excludeID *int64) (bool, error)
}
```

**Key Points:**
- Always include `context.Context` as the first parameter
- Use specific error types for better error handling
- Include both basic CRUD and domain-specific operations
- Repository interface focuses on data access, Service interface focuses on business logic

### Step 4: Implement Repository (`repository.go`)

Create the `repository.go` file to implement data access using the app-wide database interface and transactions:

```go
package productregistration

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"tushartemplategin/pkg/interfaces"
)

// ProductRepository implements the Repository interface
type ProductRepository struct {
	db     interfaces.Database     // note: Database, not DBInterface
	logger interfaces.Logger
}

// NewProductRepository creates a new product repository
func NewProductRepository(db interfaces.Database, log interfaces.Logger) Repository {
	return &ProductRepository{db: db, logger: log}
}

// Create creates a new product in the database
func (r *ProductRepository) Create(ctx context.Context, product *ProductRegistration) (*ProductRegistration, error) {
	const query = `
		INSERT INTO products (name, description, category, price, sku, stock, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	if err := r.db.WithTransaction(ctx, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, query,
			product.Name, product.Description, product.Category, product.Price,
			product.SKU, product.Stock, product.IsActive, product.CreatedAt, product.UpdatedAt,
		).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
	}); err != nil {
		r.logger.Error(ctx, "Failed to create product", interfaces.Fields{"error": err.Error(), "sku": product.SKU})
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	return product, nil
}
```

**Key Points:**
- Use `interfaces.Database` and wrap queries in `WithTransaction(ctx, func(tx *sql.Tx) error { ... })`
- Always use parameterized queries to prevent SQL injection
- Include proper error handling and logging
- Use context for cancellation and timeouts

### Step 5: Implement Service (`service.go`)

Create the `service.go` file to implement business logic:

```go
package productregistration

import (
	"context"
	"fmt"

	"tushartemplategin/pkg/interfaces"
)

// ProductService implements the Service interface
type ProductService struct {
	repo   Repository
	logger interfaces.Logger
}

// NewProductService creates a new product service
func NewProductService(repo Repository, log interfaces.Logger) Service {
	return &ProductService{repo: repo, logger: log}
}

// CreateProduct creates a new product with business logic validation
func (s *ProductService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductRegistration, error) {
	// business validation ...
	return s.repo.Create(ctx, &ProductRegistration{/* map fields */})
}
```

**Key Points:**
- Implement business logic validation
- Handle errors appropriately
- Log important operations
- Keep business logic separate from data access

### Step 6: Implement Routes (`routes.go`)

Create the `routes.go` file to implement HTTP handlers:

```go
// RegisterRoutes registers all domain routes
func RegisterRoutes(router *gin.RouterGroup) {
	productGroup := router.Group("/products")
	{
		productGroup.POST("", createProductHandler)
		// ... other routes
	}
}
```

### Step 7: Add Error Constants

Update `pkg/constants/http.go` to include domain-specific error messages:

```go
// Product Registration Error Messages
const (
	ERROR_PRODUCT_NOT_FOUND     = "Product not found"
	ERROR_PRODUCT_SKU_EXISTS    = "Product with this SKU already exists"
	ERROR_PRODUCT_CREATE_FAILED = "Failed to create product"
	ERROR_PRODUCT_UPDATE_FAILED = "Failed to update product"
	ERROR_PRODUCT_DELETE_FAILED = "Failed to delete product"
	ERROR_PRODUCT_LIST_FAILED   = "Failed to list products"
	ERROR_INVALID_STOCK         = "Invalid stock quantity"
)
```

### Step 8: Register Domain in Main Application

Update `cmd/server/main.go` to register your new domain:

1) Add import:
```go
import (
	// ... existing imports
	"tushartemplategin/internal/domains/productregistration"
)
```

2) Pass `db` into setup and adjust signature:
```go
// in main()
router = setupDomainsAndMiddleware(router, appLogger, db)

// signature
func setupDomainsAndMiddleware(router *gin.Engine, appLogger logger.Logger, db interfaces.Database) *gin.Engine {
    // create product repository using db
}
```

3) Inside setup, wire repository and service, and expose service on the Gin context:
```go
productRepo := productregistration.NewProductRepository(db, appLogger)
productService := productregistration.NewProductService(productRepo, appLogger)
router.Use(func(c *gin.Context) { c.Set("productService", productService); c.Next() })
```

4) Register routes:
```go
productregistration.RegisterRoutes(api)
```

## Database Setup

### Step 1: Create Migration Script

Create a migration script in `scripts/migrations/`:

```sql
-- Migration: Create products table
-- Description: Creates the products table for the ProductRegistration domain
-- Version: 001
-- Date: 2024-01-01

CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100) NOT NULL,
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    sku VARCHAR(50) NOT NULL UNIQUE,
    stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);
CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);
CREATE INDEX IF NOT EXISTS idx_products_is_active ON products(is_active);
CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at);

-- Create trigger for automatic updated_at updates
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_products_updated_at 
    BEFORE UPDATE ON products 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
```

### Step 2: Run Migration

Execute the migration script against your database:

```bash
psql -h localhost -U your_username -d your_database -f scripts/migrations/001_create_products_table.sql
```

## Testing

### Step 1: Test the Endpoints

Use curl or a tool like Postman to test your endpoints:

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product",
    "description": "A test product",
    "category": "Electronics",
    "price": 99.99,
    "sku": "TEST-001",
    "stock": 100,
    "is_active": true
  }'
```

## Best Practices

- Use `interfaces.Database` with transactions for DB operations
- Validate input at API layer; enforce business rules in service layer
- Prefer parameterized queries, proper indexing, and pagination
- Keep domain wiring in `setupDomainsAndMiddleware` and pass `db`

This guide provides a comprehensive framework for adding new domains to your application. Follow these patterns consistently to maintain code quality and architecture integrity.
