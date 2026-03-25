package wjobs

import (
	"context"
	"log/slog"
	"time"

	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type updateSystemInfoJob struct {
	logger *slog.Logger
	runner pworkers.DownloadTaskRunner
}

func NewUpdateSystemInfoJob(logger *slog.Logger, runner pworkers.DownloadTaskRunner) *updateSystemInfoJob {
	return &updateSystemInfoJob{
		logger: logger,
		runner: runner,
	}
}

func (j *updateSystemInfoJob) Execute(ctx context.Context) error {
	startTime := time.Now()
	j.runner.UpdateSystemInfo()
	elapsed := time.Since(startTime)

	j.logger.Debug(
		"Job done",
		"name", "UpdateSystemInfo",
		"elapsed", uformat.DurationFormat(elapsed),
	)

	return nil
}
