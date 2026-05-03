package wjobs

import (
	"context"
	"log/slog"
	"time"

	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type downloaderMigrationsJob struct {
	logger *slog.Logger
	runner pworkers.MigrationsRunner
}

func NewDownloaderMigrationsJob(logger *slog.Logger, runner pworkers.MigrationsRunner) *downloaderMigrationsJob {
	return &downloaderMigrationsJob{
		logger: logger,
		runner: runner,
	}
}

func (j *downloaderMigrationsJob) Execute(ctx context.Context) error {
	startTime := time.Now()
	err := j.runner.ExecuteMigrations(ctx)
	elapsed := time.Since(startTime)

	j.logger.Debug(
		"Job done",
		"name", "DownloaderMigrationsJob",
		"elapsed", uformat.DurationFormat(elapsed),
	)

	return err
}
