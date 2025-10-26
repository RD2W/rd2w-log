package config

import (
	"os"
	"time"

	"github.com/rd2w/rd2w-log/internal/shared/config"
)

// AuthConfig содержит конфигурацию для сервиса аутентификации
type AuthConfig struct {
	config.BaseConfig

	// Graceful Shutdown
	Shutdown config.GracefulShutdownConfig `yaml:"shutdown"`

	// Базы данных
	Postgres config.PostgresConfig `yaml:"postgres"`
	Redis    config.RedisConfig    `yaml:"redis"`

	// JWT
	JWT JWTConfig `yaml:"jwt"`

	// Безопасность
	Security SecurityConfig `yaml:"security"`

	// Специфичные для Auth Service настройки
	PasswordMinLength int           `yaml:"password_min_length" env:"PASSWORD_MIN_LENGTH"`
	MaxLoginAttempts  int           `yaml:"max_login_attempts" env:"MAX_LOGIN_ATTEMPTS"`
	LoginBlockTime    time.Duration `yaml:"login_block_time" env:"LOGIN_BLOCK_TIME"`
}

// JWTConfig расширяет базовый JWTConfig
type JWTConfig struct {
	config.JWTConfig
	Issuer string `yaml:"issuer" env:"JWT_ISSUER"`
}

// SecurityConfig содержит настройки безопасности
type SecurityConfig struct {
	BcryptCost           int           `yaml:"bcrypt_cost" env:"BCRYPT_COST"`
	TokenCleanupInterval time.Duration `yaml:"token_cleanup_interval" env:"TOKEN_CLEANUP_INTERVAL"`
}

// MustLoad загружает конфигурацию для Auth Service
func MustLoad() *AuthConfig {
	cfg := &AuthConfig{}
	config.MustLoad(cfg, config.LoadOptions{
		ConfigPath: getEnv("AUTH_CONFIG_PATH", "./configs/auth-config.yaml"),
		EnvPrefix:  "AUTH",
	})
	setAuthDefaults(cfg)
	return cfg
}

// Load загружает конфигурацию с обработкой ошибок
func Load() (*AuthConfig, error) {
	cfg := &AuthConfig{}
	err := config.Load(cfg, config.LoadOptions{
		ConfigPath: getEnv("AUTH_CONFIG_PATH", "./configs/auth-config.yaml"),
		EnvPrefix:  "AUTH",
	})
	if err != nil {
		return nil, err
	}
	setAuthDefaults(cfg)
	return cfg, nil
}

func setAuthDefaults(cfg *AuthConfig) {
	if cfg.AppName == "" {
		cfg.AppName = "auth-service"
	}
	if cfg.GRPCPort == "" {
		cfg.GRPCPort = "50051"
	}
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = "auth-secret-key-change-in-production"
	}
	if cfg.JWT.AccessTokenTTL == 0 {
		cfg.JWT.AccessTokenTTL = 15 * time.Minute
	}
	if cfg.JWT.RefreshTokenTTL == 0 {
		cfg.JWT.RefreshTokenTTL = 7 * 24 * time.Hour
	}
	if cfg.JWT.Issuer == "" {
		cfg.JWT.Issuer = "ham-radio-auth"
	}
	if cfg.PasswordMinLength == 0 {
		cfg.PasswordMinLength = 8
	}
	if cfg.MaxLoginAttempts == 0 {
		cfg.MaxLoginAttempts = 5
	}
	if cfg.LoginBlockTime == 0 {
		cfg.LoginBlockTime = 30 * time.Minute
	}

	// PostgreSQL defaults
	if cfg.Postgres.Host == "" {
		cfg.Postgres.Host = "localhost"
	}
	if cfg.Postgres.Port == "" {
		cfg.Postgres.Port = "5432"
	}
	if cfg.Postgres.User == "" {
		cfg.Postgres.User = "hamuser"
	}
	if cfg.Postgres.Password == "" {
		cfg.Postgres.Password = "hampass"
	}
	if cfg.Postgres.DBName == "" {
		cfg.Postgres.DBName = "hamradio"
	}
	if cfg.Postgres.SSLMode == "" {
		cfg.Postgres.SSLMode = "disable"
	}

	// Redis defaults
	if cfg.Redis.URL == "" {
		cfg.Redis.URL = "redis://localhost:6379"
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
