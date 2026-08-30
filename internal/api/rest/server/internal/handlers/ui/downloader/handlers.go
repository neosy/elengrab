package downloader

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/mappers"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/validators"
	httptemplates "github.com/neosy/elengrab/internal/api/rest/server/templates"
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

	templates  *httptemplates.Templates
	assets     *assets.Assets
	assetPaths paths.AssetPaths

	// Storages
	downloadsStorage pstorage.DownloadsStorage

	// Usecases
	authWeb    authweb.AuthWeb
	downloader downloader.DownloaderAPI
	linkWeb    linkweb.LinkWeb
	thumbnail  thumbnail.Thumbnail

	// Options
	appMode         dtypes.AppMode
	baseURL         string
	shortLinkPrefix string
}

func NewDownloaderHandlers(
	logger *slog.Logger,

	templates *httptemplates.Templates,
	assets *assets.Assets,
	assetPaths paths.AssetPaths,

	// Storages
	downloadsStorage pstorage.DownloadsStorage,

	usecases *usecases.Usecases,
	appMode dtypes.AppMode,
	baseURL string,
	shortLinkPrefix string,
) *DownloaderHandlers {
	return &DownloaderHandlers{
		logger:     logger,
		mappers:    mappers.NewMappers(),
		validators: validators.NewValidators(),

		templates:  templates,
		assets:     assets,
		assetPaths: assetPaths,

		// Storages
		downloadsStorage: downloadsStorage,

		// Usecases
		authWeb:    usecases.AuthWeb,
		downloader: usecases.DownloaderAPI,
		linkWeb:    usecases.LinkWeb,
		thumbnail:  usecases.Thumbnail,

		// Options
		appMode:         appMode,
		baseURL:         baseURL,
		shortLinkPrefix: shortLinkPrefix,
	}
}
