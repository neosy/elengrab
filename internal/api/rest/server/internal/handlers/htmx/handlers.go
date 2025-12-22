package htmxh

import (
	"html/template"

	grabberh "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/htmx/grabber"
	statich "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/htmx/static"
	"github.com/neosy/elengrab/internal/app/usecases"
)

type HTMXHandlers struct {
	Static  *statich.StaticHandlers
	Grabber *grabberh.GrabberHandlers
}

func NewHTMXHandlers(
	usecases *usecases.Usecases,
	templates *template.Template,
	assetsDir string,
	downloadsDir string,
) *HTMXHandlers {
	return &HTMXHandlers{
		Static:  statich.NewStaticHandlers(assetsDir),
		Grabber: grabberh.NewGrabberHandlers(usecases, templates, assetsDir, downloadsDir),
	}
}
