# Constants Package

This package contains common constants used throughout the application, including HTTP status codes and error messages.

## Usage

### HTTP Status Codes

```go
import "github.com/tushartemplategin/pkg/constants"

// Instead of hardcoded values
c.JSON(500, gin.H{"error": "Internal Server Error"})

// Use constants
c.JSON(constants.StatusInternalServerError, gin.H{"error": constants.ERROR_500_INTERNAL_SERVER_ERROR})
```

### Available Constants

#### HTTP Status Codes
- `StatusOK` (200)
- `StatusCreated` (201)
- `StatusBadRequest` (400)
- `StatusUnauthorized` (401)
- `StatusNotFound` (404)
- `StatusInternalServerError` (500)
- And many more...

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
3. **Type Safety**: Prevents typos in status codes
4. **Easy Updates**: Change error messages in one place
5. **Documentation**: Self-documenting code with meaningful constant names

## Adding New Constants

When adding new constants to this package:

1. Group related constants together
2. Use descriptive names that clearly indicate their purpose
3. Add comments for complex constants
4. Follow the existing naming convention
5. Update this README if adding new categories
