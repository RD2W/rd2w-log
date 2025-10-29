package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type slogLogger struct {
	logger *slog.Logger
}

// New creates a new structured logger using slog
func New(level string) (Logger, error) {
	// Преобразуем строковый уровень в slog level
	var slogLevel slog.Level
	switch level {
	case "debug":
		slogLevel = slog.LevelDebug
	case "info":
		slogLevel = slog.LevelInfo
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	// Настройка handler
	opts := &slog.HandlerOptions{
		Level: slogLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Кастомные преобразования атрибутов
			if a.Key == slog.TimeKey {
				// Форматируем время в ISO8601
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(time.RFC3339))
				}
			}
			return a
		},
		AddSource: slogLevel == slog.LevelDebug, // Добавляем source только для debug
	}

	// Создаем JSON handler для production
	handler := slog.NewJSONHandler(os.Stdout, opts)

	return &slogLogger{
		logger: slog.New(handler),
	}, nil
}

// NewDevelopment creates a development logger with human-readable output
func NewDevelopment() (Logger, error) {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Упрощаем вывод для development
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format("15:04:05.000"))
				}
			}
			return a
		},
		AddSource: true,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)

	return &slogLogger{
		logger: slog.New(handler),
	}, nil
}

// NewTest creates a logger suitable for testing (minimal output)
func NewTest() (Logger, error) {
	opts := &slog.HandlerOptions{
		Level: slog.LevelError, // Только ошибки в тестах
	}

	handler := slog.NewTextHandler(os.Stdout, opts)

	return &slogLogger{
		logger: slog.New(handler),
	}, nil
}

func (l *slogLogger) Debug(msg string, args ...any) {
	l.log(logLevelDebug, msg, args...)
}

func (l *slogLogger) Info(msg string, args ...any) {
	l.log(logLevelInfo, msg, args...)
}

func (l *slogLogger) Warn(msg string, args ...any) {
	l.log(logLevelWarn, msg, args...)
}

func (l *slogLogger) Error(msg string, args ...any) {
	l.log(logLevelError, msg, args...)
}

func (l *slogLogger) Fatal(msg string, args ...any) {
	l.log(logLevelError, msg, args...)
	os.Exit(1)
}

func (l *slogLogger) With(args ...any) Logger {
	return &slogLogger{
		logger: l.logger.With(args...),
	}
}

func (l *slogLogger) WithGroup(name string) Logger {
	return &slogLogger{
		logger: l.logger.WithGroup(name),
	}
}

// Log выполняет логирование с учетом уровня
func (l *slogLogger) log(level logLevel, msg string, args ...any) {
	// Добавляем caller информацию для debug уровня
	if level == logLevelDebug {
		// Получаем информацию о caller
		pc, file, line, ok := runtime.Caller(2) // Пропускаем 2 фрейма
		if ok {
			caller := slog.String("caller", formatCaller(pc, file, line))
			args = append([]any{caller}, args...)
		}
	}

	switch level {
	case logLevelDebug:
		l.logger.Debug(msg, args...)
	case logLevelInfo:
		l.logger.Info(msg, args...)
	case logLevelWarn:
		l.logger.Warn(msg, args...)
	case logLevelError:
		l.logger.Error(msg, args...)
	}
}

// formatCaller форматирует информацию о caller
func formatCaller(pc uintptr, file string, line int) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		// Если не можем получить функцию, используем хотя бы файл и строку
		return fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	funcName := fn.Name()
	// Оставляем только имя функции (без пути к пакету)
	if idx := strings.LastIndex(funcName, "/"); idx != -1 {
		funcName = funcName[idx+1:]
	}
	// Убираем возможный путь к пакету до имени функции (до ".")
	if idx := strings.LastIndex(funcName, "."); idx != -1 {
		funcName = funcName[idx+1:]
	}

	// Теперь используем ВСЕ доступные данные
	return fmt.Sprintf("%s(%s:%d)", funcName, filepath.Base(file), line)
}

// logLevel представляет внутренний уровень логирования
type logLevel int

const (
	logLevelDebug logLevel = iota
	logLevelInfo
	logLevelWarn
	logLevelError
)
