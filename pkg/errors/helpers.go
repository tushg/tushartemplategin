package errors

import (
	"context"
	"database/sql"
	"net/http"
)

// ErrorHelper provides utility functions for error handling
type ErrorHelper struct{}

// NewErrorHelper creates a new error helper instance
func NewErrorHelper() *ErrorHelper {
	return &ErrorHelper{}
}

// ===== DATABASE ERROR HELPERS =====

// HandleDatabaseError converts database errors to AppError
func (eh *ErrorHelper) HandleDatabaseError(ctx context.Context, operation, component string, err error) *AppError {
	if err == nil {
		return nil
	}

	switch err {
	case sql.ErrNoRows:
		return NewInternalErrorWithDetails(
			DB_QUERY_FAILED,
			"Record not found",
			"Database query returned no rows",
			component,
		)
	case sql.ErrConnDone:
		return NewSystemError(
			DB_CONNECTION_LOST,
			"Database connection lost",
			component,
		)
	case sql.ErrTxDone:
		return NewSystemError(
			DB_TRANSACTION_FAILED,
			"Transaction already completed",
			component,
		)
	default:
		return WrapError(
			ctx,
			DB_QUERY_FAILED,
			"Database operation failed",
			component,
			err,
		)
	}
}

// HandleConstraintError handles database constraint violations
func (eh *ErrorHelper) HandleConstraintError(ctx context.Context, component string, err error) *AppError {
	return WrapError(
		ctx,
		DB_CONSTRAINT_VIOLATION,
		"Database constraint violation",
		component,
		err,
	)
}

// ===== HTTP ERROR HELPERS =====

// HandleHTTPError converts HTTP errors to AppError
func (eh *ErrorHelper) HandleHTTPError(ctx context.Context, statusCode int, component, serviceName, endpoint string, responseBody string) *AppError {
	switch statusCode {
	case http.StatusBadRequest:
		return WrapExternalError(
			ctx,
			HTTP_BAD_REQUEST,
			"Bad request to external service",
			component,
			serviceName,
			endpoint,
			statusCode,
			responseBody,
			nil,
		)
	case http.StatusUnauthorized:
		return WrapExternalError(
			ctx,
			HTTP_UNAUTHORIZED,
			"Unauthorized access to external service",
			component,
			serviceName,
			endpoint,
			statusCode,
			responseBody,
			nil,
		)
	case http.StatusForbidden:
		return WrapExternalError(
			ctx,
			HTTP_FORBIDDEN,
			"Forbidden access to external service",
			component,
			serviceName,
			endpoint,
			statusCode,
			responseBody,
			nil,
		)
	case http.StatusNotFound:
		return WrapExternalError(
			ctx,
			HTTP_NOT_FOUND,
			"Resource not found in external service",
			component,
			serviceName,
			endpoint,
			statusCode,
			responseBody,
			nil,
		)
	case http.StatusTooManyRequests:
		return WrapExternalError(
			ctx,
			HTTP_TOO_MANY_REQUESTS,
			"Rate limit exceeded for external service",
			component,
			serviceName,
			endpoint,
			statusCode,
			responseBody,
			nil,
		)
	case http.StatusInternalServerError:
		return WrapExternalError(
			ctx,
			HTTP_INTERNAL_SERVER_ERROR,
			"External service internal error",
			component,
			serviceName,
			endpoint,
			statusCode,
			responseBody,
			nil,
		)
	case http.StatusBadGateway:
		return WrapExternalError(
			ctx,
			HTTP_BAD_GATEWAY,
			"Bad gateway to external service",
			component,
			serviceName,
			endpoint,
			statusCode,
			responseBody,
			nil,
		)
	case http.StatusServiceUnavailable:
		return WrapExternalError(
			ctx,
			HTTP_SERVICE_UNAVAILABLE,
			"External service unavailable",
			component,
			serviceName,
			endpoint,
			statusCode,
			responseBody,
			nil,
		)
	case http.StatusGatewayTimeout:
		return WrapExternalError(
			ctx,
			HTTP_GATEWAY_TIMEOUT,
			"Gateway timeout to external service",
			component,
			serviceName,
			endpoint,
			statusCode,
			responseBody,
			nil,
		)
	default:
		return WrapExternalError(
			ctx,
			HTTP_INTERNAL_SERVER_ERROR,
			"Unknown HTTP error from external service",
			component,
			serviceName,
			endpoint,
			statusCode,
			responseBody,
			nil,
		)
	}
}

// ===== VALIDATION ERROR HELPERS =====

// HandleValidationError creates a validation error
func (eh *ErrorHelper) HandleValidationError(ctx context.Context, field, message, component string) *AppError {
	return NewValidationError(
		VALIDATION_REQUIRED_FIELD,
		message,
		component,
	)
}

// HandleDuplicateValueError creates a duplicate value error
func (eh *ErrorHelper) HandleDuplicateValueError(ctx context.Context, field, value, component string) *AppError {
	return NewValidationError(
		VALIDATION_DUPLICATE_VALUE,
		"Duplicate value for field: "+field,
		component,
	)
}

// ===== BUSINESS LOGIC ERROR HELPERS =====

// HandleBusinessLogicError creates a business logic error
func (eh *ErrorHelper) HandleBusinessLogicError(ctx context.Context, code, message, component string) *AppError {
	return NewBusinessLogicError(code, message, component)
}

// HandleNotFoundError creates a not found error
func (eh *ErrorHelper) HandleNotFoundError(ctx context.Context, resource, identifier, component string) *AppError {
	return NewInternalErrorWithDetails(
		PRODUCT_NOT_FOUND,
		resource+" not found",
		"Resource with identifier '"+identifier+"' not found",
		component,
	)
}

// ===== EXTERNAL API ERROR HELPERS =====

// HandlePaymentAPIError handles payment service errors
func (eh *ErrorHelper) HandlePaymentAPIError(ctx context.Context, statusCode int, responseBody string, err error) *AppError {
	return WrapExternalError(
		ctx,
		PAYMENT_API_ERROR,
		"Payment service error",
		ComponentPaymentService,
		"payment-service",
		"/payments",
		statusCode,
		responseBody,
		err,
	)
}

// HandleInventoryAPIError handles inventory service errors
func (eh *ErrorHelper) HandleInventoryAPIError(ctx context.Context, statusCode int, responseBody string, err error) *AppError {
	return WrapExternalError(
		ctx,
		INVENTORY_API_ERROR,
		"Inventory service error",
		ComponentInventoryService,
		"inventory-service",
		"/inventory",
		statusCode,
		responseBody,
		err,
	)
}

// HandleNotificationAPIError handles notification service errors
func (eh *ErrorHelper) HandleNotificationAPIError(ctx context.Context, statusCode int, responseBody string, err error) *AppError {
	return WrapExternalError(
		ctx,
		NOTIFICATION_API_ERROR,
		"Notification service error",
		ComponentNotificationService,
		"notification-service",
		"/notifications",
		statusCode,
		responseBody,
		err,
	)
}

// ===== UTILITY FUNCTIONS =====

// IsRetryableError checks if an error is retryable
func (eh *ErrorHelper) IsRetryableError(err error) bool {
	return IsRetryable(err)
}

// GetErrorCode extracts error code from error
func (eh *ErrorHelper) GetErrorCode(err error) string {
	if appErr := GetAppError(err); appErr != nil {
		return appErr.Code
	}
	return SYSTEM_INTERNAL_ERROR
}

// GetErrorSeverity extracts error severity from error
func (eh *ErrorHelper) GetErrorSeverity(err error) ErrorSeverity {
	return GetSeverity(err)
}

// LogError logs an error with appropriate level
func (eh *ErrorHelper) LogError(ctx context.Context, err error, logger interface{}) {
	// This would integrate with your existing logger
	// Implementation depends on your logger interface
}
