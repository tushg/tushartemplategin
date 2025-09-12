package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"tushartemplategin/pkg/interfaces"
)

// Logger defines the interface for logging operations
type Logger interface {
	Debug(ctx context.Context, msg string, fields interfaces.Fields)            // Log debug message
	Info(ctx context.Context, msg string, fields interfaces.Fields)             // Log info message
	Warn(ctx context.Context, msg string, fields interfaces.Fields)             // Log warning message
	Error(ctx context.Context, msg string, fields interfaces.Fields)            // Log error message
	Fatal(ctx context.Context, msg string, err error, fields interfaces.Fields) // Log fatal message and exit
}

// logger implements the Logger interface using zap
type logger struct {
	zapLogger *zap.Logger // Underlying zap logger instance
}

// customTimeEncoder creates timestamps in YYYY-MM-DDTHH:MM:SS.ssssssZ format
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	// Format: YYYY-MM-DDTHH:MM:SS.ssssssZ
	// Example: 2024-01-15T10:30:45.123456Z
	enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.000000Z"))
}

// generateTimestampBasedFileName creates a timestamp-based log file name
// Format: {prefix}_{YYYY-MM-DD_HH-MM-SS}.log
func generateTimestampBasedFileName(logDirectory string) string {
	timestamp := time.Now().Format(TimestampFormat)
	fileName := fmt.Sprintf("%s_%s%s", LogFilePrefix, timestamp, LogFileExtension)
	return filepath.Join(logDirectory, fileName)
}

// ensureLogDirectory creates the log directory if it doesn't exist
func ensureLogDirectory(logDirectory string) error {
	if logDirectory == "" {
		return fmt.Errorf("log directory cannot be empty")
	}

	// Create directory with appropriate permissions (755)
	if err := os.MkdirAll(logDirectory, 0755); err != nil {
		return fmt.Errorf("failed to create log directory '%s': %w", logDirectory, err)
	}

	return nil
}

// getLogFilePath determines the log file path based on configuration
// Uses filePath as directory and generates timestamp-based filename
func getLogFilePath(config *Config) (string, error) {
	if config.FilePath == "" {
		return "", fmt.Errorf("filePath must be specified for file output")
	}

	// Ensure directory exists
	if err := ensureLogDirectory(config.FilePath); err != nil {
		return "", err
	}

	// Generate timestamp-based file name
	return generateTimestampBasedFileName(config.FilePath), nil
}

// NewLogger creates a new logger instance with the given configuration
func NewLogger(config *Config) (Logger, error) {
	// Convert string log level to zapcore.Level
	var level zapcore.Level
	switch config.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel // Default to info level if invalid
	}

	// Determine output destination (file or stdout)
	var output io.Writer
	if config.Output == "file" {
		// Get timestamp-based log file path
		logFilePath, err := getLogFilePath(config)
		if err != nil {
			return nil, fmt.Errorf("failed to determine log file path: %w", err)
		}

		// Use lumberjack for log rotation and file management
		output = &lumberjack.Logger{
			Filename:   logFilePath,       // Timestamp-based log file path
			MaxSize:    config.MaxSize,    // Max file size in MB
			MaxBackups: config.MaxBackups, // Max number of backup files
			MaxAge:     config.MaxAge,     // Max age of log files in days
			Compress:   config.Compress,   // Whether to compress old files
		}
	} else {
		output = os.Stdout // Default to stdout
	}

	// Configure encoder settings for structured logging
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"                     // Key for timestamp field
	encoderConfig.EncodeTime = customTimeEncoder            // Custom time format: YYYY-MM-DDTHH:MM:SS.ssssssZ
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // Capital level names

	// Choose encoder based on format (JSON or console)
	var encoder zapcore.Encoder
	if config.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig) // JSON format for machine parsing
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig) // Console format for human reading
	}

	// Create the core logger with encoder, output, and level
	core := zapcore.NewCore(encoder, zapcore.AddSync(output), level)

	// Create zap logger with additional options
	var options []zap.Option

	// Add caller information if enabled
	if config.AddCaller {
		options = append(options,
			zap.AddCaller(),      // Enable caller info
			zap.AddCallerSkip(1), // Skip 1 level to get actual caller
		)
	}

	// Add stack traces if enabled (for error level and above)
	if config.AddStack {
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	zapLogger := zap.New(core, options...)

	return &logger{zapLogger: zapLogger}, nil
}

// Debug logs a debug message with optional fields
func (l *logger) Debug(ctx context.Context, msg string, fields interfaces.Fields) {
	l.zapLogger.Debug(msg, convertFields(fields)...)
}

// Info logs an info message with optional fields
func (l *logger) Info(ctx context.Context, msg string, fields interfaces.Fields) {
	l.zapLogger.Info(msg, convertFields(fields)...)
}

// Warn logs a warning message with optional fields
func (l *logger) Warn(ctx context.Context, msg string, fields interfaces.Fields) {
	l.zapLogger.Warn(msg, convertFields(fields)...)
}

// Error logs an error message with optional fields
func (l *logger) Error(ctx context.Context, msg string, fields interfaces.Fields) {
	l.zapLogger.Error(msg, convertFields(fields)...)
}

// Fatal logs a fatal message and exits the program
func (l *logger) Fatal(ctx context.Context, msg string, err error, fields interfaces.Fields) {
	if err != nil {
		fields["error"] = err.Error() // Add error message to fields
	}
	l.zapLogger.Fatal(msg, convertFields(fields)...)
}

// convertFields converts our Fields type to zap.Field slice
func convertFields(fields interfaces.Fields) []zap.Field {
	var zapFields []zap.Field
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}
