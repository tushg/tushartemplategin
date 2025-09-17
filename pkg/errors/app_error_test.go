package errors

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	// Test basic error creation
	err := New(ErrCodeBadRequest, "Invalid request", http.StatusBadRequest)

	if err.Code != ErrCodeBadRequest {
		t.Errorf("Expected code %s, got %s", ErrCodeBadRequest, err.Code)
	}

	if err.Message != "Invalid request" {
		t.Errorf("Expected message 'Invalid request', got '%s'", err.Message)
	}

	if err.HTTPStatus != http.StatusBadRequest {
		t.Errorf("Expected HTTP status %d, got %d", http.StatusBadRequest, err.HTTPStatus)
	}

	if err.Details != "" {
		t.Errorf("Expected empty details, got '%s'", err.Details)
	}

	if err.Fields == nil {
		t.Error("Expected Fields to be initialized")
	}
}

func TestNewWithDetails(t *testing.T) {
	// Test error creation with details
	err := NewWithDetails(ErrCodeNotFound, "Resource not found", "The requested resource was not found", http.StatusNotFound)

	if err.Code != ErrCodeNotFound {
		t.Errorf("Expected code %s, got %s", ErrCodeNotFound, err.Code)
	}

	if err.Message != "Resource not found" {
		t.Errorf("Expected message 'Resource not found', got '%s'", err.Message)
	}

	if err.Details != "The requested resource was not found" {
		t.Errorf("Expected details 'The requested resource was not found', got '%s'", err.Details)
	}

	if err.HTTPStatus != http.StatusNotFound {
		t.Errorf("Expected HTTP status %d, got %d", http.StatusNotFound, err.HTTPStatus)
	}

	if err.Fields == nil {
		t.Error("Expected Fields to be initialized")
	}
}

func TestErrorInterface(t *testing.T) {
	// Test that AppError implements the error interface
	err := New(ErrCodeInternalServer, "Internal server error", http.StatusInternalServerError)

	errorString := err.Error()
	expected := "Internal server error"

	if errorString != expected {
		t.Errorf("Expected error string '%s', got '%s'", expected, errorString)
	}
}

func TestErrorWithDetails(t *testing.T) {
	// Test error string with details
	err := NewWithDetails(ErrCodeNotFound, "Resource not found", "The requested resource was not found", http.StatusNotFound)

	errorString := err.Error()
	expected := "Resource not found: The requested resource was not found"

	if errorString != expected {
		t.Errorf("Expected error string '%s', got '%s'", expected, errorString)
	}
}

func TestIsAppError(t *testing.T) {
	// Test IsAppError function
	appErr := New(ErrCodeBadRequest, "Bad request", http.StatusBadRequest)
	regularErr := fmt.Errorf("regular error")

	if !IsAppError(appErr) {
		t.Error("Expected IsAppError to return true for AppError")
	}

	if IsAppError(regularErr) {
		t.Error("Expected IsAppError to return false for regular error")
	}
}

func TestGetAppError(t *testing.T) {
	// Test GetAppError function
	appErr := New(ErrCodeBadRequest, "Bad request", http.StatusBadRequest)
	regularErr := fmt.Errorf("regular error")

	retrievedErr := GetAppError(appErr)
	if retrievedErr != appErr {
		t.Error("Expected GetAppError to return the same AppError")
	}

	retrievedErr = GetAppError(regularErr)
	if retrievedErr != nil {
		t.Error("Expected GetAppError to return nil for regular error")
	}
}
