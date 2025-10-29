package logger

import "context"

// contextKey является приватным типом для ключей контекста
type contextKey string

const (
	loggerKey contextKey = "logger"
)

// FromContext извлекает логгер из контекста или возвращает дефолтный
func FromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerKey).(Logger); ok {
		return logger
	}
	return &noopLogger{}
}

// WithContext добавляет логгер в контекст
func WithContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// WithContextFields добавляет поля к логгеру в контексте
func WithContextFields(ctx context.Context, args ...any) context.Context {
	logger := FromContext(ctx)
	return WithContext(ctx, logger.With(args...))
}

// noopLogger является логгером-заглушкой
type noopLogger struct{}

func (n *noopLogger) Debug(msg string, args ...any) {}
func (n *noopLogger) Info(msg string, args ...any)  {}
func (n *noopLogger) Warn(msg string, args ...any)  {}
func (n *noopLogger) Error(msg string, args ...any) {}
func (n *noopLogger) Fatal(msg string, args ...any) {}
func (n *noopLogger) With(args ...any) Logger       { return n }
func (n *noopLogger) WithGroup(name string) Logger  { return n }
