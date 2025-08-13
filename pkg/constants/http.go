package constants

// HTTP Status Codes
const (
	// 2xx Success
	StatusOK        = 200
	StatusCreated   = 201
	StatusAccepted  = 202
	StatusNoContent = 204

	// 4xx Client Errors
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusMethodNotAllowed    = 405
	StatusConflict            = 409
	StatusUnprocessableEntity = 422
	StatusTooManyRequests     = 429

	// 5xx Server Errors
	StatusInternalServerError = 500
	StatusNotImplemented      = 501
	StatusBadGateway          = 502
	StatusServiceUnavailable  = 503
	StatusGatewayTimeout      = 504
)

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
