package wjobs

import (
	"context"
	"log/slog"

	"github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewUpdateDBMetricsJob(logger *slog.Logger, runner pworkers.DBMMetricsRunner) workers.Job {
	run := func(context.Context) error {
		return runner.UpdateMetrics()
	}

	return workers.NewJob(
		"UpdateDBMetrics",
		workers.WrapJobExecute(logger, run),
	)
}
