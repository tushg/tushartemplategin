package productregistration

import (
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
		middleware.HandleAppError(c, errors.NewBadRequestWithDetails("Invalid request body", err.Error()))
		return
	}

	// Create product through service layer
	product, err := productService.CreateProduct(ctx, &req)
	if err != nil {
		// Check if it's already an AppError
		if appErr := errors.GetAppError(err); appErr != nil {
			middleware.HandleAppError(c, appErr)
		} else {
			middleware.HandleAppError(c, errors.NewInternalServerErrorWithError("Failed to create product", err))
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
		middleware.HandleAppError(c, errors.NewBadRequestWithDetails("Invalid query parameters", err.Error()))
		return
	}

	// List products through service layer
	response, err := productService.ListProducts(ctx, &req)
	if err != nil {
		if appErr := errors.GetAppError(err); appErr != nil {
			middleware.HandleAppError(c, appErr)
		} else {
			middleware.HandleAppError(c, errors.NewInternalServerErrorWithError("Failed to list products", err))
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
		middleware.HandleAppError(c, errors.NewBadRequestWithDetails("Invalid product ID", "Product ID must be a valid integer"))
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
				middleware.HandleAppError(c, errors.NewProductNotFound(id))
			} else {
				middleware.HandleAppError(c, errors.NewInternalServerErrorWithError("Failed to get product", err))
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
		middleware.HandleAppError(c, errors.NewBadRequestWithDetails("Invalid product ID", "Product ID must be a valid integer"))
		return
	}

	// Parse request body
	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleAppError(c, errors.NewBadRequestWithDetails("Invalid request body", err.Error()))
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
				middleware.HandleAppError(c, errors.NewProductNotFound(id))
			} else {
				middleware.HandleAppError(c, errors.NewInternalServerErrorWithError("Failed to update product", err))
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
		middleware.HandleAppError(c, errors.NewBadRequestWithDetails("Invalid product ID", "Product ID must be a valid integer"))
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
				middleware.HandleAppError(c, errors.NewProductNotFound(id))
			} else {
				middleware.HandleAppError(c, errors.NewInternalServerErrorWithError("Failed to delete product", err))
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
		middleware.HandleAppError(c, errors.NewBadRequestWithDetails("SKU parameter is required", "SKU cannot be empty"))
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
				middleware.HandleAppError(c, errors.NewNotFoundWithDetails("Product not found", "Product with SKU '"+sku+"' not found").WithField("sku", sku))
			} else {
				middleware.HandleAppError(c, errors.NewInternalServerErrorWithError("Failed to get product by SKU", err))
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
		middleware.HandleAppError(c, errors.NewBadRequestWithDetails("Invalid product ID", "Product ID must be a valid integer"))
		return
	}

	// Parse request body for stock update
	var req struct {
		Stock int `json:"stock" binding:"required,min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleAppError(c, errors.NewBadRequestWithDetails("Invalid request body", err.Error()))
		return
	}

	// Validate stock value
	if req.Stock < 0 {
		middleware.HandleAppError(c, errors.NewInvalidStock(req.Stock))
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
				middleware.HandleAppError(c, errors.NewProductNotFound(id))
			} else {
				middleware.HandleAppError(c, errors.NewInternalServerErrorWithError("Failed to update stock", err))
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
