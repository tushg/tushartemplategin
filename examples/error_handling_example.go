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

	// Example 1: Validation error
	if id == "" {
		middleware.HandleAppError(c, errors.NewBadRequestWithDetails("Missing ID parameter", "ID is required"))
		return
	}

	// Example 2: Business logic error
	if id == "notfound" {
		middleware.HandleAppError(c, errors.NewProductNotFound(123).WithField("requested_id", id))
		return
	}

	// Example 3: Conflict error
	if id == "conflict" {
		middleware.HandleAppError(c, errors.NewProductSKUExists("CONFLICT-SKU").WithField("requested_id", id))
		return
	}

	// Example 4: Validation error with custom fields
	if id == "invalid" {
		middleware.HandleAppError(c, errors.NewValidationErrorWithDetails("Invalid ID format", "ID must be numeric").
			WithField("provided_id", id).
			WithField("expected_format", "numeric"))
		return
	}

	// Example 5: Internal server error
	if id == "error" {
		middleware.HandleAppError(c, errors.NewInternalServerErrorWithError("Something went wrong",
			fmt.Errorf("database connection failed")))
		return
	}

	// Example 6: Success case
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
		return errors.NewDatabaseError("query products", fmt.Errorf("connection timeout"))
	case "not_found":
		return errors.NewProductNotFound(123)
	case "validation":
		return errors.NewValidationErrorWithDetails("Invalid input", "ID must be positive integer")
	case "conflict":
		return errors.NewProductSKUExists("EXISTING-SKU")
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
	return errors.NewInternalServerErrorWithError("Service operation failed", err)
}
