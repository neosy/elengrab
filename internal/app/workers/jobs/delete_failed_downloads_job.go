package wjobs

import (
	"log/slog"

	nworkers "github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewDeleteFailedDownloadsJob(logger *slog.Logger, runner pworkers.DownloadMaintenanceRunner) nworkers.Job {
	return nworkers.NewJob(
		"DeleteFailedDownloads",
		nworkers.MakeTimedJobExecute(logger, runner.DeleteFailedDownloads),
	)
}
