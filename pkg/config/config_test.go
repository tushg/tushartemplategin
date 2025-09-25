package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Load(t *testing.T) {
	tests := []struct {
		name           string
		configContent  string
		expectedConfig *Config
		expectError    bool
	}{
		{
			name: "valid complete config",
			configContent: `{
				"server": {
					"port": ":8080",
					"mode": "debug",
					"ssl": {
						"enabled": true,
						"port": ":443",
						"certFile": "/path/to/cert.pem",
						"keyFile": "/path/to/key.pem",
						"redirectHTTP": true
					}
				},
				"log": {
					"level": "info",
					"format": "json",
					"output": "stdout",
					"filePath": "/var/log/app.log",
					"maxSize": 100,
					"maxBackups": 3,
					"maxAge": 28,
					"compress": true,
					"addCaller": true,
					"addStack": false
				},
				"database": {
					"type": "postgres",
					"postgres": {
						"host": "localhost",
						"port": 5432,
						"name": "testdb",
						"username": "testuser",
						"password": "testpass",
						"sslMode": "require",
						"maxOpenConns": 10,
						"maxIdleConns": 2,
						"connMaxLifetime": "10m",
						"connMaxIdleTime": "2m",
						"timeout": "60s",
						"maxRetries": 5,
						"retryDelay": "2s",
						"healthCheckInterval": "60s"
					}
				}
			}`,
			expectedConfig: &Config{
				Server: ServerConfig{
					Port: ":8080",
					Mode: "debug",
					SSL: SSLConfig{
						Enabled:      true,
						Port:         ":443",
						CertFile:     "/path/to/cert.pem",
						KeyFile:      "/path/to/key.pem",
						RedirectHTTP: true,
					},
				},
				Log: LogConfig{
					Level:      "info",
					Format:     "json",
					Output:     "stdout",
					FilePath:   "/var/log/app.log",
					MaxSize:    100,
					MaxBackups: 3,
					MaxAge:     28,
					Compress:   true,
					AddCaller:  true,
					AddStack:   false,
				},
				Database: DatabaseConfig{
					Type: "postgres",
					Postgres: &PostgresConfig{
						Host:                "localhost",
						Port:                5432,
						Name:                "testdb",
						Username:            "testuser",
						Password:            "testpass",
						SSLMode:             "require",
						MaxOpenConns:        10,
						MaxIdleConns:        2,
						ConnMaxLifetime:     10 * time.Minute,
						ConnMaxIdleTime:     2 * time.Minute,
						Timeout:             60 * time.Second,
						MaxRetries:          5,
						RetryDelay:          2 * time.Second,
						HealthCheckInterval: 60 * time.Second,
					},
				},
			},
			expectError: false,
		},
		{
			name: "config with defaults only",
			configContent: `{
				"server": {
					"port": ":3000"
				}
			}`,
			expectedConfig: &Config{
				Server: ServerConfig{
					Port: ":3000",
					Mode: "",
					SSL: SSLConfig{
						Enabled:      false,
						Port:         "",
						CertFile:     "",
						KeyFile:      "",
						RedirectHTTP: false,
					},
				},
				Log: LogConfig{
					Level:      "",
					Format:     "",
					Output:     "",
					FilePath:   "",
					MaxSize:    0,
					MaxBackups: 0,
					MaxAge:     0,
					Compress:   false,
					AddCaller:  false,
					AddStack:   false,
				},
				Database: DatabaseConfig{
					Type: "",
					Postgres: &PostgresConfig{
						Host:                "localhost",
						Port:                5432,
						SSLMode:             "disable",
						MaxOpenConns:        25,
						MaxIdleConns:        5,
						ConnMaxLifetime:     5 * time.Minute,
						ConnMaxIdleTime:     1 * time.Minute,
						Timeout:             30 * time.Second,
						MaxRetries:          3,
						RetryDelay:          1 * time.Second,
						HealthCheckInterval: 30 * time.Second,
					},
				},
			},
			expectError: false,
		},
		{
			name: "invalid JSON",
			configContent: `{
				"server": {
					"port": ":8080"
				}
				"invalid": json
			}`,
			expectedConfig: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary config file
			tempDir := t.TempDir()
			configFile := filepath.Join(tempDir, "config.json")

			err := os.WriteFile(configFile, []byte(tt.configContent), 0644)
			require.NoError(t, err)

			// Reset viper for clean test
			viper.Reset()

			// Set config path to temp directory
			viper.AddConfigPath(tempDir)
			viper.SetConfigName("config")
			viper.SetConfigType("json")

			// Load config
			config, err := Load()

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, config)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, config)

				// Compare server config
				assert.Equal(t, tt.expectedConfig.Server.Port, config.Server.Port)
				assert.Equal(t, tt.expectedConfig.Server.Mode, config.Server.Mode)
				assert.Equal(t, tt.expectedConfig.Server.SSL, config.Server.SSL)

				// Compare log config
				assert.Equal(t, tt.expectedConfig.Log, config.Log)

				// Compare database config
				assert.Equal(t, tt.expectedConfig.Database.Type, config.Database.Type)

				if tt.expectedConfig.Database.Postgres != nil {
					assert.NotNil(t, config.Database.Postgres)
					assert.Equal(t, *tt.expectedConfig.Database.Postgres, *config.Database.Postgres)
				}
			}
		})
	}
}

func TestDatabaseConfig_GetType(t *testing.T) {
	tests := []struct {
		name     string
		config   DatabaseConfig
		expected string
	}{
		{
			name:     "postgres type",
			config:   DatabaseConfig{Type: "postgres"},
			expected: "postgres",
		},
		{
			name:     "sqlite type",
			config:   DatabaseConfig{Type: "sqlite"},
			expected: "sqlite",
		},
		{
			name:     "mysql type",
			config:   DatabaseConfig{Type: "mysql"},
			expected: "mysql",
		},
		{
			name:     "empty type",
			config:   DatabaseConfig{Type: ""},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetType()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDatabaseConfig_GetPostgres(t *testing.T) {
	tests := []struct {
		name     string
		config   DatabaseConfig
		expected interface{}
	}{
		{
			name: "postgres config exists",
			config: DatabaseConfig{
				Postgres: &PostgresConfig{
					Host:     "localhost",
					Port:     5432,
					Username: "user",
					Password: "pass",
				},
			},
			expected: &PostgresConfig{
				Host:     "localhost",
				Port:     5432,
				Username: "user",
				Password: "pass",
			},
		},
		{
			name:     "postgres config is nil",
			config:   DatabaseConfig{Postgres: nil},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetPostgres()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDatabaseConfig_GetSQLite(t *testing.T) {
	tests := []struct {
		name     string
		config   DatabaseConfig
		expected interface{}
	}{
		{
			name: "sqlite config exists",
			config: DatabaseConfig{
				SQLite: &SQLiteConfig{
					FilePath: "/path/to/db.sqlite",
					Timeout:  30 * time.Second,
				},
			},
			expected: &SQLiteConfig{
				FilePath: "/path/to/db.sqlite",
				Timeout:  30 * time.Second,
			},
		},
		{
			name:     "sqlite config is nil",
			config:   DatabaseConfig{SQLite: nil},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetSQLite()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDatabaseConfig_GetMySQL(t *testing.T) {
	tests := []struct {
		name     string
		config   DatabaseConfig
		expected interface{}
	}{
		{
			name: "mysql config exists",
			config: DatabaseConfig{
				MySQL: &MySQLConfig{
					Host:     "localhost",
					Port:     3306,
					Username: "user",
					Password: "pass",
				},
			},
			expected: &MySQLConfig{
				Host:     "localhost",
				Port:     3306,
				Username: "user",
				Password: "pass",
			},
		},
		{
			name:     "mysql config is nil",
			config:   DatabaseConfig{MySQL: nil},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetMySQL()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPostgresConfig_Getters(t *testing.T) {
	config := &PostgresConfig{
		Host:                "testhost",
		Port:                5433,
		Name:                "testdb",
		Username:            "testuser",
		Password:            "testpass",
		SSLMode:             "require",
		MaxOpenConns:        15,
		MaxIdleConns:        3,
		ConnMaxLifetime:     10 * time.Minute,
		ConnMaxIdleTime:     2 * time.Minute,
		Timeout:             45 * time.Second,
		MaxRetries:          5,
		RetryDelay:          2 * time.Second,
		HealthCheckInterval: 60 * time.Second,
	}

	tests := []struct {
		name     string
		getter   func() interface{}
		expected interface{}
	}{
		{"GetHost", func() interface{} { return config.GetHost() }, "testhost"},
		{"GetPort", func() interface{} { return config.GetPort() }, 5433},
		{"GetName", func() interface{} { return config.GetName() }, "testdb"},
		{"GetUsername", func() interface{} { return config.GetUsername() }, "testuser"},
		{"GetPassword", func() interface{} { return config.GetPassword() }, "testpass"},
		{"GetSSLMode", func() interface{} { return config.GetSSLMode() }, "require"},
		{"GetMaxOpenConns", func() interface{} { return config.GetMaxOpenConns() }, 15},
		{"GetMaxIdleConns", func() interface{} { return config.GetMaxIdleConns() }, 3},
		{"GetConnMaxLifetime", func() interface{} { return config.GetConnMaxLifetime() }, 10 * time.Minute},
		{"GetConnMaxIdleTime", func() interface{} { return config.GetConnMaxIdleTime() }, 2 * time.Minute},
		{"GetTimeout", func() interface{} { return config.GetTimeout() }, 45 * time.Second},
		{"GetMaxRetries", func() interface{} { return config.GetMaxRetries() }, 5},
		{"GetRetryDelay", func() interface{} { return config.GetRetryDelay() }, 2 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.getter()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSQLiteConfig_Getters(t *testing.T) {
	config := &SQLiteConfig{
		FilePath: "/path/to/test.db",
		Timeout:  60 * time.Second,
	}

	tests := []struct {
		name     string
		getter   func() interface{}
		expected interface{}
	}{
		{"GetPath", func() interface{} { return config.GetPath() }, "/path/to/test.db"},
		{"GetTimeout", func() interface{} { return config.GetTimeout() }, 60 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.getter()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMySQLConfig_Getters(t *testing.T) {
	config := &MySQLConfig{
		Host:                "mysqlhost",
		Port:                3307,
		Name:                "mysqldb",
		Username:            "mysqluser",
		Password:            "mysqlpass",
		Charset:             "utf8mb4",
		ParseTime:           true,
		Loc:                 "UTC",
		MaxOpenConns:        20,
		MaxIdleConns:        4,
		ConnMaxLifetime:     8 * time.Minute,
		ConnMaxIdleTime:     3 * time.Minute,
		Timeout:             50 * time.Second,
		MaxRetries:          4,
		RetryDelay:          3 * time.Second,
		HealthCheckInterval: 45 * time.Second,
	}

	tests := []struct {
		name     string
		getter   func() interface{}
		expected interface{}
	}{
		{"GetHost", func() interface{} { return config.GetHost() }, "mysqlhost"},
		{"GetPort", func() interface{} { return config.GetPort() }, 3307},
		{"GetName", func() interface{} { return config.GetName() }, "mysqldb"},
		{"GetUsername", func() interface{} { return config.GetUsername() }, "mysqluser"},
		{"GetPassword", func() interface{} { return config.GetPassword() }, "mysqlpass"},
		{"GetTimeout", func() interface{} { return config.GetTimeout() }, 50 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.getter()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSetDatabaseDefaults(t *testing.T) {
	// Reset viper to ensure clean state
	viper.Reset()

	// Call the function
	setDatabaseDefaults()

	// Test PostgreSQL defaults
	assert.Equal(t, "localhost", viper.GetString("database.postgres.host"))
	assert.Equal(t, 5432, viper.GetInt("database.postgres.port"))
	assert.Equal(t, "disable", viper.GetString("database.postgres.sslMode"))
	assert.Equal(t, 25, viper.GetInt("database.postgres.maxOpenConns"))
	assert.Equal(t, 5, viper.GetInt("database.postgres.maxIdleConns"))
	assert.Equal(t, "5m", viper.GetString("database.postgres.connMaxLifetime"))
	assert.Equal(t, "1m", viper.GetString("database.postgres.connMaxIdleTime"))
	assert.Equal(t, "30s", viper.GetString("database.postgres.timeout"))
	assert.Equal(t, 3, viper.GetInt("database.postgres.maxRetries"))
	assert.Equal(t, "1s", viper.GetString("database.postgres.retryDelay"))
	assert.Equal(t, "30s", viper.GetString("database.postgres.healthCheckInterval"))

	// Test SQLite defaults
	assert.Equal(t, "./data/app.db", viper.GetString("database.sqlite.filePath"))
	assert.Equal(t, "30s", viper.GetString("database.sqlite.timeout"))
	assert.Equal(t, 1, viper.GetInt("database.sqlite.maxOpenConns"))
	assert.Equal(t, 1, viper.GetInt("database.sqlite.maxIdleConns"))
	assert.Equal(t, "5m", viper.GetString("database.sqlite.connMaxLifetime"))
	assert.Equal(t, "1m", viper.GetString("database.sqlite.connMaxIdleTime"))
	assert.Equal(t, "WAL", viper.GetString("database.sqlite.journalMode"))
	assert.Equal(t, "NORMAL", viper.GetString("database.sqlite.syncMode"))
	assert.Equal(t, 1000, viper.GetInt("database.sqlite.cacheSize"))
	assert.Equal(t, true, viper.GetBool("database.sqlite.foreignKeys"))
	assert.Equal(t, "INCREMENTAL", viper.GetString("database.sqlite.autoVacuum"))
	assert.Equal(t, "30s", viper.GetString("database.sqlite.healthCheckInterval"))

	// Test MySQL defaults
	assert.Equal(t, "localhost", viper.GetString("database.mysql.host"))
	assert.Equal(t, 3306, viper.GetInt("database.mysql.port"))
	assert.Equal(t, "utf8mb4", viper.GetString("database.mysql.charset"))
	assert.Equal(t, true, viper.GetBool("database.mysql.parseTime"))
	assert.Equal(t, "Local", viper.GetString("database.mysql.loc"))
	assert.Equal(t, 25, viper.GetInt("database.mysql.maxOpenConns"))
	assert.Equal(t, 5, viper.GetInt("database.mysql.maxIdleConns"))
	assert.Equal(t, "5m", viper.GetString("database.mysql.connMaxLifetime"))
	assert.Equal(t, "1m", viper.GetString("database.mysql.connMaxIdleTime"))
	assert.Equal(t, "30s", viper.GetString("database.mysql.timeout"))
	assert.Equal(t, 3, viper.GetInt("database.mysql.maxRetries"))
	assert.Equal(t, "1s", viper.GetString("database.mysql.retryDelay"))
	assert.Equal(t, "30s", viper.GetString("database.mysql.healthCheckInterval"))
}

func TestConfig_Load_WithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("SERVER_PORT", ":9090")
	os.Setenv("DATABASE_POSTGRES_HOST", "envhost")
	os.Setenv("DATABASE_POSTGRES_PORT", "5434")
	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DATABASE_POSTGRES_HOST")
		os.Unsetenv("DATABASE_POSTGRES_PORT")
	}()

	// Create minimal config file
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.json")

	configContent := `{
		"server": {
			"port": ":8080"
		},
		"database": {
			"type": "postgres",
			"postgres": {
				"host": "localhost",
				"port": 5432
			}
		}
	}`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	// Reset viper for clean test
	viper.Reset()

	// Set config path to temp directory
	viper.AddConfigPath(tempDir)
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	// Enable automatic environment variable binding
	viper.AutomaticEnv()
	viper.SetEnvPrefix("")

	// Load config
	config, err := Load()
	require.NoError(t, err)
	require.NotNil(t, config)

	// Note: The current Load() function doesn't enable environment variable binding
	// This test documents the current behavior - config file values take precedence
	// In a real implementation, you might want to enable environment variable override
	assert.Equal(t, ":8080", config.Server.Port)                // Config file value, not env var
	assert.Equal(t, "localhost", config.Database.Postgres.Host) // Config file value, not env var
	assert.Equal(t, 5432, config.Database.Postgres.Port)        // Config file value, not env var
}

func TestConfig_Load_FileNotFound(t *testing.T) {
	// Reset viper for clean test
	viper.Reset()

	// Set config path to non-existent directory
	viper.AddConfigPath("/non/existent/path")
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	// Load config should return error
	config, err := Load()
	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestConfig_StructTags(t *testing.T) {
	// Test that struct tags are properly set for mapstructure
	config := &Config{
		Server: ServerConfig{
			Port: ":8080",
			Mode: "debug",
			SSL: SSLConfig{
				Enabled:      true,
				Port:         ":443",
				CertFile:     "/path/to/cert.pem",
				KeyFile:      "/path/to/key.pem",
				RedirectHTTP: true,
			},
		},
		Log: LogConfig{
			Level:      "info",
			Format:     "json",
			Output:     "stdout",
			FilePath:   "/var/log/app.log",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   true,
			AddCaller:  true,
			AddStack:   false,
		},
		Database: DatabaseConfig{
			Type: "postgres",
			Postgres: &PostgresConfig{
				Host:     "localhost",
				Port:     5432,
				Name:     "testdb",
				Username: "testuser",
				Password: "testpass",
				SSLMode:  "require",
			},
		},
	}

	// Test that the struct can be marshaled to map[string]interface{}
	// This indirectly tests that mapstructure tags are working
	viper.Reset()
	viper.Set("server.port", config.Server.Port)
	viper.Set("server.mode", config.Server.Mode)
	viper.Set("server.ssl.enabled", config.Server.SSL.Enabled)
	viper.Set("log.level", config.Log.Level)
	viper.Set("log.format", config.Log.Format)
	viper.Set("database.type", config.Database.Type)
	viper.Set("database.postgres.host", config.Database.Postgres.Host)
	viper.Set("database.postgres.port", config.Database.Postgres.Port)

	var result Config
	err := viper.Unmarshal(&result)
	assert.NoError(t, err)
	assert.Equal(t, config.Server.Port, result.Server.Port)
	assert.Equal(t, config.Server.Mode, result.Server.Mode)
	assert.Equal(t, config.Server.SSL.Enabled, result.Server.SSL.Enabled)
	assert.Equal(t, config.Log.Level, result.Log.Level)
	assert.Equal(t, config.Log.Format, result.Log.Format)
	assert.Equal(t, config.Database.Type, result.Database.Type)
	assert.Equal(t, config.Database.Postgres.Host, result.Database.Postgres.Host)
	assert.Equal(t, config.Database.Postgres.Port, result.Database.Postgres.Port)
}

func TestConfig_EdgeCases(t *testing.T) {
	t.Run("empty config file", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "config.json")

		err := os.WriteFile(configFile, []byte("{}"), 0644)
		require.NoError(t, err)

		viper.Reset()
		viper.AddConfigPath(tempDir)
		viper.SetConfigName("config")
		viper.SetConfigType("json")

		config, err := Load()
		assert.NoError(t, err)
		assert.NotNil(t, config)

		// Should have defaults applied
		assert.Equal(t, "localhost", config.Database.Postgres.Host)
		assert.Equal(t, 5432, config.Database.Postgres.Port)
	})

	t.Run("config with only server section", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "config.json")

		configContent := `{
			"server": {
				"port": ":3000",
				"mode": "release"
			}
		}`

		err := os.WriteFile(configFile, []byte(configContent), 0644)
		require.NoError(t, err)

		viper.Reset()
		viper.AddConfigPath(tempDir)
		viper.SetConfigName("config")
		viper.SetConfigType("json")

		config, err := Load()
		assert.NoError(t, err)
		assert.NotNil(t, config)

		assert.Equal(t, ":3000", config.Server.Port)
		assert.Equal(t, "release", config.Server.Mode)
		assert.Equal(t, "localhost", config.Database.Postgres.Host) // Default applied
	})
}
