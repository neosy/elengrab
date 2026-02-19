package wjobs

import (
	"context"
	"log/slog"
	"time"

	pworkers "github.com/neosy/elengrab/internal/ports/workers"
	uformat "github.com/neosy/elengrab/pkg/utils/format"
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
