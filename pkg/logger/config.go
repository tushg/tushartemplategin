package logger

// Config contains logger configuration settings
type Config struct {
	Level      string // Log level (debug, info, warn, error, fatal)
	Format     string // Log format (json, console)
	Output     string // Output destination (stdout, file)
	FilePath   string // Log file path (if output is file)
	MaxSize    int    // Maximum log file size in MB
	MaxBackups int    // Maximum number of backup files
	MaxAge     int    // Maximum age of log files in days
	Compress   bool   // Whether to compress old log files
	AddCaller  bool   // Whether to add caller information
	AddStack   bool   // Whether to add stack traces
}

// Constants for log file naming
const (
	// LogFilePrefix is the constant prefix used for all log files
	LogFilePrefix = "app"
	// LogFileExtension is the file extension for log files
	LogFileExtension = ".log"
	// TimestampFormat is the format used for timestamp in log file names
	// Format: YYYY-MM-DD_HH-MM-SS
	TimestampFormat = "2006-01-02_15-04-05"
)

// Level represents the logging level
type Level int

// Log level constants
const (
	DebugLevel Level = iota // 0: Debug level for detailed debugging information
	InfoLevel               // 1: Info level for general information
	WarnLevel               // 2: Warn level for warning messages
	ErrorLevel              // 3: Error level for error messages
	FatalLevel              // 4: Fatal level for fatal errors (will exit program)
)
