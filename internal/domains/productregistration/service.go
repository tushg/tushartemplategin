package productregistration

import (
	"context"
	"fmt"

	"tushartemplategin/pkg/errors"
	"tushartemplategin/pkg/interfaces"
)

// ProductService implements the Service interface for product business logic
type ProductService struct {
	repo   Repository
	logger interfaces.Logger
}

// NewProductService creates a new product service
func NewProductService(repo Repository, log interfaces.Logger) Service {
	return &ProductService{
		repo:   repo,
		logger: log,
	}
}

// CreateProduct creates a new product with business logic validation
func (s *ProductService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductRegistration, error) {
	s.logger.Info(ctx, "Creating new product", interfaces.Fields{
		"name":     req.Name,
		"category": req.Category,
		"sku":      req.SKU,
	})

	// Check if SKU already exists
	exists, err := s.repo.SKUExists(ctx, req.SKU, nil)
	if err != nil {
		s.logger.Error(ctx, "Failed to check SKU existence", interfaces.Fields{
			"error": err.Error(),
			"sku":   req.SKU,
		})
		return nil, errors.NewDatabaseError("check SKU existence", err)
	}

	if exists {
		s.logger.Warn(ctx, "Product creation failed: SKU already exists", interfaces.Fields{
			"sku": req.SKU,
		})
		return nil, errors.NewProductSKUExists(req.SKU)
	}

	// Create product entity
	product := &ProductRegistration{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Price:       req.Price,
		SKU:         req.SKU,
		Stock:       req.Stock,
		IsActive:    req.IsActive,
	}

	// Save to repository
	createdProduct, err := s.repo.Create(ctx, product)
	if err != nil {
		s.logger.Error(ctx, "Failed to create product in repository", interfaces.Fields{
			"error": err.Error(),
			"sku":   req.SKU,
		})
		return nil, errors.NewDatabaseError("create product", err)
	}

	s.logger.Info(ctx, "Product created successfully", interfaces.Fields{
		"id":  createdProduct.ID,
		"sku": createdProduct.SKU,
	})

	return createdProduct, nil
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(ctx context.Context, id int64) (*ProductRegistration, error) {
	s.logger.Info(ctx, "Getting product by ID", interfaces.Fields{
		"id": id,
	})

	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to get product", interfaces.Fields{
			"error": err.Error(),
			"id":    id,
		})
		return nil, errors.NewDatabaseError("get product by ID", err)
	}

	if product == nil {
		s.logger.Warn(ctx, "Product not found", interfaces.Fields{
			"id": id,
		})
		return nil, errors.NewProductNotFound(id)
	}

	s.logger.Info(ctx, "Product retrieved successfully", interfaces.Fields{
		"id":  product.ID,
		"sku": product.SKU,
	})

	return product, nil
}

// UpdateProduct updates an existing product with business logic validation
func (s *ProductService) UpdateProduct(ctx context.Context, id int64, req *UpdateProductRequest) (*ProductRegistration, error) {
	s.logger.Info(ctx, "Updating product", interfaces.Fields{
		"id": id,
	})

	// Check if product exists
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to check product existence", interfaces.Fields{
			"error": err.Error(),
			"id":    id,
		})
		return nil, fmt.Errorf("failed to validate product: %w", err)
	}

	if !exists {
		s.logger.Warn(ctx, "Product update failed: product not found", interfaces.Fields{
			"id": id,
		})
		return nil, fmt.Errorf("product with id %d not found", id)
	}

	// Get existing product
	existingProduct, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to get existing product", interfaces.Fields{
			"error": err.Error(),
			"id":    id,
		})
		return nil, fmt.Errorf("failed to get existing product: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		existingProduct.Name = *req.Name
	}
	if req.Description != nil {
		existingProduct.Description = *req.Description
	}
	if req.Category != nil {
		existingProduct.Category = *req.Category
	}
	if req.Price != nil {
		existingProduct.Price = *req.Price
	}
	if req.SKU != nil {
		// Check if new SKU already exists (excluding current product)
		exists, err := s.repo.SKUExists(ctx, *req.SKU, &id)
		if err != nil {
			s.logger.Error(ctx, "Failed to check SKU existence", interfaces.Fields{
				"error": err.Error(),
				"sku":   *req.SKU,
			})
			return nil, fmt.Errorf("failed to validate SKU: %w", err)
		}

		if exists {
			s.logger.Warn(ctx, "Product update failed: SKU already exists", interfaces.Fields{
				"sku": *req.SKU,
			})
			return nil, fmt.Errorf("product with SKU %s already exists", *req.SKU)
		}
		existingProduct.SKU = *req.SKU
	}
	if req.Stock != nil {
		existingProduct.Stock = *req.Stock
	}
	if req.IsActive != nil {
		existingProduct.IsActive = *req.IsActive
	}

	// Save updated product
	updatedProduct, err := s.repo.Update(ctx, id, existingProduct)
	if err != nil {
		s.logger.Error(ctx, "Failed to update product in repository", interfaces.Fields{
			"error": err.Error(),
			"id":    id,
		})
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	s.logger.Info(ctx, "Product updated successfully", interfaces.Fields{
		"id":  updatedProduct.ID,
		"sku": updatedProduct.SKU,
	})

	return updatedProduct, nil
}

// DeleteProduct removes a product
func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	s.logger.Info(ctx, "Deleting product", interfaces.Fields{
		"id": id,
	})

	// Check if product exists
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to check product existence", interfaces.Fields{
			"error": err.Error(),
			"id":    id,
		})
		return fmt.Errorf("failed to validate product: %w", err)
	}

	if !exists {
		s.logger.Warn(ctx, "Product deletion failed: product not found", interfaces.Fields{
			"id": id,
		})
		return fmt.Errorf("product with id %d not found", id)
	}

	// Delete from repository
	err = s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to delete product from repository", interfaces.Fields{
			"error": err.Error(),
			"id":    id,
		})
		return fmt.Errorf("failed to delete product: %w", err)
	}

	s.logger.Info(ctx, "Product deleted successfully", interfaces.Fields{
		"id": id,
	})

	return nil
}

// ListProducts retrieves a list of products with pagination and filtering
func (s *ProductService) ListProducts(ctx context.Context, req *ProductListRequest) (*ProductListResponse, error) {
	s.logger.Info(ctx, "Listing products", interfaces.Fields{
		"page":     req.Page,
		"limit":    req.Limit,
		"category": req.Category,
		"search":   req.Search,
	})

	products, total, err := s.repo.List(ctx, req)
	if err != nil {
		s.logger.Error(ctx, "Failed to list products", interfaces.Fields{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	// Convert to response format
	productList := make([]ProductRegistration, len(products))
	for i, product := range products {
		productList[i] = *product
	}

	response := &ProductListResponse{
		Products: productList,
		Total:    total,
		Page:     req.Page,
		Limit:    req.Limit,
	}

	s.logger.Info(ctx, "Products listed successfully", interfaces.Fields{
		"count": len(products),
		"total": total,
	})

	return response, nil
}

// GetProductBySKU retrieves a product by its SKU
func (s *ProductService) GetProductBySKU(ctx context.Context, sku string) (*ProductRegistration, error) {
	s.logger.Info(ctx, "Getting product by SKU", interfaces.Fields{
		"sku": sku,
	})

	product, err := s.repo.GetBySKU(ctx, sku)
	if err != nil {
		s.logger.Error(ctx, "Failed to get product by SKU", interfaces.Fields{
			"error": err.Error(),
			"sku":   sku,
		})
		return nil, err
	}

	s.logger.Info(ctx, "Product retrieved by SKU successfully", interfaces.Fields{
		"id":  product.ID,
		"sku": product.SKU,
	})

	return product, nil
}

// UpdateStock updates the stock quantity of a product
func (s *ProductService) UpdateStock(ctx context.Context, id int64, stock int) error {
	s.logger.Info(ctx, "Updating product stock", interfaces.Fields{
		"id":    id,
		"stock": stock,
	})

	// Check if product exists
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to check product existence", interfaces.Fields{
			"error": err.Error(),
			"id":    id,
		})
		return fmt.Errorf("failed to validate product: %w", err)
	}

	if !exists {
		s.logger.Warn(ctx, "Stock update failed: product not found", interfaces.Fields{
			"id": id,
		})
		return fmt.Errorf("product with id %d not found", id)
	}

	// Validate stock quantity
	if stock < 0 {
		s.logger.Warn(ctx, "Stock update failed: invalid stock quantity", interfaces.Fields{
			"id":    id,
			"stock": stock,
		})
		return fmt.Errorf("stock quantity cannot be negative")
	}

	// Update stock
	err = s.repo.UpdateStock(ctx, id, stock)
	if err != nil {
		s.logger.Error(ctx, "Failed to update product stock in repository", interfaces.Fields{
			"error": err.Error(),
			"id":    id,
		})
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	s.logger.Info(ctx, "Product stock updated successfully", interfaces.Fields{
		"id":    id,
		"stock": stock,
	})

	return nil
}
