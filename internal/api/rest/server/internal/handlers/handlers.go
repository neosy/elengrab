package handlers

import (
	"html/template"
	"log/slog"

	apihandlers "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/api"
	statich "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/static"
	uih "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui"
	"github.com/neosy/elengrab/internal/app/usecases"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/ports/persistence"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type Dependencies struct {
	// Storages
	DownloadsStorage pstorage.DownloadsStorage

	// Caches
	AssetFileCacheRepository persistence.AssetFileCacheRepository

	Usecases  *usecases.Usecases
	Templates *template.Template

	// Options
	AppMode         dtypes.AppMode
	BaseURL         string
	ShortLinkPrefix string
	AssetsDir       string
}

type Handlers struct {
	Static *statich.StaticHandlers
	UI     *uih.UIHandlers
	API    *apihandlers.APIHandlers
}

func New(logger *slog.Logger, deps *Dependencies) *Handlers {
	return &Handlers{
		Static: statich.NewStaticHandlers(
			deps.AssetsDir,
			deps.AssetFileCacheRepository,
			deps.Usecases.Thumbnail,
			deps.Usecases.Downloader,
		),
		UI: uih.NewUIHandlers(
			logger,
			deps.DownloadsStorage,
			deps.Usecases,
			deps.Templates,
			deps.AppMode,
			deps.BaseURL,
			deps.ShortLinkPrefix,
			deps.AssetsDir,
		),
		API: apihandlers.NewAPIHandlers(deps.Usecases),
	}
}
