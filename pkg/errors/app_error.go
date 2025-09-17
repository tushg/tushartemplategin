package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorCode represents a standardized error code
type ErrorCode string

// Predefined error codes
const (
	// General errors
	ErrCodeInternalServer      ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrCodeBadRequest          ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized        ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden           ErrorCode = "FORBIDDEN"
	ErrCodeNotFound            ErrorCode = "NOT_FOUND"
	ErrCodeConflict            ErrorCode = "CONFLICT"
	ErrCodeUnprocessableEntity ErrorCode = "UNPROCESSABLE_ENTITY"
	ErrCodeTooManyRequests     ErrorCode = "TOO_MANY_REQUESTS"
	ErrCodeServiceUnavailable  ErrorCode = "SERVICE_UNAVAILABLE"

	// Validation errors
	ErrCodeValidationFailed ErrorCode = "VALIDATION_FAILED"
	ErrCodeInvalidParameter ErrorCode = "INVALID_PARAMETER"
	ErrCodeMissingParameter ErrorCode = "MISSING_PARAMETER"

	// Business logic errors
	ErrCodeProductNotFound     ErrorCode = "PRODUCT_NOT_FOUND"
	ErrCodeProductSKUExists    ErrorCode = "PRODUCT_SKU_EXISTS"
	ErrCodeProductCreateFailed ErrorCode = "PRODUCT_CREATE_FAILED"
	ErrCodeProductUpdateFailed ErrorCode = "PRODUCT_UPDATE_FAILED"
	ErrCodeProductDeleteFailed ErrorCode = "PRODUCT_DELETE_FAILED"
	ErrCodeInvalidStock        ErrorCode = "INVALID_STOCK"

	// Database errors
	ErrCodeDatabaseConnection  ErrorCode = "DATABASE_CONNECTION_ERROR"
	ErrCodeDatabaseQuery       ErrorCode = "DATABASE_QUERY_ERROR"
	ErrCodeDatabaseTransaction ErrorCode = "DATABASE_TRANSACTION_ERROR"
)

// AppError represents a structured application error
type AppError struct {
	Code       ErrorCode              `json:"code"`              // Standardized error code
	Message    string                 `json:"message"`           // Human-readable error message
	Details    string                 `json:"details,omitempty"` // Additional error details
	HTTPStatus int                    `json:"-"`                 // HTTP status code (not serialized)
	Fields     map[string]interface{} `json:"fields,omitempty"`  // Additional context fields
	Err        error                  `json:"-"`                 // Original error (not serialized)
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

// Unwrap returns the underlying error for error unwrapping
func (e *AppError) Unwrap() error {
	return e.Err
}

// ToJSON converts the error to JSON format
func (e *AppError) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// New creates a new AppError with the given code, message, and HTTP status
func New(code ErrorCode, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Fields:     make(map[string]interface{}),
	}
}

// NewWithDetails creates a new AppError with additional details
func NewWithDetails(code ErrorCode, message, details string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		HTTPStatus: httpStatus,
		Fields:     make(map[string]interface{}),
	}
}

// NewWithError creates a new AppError wrapping an existing error
func NewWithError(code ErrorCode, message string, httpStatus int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Err:        err,
		Fields:     make(map[string]interface{}),
	}
}

// WithField adds a field to the error context
func (e *AppError) WithField(key string, value interface{}) *AppError {
	if e.Fields == nil {
		e.Fields = make(map[string]interface{})
	}
	e.Fields[key] = value
	return e
}

// WithFields adds multiple fields to the error context
func (e *AppError) WithFields(fields map[string]interface{}) *AppError {
	if e.Fields == nil {
		e.Fields = make(map[string]interface{})
	}
	for k, v := range fields {
		e.Fields[k] = v
	}
	return e
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// Predefined error constructors for common scenarios

// NewBadRequest creates a 400 Bad Request error
func NewBadRequest(message string) *AppError {
	return New(ErrCodeBadRequest, message, http.StatusBadRequest)
}

// NewBadRequestWithDetails creates a 400 Bad Request error with details
func NewBadRequestWithDetails(message, details string) *AppError {
	return NewWithDetails(ErrCodeBadRequest, message, details, http.StatusBadRequest)
}

// NewUnauthorized creates a 401 Unauthorized error
func NewUnauthorized(message string) *AppError {
	return New(ErrCodeUnauthorized, message, http.StatusUnauthorized)
}

// NewForbidden creates a 403 Forbidden error
func NewForbidden(message string) *AppError {
	return New(ErrCodeForbidden, message, http.StatusForbidden)
}

// NewNotFound creates a 404 Not Found error
func NewNotFound(message string) *AppError {
	return New(ErrCodeNotFound, message, http.StatusNotFound)
}

// NewNotFoundWithDetails creates a 404 Not Found error with details
func NewNotFoundWithDetails(message, details string) *AppError {
	return NewWithDetails(ErrCodeNotFound, message, details, http.StatusNotFound)
}

// NewConflict creates a 409 Conflict error
func NewConflict(message string) *AppError {
	return New(ErrCodeConflict, message, http.StatusConflict)
}

// NewConflictWithDetails creates a 409 Conflict error with details
func NewConflictWithDetails(message, details string) *AppError {
	return NewWithDetails(ErrCodeConflict, message, details, http.StatusConflict)
}

// NewUnprocessableEntity creates a 422 Unprocessable Entity error
func NewUnprocessableEntity(message string) *AppError {
	return New(ErrCodeUnprocessableEntity, message, http.StatusUnprocessableEntity)
}

// NewInternalServerError creates a 500 Internal Server Error
func NewInternalServerError(message string) *AppError {
	return New(ErrCodeInternalServer, message, http.StatusInternalServerError)
}

// NewInternalServerErrorWithError creates a 500 Internal Server Error wrapping an existing error
func NewInternalServerErrorWithError(message string, err error) *AppError {
	return NewWithError(ErrCodeInternalServer, message, http.StatusInternalServerError, err)
}

// NewValidationError creates a 422 Validation Failed error
func NewValidationError(message string) *AppError {
	return New(ErrCodeValidationFailed, message, http.StatusUnprocessableEntity)
}

// NewValidationErrorWithDetails creates a 422 Validation Failed error with details
func NewValidationErrorWithDetails(message, details string) *AppError {
	return NewWithDetails(ErrCodeValidationFailed, message, details, http.StatusUnprocessableEntity)
}

// Business logic error constructors

// NewProductNotFound creates a product not found error
func NewProductNotFound(id interface{}) *AppError {
	return NewNotFoundWithDetails("Product not found", fmt.Sprintf("Product with ID %v not found", id)).
		WithField("product_id", id)
}

// NewProductSKUExists creates a product SKU exists error
func NewProductSKUExists(sku string) *AppError {
	return NewConflictWithDetails("Product SKU already exists", fmt.Sprintf("Product with SKU '%s' already exists", sku)).
		WithField("sku", sku)
}

// NewInvalidStock creates an invalid stock error
func NewInvalidStock(stock int) *AppError {
	return NewValidationErrorWithDetails("Invalid stock quantity", fmt.Sprintf("Stock quantity %d is invalid", stock)).
		WithField("stock", stock)
}

// NewDatabaseError creates a database error
func NewDatabaseError(operation string, err error) *AppError {
	return NewInternalServerErrorWithError("Database operation failed", err).
		WithField("operation", operation)
}

// Helper functions

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError extracts AppError from an error, returns nil if not an AppError
func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return nil
}

// WrapError wraps a standard Go error into an AppError
func WrapError(err error, code ErrorCode, message string, httpStatus int) *AppError {
	return NewWithError(code, message, httpStatus, err)
}
