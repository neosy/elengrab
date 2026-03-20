package wjobs

import (
	"context"
	"log/slog"
	"time"

	pworkers "github.com/neosy/elengrab/internal/ports/workers"
	uformat "github.com/neosy/elengrab/pkg/utils/format"
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
