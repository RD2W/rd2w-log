package logger

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// captureOutput захватывает stdout во время выполнения функции
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	err := w.Close()
	if err != nil {
		return ""
	}
	os.Stdout = old

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	if err != nil {
		return ""
	}
	return buf.String()
}

func TestJSONOutput(t *testing.T) {
	output := captureOutput(func() {
		logger, _ := New("info")
		logger.Info("test message", "key", "value", "number", 42)
	})

	if !strings.Contains(output, `"msg":"test message"`) {
		t.Errorf("JSON output missing message: %s", output)
	}

	// Проверяем валидность JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(output), &jsonData); err != nil {
		t.Errorf("Invalid JSON output: %v, output: %s", err, output)
	}
}

func TestDebugWithCaller(t *testing.T) {
	output := captureOutput(func() {
		logger, _ := New("debug")
		logger.Debug("debug test", "user", "testuser")
	})

	// Debug должен содержать caller информацию
	if !strings.Contains(output, "caller") {
		t.Errorf("Debug output should contain caller info: %s", output)
	}
}
