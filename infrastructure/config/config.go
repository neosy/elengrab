package iconfig

import (
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/neosy/elengrab/pkg/nconfig"
)

// Основные настройки
type Config struct {
	AppName string `env:"APP_NAME" envDefault:"Elengrab"`
	// Application configuration
	AppConfig nconfig.AppConfig `envPrefix:""`

	HTMXServer HTMXServerConfig `envPrefix:"HTTP_SERVER_"`
	Elengrab   ElengrabConfig   `envPrefix:"ELENGRAB_"`
	SQLite     SQLiteConfig     `envPrefix:"SQLITE_"`
	Redis      RedisConfig      `envPrefix:"REDIS_"`
}

// Настройки HTMX сервера
type HTMXServerConfig struct {
	Address  string `env:"ADDRESS" envDefault:""`
	Port     string `env:"PORT" envDefault:"8080"`
	Compress bool   `env:"COMPRESS" envDefault:"true"`
}

type ElengrabConfig struct {
	AssetsDir        string                    `env:"ASSETS_DIR" envDefault:"/app_n/assets"`
	DownloaderBinDir string                    `env:"DOWNLOADER_BIN_DIR" envDefault:"/usr/local/bin"`
	DownloadsDir     string                    `env:"DOWNLOADS_DIR" envDefault:"/app_n/downloads"`
	DownloadWorkers  int                       `env:"DOWNLOAD_WORKERS" envDefault:"3"`
	LoadHistory      bool                      `env:"LOAD_HISTORY" envDefault:"true"`
	Maintenance      ElengrabMaintenanceConfig `envPrefix:"MAINTENANCE_"`
}

type ElengrabMaintenanceConfig struct {
	IntervalUpdateHash            time.Duration `env:"INTERVAL_UPDATE_HASH" envDefault:"8h"`
	IntervalDeleteDuplicates      time.Duration `env:"INTERVAL_DELETE_DUPLICATES" envDefault:"1h"`
	IntervalDeleteMissingFiles    time.Duration `env:"INTERVAL_DELETE_MISSING_FILES" envDefault:"30m"`
	IntervalDeleteFailedDownloads time.Duration `env:"INTERVAL_DELETE_FAILED_DOWNLOADS" envDefault:"1h"`

	// EnableMoveUnmatchedFiles controls whether the periodic
	// moveUnmatchedFiles operation is allowed. Default is false (disabled).
	EnableMoveUnmatchedFiles bool `env:"ENABLE_MOVE_UNMATCHED_FILES" envDefault:"false"`

	// DatabaseBackupsKeep defines how many of the latest backup files to keep.
	// If the value is 0, old backup files will not be cleaned up.
	DatabaseBackupsKeep int `env:"DATABASE_BACKUPS_KEEP" envDefault:"7"`
}

type SQLiteConfig struct {
	DataDir    string `env:"DATA_DIR" envDefault:"/app_n/sqlite/data"`
	BackupsDir string `env:"BACKUPS_DIR" envDefault:"/app_n/sqlite/backups"`
}

type RedisConfig struct {
	// SIMPLE(Обычный редис) or SENTINEL
	ConnectionType     string   `env:"CONNECTION_TYPE" envDefault:"SIMPLE"`
	Addresses          []string `env:"ADDRESSES" envSeparator:"," envDefault:"localhost:6379"`
	SentinelMasterName string   `env:"SENTINEL_MASTER_NAME" envDefault:"mymaster"`
	PrefixKey          string   `env:"PREFIX_KEY" envDefault:"microservice"`
}

// Создание объекта Config
func New(appName, appVersion *string) (*Config, error) {
	parseFlag()

	c := &Config{}

	if err := c.load(); err != nil {
		return nil, err
	}

	if appName != nil {
		c.AppName = *appName
	}
	if appVersion != nil {
		c.AppConfig.Version = *appVersion
	}

	return c, nil
}

// Load config from environment variables
func (config *Config) load() error {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables only")
	}
	if err := env.Parse(config); err != nil {
		return fmt.Errorf("config load(). Read configuration error: %w", err)
	}

	return nil
}
