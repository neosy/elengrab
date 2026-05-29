package uih

import (
	"html/template"
	"log/slog"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
	adminhandlers "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin_handlers"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/icons"
	dlhandlers "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader_handlers"
	"github.com/neosy/elengrab/internal/app/usecases"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type UIHandlers struct {
	Admin      *adminhandlers.AdminHandlers
	Downloader *dlhandlers.DownloaderHandlers
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
	assetFolders := assets.NewFolderPaths(assetsDir)

	icons.InitDir(assetFolders.Icons())

	return &UIHandlers{
		Admin: adminhandlers.NewAdminHandlers(
			logger,

			templates,
			assetFolders,

			// usecases
			usecases,

			// options
			appMode,
			baseURL,
		),
		Downloader: dlhandlers.NewDownloaderHandlers(
			logger,

			templates,
			assetFolders,

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
