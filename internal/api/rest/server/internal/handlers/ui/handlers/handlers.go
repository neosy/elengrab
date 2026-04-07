package handlers

import (
	"html/template"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/mappers"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/validators"
	"github.com/neosy/elengrab/internal/app/usecases"
	authweb "github.com/neosy/elengrab/internal/app/usecases/auth_web"
	"github.com/neosy/elengrab/internal/app/usecases/downloader"
	linkweb "github.com/neosy/elengrab/internal/app/usecases/link_web"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type DownloaderHandlers struct {
	mappers    *mappers.Mappers
	validators *validators.Validators

	templates    *template.Template
	assetFolders assets.FolderPaths

	// usecases
	authWeb    *authweb.AuthWeb
	downloader *downloader.Downloader
	linkWeb    *linkweb.LinkWeb

	// Options
	appMode         dtypes.AppMode
	baseURL         string
	shortLinkPrefix string
	assetsDir       string
	downloadsDir    string
}

func NewDownloaderHandlers(
	templates *template.Template,
	usecases *usecases.Usecases,
	appMode dtypes.AppMode,
	baseURL string,
	shortLinkPrefix string,
	assetsDir string,
	downloadsDir string,
) *DownloaderHandlers {
	return &DownloaderHandlers{
		mappers:    mappers.NewMappers(),
		validators: validators.NewValidators(),

		templates:    templates,
		assetFolders: assets.NewFolderPaths(assetsDir),

		// usecases
		authWeb:    usecases.AuthWeb,
		downloader: usecases.Downloader,
		linkWeb:    usecases.LinkWeb,

		// Options
		appMode:         appMode,
		baseURL:         baseURL,
		shortLinkPrefix: shortLinkPrefix,
		assetsDir:       assetsDir,
		downloadsDir:    downloadsDir,
	}
}
