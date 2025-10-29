package logger

import (
	"strings"
	"testing"
)

func TestLogLevels(t *testing.T) {
	levels := []struct {
		name  string
		level string
		debug bool
	}{
		{"debug", "debug", true},
		{"info", "info", false},
		{"warn", "warn", false},
		{"error", "error", false},
	}

	for _, tt := range levels {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				logger, _ := New(tt.level)
				logger.Debug("debug message")
			})

			hasDebug := strings.Contains(output, "debug message")
			if tt.debug && !hasDebug {
				t.Errorf("Debug level should output debug messages")
			}
			if !tt.debug && hasDebug {
				t.Errorf("%s level should not output debug messages", tt.level)
			}
		})
	}
}
