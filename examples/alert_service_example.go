package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"tushartemplategin/internal/domains/messagecatalog"
	"tushartemplategin/pkg/config"
	"tushartemplategin/pkg/interfaces"
	"tushartemplategin/pkg/logger"
)

// AlertService represents an example alert service
type AlertService struct {
	messageCatalog messagecatalog.Service
	logger         interfaces.Logger
}

// NewAlertService creates a new alert service
func NewAlertService(messageCatalog messagecatalog.Service, logger interfaces.Logger) *AlertService {
	return &AlertService{
		messageCatalog: messageCatalog,
		logger:         logger,
	}
}

// ProcessAlert processes an alert and returns formatted response
func (s *AlertService) ProcessAlert(ctx context.Context, alertCode string, parameters map[string]interface{}) (*AlertResponse, error) {
	s.logger.Info(ctx, "Processing alert", interfaces.Fields{
		"alert_code": alertCode,
		"parameters": parameters,
	})

	// Get message from catalog
	req := &messagecatalog.MessageRequest{
		MessageCode: alertCode,
		CatalogName: "alert",
		Language:    "en-US", // Could be configurable
		Parameters:  parameters,
	}

	message, err := s.messageCatalog.GetMessage(ctx, req)
	if err != nil {
		s.logger.Error(ctx, "Failed to get alert message", interfaces.Fields{
			"alert_code": alertCode,
			"error":      err.Error(),
		})
		return nil, err
	}

	// Convert to alert-specific response
	response := &AlertResponse{
		Code:        message.MessageCode,
		Category:    message.Category,
		Severity:    message.Severity,
		Component:   message.Component,
		Message:     message.FormattedMessage,
		Description: message.DetailedDescription,
		Action:      message.ResponseAction,
		Language:    message.Language,
		Timestamp:   time.Now(),
	}

	s.logger.Info(ctx, "Alert processed successfully", interfaces.Fields{
		"alert_code": alertCode,
		"category":   message.Category,
		"severity":   message.Severity,
	})

	return response, nil
}

// GetAlertMessage gets a specific alert message
func (s *AlertService) GetAlertMessage(ctx context.Context, code, language string, parameters map[string]interface{}) (*AlertMessage, error) {
	req := &messagecatalog.MessageRequest{
		MessageCode: code,
		CatalogName: "alert",
		Language:    language,
		Parameters:  parameters,
	}

	message, err := s.messageCatalog.GetMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	return &AlertMessage{
		Code:        message.MessageCode,
		Category:    message.Category,
		Severity:    message.Severity,
		Message:     message.FormattedMessage,
		Description: message.DetailedDescription,
		Action:      message.ResponseAction,
	}, nil
}

// ListAlertsByCategory gets all alerts in a category
func (s *AlertService) ListAlertsByCategory(ctx context.Context, category, language string) ([]*AlertMessage, error) {
	messages, err := s.messageCatalog.GetMessagesByCategory(ctx, category, "alert", language)
	if err != nil {
		return nil, err
	}

	var alerts []*AlertMessage
	for _, message := range messages {
		alerts = append(alerts, &AlertMessage{
			Code:        message.MessageCode,
			Category:    message.Category,
			Severity:    message.Severity,
			Message:     message.FormattedMessage,
			Description: message.DetailedDescription,
			Action:      message.ResponseAction,
		})
	}

	return alerts, nil
}

// AlertResponse represents the response for an alert
type AlertResponse struct {
	Code        string    `json:"code"`
	Category    string    `json:"category"`
	Severity    string    `json:"severity"`
	Component   string    `json:"component"`
	Message     string    `json:"message"`
	Description string    `json:"description"`
	Action      string    `json:"action"`
	Language    string    `json:"language"`
	Timestamp   time.Time `json:"timestamp"`
}

// AlertMessage represents a basic alert message
type AlertMessage struct {
	Code        string `json:"code"`
	Category    string `json:"category"`
	Severity    string `json:"severity"`
	Message     string `json:"message"`
	Description string `json:"description"`
	Action      string `json:"action"`
}

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

	// Create alert service
	alertService := NewAlertService(messageCatalogService, appLogger)

	ctx := context.Background()

	fmt.Println("=== Alert Service Example ===\n")

	// Example 1: Process a registration alert
	fmt.Println("1. Processing registration alert:")
	alert1, err := alertService.ProcessAlert(ctx, "ABC0001", map[string]interface{}{
		"exp_date": "2024-12-31",
	})
	if err != nil {
		log.Printf("Error processing alert: %v", err)
	} else {
		fmt.Printf("   Code: %s\n", alert1.Code)
		fmt.Printf("   Category: %s\n", alert1.Category)
		fmt.Printf("   Severity: %s\n", alert1.Severity)
		fmt.Printf("   Component: %s\n", alert1.Component)
		fmt.Printf("   Message: %s\n", alert1.Message)
		fmt.Printf("   Description: %s\n", alert1.Description)
		fmt.Printf("   Action: %s\n", alert1.Action)
		fmt.Printf("   Language: %s\n", alert1.Language)
		fmt.Printf("   Timestamp: %s\n", alert1.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Println()
	}

	// Example 2: Process an authentication alert
	fmt.Println("2. Processing authentication alert:")
	alert2, err := alertService.ProcessAlert(ctx, "ABC0002", map[string]interface{}{
		"username":  "john.doe",
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		log.Printf("Error processing alert: %v", err)
	} else {
		fmt.Printf("   Code: %s\n", alert2.Code)
		fmt.Printf("   Category: %s\n", alert2.Category)
		fmt.Printf("   Severity: %s\n", alert2.Severity)
		fmt.Printf("   Message: %s\n", alert2.Message)
		fmt.Printf("   Description: %s\n", alert2.Description)
		fmt.Printf("   Action: %s\n", alert2.Action)
		fmt.Println()
	}

	// Example 3: Process a security alert
	fmt.Println("3. Processing security alert:")
	alert3, err := alertService.ProcessAlert(ctx, "ABC0005", map[string]interface{}{
		"ip_address":    "192.168.1.100",
		"timestamp":     time.Now().Format("2006-01-02 15:04:05"),
		"activity_type": "brute force attack",
	})
	if err != nil {
		log.Printf("Error processing alert: %v", err)
	} else {
		fmt.Printf("   Code: %s\n", alert3.Code)
		fmt.Printf("   Category: %s\n", alert3.Category)
		fmt.Printf("   Severity: %s\n", alert3.Severity)
		fmt.Printf("   Message: %s\n", alert3.Message)
		fmt.Printf("   Description: %s\n", alert3.Description)
		fmt.Printf("   Action: %s\n", alert3.Action)
		fmt.Println()
	}

	// Example 4: Get alert in French
	fmt.Println("4. Getting alert in French:")
	alertFr, err := alertService.GetAlertMessage(ctx, "ABC0001", "fr-FR", map[string]interface{}{
		"exp_date": "2024-12-31",
	})
	if err != nil {
		log.Printf("Error getting French alert: %v", err)
	} else {
		fmt.Printf("   Code: %s\n", alertFr.Code)
		fmt.Printf("   Category: %s\n", alertFr.Category)
		fmt.Printf("   Severity: %s\n", alertFr.Severity)
		fmt.Printf("   Message: %s\n", alertFr.Message)
		fmt.Printf("   Description: %s\n", alertFr.Description)
		fmt.Printf("   Action: %s\n", alertFr.Action)
		fmt.Println()
	}

	// Example 5: List alerts by category
	fmt.Println("5. Listing alerts by category (Registration):")
	registrationAlerts, err := alertService.ListAlertsByCategory(ctx, "Registration", "en-US")
	if err != nil {
		log.Printf("Error listing alerts by category: %v", err)
	} else {
		for i, alert := range registrationAlerts {
			fmt.Printf("   %d. %s - %s (%s)\n", i+1, alert.Code, alert.Message, alert.Severity)
		}
		fmt.Println()
	}

	// Example 6: List critical alerts
	fmt.Println("6. Listing critical alerts:")
	criticalAlerts, err := alertService.ListAlertsByCategory(ctx, "Security", "en-US")
	if err != nil {
		log.Printf("Error listing critical alerts: %v", err)
	} else {
		for i, alert := range criticalAlerts {
			fmt.Printf("   %d. %s - %s (%s)\n", i+1, alert.Code, alert.Message, alert.Severity)
		}
		fmt.Println()
	}

	fmt.Println("=== Alert Service Example completed ===")
}
