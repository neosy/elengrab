package ui

import (
	"html/template"
	"log/slog"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
	adminhandlers "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	dlhandlers "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader"
	"github.com/neosy/elengrab/internal/app/usecases"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type Handlers struct {
	Admin      *adminhandlers.AdminHandlers
	Downloader *dlhandlers.DownloaderHandlers
}

func NewHandlers(
	logger *slog.Logger,

	// Storages
	downloadsStorage pstorage.DownloadsStorage,

	assets *assets.Assets,

	usecases *usecases.Usecases,
	templates *template.Template,

	// Options
	appMode dtypes.AppMode,
	baseURL string,
	shortLinkPrefix string,
) *Handlers {
	icons.InitDir(assets.FolderPaths().Icons())
	assetPaths := paths.NewAssetPaths(assets)

	return &Handlers{
		Admin: adminhandlers.NewAdminHandlers(
			logger,

			templates,
			assets,
			assetPaths,

			// usecases
			usecases,

			// options
			appMode,
			baseURL,
		),
		Downloader: dlhandlers.NewDownloaderHandlers(
			logger,

			templates,
			assets,
			assetPaths,

			// Storages
			downloadsStorage,

			// Usecases
			usecases,

			// Options
			appMode,
			baseURL,
			shortLinkPrefix,
		),
	}
}
