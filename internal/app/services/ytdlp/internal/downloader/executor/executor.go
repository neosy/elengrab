package executor

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/dto"
	formatcache "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/format_cache"
)

type Executor struct {
	logger *slog.Logger

	formatCache *formatcache.FormatCache

	// parameters
	ytDlpPath      string
	serviceOptions *dto.Options
}

func NewExecutor(
	logger *slog.Logger,
	formatCache *formatcache.FormatCache,
	ytDlpPath string,
	serviceOptions *dto.Options,
) *Executor {
	return &Executor{
		logger:         logger,
		formatCache:    formatCache,
		ytDlpPath:      ytDlpPath,
		serviceOptions: serviceOptions,
	}
}
