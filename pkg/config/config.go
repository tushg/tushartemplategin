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

// DatabaseConfig contains database configuration with support for multiple database types
type DatabaseConfig struct {
	Type     string          `mapstructure:"type"`     // Database type: postgres, sqlite, mysql
	Postgres *PostgresConfig `mapstructure:"postgres"` // PostgreSQL-specific configuration
	SQLite   *SQLiteConfig   `mapstructure:"sqlite"`   // SQLite-specific configuration
	MySQL    *MySQLConfig    `mapstructure:"mysql"`    // MySQL-specific configuration
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

// SQLiteConfig contains SQLite-specific configuration
type SQLiteConfig struct {
	FilePath            string        `mapstructure:"filePath"`            // Database file path
	Timeout             time.Duration `mapstructure:"timeout"`             // Connection timeout
	MaxOpenConns        int           `mapstructure:"maxOpenConns"`        // Max open connections (SQLite limitation: 1 for writes)
	MaxIdleConns        int           `mapstructure:"maxIdleConns"`        // Max idle connections
	ConnMaxLifetime     time.Duration `mapstructure:"connMaxLifetime"`     // Connection max lifetime
	ConnMaxIdleTime     time.Duration `mapstructure:"connMaxIdleTime"`     // Connection max idle time
	JournalMode         string        `mapstructure:"journalMode"`         // Journal mode (WAL, DELETE, TRUNCATE, PERSIST, MEMORY, OFF)
	SyncMode            string        `mapstructure:"syncMode"`            // Sync mode (OFF, NORMAL, FULL, EXTRA)
	CacheSize           int           `mapstructure:"cacheSize"`           // Cache size in pages
	ForeignKeys         bool          `mapstructure:"foreignKeys"`         // Enable foreign key constraints
	AutoVacuum          string        `mapstructure:"autoVacuum"`          // Auto vacuum mode (NONE, INCREMENTAL, FULL)
	HealthCheckInterval time.Duration `mapstructure:"healthCheckInterval"` // Health check interval
}

// MySQLConfig contains MySQL-specific configuration
type MySQLConfig struct {
	Host                string        `mapstructure:"host"`
	Port                int           `mapstructure:"port"`
	Name                string        `mapstructure:"name"`
	Username            string        `mapstructure:"username"`
	Password            string        `mapstructure:"password"`
	Charset             string        `mapstructure:"charset"`   // Character set
	ParseTime           bool          `mapstructure:"parseTime"` // Parse time values
	Loc                 string        `mapstructure:"loc"`       // Location for time parsing
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
	setDatabaseDefaults()

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

// setDatabaseDefaults sets production-ready defaults for all database types
func setDatabaseDefaults() {
	// PostgreSQL defaults
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

	// SQLite defaults
	viper.SetDefault("database.sqlite.filePath", "./data/app.db")
	viper.SetDefault("database.sqlite.timeout", "30s")
	viper.SetDefault("database.sqlite.maxOpenConns", 1) // SQLite limitation: only 1 writer
	viper.SetDefault("database.sqlite.maxIdleConns", 1)
	viper.SetDefault("database.sqlite.connMaxLifetime", "5m")
	viper.SetDefault("database.sqlite.connMaxIdleTime", "1m")
	viper.SetDefault("database.sqlite.journalMode", "WAL") // Write-Ahead Logging for better concurrency
	viper.SetDefault("database.sqlite.syncMode", "NORMAL")
	viper.SetDefault("database.sqlite.cacheSize", 1000)
	viper.SetDefault("database.sqlite.foreignKeys", true)
	viper.SetDefault("database.sqlite.autoVacuum", "INCREMENTAL")
	viper.SetDefault("database.sqlite.healthCheckInterval", "30s")

	// MySQL defaults
	viper.SetDefault("database.mysql.host", "localhost")
	viper.SetDefault("database.mysql.port", 3306)
	viper.SetDefault("database.mysql.charset", "utf8mb4")
	viper.SetDefault("database.mysql.parseTime", true)
	viper.SetDefault("database.mysql.loc", "Local")
	viper.SetDefault("database.mysql.maxOpenConns", 25)
	viper.SetDefault("database.mysql.maxIdleConns", 5)
	viper.SetDefault("database.mysql.connMaxLifetime", "5m")
	viper.SetDefault("database.mysql.connMaxIdleTime", "1m")
	viper.SetDefault("database.mysql.timeout", "30s")
	viper.SetDefault("database.mysql.maxRetries", 3)
	viper.SetDefault("database.mysql.retryDelay", "1s")
	viper.SetDefault("database.mysql.healthCheckInterval", "30s")
}
