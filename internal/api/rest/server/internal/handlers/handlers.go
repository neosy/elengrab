package handlers

import (
	"html/template"
	"log/slog"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/api"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/static"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui"
	"github.com/neosy/elengrab/internal/app/usecases"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	appenv "github.com/neosy/elengrab/internal/pkg/config/app_env"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type Dependencies struct {
	// Storages
	DownloadsStorage pstorage.DownloadsStorage

	// Assets
	Assets *assets.Assets

	Usecases  *usecases.Usecases
	Templates *template.Template

	// Options
	AppEnv          appenv.AppEnv
	AppMode         dtypes.AppMode
	BaseURL         string
	ShortLinkPrefix string
}

type Handlers struct {
	Static *static.StaticHandlers
	UI     *ui.Handlers
	API    *api.Handlers
}

func New(logger *slog.Logger, deps *Dependencies) *Handlers {
	return &Handlers{
		Static: static.NewStaticHandlers(
			deps.Assets,
			deps.Usecases.Thumbnail,
			deps.Usecases.Downloader,
			deps.AppEnv,
		),
		UI: ui.NewHandlers(
			logger,
			deps.DownloadsStorage,
			deps.Assets,
			deps.Usecases,
			deps.Templates,
			deps.AppMode,
			deps.BaseURL,
			deps.ShortLinkPrefix,
		),
		API: api.NewHandlers(deps.Usecases),
	}
}
