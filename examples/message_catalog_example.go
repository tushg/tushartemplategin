package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"tushartemplategin/internal/domains/messagecatalog"
	"tushartemplategin/pkg/config"
	"tushartemplategin/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create logger
	logConfig := &logger.Config{
		Level:      cfg.Log.Level,
		Format:     cfg.Log.Format,
		Output:     cfg.Log.Output,
		FilePath:   cfg.Log.FilePath,
		MaxSize:    cfg.Log.MaxSize,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAge,
		Compress:   cfg.Log.Compress,
		AddCaller:  cfg.Log.AddCaller,
		AddStack:   cfg.Log.AddStack,
	}

	appLogger, err := logger.NewLogger(logConfig)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Create message catalog service
	messageCatalogService := messagecatalog.NewMessageCatalogService(cfg.GetMessageCatalog(), appLogger)

	ctx := context.Background()

	fmt.Println("=== Message Catalog Service Example ===\n")

	// Example 1: Get a single message
	fmt.Println("1. Getting a single message:")
	req := &messagecatalog.MessageRequest{
		MessageCode: "ABC0001",
		CatalogName: "alert",
		Language:    "en-US",
		Parameters: map[string]interface{}{
			"exp_date": "2024-12-31",
		},
	}

	message, err := messageCatalogService.GetMessage(ctx, req)
	if err != nil {
		log.Printf("Error getting message: %v", err)
	} else {
		fmt.Printf("   Message Code: %s\n", message.MessageCode)
		fmt.Printf("   Category: %s\n", message.Category)
		fmt.Printf("   Severity: %s\n", message.Severity)
		fmt.Printf("   Message: %s\n", message.Message)
		fmt.Printf("   Formatted Message: %s\n", message.FormattedMessage)
		fmt.Printf("   Response Action: %s\n", message.ResponseAction)
		fmt.Println()
	}

	// Example 2: Get message in French
	fmt.Println("2. Getting message in French:")
	reqFr := &messagecatalog.MessageRequest{
		MessageCode: "ABC0001",
		CatalogName: "alert",
		Language:    "fr-FR",
		Parameters: map[string]interface{}{
			"exp_date": "2024-12-31",
		},
	}

	messageFr, err := messageCatalogService.GetMessage(ctx, reqFr)
	if err != nil {
		log.Printf("Error getting French message: %v", err)
	} else {
		fmt.Printf("   Message: %s\n", messageFr.Message)
		fmt.Printf("   Formatted Message: %s\n", messageFr.FormattedMessage)
		fmt.Println()
	}

	// Example 3: Get messages by category
	fmt.Println("3. Getting messages by category (Registration):")
	messages, err := messageCatalogService.GetMessagesByCategory(ctx, "Registration", "alert", "en-US")
	if err != nil {
		log.Printf("Error getting messages by category: %v", err)
	} else {
		for i, msg := range messages {
			fmt.Printf("   %d. %s - %s\n", i+1, msg.MessageCode, msg.Message)
		}
		fmt.Println()
	}

	// Example 4: Get messages by severity
	fmt.Println("4. Getting messages by severity (CRITICAL):")
	criticalMessages, err := messageCatalogService.GetMessagesBySeverity(ctx, "CRITICAL", "alert", "en-US")
	if err != nil {
		log.Printf("Error getting messages by severity: %v", err)
	} else {
		for i, msg := range criticalMessages {
			fmt.Printf("   %d. %s - %s\n", i+1, msg.MessageCode, msg.Message)
		}
		fmt.Println()
	}

	// Example 5: Get audit message
	fmt.Println("5. Getting audit message:")
	auditReq := &messagecatalog.MessageRequest{
		MessageCode: "AUD0001",
		CatalogName: "audit",
		Language:    "en-US",
		Parameters: map[string]interface{}{
			"username":         "john.doe",
			"action":           "logged in",
			"timestamp":        time.Now().Format("2006-01-02 15:04:05"),
			"ip_address":       "192.168.1.100",
			"session_duration": "2h 30m",
		},
	}

	auditMessage, err := messageCatalogService.GetMessage(ctx, auditReq)
	if err != nil {
		log.Printf("Error getting audit message: %v", err)
	} else {
		fmt.Printf("   Event Code: %s\n", auditMessage.MessageCode)
		fmt.Printf("   Category: %s\n", auditMessage.Category)
		fmt.Printf("   Severity: %s\n", auditMessage.Severity)
		fmt.Printf("   Message: %s\n", auditMessage.Message)
		fmt.Printf("   Formatted Message: %s\n", auditMessage.FormattedMessage)
		fmt.Printf("   Response Action: %s\n", auditMessage.ResponseAction)
		fmt.Println()
	}

	// Example 6: List available catalogs
	fmt.Println("6. Available catalogs:")
	catalogs, err := messageCatalogService.ListAvailableCatalogs(ctx)
	if err != nil {
		log.Printf("Error listing catalogs: %v", err)
	} else {
		for i, catalog := range catalogs {
			fmt.Printf("   %d. %s\n", i+1, catalog)
		}
		fmt.Println()
	}

	// Example 7: Get catalog info
	fmt.Println("7. Alert catalog info:")
	catalogInfo, err := messageCatalogService.GetCatalogInfo(ctx, "alert")
	if err != nil {
		log.Printf("Error getting catalog info: %v", err)
	} else {
		fmt.Printf("   Name: %s\n", catalogInfo.Name)
		fmt.Printf("   Path: %s\n", catalogInfo.Path)
		fmt.Printf("   Enabled: %t\n", catalogInfo.Enabled)
		fmt.Printf("   Message Count: %d\n", catalogInfo.MessageCount)
		fmt.Printf("   Languages: %v\n", catalogInfo.Languages)
		fmt.Printf("   Last Reloaded: %s\n", catalogInfo.LastReloaded.Format("2006-01-02 15:04:05"))
		fmt.Println()
	}

	// Example 8: Get catalog statistics
	fmt.Println("8. Catalog statistics:")
	stats, err := messageCatalogService.GetCatalogStats(ctx)
	if err != nil {
		log.Printf("Error getting catalog stats: %v", err)
	} else {
		fmt.Printf("   Total Catalogs: %d\n", stats.TotalCatalogs)
		fmt.Printf("   Total Messages: %d\n", stats.TotalMessages)
		fmt.Printf("   Languages Count: %d\n", stats.LanguagesCount)
		fmt.Printf("   Messages by Catalog: %v\n", stats.MessagesByCatalog)
		fmt.Printf("   Last Reloaded: %s\n", stats.LastReloaded.Format("2006-01-02 15:04:05"))
		fmt.Println()
	}

	// Example 9: Health check
	fmt.Println("9. Health check:")
	err = messageCatalogService.HealthCheck(ctx)
	if err != nil {
		log.Printf("Health check failed: %v", err)
	} else {
		fmt.Println("   Service is healthy!")
		fmt.Println()
	}

	fmt.Println("=== Example completed ===")
}
