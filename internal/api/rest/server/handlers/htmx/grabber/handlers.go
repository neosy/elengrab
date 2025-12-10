package grabberh

import (
	"html/template"

	"github.com/neosy/elengrab/internal/api/rest/server/handlers/htmx/grabber/mappers"
	"github.com/neosy/elengrab/internal/api/rest/server/handlers/htmx/grabber/validators"
	"github.com/neosy/elengrab/internal/app/usecases"
)

type GrabberHandlers struct {
	mappers    *mappers.Mappers
	validators *validators.Validators
	usecases   *usecases.Usecases

	templates *template.Template

	// Options
	assetsDir string
}

func NewGrabberHandlers(usecases *usecases.Usecases, templates *template.Template, assetsDir string) *GrabberHandlers {
	return &GrabberHandlers{
		mappers:    mappers.NewMappers(),
		validators: validators.NewValidators(),
		usecases:   usecases,
		templates:  templates,
		assetsDir:  assetsDir,
	}
}
