package wjobs

import (
	"context"
	"log/slog"
	"time"

	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type deleteDuplicatesJob struct {
	logger *slog.Logger
	runner pworkers.DownloadMaintenanceRunner
}

func NewDeleteDuplicatesJob(logger *slog.Logger, runner pworkers.DownloadMaintenanceRunner) *deleteDuplicatesJob {
	return &deleteDuplicatesJob{
		logger: logger,
		runner: runner,
	}
}

func (j *deleteDuplicatesJob) Execute(ctx context.Context) error {
	startTime := time.Now()
	err := j.runner.DeleteDuplicates(ctx)
	elapsed := time.Since(startTime)

	j.logger.Debug(
		"Job done",
		"name", "DeleteDuplicates",
		"elapsed", uformat.DurationFormat(elapsed),
	)

	return err
}
