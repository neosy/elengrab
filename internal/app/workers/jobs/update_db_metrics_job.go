package wjobs

import (
	"context"
	"log/slog"

	nworkers "github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewUpdateDBMetricsJob(logger *slog.Logger, runner pworkers.DBMMetricsRunner) nworkers.Job {
	run := func(context.Context) error {
		return runner.UpdateMetrics()
	}

	return nworkers.NewJob(
		"UpdateDBMetrics",
		nworkers.MakeTimedJobExecute(logger, run),
	)
}
