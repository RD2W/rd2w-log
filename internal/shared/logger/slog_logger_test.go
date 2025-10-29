package logger

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		level   string
		wantErr bool
	}{
		{"debug level", "debug", false},
		{"info level", "info", false},
		{"warn level", "warn", false},
		{"error level", "error", false},
		{"invalid level", "invalid", false}, // Должен использовать info по умолчанию
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if logger == nil {
				t.Error("New() returned nil logger")
			}
		})
	}
}

func TestSlogLogger_Levels(t *testing.T) {
	// Тестируем разные уровни логирования
	var _ bytes.Buffer
	// Здесь нужно временно подменить os.Stdout для захвата вывода

	logger, _ := NewTest()

	// Эти вызовы не должны паниковать
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")
}

func TestSlogLogger_With(t *testing.T) {
	logger, _ := NewTest()

	childLogger := logger.With("service", "test", "version", "1.0")

	if childLogger == logger {
		t.Errorf("With() should return new logger instance")
	}
}

func TestSlogLogger_WithGroup(t *testing.T) {
	logger, _ := NewTest()

	groupLogger := logger.WithGroup("http")

	if groupLogger == logger {
		t.Errorf("WithGroup() should return new logger instance")
	}
}

func TestNewDevelopment(t *testing.T) {
	logger, err := NewDevelopment()
	if err != nil {
		t.Errorf("NewDevelopment() unexpected error: %v", err)
	}
	if logger == nil {
		t.Error("NewDevelopment() returned nil logger")
	}
}

func TestNewTest(t *testing.T) {
	logger, err := NewTest()
	if err != nil {
		t.Errorf("NewTest() unexpected error: %v", err)
	}
	if logger == nil {
		t.Error("NewTest() returned nil logger")
	}
}
