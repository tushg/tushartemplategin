package config

import (
	"github.com/spf13/viper"
	"time"
	"tushartemplategin/pkg/interfaces"
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

	// SSL/TLS Configuration
	SSL SSLConfig `mapstructure:"ssl"` // SSL/TLS configuration
}

// SSLConfig contains SSL/TLS configuration settings
type SSLConfig struct {
	Enabled      bool   `mapstructure:"enabled"`      // Enable SSL/TLS
	Port         string `mapstructure:"port"`         // SSL port (e.g., ":443")
	CertFile     string `mapstructure:"certFile"`     // Path to SSL certificate file
	KeyFile      string `mapstructure:"keyFile"`      // Path to SSL private key file
	RedirectHTTP bool   `mapstructure:"redirectHTTP"` // Redirect HTTP to HTTPS
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

// GetType returns the database type
func (dc *DatabaseConfig) GetType() string {
	return dc.Type
}

// GetPostgres returns the PostgreSQL configuration
func (dc *DatabaseConfig) GetPostgres() interfaces.PostgresConfig {
	if dc.Postgres == nil {
		return nil
	}
	return dc.Postgres
}

// GetSQLite returns the SQLite configuration
func (dc *DatabaseConfig) GetSQLite() interfaces.SQLiteConfig {
	if dc.SQLite == nil {
		return nil
	}
	return dc.SQLite
}

// GetMySQL returns the MySQL configuration
func (dc *DatabaseConfig) GetMySQL() interfaces.MySQLConfig {
	if dc.MySQL == nil {
		return nil
	}
	return dc.MySQL
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

// GetHost returns the host
func (pc *PostgresConfig) GetHost() string { return pc.Host }

// GetPort returns the port
func (pc *PostgresConfig) GetPort() int { return pc.Port }

// GetName returns the database name
func (pc *PostgresConfig) GetName() string { return pc.Name }

// GetUsername returns the username
func (pc *PostgresConfig) GetUsername() string { return pc.Username }

// GetPassword returns the password
func (pc *PostgresConfig) GetPassword() string { return pc.Password }

// GetSSLMode returns the SSL mode
func (pc *PostgresConfig) GetSSLMode() string { return pc.SSLMode }

// GetMaxRetries returns the maximum retries
func (pc *PostgresConfig) GetMaxRetries() int { return pc.MaxRetries }

// GetRetryDelay returns the retry delay
func (pc *PostgresConfig) GetRetryDelay() time.Duration { return pc.RetryDelay }

// GetTimeout returns the timeout
func (pc *PostgresConfig) GetTimeout() time.Duration { return pc.Timeout }

// GetMaxOpenConns returns the maximum open connections
func (pc *PostgresConfig) GetMaxOpenConns() int { return pc.MaxOpenConns }

// GetMaxIdleConns returns the maximum idle connections
func (pc *PostgresConfig) GetMaxIdleConns() int { return pc.MaxIdleConns }

// GetConnMaxLifetime returns the connection max lifetime
func (pc *PostgresConfig) GetConnMaxLifetime() time.Duration { return pc.ConnMaxLifetime }

// GetConnMaxIdleTime returns the connection max idle time
func (pc *PostgresConfig) GetConnMaxIdleTime() time.Duration { return pc.ConnMaxIdleTime }

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

// GetPath returns the file path
func (sc *SQLiteConfig) GetPath() string { return sc.FilePath }

// GetTimeout returns the timeout
func (sc *SQLiteConfig) GetTimeout() time.Duration { return sc.Timeout }

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

// GetHost returns the host
func (mc *MySQLConfig) GetHost() string { return mc.Host }

// GetPort returns the port
func (mc *MySQLConfig) GetPort() int { return mc.Port }

// GetName returns the database name
func (mc *MySQLConfig) GetName() string { return mc.Name }

// GetUsername returns the username
func (mc *MySQLConfig) GetUsername() string { return mc.Username }

// GetPassword returns the password
func (mc *MySQLConfig) GetPassword() string { return mc.Password }

// GetTimeout returns the timeout
func (mc *MySQLConfig) GetTimeout() time.Duration { return mc.Timeout }

// Load reads configuration from config files and environment variables
// Returns a Config struct or an error if configuration cannot be loaded
func Load() (*Config, error) {
	// Set configuration file name and type
	viper.SetConfigName("config")
	viper.SetConfigType("json") // Changed from "yaml" to "json"

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
