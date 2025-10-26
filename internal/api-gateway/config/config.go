package config

import (
	"os"
	"strings"
	"time"

	"github.com/rd2w/rd2w-log/internal/shared/config"
)

// GatewayConfig содержит конфигурацию для API Gateway
type GatewayConfig struct {
	config.BaseConfig

	// Graceful Shutdown
	Shutdown config.GracefulShutdownConfig `yaml:"shutdown"`

	// Сервисы
	Services config.ServicesConfig `yaml:"services"`

	// Middleware
	Middleware MiddlewareConfig `yaml:"middleware"`

	// gRPC клиенты
	GRPCClients GRPCClientsConfig `yaml:"grpc_clients"`

	// Специфичные для Gateway настройки
	CORSAllowedOrigins []string `yaml:"cors_allowed_origins" env:"CORS_ALLOWED_ORIGINS"`
	RateLimitPerMinute int      `yaml:"rate_limit_per_minute" env:"RATE_LIMIT_PER_MINUTE"`
	EnableSwagger      bool     `yaml:"enable_swagger" env:"ENABLE_SWAGGER"`
	Timeout            int      `yaml:"timeout" env:"TIMEOUT"`
}

// AuthMiddlewareConfig содержит настройки middleware аутентификации
type AuthMiddlewareConfig struct {
	Enabled      bool     `yaml:"enabled"`
	HeaderName   string   `yaml:"header_name"`
	ExcludePaths []string `yaml:"exclude_paths"`
}

// MiddlewareConfig содержит все настройки middleware
type MiddlewareConfig struct {
	Auth AuthMiddlewareConfig `yaml:"auth"`
}

// GRPCClientsConfig содержит настройки gRPC клиентов
type GRPCClientsConfig struct {
	Timeout             time.Duration `yaml:"timeout"`
	Retries             int           `yaml:"retries"`
	HealthCheckInterval time.Duration `yaml:"health_check_interval"`
}

// MustLoad загружает конфигурацию для API Gateway
func MustLoad() *GatewayConfig {
	cfg := &GatewayConfig{}
	config.MustLoad(cfg, config.LoadOptions{
		ConfigPath: getEnv("GATEWAY_CONFIG_PATH", "./configs/gateway-config.yaml"),
		EnvPrefix:  "GATEWAY",
	})
	setGatewayDefaults(cfg)
	return cfg
}

func setGatewayDefaults(cfg *GatewayConfig) {
	if cfg.AppName == "" {
		cfg.AppName = "api-gateway"
	}
	if cfg.HTTPPort == "" {
		cfg.HTTPPort = "8080"
	}
	if cfg.RateLimitPerMinute == 0 {
		cfg.RateLimitPerMinute = 100
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30
	}
	if len(cfg.CORSAllowedOrigins) == 0 {
		cfg.CORSAllowedOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	}

	// Services defaults
	if cfg.Services.AuthServiceURL == "" {
		cfg.Services.AuthServiceURL = "localhost:50051"
	}
	if cfg.Services.QSOServiceURL == "" {
		cfg.Services.QSOServiceURL = "localhost:50052"
	}
	if cfg.Services.AnalyticsServiceURL == "" {
		cfg.Services.AnalyticsServiceURL = "localhost:50053"
	}

	// Middleware defaults
	if cfg.Middleware.Auth.HeaderName == "" {
		cfg.Middleware.Auth.HeaderName = "Authorization"
	}
	if len(cfg.Middleware.Auth.ExcludePaths) == 0 {
		cfg.Middleware.Auth.ExcludePaths = []string{
			"/api/v1/auth/login",
			"/api/v1/auth/register",
			"/api/v1/health",
			"/swagger/*",
		}
	}

	// gRPC Clients defaults
	if cfg.GRPCClients.Timeout == 0 {
		cfg.GRPCClients.Timeout = 30 * time.Second
	}
	if cfg.GRPCClients.Retries == 0 {
		cfg.GRPCClients.Retries = 3
	}
	if cfg.GRPCClients.HealthCheckInterval == 0 {
		cfg.GRPCClients.HealthCheckInterval = 30 * time.Second
	}
}

// GetCORSAllowedOrigins возвращает CORS origins как slice
func (g *GatewayConfig) GetCORSAllowedOrigins() []string {
	if len(g.CORSAllowedOrigins) == 0 {
		return []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	}

	// Если это строка с разделителями, разбиваем ее
	if len(g.CORSAllowedOrigins) == 1 && strings.Contains(g.CORSAllowedOrigins[0], ",") {
		return strings.Split(g.CORSAllowedOrigins[0], ",")
	}

	return g.CORSAllowedOrigins
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
