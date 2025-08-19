package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Server wraps the HTTP server and provides methods for lifecycle management
type Server struct {
	httpServer  *http.Server // Underlying HTTP server
	httpsServer *http.Server // Underlying HTTPS server
	router      *gin.Engine  // Gin router for handling requests
	sslConfig   SSLConfig    // SSL/TLS configuration
}

// SSLConfig contains SSL/TLS configuration
type SSLConfig struct {
	Enabled      bool   // Enable SSL/TLS
	Port         string // SSL port
	CertFile     string // Path to SSL certificate file
	KeyFile      string // Path to SSL private key file
	RedirectHTTP bool   // Redirect HTTP to HTTPS
}

// New creates a new server instance with the given port and router
func New(port string, router *gin.Engine, sslConfig SSLConfig) *Server {
	server := &Server{
		httpServer: &http.Server{
			Addr:         port,             // Server address (e.g., ":8080")
			Handler:      router,           // Gin router to handle requests
			ReadTimeout:  15 * time.Second, // Timeout for reading requests
			WriteTimeout: 15 * time.Second, // Timeout for writing responses
			IdleTimeout:  60 * time.Second, // Timeout for idle connections
		},
		router:    router,
		sslConfig: sslConfig,
	}

	// If SSL is enabled, create HTTPS server
	if sslConfig.Enabled {
		server.httpsServer = &http.Server{
			Addr:         sslConfig.Port,   // SSL port (e.g., ":443")
			Handler:      router,           // Gin router to handle requests
			ReadTimeout:  15 * time.Second, // Timeout for reading requests
			WriteTimeout: 15 * time.Second, // Timeout for writing responses
			IdleTimeout:  60 * time.Second, // Timeout for idle connections
		}
	}

	return server
}

// ListenAndServe starts the HTTP server and begins accepting requests
func (s *Server) ListenAndServe() error {
	if s.sslConfig.Enabled {
		// Start both HTTP and HTTPS servers
		go func() {
			// Start HTTPS server
			if err := s.ListenAndServeTLS(); err != nil {
				panic(fmt.Sprintf("HTTPS server failed to start: %v", err))
			}
		}()

		// Start HTTP server (for redirects if enabled)
		if s.sslConfig.RedirectHTTP {
			return s.startHTTPRedirectServer()
		}
	}

	return s.httpServer.ListenAndServe()
}

// ListenAndServeTLS starts the HTTPS server with TLS
func (s *Server) ListenAndServeTLS() error {
	if s.httpsServer == nil {
		return fmt.Errorf("HTTPS server not configured")
	}

	// Load TLS configuration
	tlsConfig, err := s.loadTLSConfig()
	if err != nil {
		return fmt.Errorf("failed to load TLS configuration: %w", err)
	}

	// Create listener with TLS
	listener, err := tls.Listen("tcp", s.httpsServer.Addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to create TLS listener: %w", err)
	}
	defer listener.Close()

	return s.httpsServer.Serve(listener)
}

// loadTLSConfig loads TLS configuration from certificate files
func (s *Server) loadTLSConfig() (*tls.Config, error) {
	// Load certificate and private key
	cert, err := tls.LoadX509KeyPair(s.sslConfig.CertFile, s.sslConfig.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate: %w", err)
	}

	// Create TLS configuration
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12, // Minimum TLS 1.2 for security
		MaxVersion:   tls.VersionTLS13, // Support up to TLS 1.3

		// Security settings
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},

		// HSTS and security headers
		NextProtos: []string{"h2", "http/1.1"}, // Support HTTP/2
	}

	return tlsConfig, nil
}

// startHTTPRedirectServer starts HTTP server that redirects to HTTPS
func (s *Server) startHTTPRedirectServer() error {
	redirectHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Redirect HTTP to HTTPS
		httpsURL := fmt.Sprintf("https://%s%s", r.Host, r.RequestURI)
		http.Redirect(w, r, httpsURL, http.StatusMovedPermanently)
	})

	redirectServer := &http.Server{
		Addr:         s.httpServer.Addr,
		Handler:      redirectHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return redirectServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server with the given context
func (s *Server) Shutdown(ctx context.Context) error {
	var errs []error

	// Shutdown HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		errs = append(errs, fmt.Errorf("HTTP server shutdown error: %w", err))
	}

	// Shutdown HTTPS server if enabled
	if s.httpsServer != nil {
		if err := s.httpsServer.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("HTTPS server shutdown error: %w", err))
		}
	}

	// Return first error if any occurred
	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}
