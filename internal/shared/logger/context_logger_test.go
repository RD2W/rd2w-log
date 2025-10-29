package logger

import (
	"context"
	"testing"
)

func TestFromContext(t *testing.T) {
	// Тест 1: Контекст без логгера возвращает noopLogger
	t.Run("returns noop logger when no logger in context", func(t *testing.T) {
		ctx := context.Background()
		logger := FromContext(ctx)

		if _, ok := logger.(*noopLogger); !ok {
			t.Errorf("Expected noopLogger, got %T", logger)
		}
	})

	// Тест 2: Контекст с логгером возвращает правильный логгер
	t.Run("returns logger from context", func(t *testing.T) {
		ctx := context.Background()
		expectedLogger, _ := NewTest()
		ctx = WithContext(ctx, expectedLogger)

		logger := FromContext(ctx)

		if logger != expectedLogger {
			t.Errorf("Expected %T, got %T", expectedLogger, logger)
		}
	})
}

func TestWithContext(t *testing.T) {
	ctx := context.Background()
	logger, _ := NewTest()

	newCtx := WithContext(ctx, logger)

	retrieved := FromContext(newCtx)
	if retrieved != logger {
		t.Errorf("WithContext failed to store logger in context")
	}
}

func TestWithContextFields(t *testing.T) {
	t.Run("adds fields to existing logger", func(t *testing.T) {
		ctx := context.Background()
		baseLogger, _ := NewTest()
		ctx = WithContext(ctx, baseLogger)

		// Добавляем поля в контекст
		ctxWithFields := WithContextFields(ctx, "key", "value", "number", 42)

		logger := FromContext(ctxWithFields)
		if logger == baseLogger {
			t.Errorf("Expected new logger with fields, got same logger")
		}
	})

	t.Run("works with noop logger", func(t *testing.T) {
		ctx := context.Background()

		ctxWithFields := WithContextFields(ctx, "key", "value")

		logger := FromContext(ctxWithFields)
		if _, ok := logger.(*noopLogger); !ok {
			t.Errorf("Expected noopLogger when adding fields to empty context")
		}
	})
}
