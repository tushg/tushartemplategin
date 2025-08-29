package health

import (
	"time"
)

// HealthStatus represents the health status of a service
type HealthStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
}

// ReadinessStatus represents the readiness status of a service
type ReadinessStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Database  string    `json:"database"`
	Service   string    `json:"service"`
}

// LivenessStatus represents the liveness status of a service
type LivenessStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}
