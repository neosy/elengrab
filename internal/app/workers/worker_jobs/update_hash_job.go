package wjobs

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewUpdateHashJob(logger *slog.Logger, runner pworkers.DownloadMaintenanceRunner) workers.Job {
	return workers.NewJob(
		"UpdateHash",
		workers.WrapJobExecute(logger, runner.UpdateHash),
	)
}
