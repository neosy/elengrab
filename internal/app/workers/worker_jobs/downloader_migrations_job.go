package wjobs

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewDownloaderMigrationsJob(logger *slog.Logger, runner pworkers.MigrationsRunner) workers.Job {
	return workers.NewJob(
		"DownloaderDeferredMigrations",
		workers.WrapJobExecute(logger, runner.RunDeferredMigrations),
	)
}
