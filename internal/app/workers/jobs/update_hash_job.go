package wjobs

import (
	"log/slog"

	nworkers "github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewUpdateHashJob(logger *slog.Logger, runner pworkers.DownloadMaintenanceRunner) nworkers.Job {
	return nworkers.NewJob(
		"UpdateHash",
		nworkers.MakeTimedJobExecute(logger, runner.UpdateHash),
	)
}
