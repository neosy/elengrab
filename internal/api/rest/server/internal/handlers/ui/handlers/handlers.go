package handlers

import (
	"html/template"
	"log/slog"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/icons"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/mappers"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/validators"
	"github.com/neosy/elengrab/internal/app/usecases"
	authweb "github.com/neosy/elengrab/internal/app/usecases/auth_web"
	"github.com/neosy/elengrab/internal/app/usecases/downloader"
	linkweb "github.com/neosy/elengrab/internal/app/usecases/link_web"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type DownloaderHandlers struct {
	logger     *slog.Logger
	mappers    *mappers.Mappers
	validators *validators.Validators

	templates    *template.Template
	assetFolders assets.FolderPaths

	// Storages
	downloadsStorage pstorage.DownloadsStorage

	// Usecases
	authWeb    *authweb.AuthWeb
	downloader *downloader.Downloader
	linkWeb    *linkweb.LinkWeb
	thumbnail  *thumbnail.Thumbnail

	// Options
	appMode         dtypes.AppMode
	baseURL         string
	shortLinkPrefix string
}

func NewDownloaderHandlers(
	logger *slog.Logger,
	templates *template.Template,

	// Storages
	downloadsStorage pstorage.DownloadsStorage,

	usecases *usecases.Usecases,
	appMode dtypes.AppMode,
	baseURL string,
	shortLinkPrefix string,
	assetsDir string,
) *DownloaderHandlers {
	assetFolders := assets.NewFolderPaths(assetsDir)

	icons.Init(assetFolders.Icons())

	return &DownloaderHandlers{
		logger:     logger,
		mappers:    mappers.NewMappers(),
		validators: validators.NewValidators(),

		templates:    templates,
		assetFolders: assetFolders,

		// Storages
		downloadsStorage: downloadsStorage,

		// Usecases
		authWeb:    usecases.AuthWeb,
		downloader: usecases.Downloader,
		linkWeb:    usecases.LinkWeb,
		thumbnail:  usecases.Thumbnail,

		// Options
		appMode:         appMode,
		baseURL:         baseURL,
		shortLinkPrefix: shortLinkPrefix,
	}
}
