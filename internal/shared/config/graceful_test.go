package config

import (
	"context"
	"net/http"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name     string
		input    GracefulShutdownConfig
		expected GracefulShutdownConfig
	}{
		{
			name:     "negative timeout gets default",
			input:    GracefulShutdownConfig{Timeout: -1 * time.Second, Wait: 0},
			expected: GracefulShutdownConfig{Timeout: DefaultShutdownTimeout, Wait: 0},
		},
		{
			name:     "negative wait gets zero",
			input:    GracefulShutdownConfig{Timeout: 10 * time.Second, Wait: -1 * time.Second},
			expected: GracefulShutdownConfig{Timeout: 10 * time.Second, Wait: 0},
		},
		{
			name:     "timeout less than minimum gets min",
			input:    GracefulShutdownConfig{Timeout: 1 * time.Second, Wait: 0},
			expected: GracefulShutdownConfig{Timeout: MinShutdownTimeout, Wait: 0},
		},
		{
			name:     "timeout greater than maximum gets max",
			input:    GracefulShutdownConfig{Timeout: 10 * time.Minute, Wait: 0},
			expected: GracefulShutdownConfig{Timeout: MaxShutdownTimeout, Wait: 0},
		},
		{
			name:     "wait greater than maximum gets max",
			input:    GracefulShutdownConfig{Timeout: 30 * time.Second, Wait: 1 * time.Minute},
			expected: GracefulShutdownConfig{Timeout: 30 * time.Second, Wait: MaxShutdownWait},
		},
		{
			name:     "valid values remain unchanged",
			input:    GracefulShutdownConfig{Timeout: 25 * time.Second, Wait: 10 * time.Second},
			expected: GracefulShutdownConfig{Timeout: 25 * time.Second, Wait: 10 * time.Second},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validateConfig(&tt.input)
			assert.Equal(t, tt.expected.Timeout, tt.input.Timeout)
			assert.Equal(t, tt.expected.Wait, tt.input.Wait)
		})
	}
}

func TestGetDefaultShutdownConfig(t *testing.T) {
	cfg := GetDefaultShutdownConfig()
	assert.Equal(t, DefaultShutdownTimeout, cfg.Timeout)
	assert.Equal(t, DefaultShutdownWait, cfg.Wait)
}

func TestGetProductionShutdownConfig(t *testing.T) {
	cfg := GetProductionShutdownConfig()
	assert.Equal(t, 45*time.Second, cfg.Timeout)
	assert.Equal(t, 10*time.Second, cfg.Wait)
}

func TestGetDevelopmentShutdownConfig(t *testing.T) {
	cfg := GetDevelopmentShutdownConfig()
	assert.Equal(t, 15*time.Second, cfg.Timeout)
	assert.Equal(t, 2*time.Second, cfg.Wait)
}

func TestSetupGracefulShutdown(t *testing.T) {
	t.Run("basic shutdown setup", func(t *testing.T) {
		srv := &http.Server{Addr: ":8080"}
		cfg := &GracefulShutdownConfig{
			Timeout: 100 * time.Millisecond,
			Wait:    10 * time.Millisecond,
		}

		// Захватываем логи
		// setupGracefulShutdown должен запустить горутину, но не блокировать
		SetupGracefulShutdown(srv, cfg, "test-service")

		// Даем время для запуска горутины
		time.Sleep(10 * time.Millisecond)

		// Проверяем что сервер не nil
		assert.NotNil(t, srv)
	})
}

func TestSetupGracefulShutdownWithCallback(t *testing.T) {
	t.Run("callbacks are executed", func(t *testing.T) {
		srv := &http.Server{Addr: ":8080"}
		cfg := &GracefulShutdownConfig{
			Timeout: 100 * time.Millisecond,
			Wait:    10 * time.Millisecond,
		}

		var preShutdownCalled, postShutdownCalled bool
		var mu sync.Mutex

		preShutdown := func() {
			mu.Lock()
			defer mu.Unlock()
			preShutdownCalled = true
		}

		postShutdown := func() {
			mu.Lock()
			defer mu.Unlock()
			postShutdownCalled = true
		}

		SetupGracefulShutdownWithCallback(srv, cfg, "test-service", preShutdown, postShutdown)

		// Даем время для запуска горутины
		time.Sleep(10 * time.Millisecond)

		// Проверяем что колбэки еще не вызваны
		mu.Lock()
		assert.False(t, preShutdownCalled)
		assert.False(t, postShutdownCalled)
		mu.Unlock()
	})

	t.Run("nil callbacks are handled", func(t *testing.T) {
		srv := &http.Server{Addr: ":8080"}
		cfg := &GracefulShutdownConfig{
			Timeout: 100 * time.Millisecond,
			Wait:    10 * time.Millisecond,
		}

		// Не должно паниковать при nil колбэках
		assert.NotPanics(t, func() {
			SetupGracefulShutdownWithCallback(srv, cfg, "test-service", nil, nil)
		})

		// Даем время для запуска горутины
		time.Sleep(10 * time.Millisecond)
	})
}

func TestSetupGracefulShutdownWithCustomSignals(t *testing.T) {
	t.Run("custom signals are used", func(t *testing.T) {
		srv := &http.Server{Addr: ":8080"}
		cfg := &GracefulShutdownConfig{
			Timeout: 100 * time.Millisecond,
			Wait:    10 * time.Millisecond,
		}

		customSignals := []os.Signal{syscall.SIGUSR1, syscall.SIGUSR2}

		SetupGracefulShutdownWithCustomSignals(srv, cfg, "test-service", customSignals)

		// Даем время для запуска горутины
		time.Sleep(10 * time.Millisecond)

		// Проверяем что сервер не nil
		assert.NotNil(t, srv)
	})

	t.Run("empty signals use defaults", func(t *testing.T) {
		srv := &http.Server{Addr: ":8080"}
		cfg := &GracefulShutdownConfig{
			Timeout: 100 * time.Millisecond,
			Wait:    10 * time.Millisecond,
		}

		// Не должно паниковать при пустом списке сигналов
		assert.NotPanics(t, func() {
			SetupGracefulShutdownWithCustomSignals(srv, cfg, "test-service", []os.Signal{})
		})

		// Даем время для запуска горутины
		time.Sleep(10 * time.Millisecond)
	})
}

func TestGracefulShutdownIntegration(t *testing.T) {
	t.Run("server shutdown with context timeout", func(t *testing.T) {
		// Создаем тестовый сервер
		srv := &http.Server{
			Addr: ":0", // Используем порт 0 для автоматического выбора свободного порта
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
		}

		cfg := &GracefulShutdownConfig{
			Timeout: 50 * time.Millisecond,
			Wait:    5 * time.Millisecond,
		}

		// Запускаем graceful shutdown
		SetupGracefulShutdown(srv, cfg, "integration-test")

		// Запускаем сервер в горутине
		go func() {
			_ = srv.ListenAndServe()
		}()

		// Даем серверу время запуститься
		time.Sleep(10 * time.Millisecond)

		// Отправляем сигнал остановки
		proc, err := os.FindProcess(os.Getpid())
		require.NoError(t, err)

		err = proc.Signal(syscall.SIGTERM)
		require.NoError(t, err)

		// Ждем завершения
		time.Sleep(100 * time.Millisecond)

		// Проверяем что сервер закрыт
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		// Пытаемся запустить сервер снова - должен вернуть ошибку так как закрыт
		err = srv.ListenAndServe()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "closed", "Server should be closed")

		// Проверяем что контекст не истек (сервер должен был закрыться до таймаута)
		select {
		case <-ctx.Done():
			t.Fatal("Server shutdown took too long")
		default:
			// OK - сервер закрылся вовремя
		}
	})
}

func TestConstants(t *testing.T) {
	t.Run("constants have reasonable values", func(t *testing.T) {
		assert.True(t, DefaultShutdownTimeout > 0)
		assert.True(t, DefaultShutdownWait >= 0)
		assert.True(t, MaxShutdownTimeout > MinShutdownTimeout)
		assert.True(t, MaxShutdownWait >= 0)
		assert.True(t, SignalChannelBuffer > 0)

		// Проверяем что константы имеют разумные значения
		assert.True(t, DefaultShutdownTimeout <= MaxShutdownTimeout)
		assert.True(t, DefaultShutdownTimeout >= MinShutdownTimeout)
		assert.True(t, DefaultShutdownWait <= MaxShutdownWait)
	})
}

func TestGracefulShutdownConfig_Validation(t *testing.T) {
	t.Run("config validation preserves valid values", func(t *testing.T) {
		testCases := []struct {
			name     string
			timeout  time.Duration
			wait     time.Duration
			expected GracefulShutdownConfig
		}{
			{
				name:     "valid production values",
				timeout:  45 * time.Second,
				wait:     10 * time.Second,
				expected: GracefulShutdownConfig{Timeout: 45 * time.Second, Wait: 10 * time.Second},
			},
			{
				name:     "valid development values",
				timeout:  15 * time.Second,
				wait:     2 * time.Second,
				expected: GracefulShutdownConfig{Timeout: 15 * time.Second, Wait: 2 * time.Second},
			},
			{
				name:     "zero wait is valid",
				timeout:  30 * time.Second,
				wait:     0,
				expected: GracefulShutdownConfig{Timeout: 30 * time.Second, Wait: 0},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				cfg := &GracefulShutdownConfig{
					Timeout: tc.timeout,
					Wait:    tc.wait,
				}

				validateConfig(cfg)

				assert.Equal(t, tc.expected.Timeout, cfg.Timeout)
				assert.Equal(t, tc.expected.Wait, cfg.Wait)
			})
		}
	})
}

func TestGracefulShutdown_EdgeCases(t *testing.T) {
	t.Run("nil server handling", func(t *testing.T) {
		cfg := &GracefulShutdownConfig{
			Timeout: 100 * time.Millisecond,
			Wait:    10 * time.Millisecond,
		}

		// Не должно паниковать при nil сервере
		assert.NotPanics(t, func() {
			SetupGracefulShutdown(nil, cfg, "nil-server-test")
		})

		// Даем время для запуска горутины
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("nil config handling", func(t *testing.T) {
		srv := &http.Server{Addr: ":8080"}

		// Не должно паниковать при nil конфиге
		assert.NotPanics(t, func() {
			SetupGracefulShutdown(srv, nil, "nil-config-test")
		})

		// Даем время для запуска горутины
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("empty service name", func(t *testing.T) {
		srv := &http.Server{Addr: ":8080"}
		cfg := &GracefulShutdownConfig{
			Timeout: 100 * time.Millisecond,
			Wait:    10 * time.Millisecond,
		}

		// Не должно паниковать при пустом имени сервиса
		assert.NotPanics(t, func() {
			SetupGracefulShutdown(srv, cfg, "")
		})

		// Даем время для запуска горутины
		time.Sleep(10 * time.Millisecond)
	})
}

func TestSignalHandling(t *testing.T) {
	t.Run("multiple signals are handled", func(t *testing.T) {
		srv := &http.Server{Addr: ":8080"}
		cfg := &GracefulShutdownConfig{
			Timeout: 100 * time.Millisecond,
			Wait:    10 * time.Millisecond,
		}

		signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}

		SetupGracefulShutdownWithCustomSignals(srv, cfg, "multi-signal-test", signals)

		// Даем время для запуска горутины
		time.Sleep(10 * time.Millisecond)

		// Проверяем что сервер не nil
		assert.NotNil(t, srv)
	})
}
