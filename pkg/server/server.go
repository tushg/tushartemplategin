package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Server wraps the HTTP server and provides methods for lifecycle management
type Server struct {
	httpServer *http.Server // Underlying HTTP server
	router     *gin.Engine  // Gin router for handling requests
}

// New creates a new server instance with the given port and router
func New(port string, router *gin.Engine) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         port,             // Server address (e.g., ":8080")
			Handler:      router,           // Gin router to handle requests
			ReadTimeout:  15 * time.Second, // Timeout for reading requests
			WriteTimeout: 15 * time.Second, // Timeout for writing responses
			IdleTimeout:  60 * time.Second, // Timeout for idle connections
		},
		router: router,
	}
}

// ListenAndServe starts the HTTP server and begins accepting requests
func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server with the given context
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
