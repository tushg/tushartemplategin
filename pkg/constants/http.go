package constants

// Error Messages
const (
	ERROR_500_INTERNAL_SERVER_ERROR = "Internal Server Error"
	ERROR_400_BAD_REQUEST           = "Bad Request"
	ERROR_401_UNAUTHORIZED          = "Unauthorized"
	ERROR_403_FORBIDDEN             = "Forbidden"
	ERROR_404_NOT_FOUND             = "Not Found"
	ERROR_409_CONFLICT              = "Conflict"
	ERROR_422_UNPROCESSABLE_ENTITY  = "Unprocessable Entity"
	ERROR_429_TOO_MANY_REQUESTS     = "Too Many Requests"
	ERROR_503_SERVICE_UNAVAILABLE   = "Service Unavailable"
)

// Health Check Error Messages
const (
	ERROR_HEALTH_STATUS_FAILED = "Failed to get health status"
	ERROR_READINESS_FAILED     = "Failed to get readiness status"
	ERROR_LIVENESS_FAILED      = "Failed to get liveness status"
)
