package config

import (
	"os"
	"testing"
	"time"

	"github.com/rd2w/rd2w-log/internal/shared/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyticsMustLoad(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		checkFn func(t *testing.T, cfg *AnalyticsConfig)
	}{
		{
			name:    "load with defaults when no yaml file",
			envVars: map[string]string{},
			checkFn: func(t *testing.T, cfg *AnalyticsConfig) {
				// Проверяем значения по умолчанию
				assert.Equal(t, "analytics-service", cfg.AppName)
				assert.Equal(t, "50053", cfg.GRPCPort)
				assert.Equal(t, false, cfg.CacheEnabled)
				assert.Equal(t, 1*time.Hour, cfg.CacheTTL)
				assert.Equal(t, 5*time.Minute, cfg.StatsUpdateInterval)
				assert.Equal(t, 10*time.Minute, cfg.ReportGenerationTimeout)
				assert.Equal(t, 10000, cfg.MaxReportSize)

				// Проверяем PostgreSQL defaults
				assert.Equal(t, "localhost", cfg.Postgres.Host)
				assert.Equal(t, "5432", cfg.Postgres.Port)
				assert.Equal(t, "hamuser", cfg.Postgres.User)
				assert.Equal(t, "hampass", cfg.Postgres.Password)
				assert.Equal(t, "hamradio", cfg.Postgres.DBName)
				assert.Equal(t, "disable", cfg.Postgres.SSLMode)
			},
		},
		{
			name: "non-existent config file uses defaults",
			envVars: map[string]string{
				"ANALYTICS_CONFIG_PATH": "./non-existent-config.yaml",
			},
			checkFn: func(t *testing.T, cfg *AnalyticsConfig) {
				// Проверяем значения по умолчанию
				assert.Equal(t, "analytics-service", cfg.AppName)
				assert.Equal(t, "50053", cfg.GRPCPort)
				assert.Equal(t, 1*time.Hour, cfg.CacheTTL)
				assert.Equal(t, 5*time.Minute, cfg.StatsUpdateInterval)
				assert.Equal(t, 10*time.Minute, cfg.ReportGenerationTimeout)
				assert.Equal(t, 10000, cfg.MaxReportSize)
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

func TestSetAnalyticsDefaults(t *testing.T) {
	tests := []struct {
		name     string
		input    *AnalyticsConfig
		expected *AnalyticsConfig
	}{
		{
			name: "all empty fields get defaults",
			input: &AnalyticsConfig{
				BaseConfig: config.BaseConfig{},
				Postgres:   config.PostgresConfig{},
			},
			expected: &AnalyticsConfig{
				BaseConfig: config.BaseConfig{
					AppName:  "analytics-service",
					GRPCPort: "50053",
				},
				CacheEnabled:            false,
				CacheTTL:                1 * time.Hour,
				StatsUpdateInterval:     5 * time.Minute,
				ReportGenerationTimeout: 10 * time.Minute,
				MaxReportSize:           10000,
				Postgres: config.PostgresConfig{
					Host:     "localhost",
					Port:     "5432",
					User:     "hamuser",
					Password: "hampass",
					DBName:   "hamradio",
					SSLMode:  "disable",
				},
			},
		},
		{
			name: "partial fields keep values",
			input: &AnalyticsConfig{
				BaseConfig: config.BaseConfig{
					AppName:  "custom-analytics",
					GRPCPort: "60053",
				},
				CacheEnabled:            true,
				CacheTTL:                2 * time.Hour,
				StatsUpdateInterval:     10 * time.Minute,
				ReportGenerationTimeout: 15 * time.Minute,
				MaxReportSize:           50000,
				Postgres: config.PostgresConfig{
					Host:     "db-server",
					Port:     "6432",
					User:     "customuser",
					Password: "custompass",
				},
			},
			expected: &AnalyticsConfig{
				BaseConfig: config.BaseConfig{
					AppName:  "custom-analytics", // Сохраняется
					GRPCPort: "60053",            // Сохраняется
				},
				CacheEnabled:            true,             // Сохраняется
				CacheTTL:                2 * time.Hour,    // Сохраняется
				StatsUpdateInterval:     10 * time.Minute, // Сохраняется
				ReportGenerationTimeout: 15 * time.Minute, // Сохраняется
				MaxReportSize:           50000,            // Сохраняется
				Postgres: config.PostgresConfig{
					Host:     "db-server",  // Сохраняется
					Port:     "6432",       // Сохраняется
					User:     "customuser", // Сохраняется
					Password: "custompass", // Сохраняется
					DBName:   "hamradio",   // По умолчанию
					SSLMode:  "disable",    // По умолчанию
				},
			},
		},
		{
			name: "zero values get defaults",
			input: &AnalyticsConfig{
				CacheTTL:                0,
				StatsUpdateInterval:     0,
				ReportGenerationTimeout: 0,
				MaxReportSize:           0,
				Postgres: config.PostgresConfig{
					Host: "",
					Port: "",
				},
			},
			expected: &AnalyticsConfig{
				BaseConfig: config.BaseConfig{
					AppName:  "analytics-service",
					GRPCPort: "50053",
				},
				CacheEnabled:            false,
				CacheTTL:                1 * time.Hour,
				StatsUpdateInterval:     5 * time.Minute,
				ReportGenerationTimeout: 10 * time.Minute,
				MaxReportSize:           10000,
				Postgres: config.PostgresConfig{
					Host:     "localhost",
					Port:     "5432",
					User:     "hamuser",
					Password: "hampass",
					DBName:   "hamradio",
					SSLMode:  "disable",
				},
			},
		},
		{
			name: "postgres partial configuration",
			input: &AnalyticsConfig{
				Postgres: config.PostgresConfig{
					Host:    "custom-host",
					DBName:  "custom-db",
					SSLMode: "require",
				},
			},
			expected: &AnalyticsConfig{
				BaseConfig: config.BaseConfig{
					AppName:  "analytics-service",
					GRPCPort: "50053",
				},
				CacheEnabled:            false,
				CacheTTL:                1 * time.Hour,
				StatsUpdateInterval:     5 * time.Minute,
				ReportGenerationTimeout: 10 * time.Minute,
				MaxReportSize:           10000,
				Postgres: config.PostgresConfig{
					Host:     "custom-host", // Сохраняется
					Port:     "5432",        // По умолчанию
					User:     "hamuser",     // По умолчанию
					Password: "hampass",     // По умолчанию
					DBName:   "custom-db",   // Сохраняется
					SSLMode:  "require",     // Сохраняется
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setAnalyticsDefaults(tt.input)

			assert.Equal(t, tt.expected.AppName, tt.input.AppName)
			assert.Equal(t, tt.expected.GRPCPort, tt.input.GRPCPort)
			assert.Equal(t, tt.expected.CacheEnabled, tt.input.CacheEnabled)
			assert.Equal(t, tt.expected.CacheTTL, tt.input.CacheTTL)
			assert.Equal(t, tt.expected.StatsUpdateInterval, tt.input.StatsUpdateInterval)
			assert.Equal(t, tt.expected.ReportGenerationTimeout, tt.input.ReportGenerationTimeout)
			assert.Equal(t, tt.expected.MaxReportSize, tt.input.MaxReportSize)
			assert.Equal(t, tt.expected.Postgres, tt.input.Postgres)
		})
	}
}

func TestAnalyticsGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		value        string
		defaultValue string
		expected     string
	}{
		{
			name:         "env var exists",
			key:          "ANALYTICS_TEST_VAR",
			value:        "analytics_value",
			defaultValue: "default_analytics",
			expected:     "analytics_value",
		},
		{
			name:         "env var does not exist",
			key:          "ANALYTICS_NON_EXISTENT",
			value:        "",
			defaultValue: "default_analytics",
			expected:     "default_analytics",
		},
		{
			name:         "empty env var",
			key:          "ANALYTICS_EMPTY_VAR",
			value:        "",
			defaultValue: "default_analytics",
			expected:     "default_analytics",
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

func TestAnalyticsMustLoad_NoPanic(t *testing.T) {
	// Проверяем что функция не паникует
	assert.NotPanics(t, func() {
		cfg := MustLoad()
		assert.NotNil(t, cfg)
	})
}

func TestAnalyticsConfig_Integration(t *testing.T) {
	t.Run("analytics config is properly initialized", func(t *testing.T) {
		cfg := MustLoad()

		// Проверяем что все основные поля инициализированы
		assert.NotEmpty(t, cfg.AppName)
		assert.NotEmpty(t, cfg.GRPCPort)
		assert.NotZero(t, cfg.CacheTTL)
		assert.NotZero(t, cfg.StatsUpdateInterval)
		assert.NotZero(t, cfg.ReportGenerationTimeout)
		assert.NotZero(t, cfg.MaxReportSize)

		// Проверяем что PostgreSQL настроен
		assert.NotEmpty(t, cfg.Postgres.Host)
		assert.NotEmpty(t, cfg.Postgres.Port)
		assert.NotEmpty(t, cfg.Postgres.User)
		assert.NotEmpty(t, cfg.Postgres.Password)
		assert.NotEmpty(t, cfg.Postgres.DBName)
		assert.NotEmpty(t, cfg.Postgres.SSLMode)
	})
}

func TestAnalyticsConfig_PostgresDSN(t *testing.T) {
	t.Run("postgres DSN generation", func(t *testing.T) {
		cfg := &AnalyticsConfig{
			Postgres: config.PostgresConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "hamuser",
				Password: "hampass",
				DBName:   "hamradio",
				SSLMode:  "disable",
			},
		}

		dsn := cfg.Postgres.GetDSN()
		expected := "host=localhost port=5432 user=hamuser password=hampass dbname=hamradio sslmode=disable"
		assert.Equal(t, expected, dsn)
	})

	t.Run("postgres DSN with custom values", func(t *testing.T) {
		cfg := &AnalyticsConfig{
			Postgres: config.PostgresConfig{
				Host:     "production-db",
				Port:     "6432",
				User:     "produser",
				Password: "prodpass",
				DBName:   "production_db",
				SSLMode:  "require",
			},
		}

		dsn := cfg.Postgres.GetDSN()
		expected := "host=production-db port=6432 user=produser password=prodpass dbname=production_db sslmode=require"
		assert.Equal(t, expected, dsn)
	})
}

func TestAnalyticsConfig_TimeDurations(t *testing.T) {
	t.Run("time durations are valid", func(t *testing.T) {
		cfg := MustLoad()

		// Проверяем что временные интервалы положительные
		assert.True(t, cfg.CacheTTL > 0)
		assert.True(t, cfg.StatsUpdateInterval > 0)
		assert.True(t, cfg.ReportGenerationTimeout > 0)

		// Проверяем разумные пределы
		assert.True(t, cfg.CacheTTL <= 24*time.Hour, "Cache TTL should be reasonable")
		assert.True(t, cfg.StatsUpdateInterval <= time.Hour, "Stats update interval should be reasonable")
		assert.True(t, cfg.ReportGenerationTimeout <= 30*time.Minute, "Report generation timeout should be reasonable")
	})
}
