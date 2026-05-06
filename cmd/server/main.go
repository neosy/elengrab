package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"

	"github.com/neosy/elengrab/cmd/server/infra"
	"github.com/neosy/elengrab/internal/app"
	"github.com/neosy/elengrab/internal/bootstrap"
	iconfig "github.com/neosy/elengrab/internal/config"
)

func main() {
	// Version output at startup
	fmt.Fprintf(os.Stderr, "%s v%s\n", iconfig.AppName, iconfig.AppVersion)

	// Initialize configuration and logger
	cfg, logger := bootstrap.Initialize()

	// Ensure assets directory exists
	if err := infra.SetupAssets(cfg); err != nil {
		logger.Error("Failed to setup assets directory", "error", err)
		os.Exit(1)
	}

	// Initialize core application
	app, err := app.NewApplication(logger, cfg)
	if err != nil {
		logger.Error("Failed to initialize application", "error", err)
		os.Exit(1)
	}
	defer app.Shutdown()

	// Admin server
	infra.StartAdminHTTPServer(app.Context(), cfg.AdminServer)

	// Start background workers
	if err := app.StartBackground(); err != nil {
		logger.Error("Failed to start background jobs", "error", err)
		os.Exit(1)
	}

	// Start HTTP server
	go infra.StartHTTPServer(logger, cfg, app)

	// Wait for shutdown signal
	waitForShutdown(app)
}

func waitForShutdown(app *app.Application) {
	// Capture termination signals (Ctrl+C, SIGTERM)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for context cancellation or termination signal
	select {
	case <-app.Context().Done():
		app.Logger().ErrorContext(app.Context(), "Context completed, shutting down services...")
	case sig := <-sigChan:
		app.Logger().ErrorContext(app.Context(), fmt.Sprintf("Termination signal received: %v, shutting down...", sig))
		app.Cancel()
	}
}
