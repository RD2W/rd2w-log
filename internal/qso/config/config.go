package config

import (
	"os"

	"github.com/rd2w/rd2w-log/internal/shared/config"
)

// QSOConfig содержит конфигурацию для сервиса QSO
type QSOConfig struct {
	config.BaseConfig

	// Graceful Shutdown
	Shutdown config.GracefulShutdownConfig `yaml:"shutdown"`

	// Базы данных
	Postgres config.PostgresConfig `yaml:"postgres"`
	Redis    config.RedisConfig    `yaml:"redis"`

	// Специфичные для QSO Service настройки
	MaxQSOsPerPage    int  `yaml:"max_qsos_per_page" env:"MAX_QSOS_PER_PAGE"`
	EnableCaching     bool `yaml:"enable_caching" env:"ENABLE_CACHING"`
	CacheTTL          int  `yaml:"cache_ttl" env:"CACHE_TTL"`
	ExportBatchSize   int  `yaml:"export_batch_size" env:"EXPORT_BATCH_SIZE"`
	ImportConcurrency int  `yaml:"import_concurrency" env:"IMPORT_CONCURRENCY"`

	// ADIF настройки
	ADIF ADIFConfig `yaml:"adif"`

	// Поиск
	Search SearchConfig `yaml:"search"`
}

// ADIFConfig содержит настройки для работы с ADIF форматом
type ADIFConfig struct {
	MaxFileSize     string `yaml:"max_file_size" env:"ADIF_MAX_FILE_SIZE"`
	DefaultEncoding string `yaml:"default_encoding" env:"ADIF_DEFAULT_ENCODING"`
}

// SearchConfig содержит настройки поиска QSO
type SearchConfig struct {
	MaxResults       int    `yaml:"max_results" env:"SEARCH_MAX_RESULTS"`
	DefaultTimeRange string `yaml:"default_time_range" env:"SEARCH_DEFAULT_TIME_RANGE"`
}

// MustLoad загружает конфигурацию для QSO Service
func MustLoad() *QSOConfig {
	cfg := &QSOConfig{}
	config.MustLoad(cfg, config.LoadOptions{
		ConfigPath: getEnv("QSO_CONFIG_PATH", "./configs/qso-config.yaml"),
		EnvPrefix:  "QSO",
	})
	setQSODefaults(cfg)
	return cfg
}

func setQSODefaults(cfg *QSOConfig) {
	if cfg.AppName == "" {
		cfg.AppName = "qso-service"
	}
	if cfg.GRPCPort == "" {
		cfg.GRPCPort = "50052"
	}
	if cfg.MaxQSOsPerPage == 0 {
		cfg.MaxQSOsPerPage = 100
	}
	if cfg.CacheTTL == 0 {
		cfg.CacheTTL = 300 // 5 minutes
	}
	if cfg.ExportBatchSize == 0 {
		cfg.ExportBatchSize = 1000
	}
	if cfg.ImportConcurrency == 0 {
		cfg.ImportConcurrency = 5
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
