package handlers

import (
	"html/template"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/mappers"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/validators"
	"github.com/neosy/elengrab/internal/app/usecases"
	"github.com/neosy/elengrab/internal/app/usecases/authweb"
	"github.com/neosy/elengrab/internal/app/usecases/downloader"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type DownloaderHandlers struct {
	mappers    *mappers.Mappers
	validators *validators.Validators

	templates *template.Template

	// usecases
	authWeb    *authweb.AuthWeb
	downloader *downloader.Downloader

	// Options
	appMode      dtypes.AppMode
	assetsDir    string
	downloadsDir string
}

func NewDownloaderHandlers(
	templates *template.Template,
	usecases *usecases.Usecases,
	appMode dtypes.AppMode,
	assetsDir string,
	downloadsDir string,
) *DownloaderHandlers {
	return &DownloaderHandlers{
		mappers:    mappers.NewMappers(),
		validators: validators.NewValidators(),
		templates:  templates,

		// usecases
		authWeb:    usecases.AuthWeb,
		downloader: usecases.Downloader,

		// Options
		appMode:      appMode,
		assetsDir:    assetsDir,
		downloadsDir: downloadsDir,
	}
}
