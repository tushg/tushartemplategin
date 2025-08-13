package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// Internal packages for health API
	"tushartemplategin/internal/health"

	// External packages for configuration, logging, and server
	"tushartemplategin/pkg/config"
	"tushartemplategin/pkg/logger"
	"tushartemplategin/pkg/server"

	"github.com/gin-gonic/gin"
)

func main() {
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

	// Step 4: Set Gin framework mode based on configuration
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode) // Production mode (no debug info)
	} else {
		gin.SetMode(gin.DebugMode) // Development mode (with debug info)
	}

	// Step 5: Create a new Gin router instance
	router := gin.New()
	//This following line is added to fix the DRP issue.
	//Need to investigate more on this.
	router.SetTrustedProxies(nil)

	// Step 6: Initialize the health package components
	// Create repository (data access layer)
	healthRepo := health.NewHealthRepository()

	// Create service (business logic layer)
	healthService := health.NewHealthService(healthRepo, appLogger)

	// Step 7: Add middleware to the router
	// Add health service to context so routes can access it
	router.Use(func(c *gin.Context) {
		c.Set("healthService", healthService)
		c.Next()
	})

	// Note: We're keeping it simple for now, but we can add more middleware here
	// router.Use(logger.RequestLogger(appLogger))  // Uncomment when you add request logging middleware
	// router.Use(middleware.CORS())                // Uncomment when you add CORS middleware

	// Step 8: Setup API routes using module-level route registration
	api := router.Group("/api/v1") // API version 1 group

	// Register health module routes
	// This makes the health module self-contained and responsible for its own routing
	health.RegisterRoutes(api)

	// Note: Product endpoints have been removed to keep only health API
	// You can add them back later when needed by creating product/routes.go
	// and calling product.RegisterRoutes(api)

	// Step 9: Create HTTP server instance with our router
	srv := server.New(cfg.Server.Port, router)

	// Step 10: Start the server in a background goroutine with proper coordination
	// Use channels to coordinate between server goroutine and main goroutine
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
	// This prevents the main goroutine from hanging if server fails to start
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
	// Listen for SIGINT (Ctrl+C) and SIGTERM (termination signal)
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

	// Step 15: Log successful shutdown
	appLogger.Info(context.Background(), "Server exited", logger.Fields{})
}
