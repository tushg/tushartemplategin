package productregistration

import (
	"context"
)

// Service defines the interface for product registration business logic
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

// Repository defines the interface for product registration data access
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
