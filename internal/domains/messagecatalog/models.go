package messagecatalog

import "time"

// Message represents a complete message with structure and translations
type Message struct {
	// Core structure (from messagecatelog.json)
	MessageCode string `json:"message_code"`
	Category    string `json:"category"`
	Severity    string `json:"severity"`
	Component   string `json:"component"`

	// Language-specific content (from messagecatelog-{lang}.json)
	Message             string `json:"message"`
	DetailedDescription string `json:"detailed_description"`
	ResponseAction      string `json:"response_action"`

	// Additional metadata
	Language    string                 `json:"language"`
	CatalogName string                 `json:"catalog_name"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// MessageRequest represents a request for a message
type MessageRequest struct {
	MessageCode string                 `json:"message_code" validate:"required"`
	Language    string                 `json:"language,omitempty"`
	CatalogName string                 `json:"catalog_name" validate:"required"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// MessageResponse represents the formatted response
type MessageResponse struct {
	MessageCode         string                 `json:"message_code"`
	Category            string                 `json:"category"`
	Severity            string                 `json:"severity"`
	Component           string                 `json:"component"`
	Message             string                 `json:"message"`
	DetailedDescription string                 `json:"detailed_description"`
	ResponseAction      string                 `json:"response_action"`
	Language            string                 `json:"language"`
	CatalogName         string                 `json:"catalog_name"`
	FormattedMessage    string                 `json:"formatted_message"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
	Timestamp           time.Time              `json:"timestamp"`
}

// CatalogInfo represents information about a catalog
type CatalogInfo struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Enabled      bool      `json:"enabled"`
	Languages    []string  `json:"languages"`
	MessageCount int       `json:"message_count"`
	LastReloaded time.Time `json:"last_reloaded"`
}

// CatalogStats represents statistics about all catalogs
type CatalogStats struct {
	TotalCatalogs     int            `json:"total_catalogs"`
	TotalMessages     int            `json:"total_messages"`
	LanguagesCount    int            `json:"languages_count"`
	Catalogs          []CatalogInfo  `json:"catalogs"`
	MessagesByCatalog map[string]int `json:"messages_by_catalog"`
	LastReloaded      time.Time      `json:"last_reloaded"`
}
