package logger

import (
	"context"
	"io"
	"os"

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
	if config.Output == "file" && config.FilePath != "" {
		// Use lumberjack for log rotation and file management
		output = &lumberjack.Logger{
			Filename:   config.FilePath,   // Log file path
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
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // ISO8601 time format
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
		options = append(options, zap.AddCaller())
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
