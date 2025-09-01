# ProductRegistration Domain

This domain handles product registration and management functionality including CRUD operations for products.

## Features

- Create new products
- Retrieve products by ID or SKU
- Update existing products
- Delete products
- List products with pagination and filtering
- Update product stock

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/products` | Create a new product |
| GET | `/products` | List products with pagination/filtering |
| GET | `/products/:id` | Get product by ID |
| PUT | `/products/:id` | Update product |
| DELETE | `/products/:id` | Delete product |
| GET | `/products/sku/:sku` | Get product by SKU |
| PATCH | `/products/:id/stock` | Update product stock |

## Data Model

### ProductRegistration
```go
type ProductRegistration struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Category    string    `json:"category"`
    Price       float64   `json:"price"`
    SKU         string    `json:"sku"`
    Stock       int       `json:"stock"`
    IsActive    bool      `json:"is_active"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

## Database Schema

```sql
CREATE TABLE products (
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
```

## Repository & Transactions

- Uses `interfaces.Database` and wraps queries in transactions via `WithTransaction(ctx, func(tx *sql.Tx) error { ... })`.
- Ensures consistent error handling and logging per operation.

## Wiring in main

Pass `db` to `setupDomainsAndMiddleware` and wire the domain:

```go
router = setupDomainsAndMiddleware(router, appLogger, db)

func setupDomainsAndMiddleware(router *gin.Engine, appLogger logger.Logger, db interfaces.Database) *gin.Engine {
    productRepo := productregistration.NewProductRepository(db, appLogger)
    productService := productregistration.NewProductService(productRepo, appLogger)
    router.Use(func(c *gin.Context) {
        c.Set("productService", productService)
        c.Next()
    })
    return router
}
```

## Business Rules

1. **SKU Uniqueness**: Each product must have a unique SKU
2. **Price Validation**: Product price must be >= 0
3. **Stock Validation**: Stock quantity must be >= 0
4. **Required Fields**: Name, category, price, and SKU are required
5. **Soft Delete**: Products can be deactivated instead of hard deleted

## Testing

Run the tests with:
```bash
go test ./internal/domains/productregistration/...
```

## Dependencies

- Database connection (PostgreSQL)
- Logger interface
- Gin framework for HTTP handling
