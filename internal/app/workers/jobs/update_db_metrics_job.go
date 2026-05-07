package wjobs

import (
	"context"
	"log/slog"
	"time"

	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

type updateDBMetricsJob struct {
	logger *slog.Logger
	runner pworkers.DBMMetricsRunner
}

func NewUpdateDBMetricsJob(logger *slog.Logger, runner pworkers.DBMMetricsRunner) *updateDBMetricsJob {
	return &updateDBMetricsJob{
		logger: logger,
		runner: runner,
	}
}

func (j *updateDBMetricsJob) Execute(ctx context.Context) error {
	startTime := time.Now()
	err := j.runner.UpdateMetrics()
	elapsed := time.Since(startTime)

	j.logger.Debug(
		"Job done",
		"name", "UpdateDBMetrics",
		"elapsed", uformat.DurationFormat(elapsed),
	)

	return err
}
