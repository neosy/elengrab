package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	iconfig "github.com/neosy/elengrab/infrastructure/config"
	httpsrv "github.com/neosy/elengrab/internal/api/rest/server"
	"github.com/neosy/elengrab/internal/app/usecases"
	"github.com/neosy/elengrab/internal/services"
	"github.com/neosy/elengrab/pkg/nlogger"
)

func main() {
	var err error

	cfg := iconfig.New()

	ctx, cancel := context.WithCancel(context.Background())

	// Создаем обработчик с уровнем Info, используя HandlerOptions
	handlerOptions := &slog.HandlerOptions{
		// Устанавливаем уровень логирования
		Level: nlogger.LevelToSlogLevel(cfg.AppConfig.LogLevel),
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))

	log.Printf("Logging level set to '%s'.\n", cfg.AppConfig.LogLevel)

	// Захват сигналов завершения (Ctrl+C, SIGTERM)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Services
	srvDeps := &services.Dependencies{
		BinDir:       cfg.Elengrab.BinDir,
		DownloadsDir: cfg.Elengrab.DownloadsDir,
	}
	services, err := services.New(logger, srvDeps)
	if err != nil {
		logger.Error("Failed to create services", "err", err)
		return

	}

	// Usecases
	ucDeps := &usecases.Dependencies{
		Services:     services,
		DownloadsDir: cfg.Elengrab.DownloadsDir,
	}
	uc := usecases.NewUsecases(logger, ucDeps)

	// FastHTTP server
	go func(ctx context.Context) {
		deps := &httpsrv.Dependencies{
			Usecases:  uc,
			AssetsDir: cfg.Elengrab.AssetsDir,
		}

		httpServer := httpsrv.NewServer(logger, deps)
		err = httpServer.ListenAndServe(ctx, cfg.HTMXServer.Port)
		if err != nil {
			cancel()
		}
	}(ctx)

	// Ждем сигнал завершения или отмены контекста
	select {
	case <-ctx.Done():
		logger.ErrorContext(ctx, "Context complete, shutting down services...")
	case sig := <-sigChan:
		logger.ErrorContext(ctx, fmt.Sprintf("Signal received: %v, shutting down...", sig))
		cancel()
	}

	if err != nil {
		log.Print(err)
	}
}
