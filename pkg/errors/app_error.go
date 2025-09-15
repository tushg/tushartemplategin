package errors

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// ErrorSource represents the source of the error
type ErrorSource string

const (
	ErrorSourceInternal ErrorSource = "internal"
	ErrorSourceExternal ErrorSource = "external"
	ErrorSourceSystem   ErrorSource = "system"
)

// ErrorSeverity represents the severity level of the error
type ErrorSeverity string

const (
	SeverityLow      ErrorSeverity = "low"
	SeverityMedium   ErrorSeverity = "medium"
	SeverityHigh     ErrorSeverity = "high"
	SeverityCritical ErrorSeverity = "critical"
)

// AppError represents a structured application error
type AppError struct {
	Code        string            `json:"code"`
	Message     string            `json:"message"`
	Details     string            `json:"details,omitempty"`
	Component   string            `json:"component"`
	Service     string            `json:"service,omitempty"`
	Timestamp   time.Time         `json:"timestamp"`
	Stack       string            `json:"stack,omitempty"`
	TraceID     string            `json:"trace_id,omitempty"`
	SpanID      string            `json:"span_id,omitempty"`
	Attributes  map[string]string `json:"attributes,omitempty"`
	Severity    ErrorSeverity     `json:"severity"`
	Source      ErrorSource       `json:"source"`
	Retryable   bool              `json:"retryable"`
	Cause       error             `json:"-"`
	ExternalAPI *ExternalAPIError `json:"external_api,omitempty"`
}

// ExternalAPIError represents error information from external APIs
type ExternalAPIError struct {
	ServiceName    string `json:"service_name"`
	Endpoint       string `json:"endpoint"`
	HTTPStatusCode int    `json:"http_status_code"`
	ResponseBody   string `json:"response_body,omitempty"`
	RequestID      string `json:"request_id,omitempty"`
	ErrorCode      string `json:"error_code,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

// Unwrap returns the underlying cause error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewInternalError creates a new internal application error
func NewInternalError(code, message, component string) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Component: component,
		Timestamp: time.Now().UTC(),
		Stack:     getStackTrace(),
		Source:    ErrorSourceInternal,
		Severity:  SeverityHigh,
		Retryable: false,
	}
}

// NewInternalErrorWithDetails creates a new internal error with details
func NewInternalErrorWithDetails(code, message, details, component string) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Details:   details,
		Component: component,
		Timestamp: time.Now().UTC(),
		Stack:     getStackTrace(),
		Source:    ErrorSourceInternal,
		Severity:  SeverityHigh,
		Retryable: false,
	}
}

// NewExternalAPIError creates a new external API error
func NewExternalAPIError(
	code, message, component, serviceName, endpoint string,
	httpStatus int, responseBody string, cause error,
) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Component: component,
		Service:   serviceName,
		Timestamp: time.Now().UTC(),
		Stack:     getStackTrace(),
		Source:    ErrorSourceExternal,
		Severity:  determineSeverity(httpStatus),
		Retryable: isRetryableError(httpStatus),
		Cause:     cause,
		ExternalAPI: &ExternalAPIError{
			ServiceName:    serviceName,
			Endpoint:       endpoint,
			HTTPStatusCode: httpStatus,
			ResponseBody:   responseBody,
		},
	}
}

// WrapExternalError wraps an external API error with context
func WrapExternalError(
	ctx context.Context,
	code, message, component, serviceName, endpoint string,
	httpStatus int, responseBody string, cause error,
) *AppError {
	err := NewExternalAPIError(code, message, component, serviceName, endpoint, httpStatus, responseBody, cause)

	// Add trace context if available
	if traceID := getTraceIDFromContext(ctx); traceID != "" {
		err.TraceID = traceID
	}

	return err
}

// WrapError wraps an existing error with additional context
func WrapError(ctx context.Context, code, message, component string, cause error) *AppError {
	err := &AppError{
		Code:      code,
		Message:   message,
		Component: component,
		Timestamp: time.Now().UTC(),
		Stack:     getStackTrace(),
		Source:    ErrorSourceInternal,
		Severity:  SeverityHigh,
		Retryable: false,
		Cause:     cause,
	}

	// Add trace context if available
	if traceID := getTraceIDFromContext(ctx); traceID != "" {
		err.TraceID = traceID
	}

	return err
}

// NewValidationError creates a new validation error
func NewValidationError(code, message, component string) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Component: component,
		Timestamp: time.Now().UTC(),
		Stack:     getStackTrace(),
		Source:    ErrorSourceInternal,
		Severity:  SeverityMedium,
		Retryable: false,
	}
}

// NewBusinessLogicError creates a new business logic error
func NewBusinessLogicError(code, message, component string) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Component: component,
		Timestamp: time.Now().UTC(),
		Stack:     getStackTrace(),
		Source:    ErrorSourceInternal,
		Severity:  SeverityMedium,
		Retryable: false,
	}
}

// NewSystemError creates a new system-level error
func NewSystemError(code, message, component string) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Component: component,
		Timestamp: time.Now().UTC(),
		Stack:     getStackTrace(),
		Source:    ErrorSourceSystem,
		Severity:  SeverityCritical,
		Retryable: false,
	}
}

// Helper functions

func getStackTrace() string {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

func getTraceIDFromContext(ctx context.Context) string {
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		return traceID
	}
	return ""
}

func determineSeverity(httpStatus int) ErrorSeverity {
	switch {
	case httpStatus >= 500:
		return SeverityHigh
	case httpStatus >= 400:
		return SeverityMedium
	default:
		return SeverityLow
	}
}

func isRetryableError(httpStatus int) bool {
	// Retryable status codes
	retryableCodes := map[int]bool{
		408: true, // Request Timeout
		429: true, // Too Many Requests
		500: true, // Internal Server Error
		502: true, // Bad Gateway
		503: true, // Service Unavailable
		504: true, // Gateway Timeout
	}
	return retryableCodes[httpStatus]
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError extracts AppError from error chain
func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return nil
}

// IsRetryable checks if an error is retryable
func IsRetryable(err error) bool {
	if appErr := GetAppError(err); appErr != nil {
		return appErr.Retryable
	}
	return false
}

// GetSeverity returns the severity of an error
func GetSeverity(err error) ErrorSeverity {
	if appErr := GetAppError(err); appErr != nil {
		return appErr.Severity
	}
	return SeverityMedium
}
