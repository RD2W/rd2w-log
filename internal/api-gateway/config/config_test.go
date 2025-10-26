package config

import (
	"os"
	"testing"
	"time"

	"github.com/rd2w/rd2w-log/internal/shared/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMustLoad(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		checkFn func(t *testing.T, cfg *GatewayConfig)
	}{
		{
			name:    "load with defaults when no yaml file",
			envVars: map[string]string{},
			checkFn: func(t *testing.T, cfg *GatewayConfig) {
				// Проверяем значения по умолчанию
				assert.Equal(t, "api-gateway", cfg.AppName)
				assert.Equal(t, "8080", cfg.HTTPPort)
				assert.Equal(t, 100, cfg.RateLimitPerMinute)
				assert.Equal(t, 30, cfg.Timeout)
				assert.Equal(t, false, cfg.EnableSwagger)
				assert.Equal(t, []string{"http://localhost:3000", "http://127.0.0.1:3000"}, cfg.CORSAllowedOrigins)

				// Проверяем Services по умолчанию
				assert.Equal(t, "localhost:50051", cfg.Services.AuthServiceURL)
				assert.Equal(t, "localhost:50052", cfg.Services.QSOServiceURL)
				assert.Equal(t, "localhost:50053", cfg.Services.AnalyticsServiceURL)

				// Проверяем Middleware по умолчанию
				assert.Equal(t, "Authorization", cfg.Middleware.Auth.HeaderName)
				assert.Equal(t, []string{
					"/api/v1/auth/login",
					"/api/v1/auth/register",
					"/api/v1/health",
					"/swagger/*",
				}, cfg.Middleware.Auth.ExcludePaths)

				// Проверяем GRPCClients по умолчанию
				assert.Equal(t, 30*time.Second, cfg.GRPCClients.Timeout)
				assert.Equal(t, 3, cfg.GRPCClients.Retries)
				assert.Equal(t, 30*time.Second, cfg.GRPCClients.HealthCheckInterval)
			},
		},
		{
			name: "non-existent config file uses defaults",
			envVars: map[string]string{
				"GATEWAY_CONFIG_PATH": "./non-existent-config.yaml",
			},
			checkFn: func(t *testing.T, cfg *GatewayConfig) {
				// Проверяем значения по умолчанию
				assert.Equal(t, "api-gateway", cfg.AppName)
				assert.Equal(t, "8080", cfg.HTTPPort)
				assert.Equal(t, 100, cfg.RateLimitPerMinute)
				assert.Equal(t, 30, cfg.Timeout)
				assert.Equal(t, false, cfg.EnableSwagger)
				assert.Equal(t, []string{"http://localhost:3000", "http://127.0.0.1:3000"}, cfg.CORSAllowedOrigins)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Устанавливаем env vars
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			// Вызываем тестируемую функцию
			cfg := MustLoad()

			// Проверяем результаты
			require.NotNil(t, cfg)
			tt.checkFn(t, cfg)
		})
	}
}

func TestSetGatewayDefaults(t *testing.T) {
	tests := []struct {
		name     string
		input    *GatewayConfig
		expected *GatewayConfig
	}{
		{
			name: "all empty fields get defaults",
			input: &GatewayConfig{
				BaseConfig:  config.BaseConfig{},
				Services:    config.ServicesConfig{},
				Middleware:  MiddlewareConfig{},
				GRPCClients: GRPCClientsConfig{},
			},
			expected: &GatewayConfig{
				BaseConfig: config.BaseConfig{
					AppName:  "api-gateway",
					HTTPPort: "8080",
				},
				RateLimitPerMinute: 100,
				Timeout:            30,
				CORSAllowedOrigins: []string{"http://localhost:3000", "http://127.0.0.1:3000"},
				Services: config.ServicesConfig{
					AuthServiceURL:      "localhost:50051",
					QSOServiceURL:       "localhost:50052",
					AnalyticsServiceURL: "localhost:50053",
				},
				Middleware: MiddlewareConfig{
					Auth: AuthMiddlewareConfig{
						HeaderName: "Authorization",
						ExcludePaths: []string{
							"/api/v1/auth/login",
							"/api/v1/auth/register",
							"/api/v1/health",
							"/swagger/*",
						},
					},
				},
				GRPCClients: GRPCClientsConfig{
					Timeout:             30 * time.Second,
					Retries:             3,
					HealthCheckInterval: 30 * time.Second,
				},
			},
		},
		{
			name: "partial fields keep values",
			input: &GatewayConfig{
				BaseConfig: config.BaseConfig{
					AppName:  "custom-gateway",
					HTTPPort: "9090",
				},
				RateLimitPerMinute: 500,
				Timeout:            60,
				CORSAllowedOrigins: []string{"https://custom.com"},
				Services: config.ServicesConfig{
					AuthServiceURL: "auth:60051",
				},
				GRPCClients: GRPCClientsConfig{
					Retries: 5,
				},
			},
			expected: &GatewayConfig{
				BaseConfig: config.BaseConfig{
					AppName:  "custom-gateway", // Сохраняется
					HTTPPort: "9090",           // Сохраняется
				},
				RateLimitPerMinute: 500,                            // Сохраняется
				Timeout:            60,                             // Сохраняется
				CORSAllowedOrigins: []string{"https://custom.com"}, // Сохраняется
				Services: config.ServicesConfig{
					AuthServiceURL:      "auth:60051",      // Сохраняется
					QSOServiceURL:       "localhost:50052", // По умолчанию
					AnalyticsServiceURL: "localhost:50053", // По умолчанию
				},
				Middleware: MiddlewareConfig{
					Auth: AuthMiddlewareConfig{
						HeaderName: "Authorization", // По умолчанию
						ExcludePaths: []string{
							"/api/v1/auth/login",
							"/api/v1/auth/register",
							"/api/v1/health",
							"/swagger/*",
						},
					},
				},
				GRPCClients: GRPCClientsConfig{
					Timeout:             30 * time.Second, // По умолчанию
					Retries:             5,                // Сохраняется
					HealthCheckInterval: 30 * time.Second, // По умолчанию
				},
			},
		},
		{
			name: "empty arrays get defaults",
			input: &GatewayConfig{
				CORSAllowedOrigins: []string{},
				Middleware: MiddlewareConfig{
					Auth: AuthMiddlewareConfig{
						ExcludePaths: []string{},
					},
				},
			},
			expected: &GatewayConfig{
				BaseConfig: config.BaseConfig{
					AppName:  "api-gateway",
					HTTPPort: "8080",
				},
				RateLimitPerMinute: 100,
				Timeout:            30,
				CORSAllowedOrigins: []string{"http://localhost:3000", "http://127.0.0.1:3000"},
				Services: config.ServicesConfig{
					AuthServiceURL:      "localhost:50051",
					QSOServiceURL:       "localhost:50052",
					AnalyticsServiceURL: "localhost:50053",
				},
				Middleware: MiddlewareConfig{
					Auth: AuthMiddlewareConfig{
						HeaderName: "Authorization",
						ExcludePaths: []string{
							"/api/v1/auth/login",
							"/api/v1/auth/register",
							"/api/v1/health",
							"/swagger/*",
						},
					},
				},
				GRPCClients: GRPCClientsConfig{
					Timeout:             30 * time.Second,
					Retries:             3,
					HealthCheckInterval: 30 * time.Second,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setGatewayDefaults(tt.input)

			assert.Equal(t, tt.expected.AppName, tt.input.AppName)
			assert.Equal(t, tt.expected.HTTPPort, tt.input.HTTPPort)
			assert.Equal(t, tt.expected.RateLimitPerMinute, tt.input.RateLimitPerMinute)
			assert.Equal(t, tt.expected.Timeout, tt.input.Timeout)
			assert.Equal(t, tt.expected.CORSAllowedOrigins, tt.input.CORSAllowedOrigins)
			assert.Equal(t, tt.expected.Services, tt.input.Services)
			assert.Equal(t, tt.expected.Middleware, tt.input.Middleware)
			assert.Equal(t, tt.expected.GRPCClients, tt.input.GRPCClients)
		})
	}
}

func TestGetCORSAllowedOrigins(t *testing.T) {
	tests := []struct {
		name     string
		config   *GatewayConfig
		expected []string
	}{
		{
			name: "empty origins - returns defaults",
			config: &GatewayConfig{
				CORSAllowedOrigins: []string{},
			},
			expected: []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		},
		{
			name: "single string with commas - splits correctly",
			config: &GatewayConfig{
				CORSAllowedOrigins: []string{"https://site1.com,https://site2.com,http://localhost:8080"},
			},
			expected: []string{"https://site1.com", "https://site2.com", "http://localhost:8080"},
		},
		{
			name: "multiple separate origins - returns as is",
			config: &GatewayConfig{
				CORSAllowedOrigins: []string{"https://site1.com", "https://site2.com"},
			},
			expected: []string{"https://site1.com", "https://site2.com"},
		},
		{
			name: "normal case with multiple origins",
			config: &GatewayConfig{
				CORSAllowedOrigins: []string{"https://production.com", "https://staging.com"},
			},
			expected: []string{"https://production.com", "https://staging.com"},
		},
		{
			name: "single origin without commas",
			config: &GatewayConfig{
				CORSAllowedOrigins: []string{"https://single-domain.com"},
			},
			expected: []string{"https://single-domain.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetCORSAllowedOrigins()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		value        string
		defaultValue string
		expected     string
	}{
		{
			name:         "env var exists",
			key:          "TEST_VAR",
			value:        "actual_value",
			defaultValue: "default_value",
			expected:     "actual_value",
		},
		{
			name:         "env var does not exist",
			key:          "NON_EXISTENT_VAR",
			value:        "",
			defaultValue: "default_value",
			expected:     "default_value",
		},
		{
			name:         "empty env var",
			key:          "EMPTY_VAR",
			value:        "",
			defaultValue: "default_value",
			expected:     "default_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				t.Setenv(tt.key, tt.value)
			} else {
				err := os.Unsetenv(tt.key)
				if err != nil {
					return
				}
			}

			result := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMustLoad_NoPanic(t *testing.T) {
	// Проверяем что функция не паникует
	assert.NotPanics(t, func() {
		cfg := MustLoad()
		assert.NotNil(t, cfg)
	})
}

func TestGatewayConfig_Integration(t *testing.T) {
	t.Run("config is properly initialized", func(t *testing.T) {
		cfg := MustLoad()

		// Проверяем что все основные поля инициализированы
		assert.NotEmpty(t, cfg.AppName)
		assert.NotEmpty(t, cfg.HTTPPort)
		assert.NotZero(t, cfg.RateLimitPerMinute)
		assert.NotZero(t, cfg.Timeout)

		// Проверяем что сервисы настроены
		assert.NotEmpty(t, cfg.Services.AuthServiceURL)
		assert.NotEmpty(t, cfg.Services.QSOServiceURL)
		assert.NotEmpty(t, cfg.Services.AnalyticsServiceURL)

		// Проверяем что middleware настроено
		assert.NotEmpty(t, cfg.Middleware.Auth.HeaderName)
		assert.NotEmpty(t, cfg.Middleware.Auth.ExcludePaths)

		// Проверяем что gRPC клиенты настроены
		assert.NotZero(t, cfg.GRPCClients.Timeout)
		assert.NotZero(t, cfg.GRPCClients.Retries)
		assert.NotZero(t, cfg.GRPCClients.HealthCheckInterval)
	})
}
