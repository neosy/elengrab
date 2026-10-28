package downloader

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/mappers"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/validators"
	"github.com/neosy/elengrab/internal/app/usecases"
	authweb "github.com/neosy/elengrab/internal/app/usecases/auth_web"
	"github.com/neosy/elengrab/internal/app/usecases/downloader"
	linkweb "github.com/neosy/elengrab/internal/app/usecases/link_web"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail"
	uploadmedia "github.com/neosy/elengrab/internal/app/usecases/upload_media"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/ports/persistence"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type DownloaderHandlers struct {
	logger     *slog.Logger
	mappers    *mappers.Mappers
	validators *validators.Validators

	templates  *assets.Templates
	assets     *assets.Assets
	assetPaths paths.AssetPaths

	// Storages
	downloadsStorage pstorage.DownloadsStorage

	// Caches
	uploadedFileCacheRepo persistence.UploadedFileCacheRepository
	sseEventLogCacheRepo  persistence.SSEEventLogCacheRepository

	// Usecases
	authWeb     authweb.AuthWeb
	downloader  downloader.DownloaderAPI
	linkWeb     linkweb.LinkWeb
	thumbnail   thumbnail.Thumbnail
	uploadFiles uploadmedia.UploadMedia

	// Options
	appMode         dtypes.AppMode
	baseURL         string
	shortLinkPrefix string
	maxUploadSize   int64
}

func NewDownloaderHandlers(
	logger *slog.Logger,

	templates *assets.Templates,
	assets *assets.Assets,
	assetPaths paths.AssetPaths,

	// Storages
	downloadsStorage pstorage.DownloadsStorage,

	// Caches
	uploadedFileCacheRepo persistence.UploadedFileCacheRepository,
	sseEventLogCacheRepo persistence.SSEEventLogCacheRepository,

	usecases *usecases.Usecases,
	appMode dtypes.AppMode,
	baseURL string,
	shortLinkPrefix string,
	maxUploadSize int64,
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

		// Caches
		uploadedFileCacheRepo: uploadedFileCacheRepo,
		sseEventLogCacheRepo:  sseEventLogCacheRepo,

		// Usecases
		authWeb:     usecases.AuthWeb,
		downloader:  usecases.DownloaderAPI,
		linkWeb:     usecases.LinkWeb,
		thumbnail:   usecases.Thumbnail,
		uploadFiles: usecases.UploadMedia,

		// Options
		appMode:         appMode,
		baseURL:         baseURL,
		shortLinkPrefix: shortLinkPrefix,
		maxUploadSize:   maxUploadSize,
	}
}
