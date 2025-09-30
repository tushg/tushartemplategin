package messagecatalog

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"tushartemplategin/pkg/config"
	"tushartemplategin/pkg/errors"
	"tushartemplategin/pkg/interfaces"
)

// MessageCatalogService implements the Service interface
type MessageCatalogService struct {
	config     config.MessageCatalogConfig
	logger     interfaces.Logger
	cache      map[string]map[string]*Message // language -> messageCode -> Message
	cacheMutex sync.RWMutex
	lastReload time.Time
}

// NewMessageCatalogService creates a new message catalog service
func NewMessageCatalogService(config config.MessageCatalogConfig, logger interfaces.Logger) Service {
	service := &MessageCatalogService{
		config: config,
		logger: logger,
		cache:  make(map[string]map[string]*Message),
	}

	// Load initial catalog
	if err := service.ReloadCatalog(context.Background()); err != nil {
		logger.Error(context.Background(), "Failed to load initial message catalog", interfaces.Fields{
			"error": err.Error(),
		})
	}

	return service
}

// GetMessage retrieves a message by code and language
func (s *MessageCatalogService) GetMessage(ctx context.Context, messageCode, language string) (*Message, error) {
	s.logger.Debug(ctx, "Getting message from catalog", interfaces.Fields{
		"message_code": messageCode,
		"language":     language,
	})

	// Use default language if not specified
	if language == "" {
		language = s.config.DefaultLanguage
	}

	// Check cache first
	if s.config.CacheEnabled {
		s.cacheMutex.RLock()
		if langCache, exists := s.cache[language]; exists {
			if message, exists := langCache[messageCode]; exists {
				s.cacheMutex.RUnlock()
				return message, nil
			}
		}
		s.cacheMutex.RUnlock()
	}

	// Load from file
	messages, err := s.loadMessagesFromFile(ctx, language)
	if err != nil {
		return nil, err
	}

	message, exists := messages[messageCode]
	if !exists {
		s.logger.Warn(ctx, "Message not found in catalog", interfaces.Fields{
			"message_code": messageCode,
			"language":     language,
		})
		return nil, errors.NewWithDetails(errors.ErrCodeNotFound, "Message not found",
			fmt.Sprintf("Message with code '%s' not found in language '%s'", messageCode, language), 404)
	}

	// Update cache
	if s.config.CacheEnabled {
		s.cacheMutex.Lock()
		if s.cache[language] == nil {
			s.cache[language] = make(map[string]*Message)
		}
		s.cache[language][messageCode] = message
		s.cacheMutex.Unlock()
	}

	return message, nil
}

// GetMessageByCode retrieves a message by code using default language
func (s *MessageCatalogService) GetMessageByCode(ctx context.Context, messageCode string) (*Message, error) {
	return s.GetMessage(ctx, messageCode, s.config.DefaultLanguage)
}

// GetMessagesByCategory retrieves all messages for a specific category and language
func (s *MessageCatalogService) GetMessagesByCategory(ctx context.Context, category, language string) ([]*Message, error) {
	s.logger.Debug(ctx, "Getting messages by category", interfaces.Fields{
		"category": category,
		"language": language,
	})

	if language == "" {
		language = s.config.DefaultLanguage
	}

	messages, err := s.loadMessagesFromFile(ctx, language)
	if err != nil {
		return nil, errors.NewWithError(errors.ErrCodeInternalServer, "Failed to load message catalog", 500, err)
	}

	var result []*Message
	for _, message := range messages {
		if message.Category == category {
			result = append(result, message)
		}
	}

	return result, nil
}

// GetMessagesBySeverity retrieves all messages for a specific severity and language
func (s *MessageCatalogService) GetMessagesBySeverity(ctx context.Context, severity, language string) ([]*Message, error) {
	s.logger.Debug(ctx, "Getting messages by severity", interfaces.Fields{
		"severity": severity,
		"language": language,
	})

	if language == "" {
		language = s.config.DefaultLanguage
	}

	messages, err := s.loadMessagesFromFile(ctx, language)
	if err != nil {
		return nil, errors.NewWithError(errors.ErrCodeInternalServer, "Failed to load message catalog", 500, err)
	}

	var result []*Message
	for _, message := range messages {
		if message.Severity == severity {
			result = append(result, message)
		}
	}

	return result, nil
}

// ListAvailableLanguages returns all available languages
func (s *MessageCatalogService) ListAvailableLanguages(ctx context.Context) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(s.config.CatalogPath, "*.json"))
	if err != nil {
		return nil, errors.NewWithError(errors.ErrCodeInternalServer, "Failed to scan catalog directory", 500, err)
	}

	var languages []string
	for _, file := range files {
		baseName := filepath.Base(file)
		language := strings.TrimSuffix(baseName, ".json")
		languages = append(languages, language)
	}

	return languages, nil
}

// ReloadCatalog reloads the message catalog from files
func (s *MessageCatalogService) ReloadCatalog(ctx context.Context) error {
	s.logger.Info(ctx, "Reloading message catalog", interfaces.Fields{})

	allMessages := make(map[string]map[string]*Message)
	languages, err := s.ListAvailableLanguages(ctx)
	if err != nil {
		return err
	}

	for _, language := range languages {
		messages, err := s.loadMessagesFromFile(ctx, language)
		if err != nil {
			s.logger.Warn(ctx, "Failed to load language", interfaces.Fields{
				"language": language,
				"error":    err.Error(),
			})
			continue
		}
		allMessages[language] = messages
	}

	// Update cache
	if s.config.CacheEnabled {
		s.cacheMutex.Lock()
		s.cache = allMessages
		s.lastReload = time.Now()
		s.cacheMutex.Unlock()
	}

	s.logger.Info(ctx, "Message catalog reloaded successfully", interfaces.Fields{
		"languages_count": len(allMessages),
	})

	return nil
}

// HealthCheck checks the health of the message catalog service
func (s *MessageCatalogService) HealthCheck(ctx context.Context) error {
	languages, err := s.ListAvailableLanguages(ctx)
	if err != nil {
		return fmt.Errorf("failed to list available languages: %w", err)
	}

	if len(languages) == 0 {
		return fmt.Errorf("no languages available in catalog")
	}

	s.logger.Debug(ctx, "Message catalog health check passed", interfaces.Fields{
		"languages_count": len(languages),
	})

	return nil
}

// Helper method to load messages from file
func (s *MessageCatalogService) loadMessagesFromFile(ctx context.Context, language string) (map[string]*Message, error) {
	filePath := filepath.Join(s.config.CatalogPath, fmt.Sprintf("%s.json", language))

	s.logger.Debug(ctx, "Loading messages from file", interfaces.Fields{
		"file_path": filePath,
		"language":  language,
	})

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errors.NewWithDetails(errors.ErrCodeNotFound, "Language file not found",
			fmt.Sprintf("Language file '%s' not found", filePath), 404)
	}

	// Read file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.NewWithError(errors.ErrCodeInternalServer, "Failed to read language file", 500, err)
	}

	// Parse JSON
	var rawMessages map[string]interface{}
	if err := json.Unmarshal(data, &rawMessages); err != nil {
		return nil, errors.NewWithError(errors.ErrCodeBadRequest, "Failed to parse language file", 400, err)
	}

	// Convert to Message objects
	messages := make(map[string]*Message)
	for key, value := range rawMessages {
		messageData, ok := value.(map[string]interface{})
		if !ok {
			s.logger.Warn(ctx, "Invalid message format", interfaces.Fields{
				"key": key,
			})
			continue
		}

		message := &Message{
			MessageCode:         getString(messageData, "messagecode"),
			Category:            getString(messageData, "category"),
			Severity:            getString(messageData, "severity"),
			Message:             getString(messageData, "message"),
			DetailedDescription: getString(messageData, "detailed_description"),
			ResponseAction:      getString(messageData, "response_action"),
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}

		// Parse metadata if exists
		if metadata, exists := messageData["metadata"]; exists {
			if metadataMap, ok := metadata.(map[string]interface{}); ok {
				message.Metadata = make(map[string]string)
				for k, v := range metadataMap {
					if str, ok := v.(string); ok {
						message.Metadata[k] = str
					}
				}
			}
		}

		messages[message.MessageCode] = message
	}

	s.logger.Info(ctx, "Messages loaded successfully", interfaces.Fields{
		"language": language,
		"count":    len(messages),
	})

	return messages, nil
}

// Helper function to safely extract string values from map
func getString(data map[string]interface{}, key string) string {
	if value, exists := data[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}
