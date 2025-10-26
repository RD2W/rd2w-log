package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v11"
	"gopkg.in/yaml.v3"
)

// BaseConfig содержит общие настройки для всех сервисов
type BaseConfig struct {
	AppName  string `yaml:"app_name" env:"APP_NAME"`
	Env      string `yaml:"env" env:"ENV"`
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL"`

	// Сервер
	HTTPPort string `yaml:"http_port" env:"HTTP_PORT"`
	GRPCPort string `yaml:"grpc_port" env:"GRPC_PORT"`
}

// PostgresConfig содержит настройки PostgreSQL
type PostgresConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     string `yaml:"port" env:"DB_PORT"`
	User     string `yaml:"user" env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	DBName   string `yaml:"db_name" env:"DB_NAME"`
	SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
}

// RedisConfig содержит настройки Redis
type RedisConfig struct {
	URL      string `yaml:"url" env:"REDIS_URL"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
	DB       int    `yaml:"db" env:"REDIS_DB"`
}

// JWTConfig содержит настройки JWT
type JWTConfig struct {
	Secret          string        `yaml:"secret" env:"JWT_SECRET"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl" env:"JWT_ACCESS_TTL"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env:"JWT_REFRESH_TTL"`
}

// ServicesConfig содержит адреса сервисов
type ServicesConfig struct {
	AuthServiceURL      string `yaml:"auth_service_url" env:"AUTH_SERVICE_URL"`
	QSOServiceURL       string `yaml:"qso_service_url" env:"QSO_SERVICE_URL"`
	AnalyticsServiceURL string `yaml:"analytics_service_url" env:"ANALYTICS_SERVICE_URL"`
}

// LoadOptions опции для загрузки конфигурации
type LoadOptions struct {
	ConfigPath string
	EnvPrefix  string
}

// MustLoad загружает конфигурацию или паникует при ошибке
func MustLoad(cfg interface{}, opts ...LoadOptions) {
	if err := Load(cfg, opts...); err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
}

// Load загружает конфигурацию из YAML файла и environment variables
func Load(cfg interface{}, opts ...LoadOptions) error {
	options := LoadOptions{
		ConfigPath: getEnv("CONFIG_PATH", ""),
		EnvPrefix:  getEnv("ENV_PREFIX", ""),
	}

	if len(opts) > 0 {
		options = opts[0]
	}

	// Загружаем из YAML файла, если он указан и существует
	if options.ConfigPath != "" {
		if err := loadFromYAML(options.ConfigPath, cfg); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to load YAML config: %w", err)
		}
	}

	// Перезаписываем значения из environment variables
	if err := loadFromEnv(cfg, options.EnvPrefix); err != nil {
		return fmt.Errorf("failed to load env config: %w", err)
	}

	return nil
}

// loadFromYAML загружает конфигурацию из YAML файла
func loadFromYAML(path string, cfg interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return nil
}

// loadFromEnv загружает конфигурацию из environment variables
func loadFromEnv(cfg interface{}, prefix string) error {
	if prefix == "" {
		return env.Parse(cfg)
	}

	opts := env.Options{
		Prefix:  prefix,
		TagName: "env",
	}

	return env.ParseWithOptions(cfg, opts)
}

// GetDSN возвращает DSN строку для PostgreSQL
func (p *PostgresConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode)
}

// Вспомогательные функции
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
