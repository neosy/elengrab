package wjobs

import (
	"context"
	"log/slog"

	nworkers "github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewUpdateSystemInfoJob(logger *slog.Logger, runner pworkers.DownloadTaskRunner) nworkers.Job {
	run := func(context.Context) error {
		runner.UpdateSystemInfo()
		return nil
	}

	return nworkers.NewJob(
		"UpdateSystemInfo",
		nworkers.MakeTimedJobExecute(logger, run),
	)
}
