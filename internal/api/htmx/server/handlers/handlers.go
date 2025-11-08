package handlers

import (
	grabh "github.com/neosy/elengrab/internal/api/htmx/server/handlers/grab"
	indexh "github.com/neosy/elengrab/internal/api/htmx/server/handlers/index"
	statich "github.com/neosy/elengrab/internal/api/htmx/server/handlers/static"
	"github.com/neosy/elengrab/internal/app/usecases"
)

type Dependencies struct {
	Usecases *usecases.Usecases

	// Options
	AssetsDir string
}

type handlers struct {
	Index  *indexh.IndexHandlers
	Static *statich.StaticHandlers
	Grab   *grabh.GrabHandlers
}

func New(deps *Dependencies) *handlers {
	return &handlers{
		Index:  indexh.NewIndexHandlers(deps.AssetsDir),
		Static: statich.NewStaticHandlers(deps.AssetsDir),
		Grab:   grabh.NewGrabHandlers(deps.Usecases, deps.AssetsDir),
	}
}
