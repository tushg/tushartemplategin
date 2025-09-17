package productregistration

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"tushartemplategin/pkg/errors"
	"tushartemplategin/pkg/middleware"
)

// RegisterRoutes registers all product registration-related routes to the given router group
func RegisterRoutes(router *gin.RouterGroup) {
	// Create a product registration group under the main API group
	// This will create routes like /api/v1/products, /api/v1/products/:id, etc.
	productGroup := router.Group("/products")
	{
		// Register product endpoints with their handlers
		// Each endpoint is clearly defined and easy to maintain

		// POST /products - Create a new product
		productGroup.POST("", createProductHandler)

		// GET /products - List all products with pagination and filtering
		productGroup.GET("", listProductsHandler)

		// GET /products/:id - Get a specific product by ID
		productGroup.GET("/:id", getProductHandler)

		// PUT /products/:id - Update a specific product
		productGroup.PUT("/:id", updateProductHandler)

		// DELETE /products/:id - Delete a specific product
		productGroup.DELETE("/:id", deleteProductHandler)

		// GET /products/sku/:sku - Get a product by SKU
		productGroup.GET("/sku/:sku", getProductBySKUHandler)

		// PATCH /products/:id/stock - Update product stock
		productGroup.PATCH("/:id/stock", updateStockHandler)
	}
}

// createProductHandler handles product creation requests
func createProductHandler(c *gin.Context) {
	// Get the product service from the context
	productService := c.MustGet("productService").(Service)

	ctx := c.Request.Context()

	// Parse request body
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleAppError(c, errors.NewWithDetails(errors.ErrCodeBadRequest, "Invalid request body", err.Error(), http.StatusBadRequest))
		return
	}

	// Create product through service layer
	product, err := productService.CreateProduct(ctx, &req)
	if err != nil {
		// Check if it's already an AppError
		if appErr := errors.GetAppError(err); appErr != nil {
			middleware.HandleAppError(c, appErr)
		} else {
			middleware.HandleAppError(c, errors.NewWithError(errors.ErrCodeInternalServer, "Failed to create product", http.StatusInternalServerError, err))
		}
		return
	}

	// Return created product with 201 Created
	c.JSON(http.StatusCreated, ProductResponse{Product: *product})
}

// listProductsHandler handles product listing requests
func listProductsHandler(c *gin.Context) {
	// Get the product service from the context
	productService := c.MustGet("productService").(Service)

	ctx := c.Request.Context()

	// Parse query parameters
	var req ProductListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.HandleAppError(c, errors.NewWithDetails(errors.ErrCodeBadRequest, "Invalid query parameters", err.Error(), http.StatusBadRequest))
		return
	}

	// List products through service layer
	response, err := productService.ListProducts(ctx, &req)
	if err != nil {
		if appErr := errors.GetAppError(err); appErr != nil {
			middleware.HandleAppError(c, appErr)
		} else {
			middleware.HandleAppError(c, errors.NewWithError(errors.ErrCodeInternalServer, "Failed to list products", http.StatusInternalServerError, err))
		}
		return
	}

	// Return product list with 200 OK
	c.JSON(http.StatusOK, response)
}

// getProductHandler handles getting a specific product by ID
func getProductHandler(c *gin.Context) {
	// Get the product service from the context
	productService := c.MustGet("productService").(Service)

	ctx := c.Request.Context()

	// Parse product ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.HandleAppError(c, errors.NewWithDetails(errors.ErrCodeBadRequest, "Invalid product ID", "Product ID must be a valid integer", http.StatusBadRequest))
		return
	}

	// Get product through service layer
	product, err := productService.GetProduct(ctx, id)
	if err != nil {
		if appErr := errors.GetAppError(err); appErr != nil {
			middleware.HandleAppError(c, appErr)
		} else {
			// Check if it's a "not found" error by string matching (for backward compatibility)
			if err.Error() == "product with id "+idStr+" not found" {
				middleware.HandleAppError(c, errors.NewWithDetails(errors.ErrCodeProductNotFound, "Product not found", fmt.Sprintf("Product with ID %v not found", id), http.StatusNotFound).WithField("product_id", id))
			} else {
				middleware.HandleAppError(c, errors.NewWithError(errors.ErrCodeInternalServer, "Failed to get product", http.StatusInternalServerError, err))
			}
		}
		return
	}

	// Return product with 200 OK
	c.JSON(http.StatusOK, ProductResponse{Product: *product})
}

// updateProductHandler handles product update requests
func updateProductHandler(c *gin.Context) {
	// Get the product service from the context
	productService := c.MustGet("productService").(Service)

	ctx := c.Request.Context()

	// Parse product ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.HandleAppError(c, errors.NewWithDetails(errors.ErrCodeBadRequest, "Invalid product ID", "Product ID must be a valid integer", http.StatusBadRequest))
		return
	}

	// Parse request body
	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleAppError(c, errors.NewWithDetails(errors.ErrCodeBadRequest, "Invalid request body", err.Error(), http.StatusBadRequest))
		return
	}

	// Update product through service layer
	product, err := productService.UpdateProduct(ctx, id, &req)
	if err != nil {
		if appErr := errors.GetAppError(err); appErr != nil {
			middleware.HandleAppError(c, appErr)
		} else {
			// Check if it's a "not found" error
			if err.Error() == "product with id "+idStr+" not found" {
				middleware.HandleAppError(c, errors.NewWithDetails(errors.ErrCodeProductNotFound, "Product not found", fmt.Sprintf("Product with ID %v not found", id), http.StatusNotFound).WithField("product_id", id))
			} else {
				middleware.HandleAppError(c, errors.NewWithError(errors.ErrCodeInternalServer, "Failed to update product", http.StatusInternalServerError, err))
			}
		}
		return
	}

	// Return updated product with 200 OK
	c.JSON(http.StatusOK, ProductResponse{Product: *product})
}

// deleteProductHandler handles product deletion requests
func deleteProductHandler(c *gin.Context) {
	// Get the product service from the context
	productService := c.MustGet("productService").(Service)

	ctx := c.Request.Context()

	// Parse product ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.HandleAppError(c, errors.NewWithDetails(errors.ErrCodeBadRequest, "Invalid product ID", "Product ID must be a valid integer", http.StatusBadRequest))
		return
	}

	// Delete product through service layer
	err = productService.DeleteProduct(ctx, id)
	if err != nil {
		if appErr := errors.GetAppError(err); appErr != nil {
			middleware.HandleAppError(c, appErr)
		} else {
			// Check if it's a "not found" error
			if err.Error() == "product with id "+idStr+" not found" {
				middleware.HandleAppError(c, errors.NewWithDetails(errors.ErrCodeProductNotFound, "Product not found", fmt.Sprintf("Product with ID %v not found", id), http.StatusNotFound).WithField("product_id", id))
			} else {
				middleware.HandleAppError(c, errors.NewWithError(errors.ErrCodeInternalServer, "Failed to delete product", http.StatusInternalServerError, err))
			}
		}
		return
	}

	// Return success with 204 No Content
	c.Status(http.StatusNoContent)
}

// getProductBySKUHandler handles getting a product by SKU
func getProductBySKUHandler(c *gin.Context) {
	// Get the product service from the context
	productService := c.MustGet("productService").(Service)

	ctx := c.Request.Context()

	// Parse SKU from URL parameter
	sku := c.Param("sku")
	if sku == "" {
		middleware.HandleAppError(c, errors.NewWithDetails(errors.ErrCodeBadRequest, "SKU parameter is required", "SKU cannot be empty", http.StatusBadRequest))
		return
	}

	// Get product through service layer
	product, err := productService.GetProductBySKU(ctx, sku)
	if err != nil {
		if appErr := errors.GetAppError(err); appErr != nil {
			middleware.HandleAppError(c, appErr)
		} else {
			// Check if it's a "not found" error
			if err.Error() == "product with sku "+sku+" not found" {
				middleware.HandleAppError(c, errors.NewWithDetails(errors.ErrCodeNotFound, "Product not found", "Product with SKU '"+sku+"' not found", http.StatusNotFound).WithField("sku", sku))
			} else {
				middleware.HandleAppError(c, errors.NewWithError(errors.ErrCodeInternalServer, "Failed to get product by SKU", http.StatusInternalServerError, err))
			}
		}
		return
	}

	// Return product with 200 OK
	c.JSON(http.StatusOK, ProductResponse{Product: *product})
}

// updateStockHandler handles product stock update requests
func updateStockHandler(c *gin.Context) {
	// Get the product service from the context
	productService := c.MustGet("productService").(Service)

	ctx := c.Request.Context()

	// Parse product ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.HandleAppError(c, errors.NewWithDetails(errors.ErrCodeBadRequest, "Invalid product ID", "Product ID must be a valid integer", http.StatusBadRequest))
		return
	}

	// Parse request body for stock update
	var req struct {
		Stock int `json:"stock" binding:"required,min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleAppError(c, errors.NewWithDetails(errors.ErrCodeBadRequest, "Invalid request body", err.Error(), http.StatusBadRequest))
		return
	}

	// Validate stock value
	if req.Stock < 0 {
		middleware.HandleAppError(c, errors.NewWithDetails(errors.ErrCodeInvalidStock, "Invalid stock quantity", fmt.Sprintf("Stock quantity %d is invalid", req.Stock), http.StatusUnprocessableEntity).WithField("stock", req.Stock))
		return
	}

	// Update stock through service layer
	err = productService.UpdateStock(ctx, id, req.Stock)
	if err != nil {
		if appErr := errors.GetAppError(err); appErr != nil {
			middleware.HandleAppError(c, appErr)
		} else {
			// Check if it's a "not found" error
			if err.Error() == "product with id "+idStr+" not found" {
				middleware.HandleAppError(c, errors.NewWithDetails(errors.ErrCodeProductNotFound, "Product not found", fmt.Sprintf("Product with ID %v not found", id), http.StatusNotFound).WithField("product_id", id))
			} else {
				middleware.HandleAppError(c, errors.NewWithError(errors.ErrCodeInternalServer, "Failed to update stock", http.StatusInternalServerError, err))
			}
		}
		return
	}

	// Return success with 200 OK
	c.JSON(http.StatusOK, gin.H{
		"message": "Stock updated successfully",
		"id":      id,
		"stock":   req.Stock,
	})
}
