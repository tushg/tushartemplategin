package middleware

import (
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"tushartemplategin/pkg/errors"
	"tushartemplategin/pkg/interfaces"
)

// ErrorHandlerMiddleware provides centralized error handling for all routes
// It catches panics, converts errors to structured responses, and logs errors
func ErrorHandlerMiddleware(logger interfaces.Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Defer panic recovery
		defer func() {
			if r := recover(); r != nil {
				// Log the panic with stack trace
				logger.Error(c.Request.Context(), "Panic recovered", interfaces.Fields{
					"panic":  r,
					"stack":  string(debug.Stack()),
					"path":   c.Request.URL.Path,
					"method": c.Request.Method,
				})

				// Create internal server error response
				appErr := errors.NewInternalServerError("Internal server error occurred")
				respondWithError(c, appErr)
				c.Abort()
			}
		}()

		// Continue to next handler
		c.Next()

		// Check for errors set by handlers
		if len(c.Errors) > 0 {
			// Get the last error (most recent)
			err := c.Errors.Last()

			// Convert to AppError if it's not already
			appErr := convertToAppError(err.Err)

			// Log the error
			logger.Error(c.Request.Context(), "Request failed", interfaces.Fields{
				"error":  appErr.Error(),
				"code":   appErr.Code,
				"path":   c.Request.URL.Path,
				"method": c.Request.Method,
			})

			// Respond with error
			respondWithError(c, appErr)
			c.Abort()
		}
	})
}

// convertToAppError converts a standard Go error to AppError
func convertToAppError(err error) *errors.AppError {
	// If it's already an AppError, return it
	if appErr := errors.GetAppError(err); appErr != nil {
		return appErr
	}

	// Convert standard error to internal server error
	return errors.NewInternalServerErrorWithError("An unexpected error occurred", err)
}

// respondWithError sends a structured error response
func respondWithError(c *gin.Context, appErr *errors.AppError) {
	// Set the appropriate HTTP status code
	c.JSON(appErr.HTTPStatus, gin.H{
		"error":   appErr.Code,
		"message": appErr.Message,
		"details": appErr.Details,
		"fields":  appErr.Fields,
	})
}

// HandleError is a helper function for handlers to set errors
// Usage: middleware.HandleError(c, errors.NewNotFound("Resource not found"))
func HandleError(c *gin.Context, err error) {
	c.Error(err)
}

// HandleAppError is a helper function for handlers to set AppErrors
// Usage: middleware.HandleAppError(c, errors.NewProductNotFound(123))
func HandleAppError(c *gin.Context, appErr *errors.AppError) {
	c.Error(appErr)
}

// ValidationErrorHandler handles validation errors from Gin's binding
func ValidationErrorHandler(c *gin.Context, err error) {
	// Convert Gin validation error to our AppError
	appErr := errors.NewValidationErrorWithDetails("Validation failed", err.Error())
	HandleAppError(c, appErr)
}

// NotFoundHandler handles 404 errors for undefined routes
func NotFoundHandler(logger interfaces.Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		appErr := errors.NewNotFoundWithDetails("Route not found",
			"The requested route does not exist")

		logger.Warn(c.Request.Context(), "Route not found", interfaces.Fields{
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
		})

		respondWithError(c, appErr)
	})
}

// MethodNotAllowedHandler handles 405 errors for unsupported methods
func MethodNotAllowedHandler(logger interfaces.Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		appErr := errors.NewBadRequestWithDetails("Method not allowed",
			"The HTTP method is not supported for this route")

		logger.Warn(c.Request.Context(), "Method not allowed", interfaces.Fields{
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
		})

		respondWithError(c, appErr)
	})
}
