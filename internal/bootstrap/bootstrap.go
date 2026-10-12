package bootstrap

import (
	"fmt"
	"log"
	"log/slog"

	iconfig "github.com/neosy/elengrab/internal/config"
	nlogger "github.com/neosy/elengrab/internal/pkg/logger"
)

// Initialize loads the application configuration and initializes the logger.
func Initialize() (*iconfig.Config, *slog.Logger, *LogWriter, error) {
	// Load application configuration
	cfg, err := iconfig.New()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("load configuration: %w", err)
	}

	writer, err := NewLogWriter(cfg.Elengrab.RootDir, cfg.Elengrab.Logging)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("initialize log writer: %w", err)
	}

	// Create a logger with Info level using HandlerOptions
	handlerOptions := &slog.HandlerOptions{
		Level: nlogger.LevelToSlogLevel(cfg.AppConfig.LogLevel),
	}

	logger := slog.New(slog.NewTextHandler(writer, handlerOptions))

	log.Printf("Logging level set to '%s'.\n", cfg.AppConfig.LogLevel)

	return cfg, logger, writer, nil
}
