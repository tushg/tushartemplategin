package health

import "time"

// HealthStatus represents the overall health status of the service
type HealthStatus struct {
	Status    string    `json:"status"`    // Service status (healthy, unhealthy)
	Timestamp time.Time `json:"timestamp"` // When the health check was performed
	Service   string    `json:"service"`   // Service name identifier
	Version   string    `json:"version"`   // Service version
}

// ReadinessStatus represents the readiness status for Kubernetes readiness probes
type ReadinessStatus struct {
	Status    string    `json:"status"`    // Readiness status (ready, not ready)
	Timestamp time.Time `json:"timestamp"` // When the readiness check was performed
	Database  string    `json:"database"`  // Database connection status
	Service   string    `json:"service"`   // Service name identifier
}

// LivenessStatus represents the liveness status for Kubernetes liveness probes
type LivenessStatus struct {
	Status    string    `json:"status"`    // Liveness status (alive, dead)
	Timestamp time.Time `json:"timestamp"` // When the liveness check was performed
	Service   string    `json:"service"`   // Service name identifier
}
