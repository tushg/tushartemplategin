package health

import (
	"github.com/gin-gonic/gin"
	"tushartemplategin/pkg/constants"
)

// RegisterRoutes registers all health-related routes to the given router group
// This function makes the health module self-contained and responsible for its own routing
func RegisterRoutes(router *gin.RouterGroup) {
	// Create a health group under the main API group
	// This will create routes like /api/v1/health, /api/v1/health/ready, etc.
	healthGroup := router.Group("/health")
	{
		// Register health endpoints with their handlers
		// Each endpoint is clearly defined and easy to maintain

		// GET /health - Overall health status check
		healthGroup.GET("", getHealthHandler)

		// GET /health/ready - Kubernetes readiness probe
		// Used to determine if the service is ready to receive traffic
		healthGroup.GET("/ready", getReadinessHandler)

		// GET /health/live - Kubernetes liveness probe
		// Used to determine if the service is alive and running
		healthGroup.GET("/live", getLivenessHandler)
	}
}

// getHealthHandler handles health check requests
func getHealthHandler(c *gin.Context) {
	// Get the health service from the context (we'll set this up in main.go)
	healthService := c.MustGet("healthService").(Service)

	ctx := c.Request.Context()

	// Get health status from service layer
	health, err := healthService.GetHealth(ctx)
	if err != nil {
		// Return 500 error if service fails
		c.JSON(constants.StatusInternalServerError, gin.H{
			"error": constants.ERROR_HEALTH_STATUS_FAILED,
		})
		return
	}

	// Return health status with 200 OK
	c.JSON(constants.StatusOK, health)
}

// getReadinessHandler handles readiness probe requests
func getReadinessHandler(c *gin.Context) {
	healthService := c.MustGet("healthService").(Service)

	ctx := c.Request.Context()

	readiness, err := healthService.GetReadiness(ctx)
	if err != nil {
		c.JSON(constants.StatusInternalServerError, gin.H{
			"error": constants.ERROR_READINESS_FAILED,
		})
		return
	}

	c.JSON(200, readiness)
}

// getLivenessHandler handles liveness probe requests
func getLivenessHandler(c *gin.Context) {
	healthService := c.MustGet("healthService").(Service)

	ctx := c.Request.Context()

	liveness, err := healthService.GetLiveness(ctx)
	if err != nil {
		c.JSON(constants.StatusInternalServerError, gin.H{
			"error": constants.ERROR_LIVENESS_FAILED,
		})
		return
	}

	c.JSON(constants.StatusOK, liveness)
}
