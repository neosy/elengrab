package wjobs

import (
	"log/slog"

	nworkers "github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewDeleteDuplicatesJob(logger *slog.Logger, runner pworkers.DownloadMaintenanceRunner) nworkers.Job {
	return nworkers.NewJob(
		"DeleteDuplicates",
		nworkers.MakeTimedJobExecute(logger, runner.DeleteDuplicates),
	)
}
