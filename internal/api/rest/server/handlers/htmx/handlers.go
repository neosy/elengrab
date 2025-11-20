package htmxh

import (
	grabberh "github.com/neosy/elengrab/internal/api/rest/server/handlers/htmx/grabber"
	statich "github.com/neosy/elengrab/internal/api/rest/server/handlers/htmx/static"
	"github.com/neosy/elengrab/internal/app/usecases"
)

type Handlers struct {
	Static  *statich.StaticHandlers
	Grabber *grabberh.GrabberHandlers
}

func NewHTMXHandlers(usecases *usecases.Usecases, assetsDir string) *Handlers {
	return &Handlers{
		Static:  statich.NewStaticHandlers(assetsDir),
		Grabber: grabberh.NewGrabberHandlers(usecases, assetsDir),
	}
}
