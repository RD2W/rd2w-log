package logger

import (
	"context"
	"testing"
)

func BenchmarkLogger_Info(b *testing.B) {
	logger, _ := NewTest()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message", "iteration", i, "data", "test")
	}
}

func BenchmarkLogger_WithContext(b *testing.B) {
	ctx := context.Background()
	logger, _ := NewTest()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WithContext(ctx, logger)
	}
}

func BenchmarkFromContext(b *testing.B) {
	ctx := context.Background()
	logger, _ := NewTest()
	ctx = WithContext(ctx, logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FromContext(ctx)
	}
}
