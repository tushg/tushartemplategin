package productregistration

import (
	"time"
)

// ProductRegistration represents a product registration entity
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

// CreateProductRequest represents the request payload for creating a product
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=255"`
	Description string  `json:"description" binding:"max=1000"`
	Category    string  `json:"category" binding:"required,min=1,max=100"`
	Price       float64 `json:"price" binding:"required,min=0"`
	SKU         string  `json:"sku" binding:"required,min=1,max=50"`
	Stock       int     `json:"stock" binding:"required,min=0"`
	IsActive    bool    `json:"is_active"`
}

// UpdateProductRequest represents the request payload for updating a product
type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty" binding:"omitempty,min=1,max=255"`
	Description *string  `json:"description,omitempty" binding:"omitempty,max=1000"`
	Category    *string  `json:"category,omitempty" binding:"omitempty,min=1,max=100"`
	Price       *float64 `json:"price,omitempty" binding:"omitempty,min=0"`
	SKU         *string  `json:"sku,omitempty" binding:"omitempty,min=1,max=50"`
	Stock       *int     `json:"stock,omitempty" binding:"omitempty,min=0"`
	IsActive    *bool    `json:"is_active,omitempty"`
}

// ProductListResponse represents the response for listing products
type ProductListResponse struct {
	Products []ProductRegistration `json:"products"`
	Total    int64                 `json:"total"`
	Page     int                   `json:"page"`
	Limit    int                   `json:"limit"`
}

// ProductListRequest represents the request parameters for listing products
type ProductListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	Limit    int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Category string `form:"category" binding:"omitempty"`
	IsActive *bool  `form:"is_active" binding:"omitempty"`
	Search   string `form:"search" binding:"omitempty"`
}

// ProductResponse represents the response for a single product
type ProductResponse struct {
	Product ProductRegistration `json:"product"`
}
