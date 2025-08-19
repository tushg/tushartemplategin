package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	// Internal packages for health API
	"tushartemplategin/internal/health"

	// External packages for configuration, logging, and server
	"tushartemplategin/pkg/config"
	"tushartemplategin/pkg/database"
	"tushartemplategin/pkg/logger"
	"tushartemplategin/pkg/middleware"
	"tushartemplategin/pkg/server"
)

func main() {
	// ===== CONFIGURATION SETUP =====
	// Step 1: Load application configuration from config files
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Step 2: Create logger configuration from config
	logConfig := &logger.Config{
		Level:      cfg.Log.Level,      // Log level (debug, info, warn, error, fatal)
		Format:     cfg.Log.Format,     // Log format (json, console)
		Output:     cfg.Log.Output,     // Output destination (stdout, file)
		FilePath:   cfg.Log.FilePath,   // Log file path (if output is file)
		MaxSize:    cfg.Log.MaxSize,    // Maximum log file size in MB
		MaxBackups: cfg.Log.MaxBackups, // Maximum number of backup files
		MaxAge:     cfg.Log.MaxAge,     // Maximum age of log files in days
		Compress:   cfg.Log.Compress,   // Whether to compress old log files
		AddCaller:  cfg.Log.AddCaller,  // Whether to add caller information
		AddStack:   cfg.Log.AddStack,   // Whether to add stack traces
	}

	// Step 3: Initialize the structured logger
	appLogger, err := logger.NewLogger(logConfig)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// ===== DATABASE INITIALIZATION =====
	// Step 4: Initialize database using factory pattern
	appLogger.Info(context.Background(), "Initializing database", logger.Fields{
		"type": cfg.Database.Type,
	})

	// Create database factory and instance
	dbFactory := database.NewDatabaseFactory(appLogger)
	db, err := dbFactory.CreateDatabase(&cfg.Database)
	if err != nil {
		appLogger.Error(context.Background(), "Failed to create database instance", logger.Fields{
			"error": err.Error(),
			"type":  cfg.Database.Type,
		})
		log.Fatalf("Failed to create database instance: %v", err)
	}

	// Connect to database
	ctx := context.Background()
	if err := db.Connect(ctx); err != nil {
		appLogger.Error(ctx, "Failed to connect to database", logger.Fields{"error": err.Error()})
		// Continue without database for now
	} else {
		appLogger.Info(ctx, "Successfully connected to database", logger.Fields{
			"type": cfg.Database.Type,
		})
	}

	// ===== SERVER INITIALIZATION =====
	// Step 5: Set Gin framework mode based on configuration
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode) // Production mode (no debug info)
	} else {
		gin.SetMode(gin.DebugMode) // Development mode (with debug info)
	}

	// Step 6: Create a new Gin router instance
	router := gin.New()
	router.SetTrustedProxies(nil)

	// ===== DOMAIN SETUP =====
	// Step 7: Setup domains and middleware
	router = setupDomainsAndMiddleware(router, appLogger, db)

	// Step 8: Setup API routes using module-level route registration
	api := router.Group("/api/v1") // API version 1 group

	// Register all domain routes in a clean, organized way
	registerAllRoutes(api, appLogger)

	// ===== SERVER LIFECYCLE =====
	// Step 9: Create server instance with our router and SSL configuration
	sslConfig := server.SSLConfig{
		Enabled:      cfg.Server.SSL.Enabled,
		Port:         cfg.Server.SSL.Port,
		CertFile:     cfg.Server.SSL.CertFile,
		KeyFile:      cfg.Server.SSL.KeyFile,
		RedirectHTTP: cfg.Server.SSL.RedirectHTTP,
	}

	srv := server.New(cfg.Server.Port, router, sslConfig)

	// Step 10: Start the server in a background goroutine with proper coordination
	serverErr := make(chan error, 1)
	serverStarted := make(chan bool, 1)

	go func() {
		appLogger.Info(context.Background(), "Starting TUSHAR TEMPLATE GIN...", logger.Fields{
			"port": cfg.Server.Port, // Log the port we're starting on
			"mode": cfg.Server.Mode, // Log the server mode
		})

		// Signal that server is attempting to start
		serverStarted <- true

		// Start listening for HTTP requests
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error(context.Background(), "Server failed to start", logger.Fields{
				"port":  cfg.Server.Port,
				"error": err.Error(),
			})
			// Send error to main goroutine so it can exit
			serverErr <- err
		} else {
			// Server stopped normally (not due to error)
			serverErr <- nil
		}
	}()

	// Step 11: Wait for either server to start successfully OR fail to start
	select {
	case <-serverStarted:
		appLogger.Info(context.Background(), "Server started successfully", logger.Fields{
			"port": cfg.Server.Port,
		})
	case err := <-serverErr:
		appLogger.Fatal(context.Background(), "Server failed to start, exiting", err, logger.Fields{
			"port": cfg.Server.Port,
		})
	}

	// Step 12: Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for either shutdown signal OR server error
	select {
	case <-quit:
		appLogger.Info(context.Background(), "Shutdown signal received", logger.Fields{})
	case err := <-serverErr:
		appLogger.Error(context.Background(), "Server encountered error, shutting down", logger.Fields{
			"error": err.Error(),
		})
	}

	appLogger.Info(context.Background(), "Shutting down server", logger.Fields{})

	// Step 13: Create a deadline for server shutdown (30 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // Ensure context is cancelled when function exits

	// Step 14: Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Fatal(context.Background(), "Server forced to shutdown", err, logger.Fields{})
	}

	// Step 15: Disconnect from database
	if err := db.Disconnect(ctx); err != nil {
		appLogger.Error(ctx, "Failed to disconnect from database", logger.Fields{"error": err.Error()})
	} else {
		appLogger.Info(ctx, "Database disconnected successfully", logger.Fields{})
	}

	// Step 16: Log successful shutdown
	appLogger.Info(context.Background(), "Server exited", logger.Fields{})
}

// setupDomainsAndMiddleware initializes domain-specific components and middleware
func setupDomainsAndMiddleware(router *gin.Engine, appLogger logger.Logger, db database.Database) *gin.Engine {
	ctx := context.Background()

	// ===== SECURITY MIDDLEWARE =====
	appLogger.Info(ctx, "Setting up security middleware", logger.Fields{})
	router.Use(middleware.SecurityHeaders())
	appLogger.Info(ctx, "Security middleware setup complete", logger.Fields{})

	// ===== CURRENT DOMAINS =====
	appLogger.Info(ctx, "Setting up health domain", logger.Fields{})

	// Create repository (data access layer) - NOW WITH DATABASE
	healthRepo := health.NewHealthRepository(db, appLogger)

	// Create service (business logic layer)
	healthService := health.NewHealthService(healthRepo, appLogger)

	// Add health service to context so routes can access it
	router.Use(func(c *gin.Context) {
		c.Set("healthService", healthService)
		c.Next()
	})
	appLogger.Info(ctx, "Health domain setup complete", logger.Fields{})

	appLogger.Info(ctx, "All domain setup complete", logger.Fields{})
	return router
}

// registerAllRoutes handles all domain route registrations in one organized place
func registerAllRoutes(api *gin.RouterGroup, appLogger logger.Logger) {
	ctx := context.Background()

	// ===== CURRENT DOMAINS =====
	appLogger.Info(ctx, "Registering health domain routes", logger.Fields{})
	health.RegisterRoutes(api)
	appLogger.Info(ctx, "Health domain routes registered successfully", logger.Fields{})

	appLogger.Info(ctx, "All domain routes registered successfully", logger.Fields{})
}
