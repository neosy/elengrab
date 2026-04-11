package iconfig

import (
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	nconfig "github.com/neosy/elengrab/internal/pkg/config"
)

// Basic Settings
type Config struct {
	// Global application settings, no ENV prefix.
	AppConfig nconfig.AppConfig `envPrefix:""`

	// Elengrab application configuration.
	Elengrab ElengrabConfig `envPrefix:"ELENGRAB_"`

	// HTTP server configuration
	HTTPServer HTTPServerConfig `envPrefix:"HTTP_SERVER_"`
	// SQLite storage configuration
	SQLite SQLiteConfig `envPrefix:"SQLITE_"`
	// Redis configuration
	Redis RedisConfig `envPrefix:"REDIS_"`
}

type HTTPServerConfig struct {
	Address  string `env:"ADDRESS" envDefault:""`
	Port     string `env:"PORT" envDefault:"8080"`
	Compress bool   `env:"COMPRESS" envDefault:"true"`
}

type ElengrabConfig struct {
	Mode     string `env:"MODE" envDefault:"public"`
	DemoMode bool   `env:"DEMO_MODE" envDefault:"false"`

	BaseURL         string `env:"BASE_URL" envDefault:""`
	ShortLinkPrefix string `env:"SHORT_LINK_PREFIX" envDefault:"/s"`
	ShortLinkLength uint8  `env:"SHORT_LINK_LENGTH" envDefault:"6"`

	// AdminLogin and AdminPassword are used to create the default administrator
	// at first startup if no admin exists.
	AdminLogin    string `env:"ADMIN_LOGIN" envDefault:""`
	AdminPassword string `env:"ADMIN_PASSWORD" envDefault:""`

	DownloaderBinDir string `env:"DOWNLOADER_BIN_DIR" envDefault:"/usr/local/bin"`

	RootDir      string `env:"ROOT_DIR" envDefault:""`
	AssetsDir    string `env:"ASSETS_DIR" envDefault:"assets"`
	DownloadsDir string `env:"DOWNLOADS_DIR" envDefault:"downloads"`
	// CookiesDir defines the directory where cookies are stored.
	// Default is "cookies".
	CookiesDir string `env:"COOKIES_DIR" envDefault:"cookies"`

	DownloadWorkers       uint32 `env:"DOWNLOAD_WORKERS" envDefault:"3"`
	DeleteDuplicatesScope string `env:"DELETE_DUPLICATES_SCOPE" envDefault:"per_user"`
	// YoutubeAllowCookies allow cookies when downloading YouTube videos.
	// Default is false (disabled).
	YoutubeAllowCookies bool `env:"YOUTUBE_ALLOW_COOKIES" envDefault:"false"`

	Maintenance ElengrabMaintenanceConfig `envPrefix:"MAINTENANCE_"`
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
	DataDir    string `env:"DATA_DIR" envDefault:"sqlite/data"`
	BackupsDir string `env:"BACKUPS_DIR" envDefault:"sqlite/backups"`
}

type RedisConfig struct {
	// SIMPLE(Обычный редис) or SENTINEL
	ConnectionType     string   `env:"CONNECTION_TYPE" envDefault:"simple"`
	Addresses          []string `env:"ADDRESSES" envSeparator:"," envDefault:"localhost:6379"`
	SentinelMasterName string   `env:"SENTINEL_MASTER_NAME" envDefault:"mymaster"`
	PrefixKey          string   `env:"PREFIX_KEY" envDefault:"microservice"`
}

// Создание объекта Config
func New() (*Config, error) {
	parseFlag()

	c := &Config{}

	if err := c.load(); err != nil {
		return nil, err
	}

	if c.Elengrab.RootDir == "" {
		var err error
		c.Elengrab.RootDir, err = defaultRootDir()
		if err != nil {
			return nil, err
		}
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
