package middleware

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"tushartemplategin/pkg/errors"
	"tushartemplategin/pkg/interfaces"
)

// ValidationMiddleware provides request validation for JSON payloads
func ValidationMiddleware(logger interfaces.Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Only validate POST, PUT, PATCH requests with JSON content
		if !isJSONRequest(c) {
			c.Next()
			return
		}

		// Get the request body
		body, err := c.GetRawData()
		if err != nil {
			logger.Error(c.Request.Context(), "Failed to read request body", interfaces.Fields{
				"error": err.Error(),
			})
			HandleAppError(c, errors.NewBadRequestWithDetails("Invalid request body", err.Error()))
			c.Abort()
			return
		}

		// Restore the body for the next handler
		c.Request.Body = &bodyReader{data: body}

		// Try to bind JSON to a generic map for basic validation
		var requestData map[string]interface{}
		if err := json.Unmarshal(body, &requestData); err != nil {
			logger.Warn(c.Request.Context(), "Invalid JSON in request body", interfaces.Fields{
				"error": err.Error(),
			})
			HandleAppError(c, errors.NewBadRequestWithDetails("Invalid JSON format", err.Error()))
			c.Abort()
			return
		}

		// Basic validation - check for required fields if needed
		if err := validateRequestData(c, requestData, logger); err != nil {
			HandleAppError(c, err.(*errors.AppError))
			c.Abort()
			return
		}

		c.Next()
	})
}

// isJSONRequest checks if the request is a JSON request
func isJSONRequest(c *gin.Context) bool {
	contentType := c.GetHeader("Content-Type")
	return strings.Contains(contentType, "application/json") &&
		(c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH")
}

// validateRequestData performs basic validation on request data
func validateRequestData(c *gin.Context, data map[string]interface{}, logger interfaces.Logger) error {
	// Add any custom validation logic here
	// For now, we'll just log the validation
	logger.Debug(c.Request.Context(), "Request validation passed", interfaces.Fields{
		"fields_count": len(data),
		"path":         c.Request.URL.Path,
	})

	return nil
}

// bodyReader implements io.ReadCloser to restore request body
type bodyReader struct {
	data []byte
	pos  int
}

func (r *bodyReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, nil
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

func (r *bodyReader) Close() error {
	return nil
}

// ValidateStruct validates a struct using go-playground/validator
func ValidateStruct(ctx context.Context, logger interfaces.Logger, s interface{}) error {
	validate := validator.New()

	if err := validate.Struct(s); err != nil {
		// Convert validation errors to our AppError format
		validationErrors := make(map[string]string)

		for _, err := range err.(validator.ValidationErrors) {
			field := strings.ToLower(err.Field())
			validationErrors[field] = getValidationMessage(err)
		}

		logger.Warn(ctx, "Validation failed", interfaces.Fields{
			"errors": validationErrors,
		})

		// Convert map[string]string to map[string]interface{}
		fields := make(map[string]interface{})
		for k, v := range validationErrors {
			fields[k] = v
		}
		
		return errors.NewValidationErrorWithDetails("Validation failed", 
			formatValidationErrors(validationErrors)).
			WithFields(fields)
	}

	return nil
}

// getValidationMessage returns a user-friendly validation message
func getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "Value is too small"
	case "max":
		return "Value is too large"
	case "email":
		return "Invalid email format"
	case "len":
		return "Invalid length"
	case "numeric":
		return "Must be a number"
	case "alpha":
		return "Must contain only letters"
	case "alphanum":
		return "Must contain only letters and numbers"
	default:
		return "Invalid value"
	}
}

// formatValidationErrors formats validation errors into a readable string
func formatValidationErrors(errors map[string]string) string {
	var messages []string
	for field, message := range errors {
		messages = append(messages, field+": "+message)
	}
	return strings.Join(messages, "; ")
}

// RequiredFieldsMiddleware validates that required fields are present
func RequiredFieldsMiddleware(requiredFields []string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		if !isJSONRequest(c) {
			c.Next()
			return
		}

		// Get the request body
		body, err := c.GetRawData()
		if err != nil {
			HandleAppError(c, errors.NewBadRequestWithDetails("Invalid request body", err.Error()))
			c.Abort()
			return
		}

		// Restore the body for the next handler
		c.Request.Body = &bodyReader{data: body}

		// Parse JSON
		var requestData map[string]interface{}
		if err := json.Unmarshal(body, &requestData); err != nil {
			HandleAppError(c, errors.NewBadRequestWithDetails("Invalid JSON format", err.Error()))
			c.Abort()
			return
		}

		// Check for required fields
		missingFields := []string{}
		for _, field := range requiredFields {
			if _, exists := requestData[field]; !exists {
				missingFields = append(missingFields, field)
			}
		}

		if len(missingFields) > 0 {
			HandleAppError(c, errors.NewValidationErrorWithDetails("Missing required fields",
				"Required fields: "+strings.Join(missingFields, ", ")).
				WithField("missing_fields", missingFields))
			c.Abort()
			return
		}

		c.Next()
	})
}
