package database

import (
	"testing"

	"tushartemplategin/pkg/config"
	"tushartemplategin/pkg/logger"
)

func TestDatabaseType_String(t *testing.T) {
	tests := []struct {
		name     string
		dbType   DatabaseType
		expected string
	}{
		{"PostgreSQL", DatabaseTypePostgreSQL, "postgres"},
		{"SQLite", DatabaseTypeSQLite, "sqlite"},
		{"MySQL", DatabaseTypeMySQL, "mysql"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dbType.String(); got != tt.expected {
				t.Errorf("DatabaseType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDatabaseType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		dbType   DatabaseType
		expected bool
	}{
		{"Valid PostgreSQL", DatabaseTypePostgreSQL, true},
		{"Valid SQLite", DatabaseTypeSQLite, true},
		{"Valid MySQL", DatabaseTypeMySQL, true},
		{"Invalid empty", "", false},
		{"Invalid random", "invalid_db", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dbType.IsValid(); got != tt.expected {
				t.Errorf("DatabaseType.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewDatabaseFactory(t *testing.T) {
	logger, err := logger.NewLogger(&logger.Config{Level: "info"})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	factory := NewDatabaseFactory(logger)

	if factory == nil {
		t.Error("NewDatabaseFactory() returned nil")
	}

	if factory.logger == nil {
		t.Error("DatabaseFactory logger is nil")
	}
}

func TestDatabaseFactory_CreateDatabase(t *testing.T) {
	logger, err := logger.NewLogger(&logger.Config{Level: "info"})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	factory := NewDatabaseFactory(logger)

	tests := []struct {
		name        string
		config      *config.DatabaseConfig
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Nil config",
			config:      nil,
			expectError: true,
			errorMsg:    "database configuration is required",
		},
		{
			name: "Invalid database type",
			config: &config.DatabaseConfig{
				Type: "invalid_db",
			},
			expectError: true,
			errorMsg:    "unsupported database type",
		},
		{
			name: "PostgreSQL with valid config",
			config: &config.DatabaseConfig{
				Type: "postgres",
				Postgres: &config.PostgresConfig{
					Host:     "localhost",
					Port:     5432,
					Name:     "test_db",
					Username: "test_user",
					Password: "test_pass",
				},
			},
			expectError: false,
		},
		{
			name: "PostgreSQL with nil config",
			config: &config.DatabaseConfig{
				Type: "postgres",
			},
			expectError: true,
			errorMsg:    "PostgreSQL configuration is required",
		},
		{
			name: "SQLite not implemented",
			config: &config.DatabaseConfig{
				Type: "sqlite",
				SQLite: &config.SQLiteConfig{
					FilePath: "./test.db",
				},
			},
			expectError: true,
			errorMsg:    "SQLite database implementation not yet available",
		},
		{
			name: "MySQL not implemented",
			config: &config.DatabaseConfig{
				Type: "mysql",
				MySQL: &config.MySQLConfig{
					Host:     "localhost",
					Port:     3306,
					Name:     "test_db",
					Username: "test_user",
					Password: "test_pass",
				},
			},
			expectError: true,
			errorMsg:    "MySQL database implementation not yet available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := factory.CreateDatabase(tt.config)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				if db == nil {
					t.Error("Expected database instance but got nil")
				}
			}
		})
	}
}

func TestDatabaseFactory_validatePostgreSQLConfig(t *testing.T) {
	logger, err := logger.NewLogger(&logger.Config{Level: "info"})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	factory := NewDatabaseFactory(logger)

	tests := []struct {
		name        string
		config      *config.PostgresConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid config",
			config: &config.PostgresConfig{
				Host:     "localhost",
				Port:     5432,
				Name:     "test_db",
				Username: "test_user",
				Password: "test_pass",
			},
			expectError: false,
		},
		{
			name: "Missing host",
			config: &config.PostgresConfig{
				Port:     5432,
				Name:     "test_db",
				Username: "test_user",
				Password: "test_pass",
			},
			expectError: true,
			errorMsg:    "host is required",
		},
		{
			name: "Invalid port",
			config: &config.PostgresConfig{
				Host:     "localhost",
				Port:     0,
				Name:     "test_db",
				Username: "test_user",
				Password: "test_pass",
			},
			expectError: true,
			errorMsg:    "port must be between 1 and 65535",
		},
		{
			name: "Missing database name",
			config: &config.PostgresConfig{
				Host:     "localhost",
				Port:     5432,
				Username: "test_user",
				Password: "test_pass",
			},
			expectError: true,
			errorMsg:    "database name is required",
		},
		{
			name: "Missing username",
			config: &config.PostgresConfig{
				Host:     "localhost",
				Port:     5432,
				Name:     "test_db",
				Password: "test_pass",
			},
			expectError: true,
			errorMsg:    "username is required",
		},
		{
			name: "Empty password (warning only)",
			config: &config.PostgresConfig{
				Host:     "localhost",
				Port:     5432,
				Name:     "test_db",
				Username: "test_user",
				Password: "",
			},
			expectError: false, // Password is optional but generates warning
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := factory.validatePostgreSQLConfig(tt.config)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestDatabaseFactory_validateSQLiteConfig(t *testing.T) {
	logger, err := logger.NewLogger(&logger.Config{Level: "info"})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	factory := NewDatabaseFactory(logger)

	tests := []struct {
		name        string
		config      *config.SQLiteConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid config",
			config: &config.SQLiteConfig{
				FilePath: "./test.db",
			},
			expectError: false,
		},
		{
			name:        "Missing file path",
			config:      &config.SQLiteConfig{},
			expectError: true,
			errorMsg:    "file path is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := factory.validateSQLiteConfig(tt.config)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestDatabaseFactory_getDatabaseHost(t *testing.T) {
	logger, err := logger.NewLogger(&logger.Config{Level: "info"})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	factory := NewDatabaseFactory(logger)

	tests := []struct {
		name     string
		config   *config.DatabaseConfig
		expected string
	}{
		{
			name: "PostgreSQL host",
			config: &config.DatabaseConfig{
				Type: "postgres",
				Postgres: &config.PostgresConfig{
					Host: "localhost",
					Port: 5432,
				},
			},
			expected: "localhost:5432",
		},
		{
			name: "SQLite file path",
			config: &config.DatabaseConfig{
				Type: "sqlite",
				SQLite: &config.SQLiteConfig{
					FilePath: "./test.db",
				},
			},
			expected: "./test.db",
		},
		{
			name: "Unknown type",
			config: &config.DatabaseConfig{
				Type: "unknown",
			},
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host := factory.getDatabaseHost(tt.config)
			if host != tt.expected {
				t.Errorf("getDatabaseHost() = %v, want %v", host, tt.expected)
			}
		})
	}
}

func TestDatabaseFactory_getSupportedTypes(t *testing.T) {
	logger, err := logger.NewLogger(&logger.Config{Level: "info"})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	factory := NewDatabaseFactory(logger)

	expected := "postgres, sqlite, mysql"
	got := factory.getSupportedTypes()

	if got != expected {
		t.Errorf("getSupportedTypes() = %v, want %v", got, expected)
	}
}

func TestDatabaseFactory_GetDatabaseType(t *testing.T) {
	logger, err := logger.NewLogger(&logger.Config{Level: "info"})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	factory := NewDatabaseFactory(logger)

	config := &config.DatabaseConfig{
		Type: "postgres",
	}

	expected := DatabaseTypePostgreSQL
	got := factory.GetDatabaseType(config)

	if got != expected {
		t.Errorf("GetDatabaseType() = %v, want %v", got, expected)
	}
}

func TestDatabaseType_IsDatabaseTypeSupported(t *testing.T) {
	tests := []struct {
		name     string
		dbType   DatabaseType
		input    string
		expected bool
	}{
		{"Supported PostgreSQL", DatabaseTypePostgreSQL, "postgres", true},
		{"Supported SQLite", DatabaseTypeSQLite, "sqlite", true},
		{"Supported MySQL", DatabaseTypeMySQL, "mysql", true},
		{"Unsupported type", DatabaseTypePostgreSQL, "invalid", false},
		{"Empty string", DatabaseTypePostgreSQL, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dbType.IsDatabaseTypeSupported(tt.input); got != tt.expected {
				t.Errorf("IsDatabaseTypeSupported() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
