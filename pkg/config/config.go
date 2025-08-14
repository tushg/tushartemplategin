package config

import (
	"github.com/spf13/viper"
	"time"
)

// Config represents the main application configuration structure
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`   // Server-related configuration
	Log      LogConfig      `mapstructure:"log"`      // Logging configuration
	Database DatabaseConfig `mapstructure:"database"` // Database configuration
}

// ServerConfig contains server-specific settings
type ServerConfig struct {
	Port string `mapstructure:"port"` // Server port (e.g., ":8080")
	Mode string `mapstructure:"mode"` // Server mode (debug/release)
}

// LogConfig contains logging configuration settings
type LogConfig struct {
	Level      string `mapstructure:"level"`      // Log level (debug, info, warn, error, fatal)
	Format     string `mapstructure:"format"`     // Log format (json, console)
	Output     string `mapstructure:"output"`     // Output destination (stdout, file)
	FilePath   string `mapstructure:"filePath"`   // Log file path (if output is file)
	MaxSize    int    `mapstructure:"maxSize"`    // Maximum log file size in MB
	MaxBackups int    `mapstructure:"maxBackups"` // Maximum number of backup files
	MaxAge     int    `mapstructure:"maxAge"`     // Maximum age of log files in days
	Compress   bool   `mapstructure:"compress"`   // Whether to compress old log files
	AddCaller  bool   `mapstructure:"addCaller"`  // Whether to add caller information
	AddStack   bool   `mapstructure:"addStack"`   // Whether to add stack traces
}

// DatabaseConfig contains database configuration
type DatabaseConfig struct {
	Postgres PostgresConfig `mapstructure:"postgres"`
}

// PostgresConfig contains PostgreSQL-specific configuration
type PostgresConfig struct {
	Host                string        `mapstructure:"host"`
	Port                int           `mapstructure:"port"`
	Name                string        `mapstructure:"name"`
	Username            string        `mapstructure:"username"`
	Password            string        `mapstructure:"password"`
	SSLMode             string        `mapstructure:"sslMode"`
	MaxOpenConns        int           `mapstructure:"maxOpenConns"`
	MaxIdleConns        int           `mapstructure:"maxIdleConns"`
	ConnMaxLifetime     time.Duration `mapstructure:"connMaxLifetime"`
	ConnMaxIdleTime     time.Duration `mapstructure:"connMaxIdleTime"`
	Timeout             time.Duration `mapstructure:"timeout"`
	MaxRetries          int           `mapstructure:"maxRetries"`
	RetryDelay          time.Duration `mapstructure:"retryDelay"`
	HealthCheckInterval time.Duration `mapstructure:"healthCheckInterval"`
}

// Load reads configuration from config files and environment variables
// Returns a Config struct or an error if configuration cannot be loaded
func Load() (*Config, error) {
	// Set configuration file name and type
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Add configuration file search paths
	viper.AddConfigPath("./configs") // Look in configs directory
	viper.AddConfigPath(".")         // Look in current directory

	// Set production-ready defaults
	viper.SetDefault("database.postgres.host", "localhost")
	viper.SetDefault("database.postgres.port", 5432)
	viper.SetDefault("database.postgres.sslMode", "disable")
	viper.SetDefault("database.postgres.maxOpenConns", 25)
	viper.SetDefault("database.postgres.maxIdleConns", 5)
	viper.SetDefault("database.postgres.connMaxLifetime", "5m")
	viper.SetDefault("database.postgres.connMaxIdleTime", "1m")
	viper.SetDefault("database.postgres.timeout", "30s")
	viper.SetDefault("database.postgres.maxRetries", 3)
	viper.SetDefault("database.postgres.retryDelay", "1s")
	viper.SetDefault("database.postgres.healthCheckInterval", "30s")

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Unmarshal configuration into Config struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
