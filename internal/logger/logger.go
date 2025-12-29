package logger

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Context keys for logger fields
type contextKey string

const (
	// TraceIDKey is the context key for trace ID
	TraceIDKey contextKey = "trace_id"
	// UserIDKey is the context key for user ID
	UserIDKey contextKey = "user_id"
	// RequestIDKey is the context key for request ID
	RequestIDKey contextKey = "request_id"
)

var globalLogger zerolog.Logger

// Config holds logger configuration
type Config struct {
	Level      string
	LogDir     string
	MaxSize    int  // megabytes
	MaxBackups int  // number of backups
	MaxAge     int  // days
	Compress   bool // compress rotated files
}

// DefaultConfig returns default logger configuration
func DefaultConfig() Config {
	return Config{
		Level:      "info",
		LogDir:     "logs",
		MaxSize:    100,  // 100 MB
		MaxBackups: 5,    // keep 5 backups
		MaxAge:     30,   // 30 days
		Compress:   true, // compress old logs
	}
}

// Setup initializes the global logger with file and console output
func Setup(config Config) error {
	// Parse log level
	lvl, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return err
	}

	// Setup file logger with rotation
	fileLogger := &lumberjack.Logger{
		Filename:   filepath.Join(config.LogDir, "app.log"),
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	// Setup console output with pretty formatting for development
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	// Multi-writer: write to both file and console
	multiWriter := io.MultiWriter(consoleWriter, fileLogger)

	// Configure zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Create logger
	globalLogger = zerolog.New(multiWriter).
		Level(lvl).
		With().
		Timestamp().
		Caller().
		Logger()

	// Set global logger
	log.Logger = globalLogger

	return nil
}

// SetupGlobal initializes the global logger with default config (backward compatible)
func SetupGlobal(level string) {
	config := DefaultConfig()
	config.Level = level
	if err := Setup(config); err != nil {
		log.Fatal().Err(err).Msg("Failed to setup logger")
	}
}

// New creates a new logger instance with the given level
func New(level string) zerolog.Logger {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).Level(lvl).With().Timestamp().Caller().Logger()
	return logger
}

// FromContext extracts the logger from context with all contextual fields
func FromContext(ctx context.Context) *zerolog.Logger {
	logger := globalLogger

	// Add trace ID if present
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		logger = logger.With().Str("trace_id", traceID).Logger()
	}

	// Add request ID if present
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		logger = logger.With().Str("request_id", requestID).Logger()
	}

	// Add user ID if present
	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		logger = logger.With().Str("user_id", userID).Logger()
	}

	return &logger
}

// WithTraceID adds a trace ID to the context
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// WithRequestID adds a request ID to the context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// WithUserID adds a user ID to the context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetTraceID retrieves the trace ID from context
func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// GetRequestID retrieves the request ID from context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// GetUserID retrieves the user ID from context
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}
