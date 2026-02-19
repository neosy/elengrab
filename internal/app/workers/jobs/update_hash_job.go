package wjobs

import (
	"context"
	"log/slog"
	"time"

	pworkers "github.com/neosy/elengrab/internal/ports/workers"
	uformat "github.com/neosy/elengrab/pkg/utils/format"
)

type updateHashJob struct {
	logger *slog.Logger
	runner pworkers.DownloadMaintenanceRunner
}

func NewUpdateHashJob(logger *slog.Logger, runner pworkers.DownloadMaintenanceRunner) *updateHashJob {
	return &updateHashJob{
		logger: logger,
		runner: runner,
	}
}

func (j *updateHashJob) Execute(ctx context.Context) error {
	startTime := time.Now()
	err := j.runner.UpdateHash(ctx)
	elapsed := time.Since(startTime)

	j.logger.Debug(
		"Job done",
		"name", "UpdateHash",
		"elapsed", uformat.DurationFormat(elapsed),
	)

	return err
}
