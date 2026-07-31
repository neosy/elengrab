package cachejobs

import (
	"fmt"
	"log/slog"

	"github.com/neosy/elengrab/internal/pkg/workers"
	pworkers "github.com/neosy/elengrab/internal/ports/workers"
)

func NewCleanCacheJob(logger *slog.Logger, runner pworkers.CacheRunner) workers.Job {
	return workers.NewJob(
		fmt.Sprintf("CleanCacheExpired[%s]", runner.Name()),
		workers.WrapJobExecute(logger, runner.CleanExpired),
	)
}
