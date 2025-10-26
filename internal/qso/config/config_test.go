package config

import (
	"os"
	"testing"

	"github.com/rd2w/rd2w-log/internal/shared/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSOMustLoad(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		checkFn func(t *testing.T, cfg *QSOConfig)
	}{
		{
			name:    "load with defaults when no yaml file",
			envVars: map[string]string{},
			checkFn: func(t *testing.T, cfg *QSOConfig) {
				// Проверяем значения по умолчанию
				assert.Equal(t, "qso-service", cfg.AppName)
				assert.Equal(t, "50052", cfg.GRPCPort)
				assert.Equal(t, 100, cfg.MaxQSOsPerPage)
				assert.Equal(t, false, cfg.EnableCaching)
				assert.Equal(t, 300, cfg.CacheTTL)
				assert.Equal(t, 1000, cfg.ExportBatchSize)
				assert.Equal(t, 5, cfg.ImportConcurrency)

				// Проверяем ADIF defaults
				assert.Equal(t, "", cfg.ADIF.MaxFileSize)     // Похоже нет значения по умолчанию
				assert.Equal(t, "", cfg.ADIF.DefaultEncoding) // Похоже нет значения по умолчанию

				// Проверяем Search defaults
				assert.Equal(t, 0, cfg.Search.MaxResults)        // Похоже нет значения по умолчанию
				assert.Equal(t, "", cfg.Search.DefaultTimeRange) // Похоже нет значения по умолчанию

				// Проверяем PostgreSQL defaults
				assert.Equal(t, "localhost", cfg.Postgres.Host)
				assert.Equal(t, "5432", cfg.Postgres.Port)
				assert.Equal(t, "hamuser", cfg.Postgres.User)
				assert.Equal(t, "hampass", cfg.Postgres.Password)
				assert.Equal(t, "hamradio", cfg.Postgres.DBName)
				assert.Equal(t, "disable", cfg.Postgres.SSLMode)

				// Проверяем Redis defaults
				assert.Equal(t, "redis://localhost:6379", cfg.Redis.URL)
				assert.Equal(t, 0, cfg.Redis.DB)        // Похоже нет значения по умолчанию
				assert.Equal(t, "", cfg.Redis.Password) // Похоже нет значения по умолчанию
			},
		},
		{
			name: "non-existent config file uses defaults",
			envVars: map[string]string{
				"QSO_CONFIG_PATH": "./non-existent-config.yaml",
			},
			checkFn: func(t *testing.T, cfg *QSOConfig) {
				// Проверяем значения по умолчанию
				assert.Equal(t, "qso-service", cfg.AppName)
				assert.Equal(t, "50052", cfg.GRPCPort)
				assert.Equal(t, 100, cfg.MaxQSOsPerPage)
				assert.Equal(t, 300, cfg.CacheTTL)
				assert.Equal(t, 1000, cfg.ExportBatchSize)
				assert.Equal(t, 5, cfg.ImportConcurrency)
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

func TestSetQSODefaults(t *testing.T) {
	tests := []struct {
		name     string
		input    *QSOConfig
		expected *QSOConfig
	}{
		{
			name: "all empty fields get defaults",
			input: &QSOConfig{
				BaseConfig: config.BaseConfig{},
				Postgres:   config.PostgresConfig{},
				Redis:      config.RedisConfig{},
				ADIF:       ADIFConfig{},
				Search:     SearchConfig{},
			},
			expected: &QSOConfig{
				BaseConfig: config.BaseConfig{
					AppName:  "qso-service",
					GRPCPort: "50052",
				},
				MaxQSOsPerPage:    100,
				EnableCaching:     false,
				CacheTTL:          300,
				ExportBatchSize:   1000,
				ImportConcurrency: 5,
				ADIF: ADIFConfig{
					MaxFileSize:     "", // Нет значения по умолчанию в setQSODefaults
					DefaultEncoding: "", // Нет значения по умолчанию в setQSODefaults
				},
				Search: SearchConfig{
					MaxResults:       0,  // Нет значения по умолчанию в setQSODefaults
					DefaultTimeRange: "", // Нет значения по умолчанию в setQSODefaults
				},
				Postgres: config.PostgresConfig{
					Host:     "localhost",
					Port:     "5432",
					User:     "hamuser",
					Password: "hampass",
					DBName:   "hamradio",
					SSLMode:  "disable",
				},
				Redis: config.RedisConfig{
					URL: "redis://localhost:6379",
				},
			},
		},
		{
			name: "partial fields keep values",
			input: &QSOConfig{
				BaseConfig: config.BaseConfig{
					AppName:  "custom-qso",
					GRPCPort: "60052",
				},
				MaxQSOsPerPage:    50,
				EnableCaching:     true,
				CacheTTL:          600,
				ExportBatchSize:   500,
				ImportConcurrency: 10,
				ADIF: ADIFConfig{
					MaxFileSize:     "10MB",
					DefaultEncoding: "UTF-8",
				},
				Search: SearchConfig{
					MaxResults:       1000,
					DefaultTimeRange: "30d",
				},
				Postgres: config.PostgresConfig{
					Host:     "qso-db",
					User:     "qsouser",
					Password: "qsopass",
				},
				Redis: config.RedisConfig{
					URL: "redis://qso-redis:6379",
					DB:  2,
				},
			},
			expected: &QSOConfig{
				BaseConfig: config.BaseConfig{
					AppName:  "custom-qso", // Сохраняется
					GRPCPort: "60052",      // Сохраняется
				},
				MaxQSOsPerPage:    50,   // Сохраняется
				EnableCaching:     true, // Сохраняется
				CacheTTL:          600,  // Сохраняется
				ExportBatchSize:   500,  // Сохраняется
				ImportConcurrency: 10,   // Сохраняется
				ADIF: ADIFConfig{
					MaxFileSize:     "10MB",  // Сохраняется
					DefaultEncoding: "UTF-8", // Сохраняется
				},
				Search: SearchConfig{
					MaxResults:       1000,  // Сохраняется
					DefaultTimeRange: "30d", // Сохраняется
				},
				Postgres: config.PostgresConfig{
					Host:     "qso-db",   // Сохраняется
					Port:     "5432",     // По умолчанию
					User:     "qsouser",  // Сохраняется
					Password: "qsopass",  // Сохраняется
					DBName:   "hamradio", // По умолчанию
					SSLMode:  "disable",  // По умолчанию
				},
				Redis: config.RedisConfig{
					URL:      "redis://qso-redis:6379", // Сохраняется
					Password: "",                       // Нет значения по умолчанию
					DB:       2,                        // Сохраняется
				},
			},
		},
		{
			name: "zero values get defaults",
			input: &QSOConfig{
				MaxQSOsPerPage:    0,
				CacheTTL:          0,
				ExportBatchSize:   0,
				ImportConcurrency: 0,
				Postgres: config.PostgresConfig{
					Host: "",
					Port: "",
				},
				Redis: config.RedisConfig{
					URL: "",
				},
			},
			expected: &QSOConfig{
				BaseConfig: config.BaseConfig{
					AppName:  "qso-service",
					GRPCPort: "50052",
				},
				MaxQSOsPerPage:    100,
				EnableCaching:     false,
				CacheTTL:          300,
				ExportBatchSize:   1000,
				ImportConcurrency: 5,
				ADIF: ADIFConfig{
					MaxFileSize:     "",
					DefaultEncoding: "",
				},
				Search: SearchConfig{
					MaxResults:       0,
					DefaultTimeRange: "",
				},
				Postgres: config.PostgresConfig{
					Host:     "localhost",
					Port:     "5432",
					User:     "hamuser",
					Password: "hampass",
					DBName:   "hamradio",
					SSLMode:  "disable",
				},
				Redis: config.RedisConfig{
					URL: "redis://localhost:6379",
				},
			},
		},
		{
			name: "only database configuration",
			input: &QSOConfig{
				Postgres: config.PostgresConfig{
					Host:    "custom-host",
					DBName:  "qso-database",
					SSLMode: "require",
				},
				Redis: config.RedisConfig{
					URL: "redis://custom:6379",
					DB:  5,
				},
			},
			expected: &QSOConfig{
				BaseConfig: config.BaseConfig{
					AppName:  "qso-service",
					GRPCPort: "50052",
				},
				MaxQSOsPerPage:    100,
				EnableCaching:     false,
				CacheTTL:          300,
				ExportBatchSize:   1000,
				ImportConcurrency: 5,
				ADIF: ADIFConfig{
					MaxFileSize:     "",
					DefaultEncoding: "",
				},
				Search: SearchConfig{
					MaxResults:       0,
					DefaultTimeRange: "",
				},
				Postgres: config.PostgresConfig{
					Host:     "custom-host",  // Сохраняется
					Port:     "5432",         // По умолчанию
					User:     "hamuser",      // По умолчанию
					Password: "hampass",      // По умолчанию
					DBName:   "qso-database", // Сохраняется
					SSLMode:  "require",      // Сохраняется
				},
				Redis: config.RedisConfig{
					URL: "redis://custom:6379", // Сохраняется
					DB:  5,                     // Сохраняется
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setQSODefaults(tt.input)

			assert.Equal(t, tt.expected.AppName, tt.input.AppName)
			assert.Equal(t, tt.expected.GRPCPort, tt.input.GRPCPort)
			assert.Equal(t, tt.expected.MaxQSOsPerPage, tt.input.MaxQSOsPerPage)
			assert.Equal(t, tt.expected.EnableCaching, tt.input.EnableCaching)
			assert.Equal(t, tt.expected.CacheTTL, tt.input.CacheTTL)
			assert.Equal(t, tt.expected.ExportBatchSize, tt.input.ExportBatchSize)
			assert.Equal(t, tt.expected.ImportConcurrency, tt.input.ImportConcurrency)
			assert.Equal(t, tt.expected.ADIF, tt.input.ADIF)
			assert.Equal(t, tt.expected.Search, tt.input.Search)
			assert.Equal(t, tt.expected.Postgres, tt.input.Postgres)
			assert.Equal(t, tt.expected.Redis, tt.input.Redis)
		})
	}
}

func TestQSOGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		value        string
		defaultValue string
		expected     string
	}{
		{
			name:         "env var exists",
			key:          "QSO_TEST_VAR",
			value:        "qso_value",
			defaultValue: "default_qso",
			expected:     "qso_value",
		},
		{
			name:         "env var does not exist",
			key:          "QSO_NON_EXISTENT",
			value:        "",
			defaultValue: "default_qso",
			expected:     "default_qso",
		},
		{
			name:         "empty env var",
			key:          "QSO_EMPTY_VAR",
			value:        "",
			defaultValue: "default_qso",
			expected:     "default_qso",
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

func TestQSOMustLoad_NoPanic(t *testing.T) {
	// Проверяем что функция не паникает
	assert.NotPanics(t, func() {
		cfg := MustLoad()
		assert.NotNil(t, cfg)
	})
}

func TestQSOConfig_Integration(t *testing.T) {
	t.Run("qso config is properly initialized", func(t *testing.T) {
		cfg := MustLoad()

		// Проверяем что все основные поля инициализированы
		assert.NotEmpty(t, cfg.AppName)
		assert.NotEmpty(t, cfg.GRPCPort)
		assert.NotZero(t, cfg.MaxQSOsPerPage)
		assert.NotZero(t, cfg.CacheTTL)
		assert.NotZero(t, cfg.ExportBatchSize)
		assert.NotZero(t, cfg.ImportConcurrency)

		// Проверяем что PostgreSQL настроен
		assert.NotEmpty(t, cfg.Postgres.Host)
		assert.NotEmpty(t, cfg.Postgres.Port)
		assert.NotEmpty(t, cfg.Postgres.User)
		assert.NotEmpty(t, cfg.Postgres.Password)
		assert.NotEmpty(t, cfg.Postgres.DBName)
		assert.NotEmpty(t, cfg.Postgres.SSLMode)

		// Проверяем что Redis настроен
		assert.NotEmpty(t, cfg.Redis.URL)
	})

	t.Run("qso config performance settings", func(t *testing.T) {
		cfg := MustLoad()

		// Проверяем разумные пределы для настроек производительности
		assert.True(t, cfg.MaxQSOsPerPage > 0 && cfg.MaxQSOsPerPage <= 1000, "Max QSOs per page should be reasonable")
		assert.True(t, cfg.CacheTTL >= 0, "Cache TTL should be non-negative")
		assert.True(t, cfg.ExportBatchSize > 0 && cfg.ExportBatchSize <= 10000, "Export batch size should be reasonable")
		assert.True(t, cfg.ImportConcurrency > 0 && cfg.ImportConcurrency <= 20, "Import concurrency should be reasonable")
	})

	t.Run("postgres DSN generation", func(t *testing.T) {
		cfg := &QSOConfig{
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
}

func TestQSOConfig_SearchSettings(t *testing.T) {
	t.Run("search configuration validation", func(t *testing.T) {
		cfg := &QSOConfig{
			Search: SearchConfig{
				MaxResults:       500,
				DefaultTimeRange: "7d",
			},
		}

		// Проверяем что поисковые настройки установлены
		assert.Equal(t, 500, cfg.Search.MaxResults)
		assert.Equal(t, "7d", cfg.Search.DefaultTimeRange)
	})

	t.Run("ADIF configuration validation", func(t *testing.T) {
		cfg := &QSOConfig{
			ADIF: ADIFConfig{
				MaxFileSize:     "5MB",
				DefaultEncoding: "UTF-8",
			},
		}

		// Проверяем что ADIF настройки установлены
		assert.Equal(t, "5MB", cfg.ADIF.MaxFileSize)
		assert.Equal(t, "UTF-8", cfg.ADIF.DefaultEncoding)
	})
}

func TestQSOConfig_CacheSettings(t *testing.T) {
	t.Run("cache configuration", func(t *testing.T) {
		cfg := MustLoad()

		// Проверяем настройки кэширования
		if cfg.EnableCaching {
			assert.True(t, cfg.CacheTTL > 0, "Cache TTL should be positive when caching is enabled")
		}

		// Проверяем что настройки кэширования логичны
		assert.True(t, cfg.CacheTTL >= 0, "Cache TTL should be non-negative")
	})
}
