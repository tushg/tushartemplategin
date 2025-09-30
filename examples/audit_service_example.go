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

// AuditService represents an example audit service
type AuditService struct {
	messageCatalog messagecatalog.Service
	logger         interfaces.Logger
}

// NewAuditService creates a new audit service
func NewAuditService(messageCatalog messagecatalog.Service, logger interfaces.Logger) *AuditService {
	return &AuditService{
		messageCatalog: messageCatalog,
		logger:         logger,
	}
}

// LogEvent logs an audit event and returns formatted response
func (s *AuditService) LogEvent(ctx context.Context, eventCode string, parameters map[string]interface{}) (*AuditEvent, error) {
	s.logger.Info(ctx, "Logging audit event", interfaces.Fields{
		"event_code": eventCode,
		"parameters": parameters,
	})

	// Get message from catalog
	req := &messagecatalog.MessageRequest{
		MessageCode: eventCode,
		CatalogName: "audit",
		Language:    "en-US", // Could be configurable
		Parameters:  parameters,
	}

	message, err := s.messageCatalog.GetMessage(ctx, req)
	if err != nil {
		s.logger.Error(ctx, "Failed to get audit message", interfaces.Fields{
			"event_code": eventCode,
			"error":      err.Error(),
		})
		return nil, err
	}

	// Convert to audit-specific response
	event := &AuditEvent{
		EventCode:     message.MessageCode,
		EventCategory: message.Category,
		RiskLevel:     message.Severity,
		Component:     message.Component,
		Description:   message.FormattedMessage,
		Details:       message.DetailedDescription,
		Action:        message.ResponseAction,
		Language:      message.Language,
		Timestamp:     time.Now(),
	}

	s.logger.Info(ctx, "Audit event logged successfully", interfaces.Fields{
		"event_code": eventCode,
		"category":   message.Category,
		"risk_level": message.Severity,
	})

	return event, nil
}

// GetAuditMessage gets a specific audit message
func (s *AuditService) GetAuditMessage(ctx context.Context, code, language string, parameters map[string]interface{}) (*AuditMessage, error) {
	req := &messagecatalog.MessageRequest{
		MessageCode: code,
		CatalogName: "audit",
		Language:    language,
		Parameters:  parameters,
	}

	message, err := s.messageCatalog.GetMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	return &AuditMessage{
		EventCode:     message.MessageCode,
		EventCategory: message.Category,
		RiskLevel:     message.Severity,
		Description:   message.FormattedMessage,
		Details:       message.DetailedDescription,
		Action:        message.ResponseAction,
	}, nil
}

// ListEventsByCategory gets all events in a category
func (s *AuditService) ListEventsByCategory(ctx context.Context, category, language string) ([]*AuditMessage, error) {
	messages, err := s.messageCatalog.GetMessagesByCategory(ctx, category, "audit", language)
	if err != nil {
		return nil, err
	}

	var events []*AuditMessage
	for _, message := range messages {
		events = append(events, &AuditMessage{
			EventCode:     message.MessageCode,
			EventCategory: message.Category,
			RiskLevel:     message.Severity,
			Description:   message.FormattedMessage,
			Details:       message.DetailedDescription,
			Action:        message.ResponseAction,
		})
	}

	return events, nil
}

// AuditEvent represents an audit event
type AuditEvent struct {
	EventCode     string    `json:"event_code"`
	EventCategory string    `json:"event_category"`
	RiskLevel     string    `json:"risk_level"`
	Component     string    `json:"component"`
	Description   string    `json:"description"`
	Details       string    `json:"details"`
	Action        string    `json:"action"`
	Language      string    `json:"language"`
	Timestamp     time.Time `json:"timestamp"`
}

// AuditMessage represents a basic audit message
type AuditMessage struct {
	EventCode     string `json:"event_code"`
	EventCategory string `json:"event_category"`
	RiskLevel     string `json:"risk_level"`
	Description   string `json:"description"`
	Details       string `json:"details"`
	Action        string `json:"action"`
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

	// Create audit service
	auditService := NewAuditService(messageCatalogService, appLogger)

	ctx := context.Background()

	fmt.Println("=== Audit Service Example ===\n")

	// Example 1: Log user access event
	fmt.Println("1. Logging user access event:")
	event1, err := auditService.LogEvent(ctx, "AUD0001", map[string]interface{}{
		"username":         "john.doe",
		"action":           "logged in",
		"timestamp":        time.Now().Format("2006-01-02 15:04:05"),
		"ip_address":       "192.168.1.100",
		"session_duration": "2h 30m",
	})
	if err != nil {
		log.Printf("Error logging event: %v", err)
	} else {
		fmt.Printf("   Event Code: %s\n", event1.EventCode)
		fmt.Printf("   Category: %s\n", event1.EventCategory)
		fmt.Printf("   Risk Level: %s\n", event1.RiskLevel)
		fmt.Printf("   Component: %s\n", event1.Component)
		fmt.Printf("   Description: %s\n", event1.Description)
		fmt.Printf("   Details: %s\n", event1.Details)
		fmt.Printf("   Action: %s\n", event1.Action)
		fmt.Printf("   Language: %s\n", event1.Language)
		fmt.Printf("   Timestamp: %s\n", event1.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Println()
	}

	// Example 2: Log data modification event
	fmt.Println("2. Logging data modification event:")
	event2, err := auditService.LogEvent(ctx, "AUD0002", map[string]interface{}{
		"table_name":     "users",
		"username":       "admin",
		"timestamp":      time.Now().Format("2006-01-02 15:04:05"),
		"record_count":   "150",
		"operation_type": "UPDATE",
	})
	if err != nil {
		log.Printf("Error logging event: %v", err)
	} else {
		fmt.Printf("   Event Code: %s\n", event2.EventCode)
		fmt.Printf("   Category: %s\n", event2.EventCategory)
		fmt.Printf("   Risk Level: %s\n", event2.RiskLevel)
		fmt.Printf("   Description: %s\n", event2.Description)
		fmt.Printf("   Details: %s\n", event2.Details)
		fmt.Printf("   Action: %s\n", event2.Action)
		fmt.Println()
	}

	// Example 3: Log security event
	fmt.Println("3. Logging security event:")
	event3, err := auditService.LogEvent(ctx, "AUD0004", map[string]interface{}{
		"security_event_type": "unauthorized_access_attempt",
		"timestamp":           time.Now().Format("2006-01-02 15:04:05"),
		"source_ip":           "10.0.0.50",
		"target_resource":     "/api/admin/users",
		"risk_level":          "HIGH",
	})
	if err != nil {
		log.Printf("Error logging event: %v", err)
	} else {
		fmt.Printf("   Event Code: %s\n", event3.EventCode)
		fmt.Printf("   Category: %s\n", event3.EventCategory)
		fmt.Printf("   Risk Level: %s\n", event3.RiskLevel)
		fmt.Printf("   Description: %s\n", event3.Description)
		fmt.Printf("   Details: %s\n", event3.Details)
		fmt.Printf("   Action: %s\n", event3.Action)
		fmt.Println()
	}

	// Example 4: Log configuration change event
	fmt.Println("4. Logging configuration change event:")
	event4, err := auditService.LogEvent(ctx, "AUD0005", map[string]interface{}{
		"config_section": "database",
		"username":       "admin",
		"timestamp":      time.Now().Format("2006-01-02 15:04:05"),
		"old_value":      "localhost:5432",
		"new_value":      "prod-db:5432",
	})
	if err != nil {
		log.Printf("Error logging event: %v", err)
	} else {
		fmt.Printf("   Event Code: %s\n", event4.EventCode)
		fmt.Printf("   Category: %s\n", event4.EventCategory)
		fmt.Printf("   Risk Level: %s\n", event4.RiskLevel)
		fmt.Printf("   Description: %s\n", event4.Description)
		fmt.Printf("   Details: %s\n", event4.Details)
		fmt.Printf("   Action: %s\n", event4.Action)
		fmt.Println()
	}

	// Example 5: Get audit message in French
	fmt.Println("5. Getting audit message in French:")
	auditFr, err := auditService.GetAuditMessage(ctx, "AUD0001", "fr-FR", map[string]interface{}{
		"username":         "john.doe",
		"action":           "logged in",
		"timestamp":        time.Now().Format("2006-01-02 15:04:05"),
		"ip_address":       "192.168.1.100",
		"session_duration": "2h 30m",
	})
	if err != nil {
		log.Printf("Error getting French audit message: %v", err)
	} else {
		fmt.Printf("   Event Code: %s\n", auditFr.EventCode)
		fmt.Printf("   Category: %s\n", auditFr.EventCategory)
		fmt.Printf("   Risk Level: %s\n", auditFr.RiskLevel)
		fmt.Printf("   Description: %s\n", auditFr.Description)
		fmt.Printf("   Details: %s\n", auditFr.Details)
		fmt.Printf("   Action: %s\n", auditFr.Action)
		fmt.Println()
	}

	// Example 6: List events by category
	fmt.Println("6. Listing events by category (UserAccess):")
	userAccessEvents, err := auditService.ListEventsByCategory(ctx, "UserAccess", "en-US")
	if err != nil {
		log.Printf("Error listing events by category: %v", err)
	} else {
		for i, event := range userAccessEvents {
			fmt.Printf("   %d. %s - %s (%s)\n", i+1, event.EventCode, event.Description, event.RiskLevel)
		}
		fmt.Println()
	}

	// Example 7: List security events
	fmt.Println("7. Listing security events:")
	securityEvents, err := auditService.ListEventsByCategory(ctx, "SecurityEvent", "en-US")
	if err != nil {
		log.Printf("Error listing security events: %v", err)
	} else {
		for i, event := range securityEvents {
			fmt.Printf("   %d. %s - %s (%s)\n", i+1, event.EventCode, event.Description, event.RiskLevel)
		}
		fmt.Println()
	}

	fmt.Println("=== Audit Service Example completed ===")
}
