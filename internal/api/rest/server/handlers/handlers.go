package handlers

import (
	htmxh "github.com/neosy/elengrab/internal/api/rest/server/handlers/htmx"
	"github.com/neosy/elengrab/internal/app/usecases"
)

type Dependencies struct {
	Usecases *usecases.Usecases

	// Options
	AssetsDir string
}

type handlers struct {
	HTMX *htmxh.Handlers
}

func New(deps *Dependencies) *handlers {
	return &handlers{
		HTMX: htmxh.NewHTMXHandlers(deps.Usecases, deps.AssetsDir),
	}
}
