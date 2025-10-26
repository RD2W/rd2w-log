package config

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Константы для graceful shutdown
const (
	// DefaultShutdownTimeout - время по умолчанию для ожидания завершения операций
	DefaultShutdownTimeout = 30 * time.Second

	// DefaultShutdownWait - дополнительное время ожидания по умолчанию
	DefaultShutdownWait = 5 * time.Second

	// SignalChannelBuffer - размер буфера канала для сигналов
	SignalChannelBuffer = 1

	// MaxShutdownTimeout - максимальное время ожидания shutdown
	MaxShutdownTimeout = 5 * time.Minute

	// MinShutdownTimeout - минимальное время ожидания shutdown
	MinShutdownTimeout = 5 * time.Second

	// MaxShutdownWait - максимальное дополнительное время ожидания
	MaxShutdownWait = 30 * time.Second
)

// GracefulShutdownConfig содержит настройки graceful shutdown
type GracefulShutdownConfig struct {
	Timeout time.Duration `yaml:"timeout" env:"SHUTDOWN_TIMEOUT"`
	Wait    time.Duration `yaml:"wait" env:"SHUTDOWN_WAIT"`
}

// validateConfig валидирует настройки shutdown
func validateConfig(cfg *GracefulShutdownConfig) {
	if cfg == nil {
		log.Printf("⚠️ Shutdown config is nil, using defaults")
		return
	}

	// Проверяем и корректируем timeout
	if cfg.Timeout <= 0 {
		cfg.Timeout = DefaultShutdownTimeout
	} else if cfg.Timeout > MaxShutdownTimeout {
		log.Printf("⚠️ Shutdown timeout %v exceeds maximum %v, using maximum", cfg.Timeout, MaxShutdownTimeout)
		cfg.Timeout = MaxShutdownTimeout
	} else if cfg.Timeout < MinShutdownTimeout {
		log.Printf("⚠️ Shutdown timeout %v is less than minimum %v, using minimum", cfg.Timeout, MinShutdownTimeout)
		cfg.Timeout = MinShutdownTimeout
	}

	// Проверяем и корректируем wait
	if cfg.Wait < 0 {
		cfg.Wait = 0
	} else if cfg.Wait > MaxShutdownWait {
		log.Printf("⚠️ Shutdown wait %v exceeds maximum %v, using maximum", cfg.Wait, MaxShutdownWait)
		cfg.Wait = MaxShutdownWait
	}
}

// SetupGracefulShutdown настраивает graceful shutdown для HTTP сервера
func SetupGracefulShutdown(srv *http.Server, cfg *GracefulShutdownConfig, serviceName string) {
	// Валидируем конфигурацию
	validateConfig(cfg)

	// Канал для сигналов ОС
	quit := make(chan os.Signal, SignalChannelBuffer)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Горутина для ожидания сигнала остановки
	go func() {
		<-quit
		log.Printf("🛑 [%s] Shutdown signal received", serviceName)

		// Создаем контекст с таймаутом
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
		defer cancel()

		// Graceful shutdown сервера
		log.Printf("⏳ [%s] Starting graceful shutdown...", serviceName)
		log.Printf("⏰ [%s] Waiting up to %v for requests to complete", serviceName, cfg.Timeout)

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("⚠️ [%s] Server shutdown failed: %v", serviceName, err)
		} else {
			log.Printf("✅ [%s] All requests completed successfully", serviceName)
		}

		// Дополнительное ожидание если нужно
		if cfg.Wait > 0 {
			log.Printf("💤 [%s] Waiting additional %v before exit...", serviceName, cfg.Wait)
			time.Sleep(cfg.Wait)
		}

		log.Printf("👋 [%s] Service stopped gracefully", serviceName)
	}()
}

// SetupGracefulShutdownWithCallback расширенная версия с callback'ами
func SetupGracefulShutdownWithCallback(
	srv *http.Server,
	cfg *GracefulShutdownConfig,
	serviceName string,
	preShutdown func(),
	postShutdown func(),
) {
	// Валидируем конфигурацию
	validateConfig(cfg)

	// Канал для сигналов ОС
	quit := make(chan os.Signal, SignalChannelBuffer)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Горутина для ожидания сигнала остановки
	go func() {
		<-quit
		log.Printf("🛑 [%s] Shutdown signal received", serviceName)

		// Pre-shutdown callback
		if preShutdown != nil {
			log.Printf("🔧 [%s] Running pre-shutdown tasks...", serviceName)
			preShutdown()
		}

		// Создаем контекст с таймаутом
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
		defer cancel()

		// Graceful shutdown сервера
		log.Printf("⏳ [%s] Starting graceful shutdown...", serviceName)
		log.Printf("⏰ [%s] Waiting up to %v for requests to complete", serviceName, cfg.Timeout)

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("⚠️ [%s] Server shutdown failed: %v", serviceName, err)
		} else {
			log.Printf("✅ [%s] All requests completed successfully", serviceName)
		}

		// Post-shutdown callback
		if postShutdown != nil {
			log.Printf("🔧 [%s] Running post-shutdown tasks...", serviceName)
			postShutdown()
		}

		// Дополнительное ожидание если нужно
		if cfg.Wait > 0 {
			log.Printf("💤 [%s] Waiting additional %v before exit...", serviceName, cfg.Wait)
			time.Sleep(cfg.Wait)
		}

		log.Printf("👋 [%s] Service stopped gracefully", serviceName)
	}()
}

// SetupGracefulShutdownWithCustomSignals версия с кастомными сигналами
func SetupGracefulShutdownWithCustomSignals(
	srv *http.Server,
	cfg *GracefulShutdownConfig,
	serviceName string,
	signals []os.Signal,
) {
	// Валидируем конфигурацию
	validateConfig(cfg)

	// Если сигналы не указаны, используем стандартные
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	// Канал для сигналов ОС
	quit := make(chan os.Signal, SignalChannelBuffer)
	signal.Notify(quit, signals...)

	// Горутина для ожидания сигнала остановки
	go func() {
		sig := <-quit
		log.Printf("🛑 [%s] Shutdown signal received: %v", serviceName, sig)

		// Создаем контекст с таймаутом
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
		defer cancel()

		// Graceful shutdown сервера
		log.Printf("⏳ [%s] Starting graceful shutdown...", serviceName)
		log.Printf("⏰ [%s] Waiting up to %v for requests to complete", serviceName, cfg.Timeout)

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("⚠️ [%s] Server shutdown failed: %v", serviceName, err)
		} else {
			log.Printf("✅ [%s] All requests completed successfully", serviceName)
		}

		// Дополнительное ожидание если нужно
		if cfg.Wait > 0 {
			log.Printf("💤 [%s] Waiting additional %v before exit...", serviceName, cfg.Wait)
			time.Sleep(cfg.Wait)
		}

		log.Printf("👋 [%s] Service stopped gracefully", serviceName)
	}()
}

// GetDefaultShutdownConfig возвращает конфигурацию shutdown по умолчанию
func GetDefaultShutdownConfig() GracefulShutdownConfig {
	return GracefulShutdownConfig{
		Timeout: DefaultShutdownTimeout,
		Wait:    DefaultShutdownWait,
	}
}

// GetProductionShutdownConfig возвращает конфигурацию для production
func GetProductionShutdownConfig() GracefulShutdownConfig {
	return GracefulShutdownConfig{
		Timeout: 45 * time.Second, // Больше времени для production
		Wait:    10 * time.Second,
	}
}

// GetDevelopmentShutdownConfig возвращает конфигурацию для development
func GetDevelopmentShutdownConfig() GracefulShutdownConfig {
	return GracefulShutdownConfig{
		Timeout: 15 * time.Second, // Меньше времени для development
		Wait:    2 * time.Second,
	}
}
