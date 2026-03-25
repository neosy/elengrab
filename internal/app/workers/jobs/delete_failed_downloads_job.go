package wjobs

import (
	"context"
	"log/slog"
	"time"

	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type deleteFailedDownloadsJob struct {
	logger *slog.Logger
	runner pworkers.DownloadMaintenanceRunner
}

func NewDeleteFailedDownloadsJob(logger *slog.Logger, runner pworkers.DownloadMaintenanceRunner) *deleteFailedDownloadsJob {
	return &deleteFailedDownloadsJob{
		logger: logger,
		runner: runner,
	}
}

func (j *deleteFailedDownloadsJob) Execute(ctx context.Context) error {
	startTime := time.Now()
	err := j.runner.DeleteFailedDownloads(ctx)
	elapsed := time.Since(startTime)

	j.logger.Debug(
		"Job done",
		"name", "DeleteFailedDownloads",
		"elapsed", uformat.DurationFormat(elapsed),
	)

	return err
}
