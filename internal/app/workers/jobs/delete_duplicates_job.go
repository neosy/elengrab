package wjobs

import (
	"context"
	"log/slog"
	"time"

	pworkers "github.com/neosy/elengrab/internal/ports/workers"
	uformat "github.com/neosy/elengrab/pkg/utils/format"
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
