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
	"github.com/neosy/elengrab/internal/cli"
	iconfig "github.com/neosy/elengrab/internal/config"
)

func main() {
	// Parse command-line flags before initializing the application.
	flags := cli.ParseFlags()

	// Handle flags that don't require application startup.
	if err := flags.Handle(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to handle command-line flags: %v\n", err)
		os.Exit(1)
	}

	// Exit if a command-line flag requests termination.
	if flags.ShouldExit() {
		return
	}

	// Version output at startup
	fmt.Fprintf(os.Stderr, "%s v%s\n", iconfig.AppName, iconfig.AppVersion)

	// Initialize configuration and logger
	cfg, logger, logWriter, err := bootstrap.Initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize configuration and logger: %v\n", err)
		os.Exit(1)
	}
	defer logWriter.Close()

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

	// Run required migrations before starting the web service.
	// Application startup must be aborted if any required migration fails.
	err = app.RunRequiredMigrations()
	if err != nil {
		logger.Error("failed to run required migrations", "error", err)
		os.Exit(1)
	}

	// Admin server
	infra.StartAdminHTTPServer(app.Context(), cfg.AdminServer)

	// Start background workers
	if err := app.StartBackground(); err != nil {
		logger.Error("Failed to start background jobs", "error", err)
		os.Exit(1)
	}

	// Start HTTP server
	readyHTTPServer := make(chan struct{})
	go infra.StartHTTPServer(logger, cfg, app, readyHTTPServer)

	go func() {
		select {
		case <-readyHTTPServer:
		case <-app.Context().Done():
			// The application is shutting down before the HTTP server is ready.
		}

		if err := logWriter.StopStartupFileLogging(); err != nil {
			logger.Error("Failed to close log file", "error", err)
		}
	}()

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
