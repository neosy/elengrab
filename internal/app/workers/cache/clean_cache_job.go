package cachejobs

import (
	"fmt"
	"log/slog"

	nworkers "github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewCleanCacheJob(logger *slog.Logger, runner pworkers.CacheRunner) nworkers.Job {
	return nworkers.NewJob(
		fmt.Sprintf("CleanCacheExpired[%s]", runner.Name()),
		nworkers.MakeTimedJobExecute(logger, runner.CleanExpired),
	)
}
