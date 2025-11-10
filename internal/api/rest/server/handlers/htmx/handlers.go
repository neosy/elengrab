package htmxh

import (
	grabberh "github.com/neosy/elengrab/internal/api/rest/server/handlers/htmx/grabber"
	indexh "github.com/neosy/elengrab/internal/api/rest/server/handlers/htmx/index"
	statich "github.com/neosy/elengrab/internal/api/rest/server/handlers/htmx/static"
	"github.com/neosy/elengrab/internal/app/usecases"
)

type Handlers struct {
	Index   *indexh.IndexHandlers
	Static  *statich.StaticHandlers
	Grabber *grabberh.GrabberHandlers
}

func NewHTMXHandlers(usecases *usecases.Usecases, assetsDir string) *Handlers {
	return &Handlers{
		Index:   indexh.NewIndexHandlers(assetsDir),
		Static:  statich.NewStaticHandlers(assetsDir),
		Grabber: grabberh.NewGrabberHandlers(usecases, assetsDir),
	}
}
