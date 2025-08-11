package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"tushartemplategin/pkg/config"
	"tushartemplategin/pkg/logger"
	"tushartemplategin/pkg/server"
)

func main() {
	// Load configuration
	cfg := config.Load()
	
	// Initialize logger
	logger := logger.New()
	
	// Create and start server
	srv := server.New(cfg, logger)
	
	// Start server in a goroutine
	go func() {
		logger.Info("Starting server...")
		if err := srv.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	logger.Info("Shutting down server...")
}
