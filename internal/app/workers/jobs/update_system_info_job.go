package wjobs

import (
	"context"
	"log/slog"

	"github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewUpdateSystemInfoJob(logger *slog.Logger, runner pworkers.DownloadTaskRunner) workers.Job {
	run := func(context.Context) error {
		runner.UpdateSystemInfo()
		return nil
	}

	return workers.NewJob(
		"UpdateSystemInfo",
		workers.WrapJobExecute(logger, run),
	)
}
