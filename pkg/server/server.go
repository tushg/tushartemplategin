package server

import (
	"fmt"
	"net/http"

	"tushartemplategin/internal/health"
	"tushartemplategin/pkg/config"
	"tushartemplategin/pkg/logger"
)

type Server struct {
	config *config.Config
	logger *logger.Logger
	router *http.ServeMux
}

func New(cfg *config.Config, log *logger.Logger) *Server {
	s := &Server{
		config: cfg,
		logger: log,
		router: http.NewServeMux(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// Health check endpoints
	s.router.HandleFunc("/health", health.HealthCheck)
	s.router.HandleFunc("/health/detailed", health.DetailedHealthCheck)
	
	// Root endpoint
	s.router.HandleFunc("/", s.handleRoot)
	
	// API endpoints
	s.router.HandleFunc("/api/status", s.handleAPIStatus)
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Welcome to Tushar Template Gin", "service": "tushartemplategin", "version": "1.0.0"}`)
}

func (s *Server) handleAPIStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "running", "environment": "%s", "port": %d}`, s.config.Environment, s.config.Port)
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	s.logger.Info("Starting server on port %d", s.config.Port)
	s.logger.Info("Environment: %s", s.config.Environment)
	s.logger.Info("Health check available at: http://localhost%s/health", addr)
	
	return http.ListenAndServe(addr, s.router)
}
