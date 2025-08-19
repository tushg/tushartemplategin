package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Server wraps the HTTP server and provides methods for lifecycle management
type Server struct {
	router     *gin.Engine  // Gin router for handling requests
	port       string       // HTTP port
	sslConfig  SSLConfig    // SSL/TLS configuration
	httpServer *http.Server // HTTP server for redirects (if needed)
}

// SSLConfig contains SSL/TLS configuration settings
type SSLConfig struct {
	Enabled      bool   // Enable SSL/TLS
	Port         string // SSL port (e.g., ":443")
	CertFile     string // Path to SSL certificate file
	KeyFile      string // Path to SSL private key file
	RedirectHTTP bool   // Redirect HTTP to HTTPS
}

// New creates a new server instance with the given port and router
func New(port string, router *gin.Engine, sslConfig SSLConfig) *Server {
	return &Server{
		router:    router,
		port:      port,
		sslConfig: sslConfig,
	}
}

// ListenAndServe starts the server and begins accepting requests
func (s *Server) ListenAndServe() error {
	if s.sslConfig.Enabled {
		// Start HTTP redirect server if enabled
		if s.sslConfig.RedirectHTTP {
			go s.startHTTPRedirectServer()
		}

		// Start HTTPS server using Gin's built-in TLS
		return s.router.RunTLS(s.sslConfig.Port, s.sslConfig.CertFile, s.sslConfig.KeyFile)
	}

	// Start HTTP server only
	return s.router.Run(s.port)
}

// startHTTPRedirectServer starts HTTP server that redirects to HTTPS
func (s *Server) startHTTPRedirectServer() error {
	redirectHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Redirect HTTP to HTTPS
		httpsURL := fmt.Sprintf("https://%s%s", r.Host, r.RequestURI)
		http.Redirect(w, r, httpsURL, http.StatusMovedPermanently)
	})

	s.httpServer = &http.Server{
		Addr:         s.port,
		Handler:      redirectHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server with the given context
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}
