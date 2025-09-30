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
	"text/template"
	"time"

	"tushartemplategin/pkg/config"
	"tushartemplategin/pkg/errors"
	"tushartemplategin/pkg/interfaces"
)

// MessageCatalogService implements the Service interface
type MessageCatalogService struct {
	config     config.MessageCatalogConfig
	logger     interfaces.Logger
	cache      map[string]map[string]map[string]*Message // catalog -> language -> messageCode -> Message
	cacheMutex sync.RWMutex
	lastReload map[string]time.Time
}

// NewMessageCatalogService creates a new message catalog service
func NewMessageCatalogService(config config.MessageCatalogConfig, logger interfaces.Logger) Service {
	service := &MessageCatalogService{
		config:     config,
		logger:     logger,
		cache:      make(map[string]map[string]map[string]*Message),
		lastReload: make(map[string]time.Time),
	}

	// Load only default language for all catalogs on startup
	if err := service.LoadDefaultLanguageCatalogs(context.Background()); err != nil {
		logger.Error(context.Background(), "Failed to load initial message catalogs", interfaces.Fields{
			"error": err.Error(),
		})
	}

	return service
}

// GetMessage retrieves a complete message with structure and translations
func (s *MessageCatalogService) GetMessage(ctx context.Context, req *MessageRequest) (*MessageResponse, error) {
	s.logger.Debug(ctx, "Getting message from catalog", interfaces.Fields{
		"message_code": req.MessageCode,
		"catalog_name": req.CatalogName,
		"language":     req.Language,
	})

	// Use default language if not specified
	if req.Language == "" {
		req.Language = s.config.DefaultLanguage
	}

	// Check cache first
	if s.config.CacheEnabled {
		s.cacheMutex.RLock()
		if catalogCache, exists := s.cache[req.CatalogName]; exists {
			if langCache, exists := catalogCache[req.Language]; exists {
				if message, exists := langCache[req.MessageCode]; exists {
					s.cacheMutex.RUnlock()
					return s.formatMessageResponse(message, req.Parameters), nil
				}
			}
		}
		s.cacheMutex.RUnlock()
	}

	// Load message from files
	message, err := s.loadMessageFromFiles(ctx, req.CatalogName, req.MessageCode, req.Language)
	if err != nil {
		return nil, err
	}

	// Update cache
	if s.config.CacheEnabled {
		s.cacheMutex.Lock()
		if s.cache[req.CatalogName] == nil {
			s.cache[req.CatalogName] = make(map[string]map[string]*Message)
		}
		if s.cache[req.CatalogName][req.Language] == nil {
			s.cache[req.CatalogName][req.Language] = make(map[string]*Message)
		}
		s.cache[req.CatalogName][req.Language][req.MessageCode] = message
		s.cacheMutex.Unlock()

		// Log when loading a non-default language on-demand
		if req.Language != s.config.DefaultLanguage {
			s.logger.Info(ctx, "Loaded non-default language on-demand", interfaces.Fields{
				"catalog_name": req.CatalogName,
				"language":     req.Language,
				"message_code": req.MessageCode,
			})
		}
	}

	return s.formatMessageResponse(message, req.Parameters), nil
}

// LoadDefaultLanguageCatalogs loads only the default language for all enabled catalogs
func (s *MessageCatalogService) LoadDefaultLanguageCatalogs(ctx context.Context) error {
	s.logger.Info(ctx, "Loading default language catalogs", interfaces.Fields{
		"default_language": s.config.DefaultLanguage,
	})

	for _, catalog := range s.config.Catalogs {
		if catalog.Enabled {
			if err := s.LoadCatalogDefaultLanguage(ctx, catalog.Name); err != nil {
				s.logger.Error(ctx, "Failed to load default language for catalog", interfaces.Fields{
					"catalog_name": catalog.Name,
					"error":        err.Error(),
				})
				return err
			}
		}
	}

	s.logger.Info(ctx, "Default language catalogs loaded successfully", interfaces.Fields{})
	return nil
}

// LoadCatalogDefaultLanguage loads only the default language for a specific catalog
func (s *MessageCatalogService) LoadCatalogDefaultLanguage(ctx context.Context, catalogName string) error {
	s.logger.Info(ctx, "Loading default language for catalog", interfaces.Fields{
		"catalog_name":     catalogName,
		"default_language": s.config.DefaultLanguage,
	})

	// Find catalog configuration
	var catalogConfig *config.CatalogConfig
	for _, catalog := range s.config.Catalogs {
		if catalog.Name == catalogName {
			catalogConfig = &catalog
			break
		}
	}

	if catalogConfig == nil {
		return errors.NewWithDetails(errors.ErrCodeNotFound, "Catalog not found",
			fmt.Sprintf("Catalog '%s' not found in configuration", catalogName), 404)
	}

	// Load structure file
	structureData, err := s.loadStructureFile(ctx, catalogConfig)
	if err != nil {
		return err
	}

	// Load only default language file
	languageData, err := s.loadLanguageFile(ctx, catalogConfig, s.config.DefaultLanguage)
	if err != nil {
		s.logger.Warn(ctx, "Failed to load default language file", interfaces.Fields{
			"catalog_name":     catalogName,
			"default_language": s.config.DefaultLanguage,
			"error":            err.Error(),
		})
		languageData = make(map[string]interface{})
	}

	// Initialize catalog cache with only default language
	catalogMessages := make(map[string]map[string]*Message) // language -> messageCode -> Message
	catalogMessages[s.config.DefaultLanguage] = make(map[string]*Message)

	// Combine structure and default language data
	for messageCode, messageStructure := range structureData {
		structure, ok := messageStructure.(map[string]interface{})
		if !ok {
			continue
		}

		languageContent, exists := languageData[messageCode]
		if !exists {
			languageContent = make(map[string]interface{})
		}

		// Type assert to map[string]interface{}
		languageMap, ok := languageContent.(map[string]interface{})
		if !ok {
			languageMap = make(map[string]interface{})
		}

		message := s.combineMessageData(structure, languageMap, catalogName, s.config.DefaultLanguage)
		catalogMessages[s.config.DefaultLanguage][messageCode] = message
	}

	// Update cache atomically
	s.cacheMutex.Lock()
	s.cache[catalogName] = catalogMessages
	s.lastReload[catalogName] = time.Now()
	s.cacheMutex.Unlock()

	// Count total messages for default language
	totalMessages := len(catalogMessages[s.config.DefaultLanguage])

	s.logger.Info(ctx, "Default language catalog loaded successfully", interfaces.Fields{
		"catalog_name":     catalogName,
		"default_language": s.config.DefaultLanguage,
		"message_count":    totalMessages,
	})

	return nil
}

// GetMessageByCode retrieves a message by code, catalog, and language
func (s *MessageCatalogService) GetMessageByCode(ctx context.Context, messageCode, catalogName, language string) (*MessageResponse, error) {
	req := &MessageRequest{
		MessageCode: messageCode,
		CatalogName: catalogName,
		Language:    language,
	}
	return s.GetMessage(ctx, req)
}

// GetMessagesByCategory retrieves all messages in a category
func (s *MessageCatalogService) GetMessagesByCategory(ctx context.Context, category, catalogName, language string) ([]*MessageResponse, error) {
	s.logger.Debug(ctx, "Getting messages by category", interfaces.Fields{
		"category":     category,
		"catalog_name": catalogName,
		"language":     language,
	})

	// Find catalog configuration
	var catalogConfig *config.CatalogConfig
	for _, catalog := range s.config.Catalogs {
		if catalog.Name == catalogName {
			catalogConfig = &catalog
			break
		}
	}

	if catalogConfig == nil {
		return nil, errors.NewWithDetails(errors.ErrCodeNotFound, "Catalog not found",
			fmt.Sprintf("Catalog '%s' not found in configuration", catalogName), 404)
	}

	// Load structure file to get all messages
	structureData, err := s.loadStructureFile(ctx, catalogConfig)
	if err != nil {
		return nil, err
	}

	var messages []*MessageResponse
	for messageCode, messageData := range structureData {
		messageStructure, ok := messageData.(map[string]interface{})
		if !ok {
			continue
		}

		// Check if message belongs to the requested category
		if msgCategory, exists := messageStructure["category"]; exists {
			if categoryStr, ok := msgCategory.(string); ok && categoryStr == category {
				// Load the complete message
				message, err := s.loadMessageFromFiles(ctx, catalogName, messageCode, language)
				if err != nil {
					s.logger.Warn(ctx, "Failed to load message", interfaces.Fields{
						"message_code": messageCode,
						"error":        err.Error(),
					})
					continue
				}
				messages = append(messages, s.formatMessageResponse(message, nil))
			}
		}
	}

	return messages, nil
}

// GetMessagesBySeverity retrieves all messages with a specific severity
func (s *MessageCatalogService) GetMessagesBySeverity(ctx context.Context, severity, catalogName, language string) ([]*MessageResponse, error) {
	s.logger.Debug(ctx, "Getting messages by severity", interfaces.Fields{
		"severity":     severity,
		"catalog_name": catalogName,
		"language":     language,
	})

	// Find catalog configuration
	var catalogConfig *config.CatalogConfig
	for _, catalog := range s.config.Catalogs {
		if catalog.Name == catalogName {
			catalogConfig = &catalog
			break
		}
	}

	if catalogConfig == nil {
		return nil, errors.NewWithDetails(errors.ErrCodeNotFound, "Catalog not found",
			fmt.Sprintf("Catalog '%s' not found in configuration", catalogName), 404)
	}

	// Load structure file to get all messages
	structureData, err := s.loadStructureFile(ctx, catalogConfig)
	if err != nil {
		return nil, err
	}

	var messages []*MessageResponse
	for messageCode, messageData := range structureData {
		messageStructure, ok := messageData.(map[string]interface{})
		if !ok {
			continue
		}

		// Check if message has the requested severity
		if msgSeverity, exists := messageStructure["severity"]; exists {
			if severityStr, ok := msgSeverity.(string); ok && severityStr == severity {
				// Load the complete message
				message, err := s.loadMessageFromFiles(ctx, catalogName, messageCode, language)
				if err != nil {
					s.logger.Warn(ctx, "Failed to load message", interfaces.Fields{
						"message_code": messageCode,
						"error":        err.Error(),
					})
					continue
				}
				messages = append(messages, s.formatMessageResponse(message, nil))
			}
		}
	}

	return messages, nil
}

// ReloadCatalog reloads a specific catalog
func (s *MessageCatalogService) ReloadCatalog(ctx context.Context, catalogName string) error {
	s.logger.Info(ctx, "Reloading catalog", interfaces.Fields{
		"catalog_name": catalogName,
	})

	// Find catalog configuration
	var catalogConfig *config.CatalogConfig
	for _, catalog := range s.config.Catalogs {
		if catalog.Name == catalogName {
			catalogConfig = &catalog
			break
		}
	}

	if catalogConfig == nil {
		return errors.NewWithDetails(errors.ErrCodeNotFound, "Catalog not found",
			fmt.Sprintf("Catalog '%s' not found in configuration", catalogName), 404)
	}

	// Load structure file
	structureData, err := s.loadStructureFile(ctx, catalogConfig)
	if err != nil {
		return err
	}

	// Load only default language
	catalogMessages := make(map[string]map[string]*Message) // language -> messageCode -> Message
	languageData, err := s.loadLanguageFile(ctx, catalogConfig, s.config.DefaultLanguage)
	if err != nil {
		s.logger.Warn(ctx, "Failed to load default language file", interfaces.Fields{
			"catalog_name":     catalogName,
			"default_language": s.config.DefaultLanguage,
			"error":            err.Error(),
		})
		languageData = make(map[string]interface{})
	}

	// Initialize default language cache
	catalogMessages[s.config.DefaultLanguage] = make(map[string]*Message)

	// Combine structure and default language data
	for messageCode, messageStructure := range structureData {
		structure, ok := messageStructure.(map[string]interface{})
		if !ok {
			continue
		}

		languageContent, exists := languageData[messageCode]
		if !exists {
			languageContent = make(map[string]interface{})
		}

		// Type assert to map[string]interface{}
		languageMap, ok := languageContent.(map[string]interface{})
		if !ok {
			languageMap = make(map[string]interface{})
		}

		message := s.combineMessageData(structure, languageMap, catalogName, s.config.DefaultLanguage)
		catalogMessages[s.config.DefaultLanguage][messageCode] = message
	}

	// Update cache atomically
	s.cacheMutex.Lock()
	s.cache[catalogName] = catalogMessages
	s.lastReload[catalogName] = time.Now()
	s.cacheMutex.Unlock()

	// Count total messages across all languages
	totalMessages := 0
	for _, langMessages := range catalogMessages {
		totalMessages += len(langMessages)
	}

	s.logger.Info(ctx, "Catalog reloaded successfully", interfaces.Fields{
		"catalog_name":  catalogName,
		"message_count": totalMessages,
		"languages":     len(catalogMessages),
	})

	return nil
}

// ReloadAllCatalogs reloads all enabled catalogs
func (s *MessageCatalogService) ReloadAllCatalogs(ctx context.Context) error {
	s.logger.Info(ctx, "Reloading all catalogs", interfaces.Fields{})

	for _, catalog := range s.config.Catalogs {
		if catalog.Enabled {
			if err := s.ReloadCatalog(ctx, catalog.Name); err != nil {
				s.logger.Error(ctx, "Failed to reload catalog", interfaces.Fields{
					"catalog_name": catalog.Name,
					"error":        err.Error(),
				})
				return err
			}
		}
	}

	s.logger.Info(ctx, "All catalogs reloaded successfully", interfaces.Fields{})
	return nil
}

// HealthCheck checks the health of the message catalog service
func (s *MessageCatalogService) HealthCheck(ctx context.Context) error {
	s.logger.Debug(ctx, "Performing health check", interfaces.Fields{})

	// Check if we can list available catalogs
	catalogs, err := s.ListAvailableCatalogs(ctx)
	if err != nil {
		return err
	}

	if len(catalogs) == 0 {
		return errors.NewWithDetails(errors.ErrCodeInternalServer, "No catalogs available", "No catalogs are loaded", 500)
	}

	s.logger.Debug(ctx, "Health check passed", interfaces.Fields{
		"catalog_count": len(catalogs),
	})

	return nil
}

// ListAvailableCatalogs returns a list of available catalogs
func (s *MessageCatalogService) ListAvailableCatalogs(ctx context.Context) ([]string, error) {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	var catalogs []string
	for catalogName := range s.cache {
		catalogs = append(catalogs, catalogName)
	}

	return catalogs, nil
}

// ListAvailableLanguages returns available languages for a catalog
func (s *MessageCatalogService) ListAvailableLanguages(ctx context.Context, catalogName string) ([]string, error) {
	s.logger.Debug(ctx, "Listing available languages", interfaces.Fields{
		"catalog_name": catalogName,
	})

	// Find catalog configuration
	var catalogConfig *config.CatalogConfig
	for _, catalog := range s.config.Catalogs {
		if catalog.Name == catalogName {
			catalogConfig = &catalog
			break
		}
	}

	if catalogConfig == nil {
		return nil, errors.NewWithDetails(errors.ErrCodeNotFound, "Catalog not found",
			fmt.Sprintf("Catalog '%s' not found in configuration", catalogName), 404)
	}

	// Discover available languages dynamically
	return s.discoverAvailableLanguages(ctx, catalogConfig)
}

// discoverAvailableLanguages discovers available language files for a catalog
func (s *MessageCatalogService) discoverAvailableLanguages(ctx context.Context, catalogConfig *config.CatalogConfig) ([]string, error) {
	languages := []string{}

	// Always include default language
	languages = append(languages, s.config.DefaultLanguage)

	// Scan directory for language files matching the pattern
	pattern := filepath.Join(catalogConfig.Path, strings.Replace(catalogConfig.LanguageFilePattern, "{lang}", "*", 1))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		s.logger.Warn(ctx, "Failed to scan for language files", interfaces.Fields{
			"catalog_name": catalogConfig.Name,
			"pattern":      pattern,
			"error":        err.Error(),
		})
		return languages, nil // Return at least default language
	}

	// Extract language codes from file names
	for _, match := range matches {
		filename := filepath.Base(match)
		// Extract language from filename like "messagecatelog-en-US.json"
		if strings.HasPrefix(filename, "messagecatelog-") && strings.HasSuffix(filename, ".json") {
			langCode := strings.TrimPrefix(filename, "messagecatelog-")
			langCode = strings.TrimSuffix(langCode, ".json")
			if langCode != s.config.DefaultLanguage {
				languages = append(languages, langCode)
			}
		}
	}

	s.logger.Debug(ctx, "Discovered available languages", interfaces.Fields{
		"catalog_name": catalogConfig.Name,
		"languages":    languages,
	})

	return languages, nil
}

// GetCatalogInfo returns information about a specific catalog
func (s *MessageCatalogService) GetCatalogInfo(ctx context.Context, catalogName string) (*CatalogInfo, error) {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	// Find catalog configuration
	var catalogConfig *config.CatalogConfig
	for _, catalog := range s.config.Catalogs {
		if catalog.Name == catalogName {
			catalogConfig = &catalog
			break
		}
	}

	if catalogConfig == nil {
		return nil, errors.NewWithDetails(errors.ErrCodeNotFound, "Catalog not found",
			fmt.Sprintf("Catalog '%s' not found in configuration", catalogName), 404)
	}

	messageCount := 0
	if catalogCache, exists := s.cache[catalogName]; exists {
		for _, langMessages := range catalogCache {
			messageCount += len(langMessages)
		}
	}

	lastReloaded := time.Time{}
	if reloadTime, exists := s.lastReload[catalogName]; exists {
		lastReloaded = reloadTime
	}

	// Get available languages dynamically
	availableLanguages, err := s.discoverAvailableLanguages(ctx, catalogConfig)
	if err != nil {
		s.logger.Warn(ctx, "Failed to discover available languages", interfaces.Fields{
			"catalog_name": catalogName,
			"error":        err.Error(),
		})
		availableLanguages = []string{s.config.DefaultLanguage} // Fallback to default
	}

	return &CatalogInfo{
		Name:         catalogName,
		Path:         catalogConfig.Path,
		Enabled:      catalogConfig.Enabled,
		Languages:    availableLanguages,
		MessageCount: messageCount,
		LastReloaded: lastReloaded,
	}, nil
}

// GetCatalogStats returns statistics about all catalogs
func (s *MessageCatalogService) GetCatalogStats(ctx context.Context) (*CatalogStats, error) {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	// Count total unique languages across all catalogs
	allLanguages := make(map[string]bool)
	for catalogName := range s.cache {
		// Find catalog configuration
		var catalogConfig *config.CatalogConfig
		for _, catalog := range s.config.Catalogs {
			if catalog.Name == catalogName {
				catalogConfig = &catalog
				break
			}
		}
		if catalogConfig != nil {
			availableLanguages, err := s.discoverAvailableLanguages(ctx, catalogConfig)
			if err == nil {
				for _, lang := range availableLanguages {
					allLanguages[lang] = true
				}
			}
		}
	}

	stats := &CatalogStats{
		TotalCatalogs:     len(s.cache),
		TotalMessages:     0,
		LanguagesCount:    len(allLanguages),
		Catalogs:          []CatalogInfo{},
		MessagesByCatalog: make(map[string]int),
		LastReloaded:      time.Time{},
	}

	for catalogName, catalogCache := range s.cache {
		messageCount := 0
		for _, langMessages := range catalogCache {
			messageCount += len(langMessages)
		}
		stats.TotalMessages += messageCount
		stats.MessagesByCatalog[catalogName] = messageCount

		// Get catalog info
		catalogInfo, err := s.GetCatalogInfo(ctx, catalogName)
		if err == nil {
			stats.Catalogs = append(stats.Catalogs, *catalogInfo)
		}

		// Track latest reload time
		if reloadTime, exists := s.lastReload[catalogName]; exists {
			if reloadTime.After(stats.LastReloaded) {
				stats.LastReloaded = reloadTime
			}
		}
	}

	return stats, nil
}

// loadMessageFromFiles loads and combines structure and language files
func (s *MessageCatalogService) loadMessageFromFiles(ctx context.Context, catalogName, messageCode, language string) (*Message, error) {
	// Find catalog configuration
	var catalogConfig *config.CatalogConfig
	for _, catalog := range s.config.Catalogs {
		if catalog.Name == catalogName {
			catalogConfig = &catalog
			break
		}
	}

	if catalogConfig == nil {
		return nil, errors.NewWithDetails(errors.ErrCodeNotFound, "Catalog not found",
			fmt.Sprintf("Catalog '%s' not found in configuration", catalogName), 404)
	}

	// Load structure file
	structureData, err := s.loadStructureFile(ctx, catalogConfig)
	if err != nil {
		return nil, err
	}

	// Get message structure
	messageStructure, exists := structureData[messageCode]
	if !exists {
		return nil, errors.NewWithDetails(errors.ErrCodeNotFound, "Message not found",
			fmt.Sprintf("Message '%s' not found in catalog '%s'", messageCode, catalogName), 404)
	}

	structure, ok := messageStructure.(map[string]interface{})
	if !ok {
		return nil, errors.NewWithDetails(errors.ErrCodeBadRequest, "Invalid message structure",
			fmt.Sprintf("Message '%s' has invalid structure in catalog '%s'", messageCode, catalogName), 400)
	}

	// Load language file
	languageData, err := s.loadLanguageFile(ctx, catalogConfig, language)
	if err != nil {
		s.logger.Warn(ctx, "Language content not found, using structure only", interfaces.Fields{
			"message_code": messageCode,
			"catalog_name": catalogName,
			"language":     language,
		})
		languageData = make(map[string]interface{})
	}

	// Get language-specific content
	languageContent, exists := languageData[messageCode]
	if !exists {
		languageContent = make(map[string]interface{})
	}

	// Type assert to map[string]interface{}
	languageMap, ok := languageContent.(map[string]interface{})
	if !ok {
		languageMap = make(map[string]interface{})
	}

	// Combine structure and language content
	message := s.combineMessageData(structure, languageMap, catalogName, language)

	return message, nil
}

// loadStructureFile loads the structure file for a catalog
func (s *MessageCatalogService) loadStructureFile(ctx context.Context, catalogConfig *config.CatalogConfig) (map[string]interface{}, error) {
	filePath := filepath.Join(catalogConfig.Path, catalogConfig.StructureFile)

	s.logger.Debug(ctx, "Loading structure file", interfaces.Fields{
		"file_path": filePath,
		"catalog":   catalogConfig.Name,
	})

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errors.NewWithDetails(errors.ErrCodeNotFound, "Structure file not found",
			fmt.Sprintf("Structure file '%s' not found for catalog '%s'", filePath, catalogConfig.Name), 404)
	}

	// Read file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.NewWithError(errors.ErrCodeInternalServer, "Failed to read structure file", 500, err)
	}

	// Parse JSON
	var structureData map[string]interface{}
	if err := json.Unmarshal(data, &structureData); err != nil {
		return nil, errors.NewWithError(errors.ErrCodeBadRequest, "Failed to parse structure file", 400, err)
	}

	return structureData, nil
}

// loadLanguageFile loads a language file for a catalog
func (s *MessageCatalogService) loadLanguageFile(ctx context.Context, catalogConfig *config.CatalogConfig, language string) (map[string]interface{}, error) {
	fileName := strings.ReplaceAll(catalogConfig.LanguageFilePattern, "{lang}", language)
	filePath := filepath.Join(catalogConfig.Path, fileName)

	s.logger.Debug(ctx, "Loading language file", interfaces.Fields{
		"file_path": filePath,
		"catalog":   catalogConfig.Name,
		"language":  language,
	})

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errors.NewWithDetails(errors.ErrCodeNotFound, "Language file not found",
			fmt.Sprintf("Language file '%s' not found for catalog '%s'", filePath, catalogConfig.Name), 404)
	}

	// Read file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.NewWithError(errors.ErrCodeInternalServer, "Failed to read language file", 500, err)
	}

	// Parse JSON
	var languageData map[string]interface{}
	if err := json.Unmarshal(data, &languageData); err != nil {
		return nil, errors.NewWithError(errors.ErrCodeBadRequest, "Failed to parse language file", 400, err)
	}

	return languageData, nil
}

// combineMessageData combines structure and language data into a Message object
func (s *MessageCatalogService) combineMessageData(structure, language map[string]interface{}, catalogName, languageCode string) *Message {
	message := &Message{
		MessageCode: getString(structure, "message_code"),
		Category:    getString(structure, "category"),
		Severity:    getString(structure, "severity"),
		Component:   getString(structure, "component"),

		Message:             getString(language, "message"),
		DetailedDescription: getString(language, "detailed_description"),
		ResponseAction:      getString(language, "response_action"),

		Language:    languageCode,
		CatalogName: catalogName,
		Metadata:    make(map[string]interface{}),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return message
}

// formatMessageResponse formats the message with parameters
func (s *MessageCatalogService) formatMessageResponse(message *Message, parameters map[string]interface{}) *MessageResponse {
	response := &MessageResponse{
		MessageCode:         message.MessageCode,
		Category:            message.Category,
		Severity:            message.Severity,
		Component:           message.Component,
		Message:             message.Message,
		DetailedDescription: message.DetailedDescription,
		ResponseAction:      message.ResponseAction,
		Language:            message.Language,
		CatalogName:         message.CatalogName,
		Metadata:            message.Metadata,
		Timestamp:           time.Now(),
	}

	// Apply parameter substitution if parameters provided
	if parameters != nil && len(parameters) > 0 {
		response.FormattedMessage = s.applyParameters(message.Message, parameters)
		response.DetailedDescription = s.applyParameters(message.DetailedDescription, parameters)
	} else {
		response.FormattedMessage = message.Message
	}

	return response
}

// applyParameters applies template parameters to message content
func (s *MessageCatalogService) applyParameters(content string, parameters map[string]interface{}) string {
	if content == "" {
		return content
	}

	tmpl, err := template.New("message").Parse(content)
	if err != nil {
		s.logger.Warn(context.Background(), "Failed to parse message template", interfaces.Fields{
			"content": content,
			"error":   err.Error(),
		})
		return content
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, parameters); err != nil {
		s.logger.Warn(context.Background(), "Failed to execute message template", interfaces.Fields{
			"content": content,
			"error":   err.Error(),
		})
		return content
	}

	return result.String()
}

// Helper function to safely get string values from map
func getString(data map[string]interface{}, key string) string {
	if value, exists := data[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}
