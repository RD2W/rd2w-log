package logger

// Logger defines the interface for structured logging compatible with slog
type Logger interface {
	// Debug logs a debug message with key-value pairs
	Debug(msg string, args ...any)

	// Info logs an info message with key-value pairs
	Info(msg string, args ...any)

	// Warn logs a warning message with key-value pairs
	Warn(msg string, args ...any)

	// Error logs an error message with key-value pairs
	Error(msg string, args ...any)

	// Fatal logs a fatal message and exits the application
	Fatal(msg string, args ...any)

	// With returns a new logger with additional context
	With(args ...any) Logger

	// WithGroup returns a new logger that starts a group
	WithGroup(name string) Logger
}
