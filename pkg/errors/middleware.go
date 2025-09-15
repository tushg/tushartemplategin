package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents the standardized error response format
type ErrorResponse struct {
	Error     *AppError `json:"error"`
	RequestID string    `json:"request_id,omitempty"`
	Path      string    `json:"path"`
	Method    string    `json:"method"`
	Timestamp string    `json:"timestamp"`
}

// ErrorHandlerMiddleware handles application errors and returns standardized responses
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Convert to AppError if it's not already
			appErr := convertToAppError(err)

			// Determine HTTP status code
			statusCode := determineHTTPStatusCode(appErr)

			// Create error response
			response := ErrorResponse{
				Error:     appErr,
				RequestID: getRequestID(c),
				Path:      c.Request.URL.Path,
				Method:    c.Request.Method,
				Timestamp: appErr.Timestamp.Format("2006-01-02T15:04:05.000000Z"),
			}

			// Return error response
			c.JSON(statusCode, response)
			c.Abort()
		}
	})
}

// convertToAppError converts a generic error to AppError
func convertToAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}

	// Convert generic error to AppError
	return &AppError{
		Code:      SYSTEM_INTERNAL_ERROR,
		Message:   "Internal server error",
		Details:   err.Error(),
		Component: ComponentMiddleware,
		Source:    ErrorSourceInternal,
		Severity:  SeverityHigh,
		Retryable: false,
	}
}

// determineHTTPStatusCode determines the appropriate HTTP status code for an error
func determineHTTPStatusCode(err *AppError) int {
	// Check for specific error codes first
	switch err.Code {
	// Client errors (4xx)
	case VALIDATION_REQUIRED_FIELD, VALIDATION_INVALID_FORMAT, VALIDATION_INVALID_RANGE, VALIDATION_INVALID_TYPE:
		return http.StatusBadRequest
	case AUTH_INVALID_CREDENTIALS, AUTH_TOKEN_EXPIRED, AUTH_TOKEN_INVALID:
		return http.StatusUnauthorized
	case AUTH_INSUFFICIENT_PERMISSIONS, AUTH_ACCOUNT_LOCKED, AUTH_ACCOUNT_DISABLED:
		return http.StatusForbidden
	case PRODUCT_NOT_FOUND, HTTP_NOT_FOUND:
		return http.StatusNotFound
	case VALIDATION_DUPLICATE_VALUE, PRODUCT_CREATE_DUPLICATE_SKU, HTTP_CONFLICT:
		return http.StatusConflict
	case HTTP_TOO_MANY_REQUESTS:
		return http.StatusTooManyRequests
	case VALIDATION_CONSTRAINT_VIOLATION, HTTP_UNPROCESSABLE_ENTITY:
		return http.StatusUnprocessableEntity

	// Server errors (5xx)
	case SYSTEM_INTERNAL_ERROR, SYSTEM_DATABASE_ERROR, SYSTEM_NETWORK_ERROR:
		return http.StatusInternalServerError
	case SYSTEM_TIMEOUT_ERROR, HTTP_GATEWAY_TIMEOUT:
		return http.StatusGatewayTimeout
	case HTTP_BAD_GATEWAY:
		return http.StatusBadGateway
	case HTTP_SERVICE_UNAVAILABLE:
		return http.StatusServiceUnavailable

	// External API errors
	case PAYMENT_API_ERROR, INVENTORY_API_ERROR, NOTIFICATION_API_ERROR:
		if err.Source == ErrorSourceExternal {
			return http.StatusBadGateway
		}
		return http.StatusInternalServerError

	// Default to internal server error
	default:
		return http.StatusInternalServerError
	}
}

// getRequestID extracts request ID from context or headers
func getRequestID(c *gin.Context) string {
	// Try to get from X-Request-ID header
	if requestID := c.GetHeader("X-Request-ID"); requestID != "" {
		return requestID
	}

	// Try to get from correlation ID
	if correlationID := c.GetHeader("X-Correlation-ID"); correlationID != "" {
		return correlationID
	}

	// Generate a new one if not found
	return generateRequestID()
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	// In production, you might want to use a more sophisticated ID generation
	return "req_" + generateShortID()
}

// generateShortID generates a short unique ID
func generateShortID() string {
	// Simple implementation - in production, use a proper ID generator
	return "12345678"
}
