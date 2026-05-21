package wjobs

import (
	"log/slog"

	nworkers "github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewDownloaderMigrationsJob(logger *slog.Logger, runner pworkers.MigrationsRunner) nworkers.Job {
	return nworkers.NewJob(
		"DownloaderDeferredMigrations",
		nworkers.MakeTimedJobExecute(logger, runner.RunDeferredMigrations),
	)
}
