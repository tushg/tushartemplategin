package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"tushartemplategin/pkg/errors"
	"tushartemplategin/pkg/middleware"
)

// Example demonstrating the new error handling system
func main() {
	// Create a simple Gin router for demonstration
	router := gin.New()

	// Add error handling middleware
	router.Use(middleware.ErrorHandlerMiddleware(nil)) // In real app, pass logger

	// Example route that demonstrates different error types
	router.GET("/example/:id", exampleHandler)

	// Start server
	fmt.Println("Starting error handling example server on :8080")
	router.Run(":8080")
}

func exampleHandler(c *gin.Context) {
	id := c.Param("id")

	// Example 1: Missing parameter error
	if id == "" {
		middleware.HandleAppError(c, errors.NewWithDetails(errors.ErrCodeBadRequest, "Missing ID parameter", "ID is required", http.StatusBadRequest))
		return
	}

	// Example 2: Business logic error
	if id == "notfound" {
		middleware.HandleAppError(c, errors.NewWithDetails(errors.ErrCodeProductNotFound, "Product not found", fmt.Sprintf("Product with ID %v not found", 123), http.StatusNotFound).WithField("product_id", 123))
		return
	}

	// Example 3: Conflict error
	if id == "conflict" {
		middleware.HandleAppError(c, errors.NewWithDetails(errors.ErrCodeProductSKUExists, "Product SKU already exists", fmt.Sprintf("Product with SKU '%s' already exists", "CONFLICT-SKU"), http.StatusConflict).WithField("sku", "CONFLICT-SKU"))
		return
	}

	// Example 4: Internal server error
	if id == "error" {
		middleware.HandleAppError(c, errors.NewWithError(errors.ErrCodeInternalServer, "Something went wrong", http.StatusInternalServerError,
			fmt.Errorf("database connection failed")))
		return
	}

	// Example 5: Success case
	c.JSON(http.StatusOK, gin.H{
		"message": "Success!",
		"id":      id,
		"status":  "ok",
	})
}

// Example of using errors in service layer
func exampleServiceMethod(ctx context.Context, id string) error {
	// Simulate different error scenarios
	switch id {
	case "db_error":
		return errors.NewWithError(errors.ErrCodeDatabaseQuery, "Database operation failed", http.StatusInternalServerError, fmt.Errorf("connection timeout")).WithField("operation", "query products")
	case "not_found":
		return errors.NewWithDetails(errors.ErrCodeProductNotFound, "Product not found", fmt.Sprintf("Product with ID %v not found", 123), http.StatusNotFound).WithField("product_id", 123)
	case "conflict":
		return errors.NewWithDetails(errors.ErrCodeProductSKUExists, "Product SKU already exists", fmt.Sprintf("Product with SKU '%s' already exists", "EXISTING-SKU"), http.StatusConflict).WithField("sku", "EXISTING-SKU")
	default:
		return nil
	}
}

// Example of error handling in service layer
func handleServiceError(err error) *errors.AppError {
	if err == nil {
		return nil
	}

	// Check if it's already an AppError
	if appErr := errors.GetAppError(err); appErr != nil {
		return appErr
	}

	// Wrap standard error
	return errors.NewWithError(errors.ErrCodeInternalServer, "Service operation failed", http.StatusInternalServerError, err)
}
