package config

import (
	"os"
	"testing"
	"time"

	"github.com/rd2w/rd2w-log/internal/shared/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMustLoad(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		checkFn func(t *testing.T, cfg *AuthConfig)
	}{
		{
			name:    "load with defaults when no yaml file",
			envVars: map[string]string{},
			checkFn: func(t *testing.T, cfg *AuthConfig) {
				// Проверяем значения по умолчанию
				assert.Equal(t, "auth-service", cfg.AppName)
				assert.Equal(t, "50051", cfg.GRPCPort)
				assert.Equal(t, 8, cfg.PasswordMinLength)
				assert.Equal(t, 5, cfg.MaxLoginAttempts)
				assert.Equal(t, 30*time.Minute, cfg.LoginBlockTime)

				// Проверяем JWT defaults
				assert.Equal(t, "auth-secret-key-change-in-production", cfg.JWT.Secret)
				assert.Equal(t, 15*time.Minute, cfg.JWT.AccessTokenTTL)
				assert.Equal(t, 7*24*time.Hour, cfg.JWT.RefreshTokenTTL)
				assert.Equal(t, "ham-radio-auth", cfg.JWT.Issuer)

				// Проверяем Security defaults
				assert.Equal(t, 0, cfg.Security.BcryptCost)                          // Похоже нет значения по умолчанию в setAuthDefaults
				assert.Equal(t, time.Duration(0), cfg.Security.TokenCleanupInterval) // Похоже нет значения по умолчанию

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
				"AUTH_CONFIG_PATH": "./non-existent-config.yaml",
			},
			checkFn: func(t *testing.T, cfg *AuthConfig) {
				// Проверяем значения по умолчанию
				assert.Equal(t, "auth-service", cfg.AppName)
				assert.Equal(t, "50051", cfg.GRPCPort)
				assert.Equal(t, 8, cfg.PasswordMinLength)
				assert.Equal(t, 5, cfg.MaxLoginAttempts)
				assert.Equal(t, 30*time.Minute, cfg.LoginBlockTime)
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

func TestAuthLoad(t *testing.T) {
	t.Run("load with error handling", func(t *testing.T) {
		cfg, err := Load()

		require.NoError(t, err)
		require.NotNil(t, cfg)
		assert.Equal(t, "auth-service", cfg.AppName)
		assert.Equal(t, "50051", cfg.GRPCPort)
	})

	t.Run("load with custom config path", func(t *testing.T) {
		t.Setenv("AUTH_CONFIG_PATH", "./non-existent-config.yaml")

		cfg, err := Load()

		require.NoError(t, err)
		require.NotNil(t, cfg)
		// Должны использоваться значения по умолчанию
		assert.Equal(t, "auth-service", cfg.AppName)
	})
}

func TestSetAuthDefaults(t *testing.T) {
	tests := []struct {
		name     string
		input    *AuthConfig
		expected *AuthConfig
	}{
		{
			name: "all empty fields get defaults",
			input: &AuthConfig{
				BaseConfig: config.BaseConfig{},
				Postgres:   config.PostgresConfig{},
				Redis:      config.RedisConfig{},
				JWT:        JWTConfig{},
				Security:   SecurityConfig{},
			},
			expected: &AuthConfig{
				BaseConfig: config.BaseConfig{
					AppName:  "auth-service",
					GRPCPort: "50051",
				},
				PasswordMinLength: 8,
				MaxLoginAttempts:  5,
				LoginBlockTime:    30 * time.Minute,
				JWT: JWTConfig{
					JWTConfig: config.JWTConfig{
						Secret:          "auth-secret-key-change-in-production",
						AccessTokenTTL:  15 * time.Minute,
						RefreshTokenTTL: 7 * 24 * time.Hour,
					},
					Issuer: "ham-radio-auth",
				},
				Security: SecurityConfig{
					BcryptCost:           0, // Нет значения по умолчанию в setAuthDefaults
					TokenCleanupInterval: 0, // Нет значения по умолчанию в setAuthDefaults
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
			input: &AuthConfig{
				BaseConfig: config.BaseConfig{
					AppName:  "custom-auth",
					GRPCPort: "60051",
				},
				PasswordMinLength: 12,
				MaxLoginAttempts:  3,
				LoginBlockTime:    1 * time.Hour,
				JWT: JWTConfig{
					JWTConfig: config.JWTConfig{
						Secret:          "custom-secret",
						AccessTokenTTL:  30 * time.Minute,
						RefreshTokenTTL: 30 * 24 * time.Hour,
					},
					Issuer: "custom-issuer",
				},
				Security: SecurityConfig{
					BcryptCost:           12,
					TokenCleanupInterval: 1 * time.Hour,
				},
				Postgres: config.PostgresConfig{
					Host:     "auth-db",
					User:     "authuser",
					Password: "authpass",
				},
				Redis: config.RedisConfig{
					URL: "redis://auth-redis:6379",
					DB:  1,
				},
			},
			expected: &AuthConfig{
				BaseConfig: config.BaseConfig{
					AppName:  "custom-auth", // Сохраняется
					GRPCPort: "60051",       // Сохраняется
				},
				PasswordMinLength: 12,            // Сохраняется
				MaxLoginAttempts:  3,             // Сохраняется
				LoginBlockTime:    1 * time.Hour, // Сохраняется
				JWT: JWTConfig{
					JWTConfig: config.JWTConfig{
						Secret:          "custom-secret",     // Сохраняется
						AccessTokenTTL:  30 * time.Minute,    // Сохраняется
						RefreshTokenTTL: 30 * 24 * time.Hour, // Сохраняется
					},
					Issuer: "custom-issuer", // Сохраняется
				},
				Security: SecurityConfig{
					BcryptCost:           12,            // Сохраняется
					TokenCleanupInterval: 1 * time.Hour, // Сохраняется
				},
				Postgres: config.PostgresConfig{
					Host:     "auth-db",  // Сохраняется
					Port:     "5432",     // По умолчанию
					User:     "authuser", // Сохраняется
					Password: "authpass", // Сохраняется
					DBName:   "hamradio", // По умолчанию
					SSLMode:  "disable",  // По умолчанию
				},
				Redis: config.RedisConfig{
					URL:      "redis://auth-redis:6379", // Сохраняется
					Password: "",                        // Нет значения по умолчанию
					DB:       1,                         // Сохраняется
				},
			},
		},
		{
			name: "zero values get defaults",
			input: &AuthConfig{
				PasswordMinLength: 0,
				MaxLoginAttempts:  0,
				LoginBlockTime:    0,
				JWT: JWTConfig{
					JWTConfig: config.JWTConfig{
						Secret:          "",
						AccessTokenTTL:  0,
						RefreshTokenTTL: 0,
					},
					Issuer: "",
				},
				Postgres: config.PostgresConfig{
					Host: "",
					Port: "",
				},
				Redis: config.RedisConfig{
					URL: "",
				},
			},
			expected: &AuthConfig{
				BaseConfig: config.BaseConfig{
					AppName:  "auth-service",
					GRPCPort: "50051",
				},
				PasswordMinLength: 8,
				MaxLoginAttempts:  5,
				LoginBlockTime:    30 * time.Minute,
				JWT: JWTConfig{
					JWTConfig: config.JWTConfig{
						Secret:          "auth-secret-key-change-in-production",
						AccessTokenTTL:  15 * time.Minute,
						RefreshTokenTTL: 7 * 24 * time.Hour,
					},
					Issuer: "ham-radio-auth",
				},
				Security: SecurityConfig{
					BcryptCost:           0,
					TokenCleanupInterval: 0,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setAuthDefaults(tt.input)

			assert.Equal(t, tt.expected.AppName, tt.input.AppName)
			assert.Equal(t, tt.expected.GRPCPort, tt.input.GRPCPort)
			assert.Equal(t, tt.expected.PasswordMinLength, tt.input.PasswordMinLength)
			assert.Equal(t, tt.expected.MaxLoginAttempts, tt.input.MaxLoginAttempts)
			assert.Equal(t, tt.expected.LoginBlockTime, tt.input.LoginBlockTime)
			assert.Equal(t, tt.expected.JWT, tt.input.JWT)
			assert.Equal(t, tt.expected.Security, tt.input.Security)
			assert.Equal(t, tt.expected.Postgres, tt.input.Postgres)
			assert.Equal(t, tt.expected.Redis, tt.input.Redis)
		})
	}
}

func TestAuthGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		value        string
		defaultValue string
		expected     string
	}{
		{
			name:         "env var exists",
			key:          "AUTH_TEST_VAR",
			value:        "auth_value",
			defaultValue: "default_auth",
			expected:     "auth_value",
		},
		{
			name:         "env var does not exist",
			key:          "AUTH_NON_EXISTENT",
			value:        "",
			defaultValue: "default_auth",
			expected:     "default_auth",
		},
		{
			name:         "empty env var",
			key:          "AUTH_EMPTY_VAR",
			value:        "",
			defaultValue: "default_auth",
			expected:     "default_auth",
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

func TestAuthMustLoad_NoPanic(t *testing.T) {
	// Проверяем что функция не паникует
	assert.NotPanics(t, func() {
		cfg := MustLoad()
		assert.NotNil(t, cfg)
	})
}

func TestAuthConfig_Integration(t *testing.T) {
	t.Run("auth config is properly initialized", func(t *testing.T) {
		cfg := MustLoad()

		// Проверяем что все основные поля инициализированы
		assert.NotEmpty(t, cfg.AppName)
		assert.NotEmpty(t, cfg.GRPCPort)
		assert.NotZero(t, cfg.PasswordMinLength)
		assert.NotZero(t, cfg.MaxLoginAttempts)
		assert.NotZero(t, cfg.LoginBlockTime)

		// Проверяем JWT настройки
		assert.NotEmpty(t, cfg.JWT.Secret)
		assert.NotZero(t, cfg.JWT.AccessTokenTTL)
		assert.NotZero(t, cfg.JWT.RefreshTokenTTL)
		assert.NotEmpty(t, cfg.JWT.Issuer)

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

	t.Run("auth config security settings", func(t *testing.T) {
		cfg := MustLoad()

		// Проверяем разумные пределы для настроек безопасности
		assert.True(t, cfg.PasswordMinLength >= 6, "Password min length should be at least 6")
		assert.True(t, cfg.MaxLoginAttempts >= 3, "Max login attempts should be at least 3")
		assert.True(t, cfg.LoginBlockTime >= 5*time.Minute, "Login block time should be reasonable")

		// Проверяем JWT TTL
		assert.True(t, cfg.JWT.AccessTokenTTL >= 5*time.Minute, "Access token TTL should be reasonable")
		assert.True(t, cfg.JWT.RefreshTokenTTL >= time.Hour, "Refresh token TTL should be reasonable")
	})

	t.Run("postgres DSN generation", func(t *testing.T) {
		cfg := &AuthConfig{
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

func TestAuthConfig_JWTValidation(t *testing.T) {
	t.Run("JWT configuration validation", func(t *testing.T) {
		cfg := MustLoad()

		// Проверяем что JWT секрет не является дефолтным в продакшене (это должно быть в env)
		// Это просто проверка что значение установлено
		assert.NotEmpty(t, cfg.JWT.Secret)

		// Проверяем что TTL времена разумные
		assert.True(t, cfg.JWT.AccessTokenTTL > 0)
		assert.True(t, cfg.JWT.RefreshTokenTTL > cfg.JWT.AccessTokenTTL, "Refresh token should live longer than access token")
	})
}
