package config

import (
	"os"
	"time"

	"github.com/rd2w/rd2w-log/internal/shared/config"
)

// AnalyticsConfig содержит конфигурацию для сервиса аналитики
type AnalyticsConfig struct {
	config.BaseConfig

	// Graceful Shutdown
	Shutdown config.GracefulShutdownConfig `yaml:"shutdown"`

	// Базы данных
	Postgres config.PostgresConfig `yaml:"postgres"`

	// Специфичные для Analytics Service настройки
	CacheEnabled            bool          `yaml:"cache_enabled" env:"CACHE_ENABLED"`
	CacheTTL                time.Duration `yaml:"cache_ttl" env:"CACHE_TTL"`
	StatsUpdateInterval     time.Duration `yaml:"stats_update_interval" env:"STATS_UPDATE_INTERVAL"`
	ReportGenerationTimeout time.Duration `yaml:"report_generation_timeout" env:"REPORT_GENERATION_TIMEOUT"`
	MaxReportSize           int           `yaml:"max_report_size" env:"MAX_REPORT_SIZE"`
}

// MustLoad загружает конфигурацию для Analytics Service
func MustLoad() *AnalyticsConfig {
	cfg := &AnalyticsConfig{}
	config.MustLoad(cfg, config.LoadOptions{
		ConfigPath: getEnv("ANALYTICS_CONFIG_PATH", "./configs/analytics-config.yaml"),
		EnvPrefix:  "ANALYTICS",
	})
	setAnalyticsDefaults(cfg)
	return cfg
}

func setAnalyticsDefaults(cfg *AnalyticsConfig) {
	if cfg.AppName == "" {
		cfg.AppName = "analytics-service"
	}
	if cfg.GRPCPort == "" {
		cfg.GRPCPort = "50053"
	}
	if cfg.CacheTTL == 0 {
		cfg.CacheTTL = 1 * time.Hour
	}
	if cfg.StatsUpdateInterval == 0 {
		cfg.StatsUpdateInterval = 5 * time.Minute
	}
	if cfg.ReportGenerationTimeout == 0 {
		cfg.ReportGenerationTimeout = 10 * time.Minute
	}
	if cfg.MaxReportSize == 0 {
		cfg.MaxReportSize = 10000
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
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
