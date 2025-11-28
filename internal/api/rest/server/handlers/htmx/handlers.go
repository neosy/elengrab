package htmxh

import (
	"html/template"

	grabberh "github.com/neosy/elengrab/internal/api/rest/server/handlers/htmx/grabber"
	statich "github.com/neosy/elengrab/internal/api/rest/server/handlers/htmx/static"
	"github.com/neosy/elengrab/internal/app/usecases"
)

type Handlers struct {
	Static  *statich.StaticHandlers
	Grabber *grabberh.GrabberHandlers
}

func NewHTMXHandlers(usecases *usecases.Usecases, templates *template.Template, assetsDir string) *Handlers {
	return &Handlers{
		Static:  statich.NewStaticHandlers(assetsDir),
		Grabber: grabberh.NewGrabberHandlers(usecases, templates, assetsDir),
	}
}
