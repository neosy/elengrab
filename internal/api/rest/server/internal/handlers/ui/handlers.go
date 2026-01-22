package uih

import (
	"html/template"

	downloaderh "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader"
	"github.com/neosy/elengrab/internal/app/usecases"
)

type UIHandlers struct {
	Downloader *downloaderh.DownloaderHandlers
}

func NewUIHandlers(
	usecases *usecases.Usecases,
	templates *template.Template,
	assetsDir string,
	downloadsDir string,
) *UIHandlers {
	return &UIHandlers{
		Downloader: downloaderh.NewDownloaderHandlers(usecases, templates, assetsDir, downloadsDir),
	}
}
