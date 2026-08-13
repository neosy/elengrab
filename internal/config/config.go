package iconfig

import (
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/neosy/elengrab/internal/pkg/config"
)

// Config contains the application configuration.
//
//go:generate envdoc -dir=../.. -files='**/*.go' -types='Config' -output ../../environments.md
type Config struct {
	// Common application settings.
	AppConfig config.AppConfig `envPrefix:"ELENGRAB_"`

	// Admin server settings.
	AdminServer AdminServerConfig `envPrefix:"ELENGRAB_ADMIN_SERVER_"`

	// Elengrab application settings.
	Elengrab ElengrabConfig `envPrefix:"ELENGRAB_"`

	// HTTP server settings.
	HTTPServer HTTPServerConfig `envPrefix:"ELENGRAB_HTTP_SERVER_"`

	// SQLite storage settings.
	SQLite SQLiteConfig `envPrefix:"ELENGRAB_SQLITE_"`

	// Redis settings.
	// Redis RedisConfig `envPrefix:"ELENGRAB_REDIS_"`
}

// AdminServerConfig contains the configuration for the admin server.
type AdminServerConfig struct {
	// Enables the admin server.
	Enable bool `env:"ENABLE" envDefault:"false"`

	// Address on which the admin server listens.
	Address string `env:"ADDRESS" envDefault:"127.0.0.1"`

	// Port on which the admin server listens.
	Port string `env:"PORT" envDefault:"6060"`

	// Debug server configuration.
	DebugConfig AdminServerDebugConfig `envPrefix:"DEBUG_"`
}

// AdminServerDebugConfig contains configuration for admin server endpoints.
type AdminServerDebugConfig struct {
	// Enables the pprof debugging endpoints.
	EnablePprof bool `env:"PPROF" envDefault:"true"`

	// Enables the metrics endpoint.
	EnableMetrics bool `env:"METRICS" envDefault:"true"`

	// Enables the health check endpoint.
	EnableHealth bool `env:"HEALTH" envDefault:"true"`
}

// HTTPServerConfig contains the HTTP server configuration.
type HTTPServerConfig struct {
	// Address specifies the network address on which the HTTP server listens.
	Address string `env:"ADDRESS" envDefault:""`

	// Port specifies the port on which the HTTP server listens.
	Port string `env:"PORT" envDefault:"8080"`

	// Compress enables HTTP response compression.
	Compress bool `env:"COMPRESS" envDefault:"true"`
}

// ElengrabConfig contains the Elengrab application configuration.
type ElengrabConfig struct {
	// Application operating mode.
	Mode string `env:"MODE" envDefault:"public"`
	// Enables demo mode with restricted functionality.
	DemoMode bool `env:"DEMO_MODE" envDefault:"false"`

	// Base URL used to build absolute application links.
	BaseURL string `env:"BASE_URL" envDefault:""`

	// Prefix used for generated short links.
	ShortLinkPrefix string `env:"SHORT_LINK_PREFIX" envDefault:"/s"`
	// Length of generated short link identifiers.
	ShortLinkLength uint8 `env:"SHORT_LINK_LENGTH" envDefault:"6"`
	// Lifetime of generated short links in days.
	ShortLinkTTLDays uint16 `env:"SHORT_LINK_TTL_DAYS" envDefault:"180"`

	// AdminLogin is used to create the default administrator at first startup
	AdminLogin string `env:"ADMIN_LOGIN" envDefault:""`
	// AdminPassword is used to create the default administrator at first startup
	AdminPassword string `env:"ADMIN_PASSWORD" envDefault:""`

	// Directory containing downloader binaries such as yt-dlp.
	DownloaderBinDir string `env:"DOWNLOADER_BIN_DIR" envDefault:"/usr/local/bin"`

	// Root directory for application data.
	RootDir string `env:"ROOT_DIR" envDefault:""`
	// Directory containing application assets.
	AssetsDir string `env:"ASSETS_DIR" envDefault:"assets"`
	// Directory containing media files (e.g. thumbnails).
	MediaDir string `env:"MEDIA_DIR" envDefault:"media"`
	// Directory containing downloaded files.
	DownloadsDir string `env:"DOWNLOADS_DIR" envDefault:"downloads"`
	// Directory containing cookies used by the downloader.
	CookiesDir string `env:"COOKIES_DIR" envDefault:"cookies"`

	// Maximum number of concurrent media downloads.
	DownloadWorkers uint32 `env:"DOWNLOAD_WORKERS" envDefault:"3"`
	// Maximum number of concurrent background operations.
	OperationWorkers uint32 `env:"OPERATION_WORKERS" envDefault:"5"`
	// Scope used when checking media uniqueness before deleting duplicates.
	DeleteDuplicatesUniquenessScope string `env:"DELETE_DUPLICATES_UNIQUENESS_SCOPE" envDefault:"per_user"`

	// Enables the use of cookies when downloading media.
	AllowCookies bool `env:"ALLOW_COOKIES" envDefault:"false"`

	// Maintenance task configuration.
	Maintenance ElengrabMaintenanceConfig `envPrefix:"MAINTENANCE_"`
}

// ElengrabMaintenanceConfig contains configuration for periodic maintenance tasks.
type ElengrabMaintenanceConfig struct {
	// Interval for calculating and updating hashes of downloaded files.
	// File hashes are used to identify duplicate files.
	UpdateHashInterval time.Duration `env:"UPDATE_HASH_INTERVAL" envDefault:"8h"`

	// Interval for deleting duplicate media files.
	DeleteDuplicatesInterval time.Duration `env:"DELETE_DUPLICATES_INTERVAL" envDefault:"1h"`

	// Interval for deleting database records whose downloaded files are missing.
	DeleteMissingDownloadsInterval time.Duration `env:"DELETE_MISSING_DOWNLOADS_INTERVAL" envDefault:"12h"`

	// Interval for deleting failed download records.
	DeleteFailedDownloadsInterval time.Duration `env:"DELETE_FAILED_DOWNLOADS_INTERVAL" envDefault:"1h"`

	// Enables moving unmatched files during periodic maintenance.
	MoveUnmatchedFilesEnabled bool `env:"ENABLE_MOVE_UNMATCHED_FILES" envDefault:"false"`

	// Number of the latest database backups to keep.
	// If set to 0, old backup files are not deleted.
	DatabaseBackupsKeep int `env:"DATABASE_BACKUPS_KEEP" envDefault:"7"`
}

// SQLiteConfig contains the SQLite storage configuration.
type SQLiteConfig struct {
	// Directory where SQLite database files are stored.
	DataDir string `env:"DATA_DIR" envDefault:"sqlite/data"`

	// Directory where SQLite database backups are stored.
	BackupsDir string `env:"BACKUPS_DIR" envDefault:"sqlite/backups"`
}

// RedisConfig contains the Redis connection configuration.
type RedisConfig struct {
	// Redis connection type: "simple" for a direct connection or "sentinel" for Redis Sentinel.
	ConnectionType string `env:"CONNECTION_TYPE" envDefault:"simple"`

	// Redis server addresses. Multiple addresses are separated by commas.
	Addresses []string `env:"ADDRESSES" envSeparator:"," envDefault:"localhost:6379"`

	// Name of the Redis Sentinel master to monitor.
	SentinelMasterName string `env:"SENTINEL_MASTER_NAME" envDefault:"mymaster"`

	// Prefix added to Redis keys to namespace application data.
	PrefixKey string `env:"PREFIX_KEY" envDefault:"microservice"`
}

// New loads the application configuration and initializes default values.
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

	startupInfo = newStartupInfo(c)

	return c, nil
}

// load loads configuration from the .env file and environment variables.
func (config *Config) load() error {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables only")
	}

	if err := env.Parse(config); err != nil {
		return fmt.Errorf("config load(). Read configuration error: %w", err)
	}

	return nil
}
