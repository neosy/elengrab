package uih

import (
	"html/template"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers"
	"github.com/neosy/elengrab/internal/app/usecases"
)

type UIHandlers struct {
	Downloader *handlers.DownloaderHandlers
}

func NewUIHandlers(
	usecases *usecases.Usecases,
	templates *template.Template,
	assetsDir string,
	downloadsDir string,
) *UIHandlers {
	return &UIHandlers{
		Downloader: handlers.NewDownloaderHandlers(usecases, templates, assetsDir, downloadsDir),
	}
}
