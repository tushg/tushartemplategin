package errors

import (
	"encoding/json"
	"fmt"
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
// Optional fields can be left empty if not needed
func New(code ErrorCode, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Fields:     make(map[string]interface{}),
	}
}

// NewWithDetails creates a new AppError with additional details
// Optional fields can be left empty if not needed
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
// Used for wrapping database errors and other Go errors
func NewWithError(code ErrorCode, message string, httpStatus int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Err:        err,
		Fields:     make(map[string]interface{}),
	}
}

// WithField adds a field to the error context (for chaining)
func (e *AppError) WithField(key string, value interface{}) *AppError {
	if e.Fields == nil {
		e.Fields = make(map[string]interface{})
	}
	e.Fields[key] = value
	return e
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
