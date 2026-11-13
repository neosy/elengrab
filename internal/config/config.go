package iconfig

import (
	"fmt"
	"log"

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

	startupInfo = newRuntimeInfo(c)

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
