# Constants Package

This package contains common constants used throughout the application, including error messages and other application-specific constants.

## Usage

### HTTP Status Codes

```go
import (
    "net/http"
    "github.com/tushartemplategin/pkg/constants"
)

// Use standard net/http constants for status codes
c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ERROR_500_INTERNAL_SERVER_ERROR})
c.JSON(http.StatusOK, gin.H{"status": "success"})
c.JSON(http.StatusBadRequest, gin.H{"error": constants.ERROR_400_BAD_REQUEST})
```

### Available Constants

#### Error Messages
- `ERROR_500_INTERNAL_SERVER_ERROR`
- `ERROR_400_BAD_REQUEST`
- `ERROR_404_NOT_FOUND`
- And more...

#### Health Check Specific Errors
- `ERROR_HEALTH_STATUS_FAILED`
- `ERROR_READINESS_FAILED`
- `ERROR_LIVENESS_FAILED`

## Benefits

1. **Maintainability**: Centralized location for all constants
2. **Consistency**: Ensures consistent error messages across the application
3. **Standards Compliance**: Uses standard Go `net/http` constants for HTTP status codes
4. **Easy Updates**: Change error messages in one place
5. **Documentation**: Self-documenting code with meaningful constant names

## Adding New Constants

When adding new constants to this package:

1. Group related constants together
2. Use descriptive names that clearly indicate their purpose
3. Add comments for complex constants
4. Follow the existing naming convention
5. Update this README if adding new categories

## Note on HTTP Status Codes

This package no longer defines custom HTTP status code constants. Instead, use the standard Go `net/http` package constants:
- `http.StatusOK` (200)
- `http.StatusCreated` (201)
- `http.StatusBadRequest` (400)
- `http.StatusUnauthorized` (401)
- `http.StatusNotFound` (404)
- `http.StatusInternalServerError` (500)
- And many more...
