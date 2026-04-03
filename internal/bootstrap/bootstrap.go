package bootstrap

import (
	"log"
	"log/slog"
	"os"

	iconfig "github.com/neosy/elengrab/internal/config"
	"github.com/neosy/elengrab/internal/pkg/nlogger"
)

func Initialize() (*iconfig.Config, *slog.Logger) {
	// Load application configuration
	cfg, err := iconfig.New()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create a logger with Info level using HandlerOptions
	handlerOptions := &slog.HandlerOptions{
		// Set the logging level
		Level: nlogger.LevelToSlogLevel(cfg.AppConfig.LogLevel),
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))

	log.Printf("Logging level set to '%s'.\n", cfg.AppConfig.LogLevel)

	return cfg, logger
}
