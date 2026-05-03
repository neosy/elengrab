package uih

import (
	"html/template"
	"log/slog"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers"
	"github.com/neosy/elengrab/internal/app/usecases"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type UIHandlers struct {
	Downloader *handlers.DownloaderHandlers
}

func NewUIHandlers(
	logger *slog.Logger,

	// Storages
	downloadsStorage pstorage.DownloadsStorage,

	usecases *usecases.Usecases,
	templates *template.Template,

	// Options
	appMode dtypes.AppMode,
	baseURL string,
	shortLinkPrefix string,
	assetsDir string,
) *UIHandlers {
	return &UIHandlers{
		Downloader: handlers.NewDownloaderHandlers(
			logger,
			templates,

			// Storages
			downloadsStorage,

			usecases,

			// Options
			appMode,
			baseURL,
			shortLinkPrefix,
			assetsDir,
		),
	}
}
