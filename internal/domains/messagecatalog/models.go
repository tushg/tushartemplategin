package messagecatalog

import "time"

// Message represents a message in the catalog
type Message struct {
	MessageCode         string            `json:"messagecode"`
	Category            string            `json:"category"`
	Severity            string            `json:"severity"`
	Message             string            `json:"message"`
	DetailedDescription string            `json:"detailed_description"`
	ResponseAction      string            `json:"response_action"`
	Metadata            map[string]string `json:"metadata,omitempty"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"`
}

// MessageRequest represents a request to get a message
type MessageRequest struct {
	MessageCode string            `json:"message_code" validate:"required"`
	Language    string            `json:"language,omitempty"`
	Parameters  map[string]string `json:"parameters,omitempty"`
}

// MessageResponse represents a formatted message response
type MessageResponse struct {
	MessageCode         string            `json:"message_code"`
	Category            string            `json:"category"`
	Severity            string            `json:"severity"`
	Message             string            `json:"message"`
	DetailedDescription string            `json:"detailed_description"`
	ResponseAction      string            `json:"response_action"`
	FormattedMessage    string            `json:"formatted_message"`
	Metadata            map[string]string `json:"metadata,omitempty"`
	Language            string            `json:"language"`
	Timestamp           time.Time         `json:"timestamp"`
}

// CatalogStats represents statistics about the message catalog
type CatalogStats struct {
	TotalMessages      int            `json:"total_messages"`
	LanguagesCount     int            `json:"languages_count"`
	CategoriesCount    int            `json:"categories_count"`
	SeveritiesCount    int            `json:"severities_count"`
	Languages          []string       `json:"languages"`
	Categories         []string       `json:"categories"`
	Severities         []string       `json:"severities"`
	LastReloaded       time.Time      `json:"last_reloaded"`
	MessagesByLanguage map[string]int `json:"messages_by_language"`
}
