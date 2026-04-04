package uih

import (
	"html/template"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers"
	"github.com/neosy/elengrab/internal/app/usecases"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type UIHandlers struct {
	Downloader *handlers.DownloaderHandlers
}

func NewUIHandlers(
	usecases *usecases.Usecases,
	templates *template.Template,

	// options
	appMode dtypes.AppMode,
	baseURL string,
	shortLinkPrefix string,
	assetsDir string,
	downloadsDir string,
) *UIHandlers {
	return &UIHandlers{
		Downloader: handlers.NewDownloaderHandlers(
			templates,
			usecases,

			// options
			appMode,
			baseURL,
			shortLinkPrefix,
			assetsDir,
			downloadsDir,
		),
	}
}
