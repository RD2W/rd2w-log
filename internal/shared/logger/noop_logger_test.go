package logger

import (
	"testing"
)

func TestNoopLogger(t *testing.T) {
	logger := &noopLogger{}

	// Эти вызовы не должны паниковать
	logger.Debug("test")
	logger.Info("test")
	logger.Warn("test")
	logger.Error("test")

	// With должен возвращать тот же логгер
	newLogger := logger.With("key", "value")
	if newLogger != logger {
		t.Errorf("noopLogger.With() should return same instance")
	}

	// WithGroup должен возвращать тот же логгер
	groupLogger := logger.WithGroup("test")
	if groupLogger != logger {
		t.Errorf("noopLogger.WithGroup() should return same instance")
	}
}
