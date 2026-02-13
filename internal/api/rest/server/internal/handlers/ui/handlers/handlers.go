package handlers

import (
	"html/template"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/mappers"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/validators"
	"github.com/neosy/elengrab/internal/app/usecases"
)

type DownloaderHandlers struct {
	mappers    *mappers.Mappers
	validators *validators.Validators
	usecases   *usecases.Usecases

	templates *template.Template

	// Options
	assetsDir    string
	downloadsDir string
}

func NewDownloaderHandlers(
	usecases *usecases.Usecases,
	templates *template.Template,
	assetsDir string,
	downloadsDir string,
) *DownloaderHandlers {
	return &DownloaderHandlers{
		mappers:    mappers.NewMappers(),
		validators: validators.NewValidators(),
		usecases:   usecases,
		templates:  templates,

		// Options
		assetsDir:    assetsDir,
		downloadsDir: downloadsDir,
	}
}
