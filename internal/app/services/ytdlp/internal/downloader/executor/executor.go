package executor

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/app/services/ytdlp/dto"
	formatcache "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/format_cache"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type Executor struct {
	logger *slog.Logger

	// Storages
	storage pstorage.DownloadsStorage

	formatCache *formatcache.FormatCache

	// parameters
	ytDlpPath      string
	serviceOptions dto.Options
}

func NewExecutor(
	logger *slog.Logger,
	storage pstorage.DownloadsStorage,
	formatCache *formatcache.FormatCache,
	ytDlpPath string,
	serviceOptions dto.Options,
) *Executor {
	return &Executor{
		logger:         logger,
		storage:        storage,
		formatCache:    formatCache,
		ytDlpPath:      ytDlpPath,
		serviceOptions: serviceOptions,
	}
}
