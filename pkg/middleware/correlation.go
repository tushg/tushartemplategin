package middleware

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// CorrelationIDHeader is the HTTP header name for correlation ID
	CorrelationIDHeader = "X-Correlation-ID"

	// CorrelationIDKey is the context key for storing correlation ID
	CorrelationIDKey = "correlation_id"

	// TraceIDHeader is the HTTP header name for trace ID (OpenTelemetry compatible)
	TraceIDHeader = "X-Trace-ID"

	// TraceIDKey is the context key for storing trace ID
	TraceIDKey = "trace_id"
)

// Pre-compiled regex patterns for better performance
var (
	// UUID pattern: 8-4-4-4-12 format (RFC 4122 compliant)
	uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

	// OpenTelemetry trace ID pattern: exactly 32 hex chars
	traceIDPattern = regexp.MustCompile(`^[0-9a-f]{32}$`)
)

// CorrelationIDMiddleware extracts or generates correlation ID and trace ID
// and stores them in the request context for use throughout the request lifecycle
//
// Features:
// - Extracts existing correlation ID from X-Correlation-ID header
// - Generates new correlation ID if not present (GUID format)
// - Extracts trace ID from X-Trace-ID header (OpenTelemetry compatible)
// - Generates new trace ID if not present
// - Adds correlation ID to response headers
// - Stores both IDs in request context
//
// Production Considerations:
// - Thread-safe GUID generation
// - Header validation and sanitization
// - Performance optimized (minimal allocations)
// - Compatible with distributed tracing systems
func CorrelationIDMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Extract or generate correlation ID
		correlationID := extractOrGenerateCorrelationID(c)

		// Extract or generate trace ID
		traceID := extractOrGenerateTraceID(c)

		// Store in context for use throughout request lifecycle
		ctx := context.WithValue(c.Request.Context(), CorrelationIDKey, correlationID)
		ctx = context.WithValue(ctx, TraceIDKey, traceID)

		// Update request context
		c.Request = c.Request.WithContext(ctx)

		// Add correlation ID to response headers for client tracking
		c.Header(CorrelationIDHeader, correlationID)
		c.Header(TraceIDHeader, traceID)

		// Continue to next middleware/handler
		c.Next()
	})
}

// extractOrGenerateCorrelationID extracts correlation ID from header or generates new one
func extractOrGenerateCorrelationID(c *gin.Context) string {
	// Try to extract from X-Correlation-ID header
	if correlationID := c.GetHeader(CorrelationIDHeader); correlationID != "" {
		// Validate and sanitize the correlation ID
		if isValidCorrelationID(correlationID) {
			return correlationID
		}
	}

	// Generate new correlation ID (UUID format)
	return generateUUID()
}

// extractOrGenerateTraceID extracts trace ID from header or generates new one
func extractOrGenerateTraceID(c *gin.Context) string {
	// Try to extract from X-Trace-ID header (OpenTelemetry compatible)
	if traceID := c.GetHeader(TraceIDHeader); traceID != "" {
		// Validate and sanitize the trace ID
		if isValidTraceID(traceID) {
			return traceID
		}
	}

	// Generate new trace ID (OpenTelemetry format)
	return generateTraceID()
}

// generateUUID generates a standard UUID v4
func generateUUID() string {
	// Generate a new UUID v4 (random UUID)
	return uuid.New().String()
}

// generateTraceID generates an OpenTelemetry-compliant trace ID
func generateTraceID() string {
	// Generate a new UUID and convert to 32-char hex (OpenTelemetry format)
	uuid := uuid.New()
	// Remove hyphens and convert to lowercase for OpenTelemetry compatibility
	return strings.ReplaceAll(uuid.String(), "-", "")
}

// generateFallbackID generates a fallback ID using timestamp
func generateFallbackID() string {
	// This is a fallback method - in production, uuid.New() should always work
	return strings.ReplaceAll(time.Now().Format("20060102150405.000000"), ".", "")
}

// isValidCorrelationID validates correlation ID format using regex (UUID format)
func isValidCorrelationID(id string) bool {
	// Basic validation: non-empty, reasonable length
	if id == "" || len(id) > 64 {
		return false
	}

	// Trim whitespace
	id = strings.TrimSpace(id)

	// Use regex to validate UUID format (RFC 4122 compliant)
	return uuidPattern.MatchString(id)
}

// isValidTraceID validates trace ID format using regex (OpenTelemetry compatible)
func isValidTraceID(id string) bool {
	// Basic validation: non-empty, exactly 32 characters
	if id == "" || len(id) != 32 {
		return false
	}

	// Trim whitespace
	id = strings.TrimSpace(id)

	// Use regex to validate OpenTelemetry trace ID format
	// Must be exactly 32 hex characters and not all zeros
	return traceIDPattern.MatchString(id) && id != "00000000000000000000000000000000"
}
