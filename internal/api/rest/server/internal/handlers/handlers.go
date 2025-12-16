package handlers

import (
	"html/template"

	apihandlers "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/api"
	htmxh "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/htmx"
	"github.com/neosy/elengrab/internal/app/usecases"
)

type Dependencies struct {
	Usecases  *usecases.Usecases
	Templates *template.Template

	// Options
	AssetsDir string
}

type handlers struct {
	HTMX *htmxh.HTMXHandlers
	API  *apihandlers.APIHandlers
}

func New(deps *Dependencies) *handlers {
	return &handlers{
		HTMX: htmxh.NewHTMXHandlers(deps.Usecases, deps.Templates, deps.AssetsDir),
		API:  apihandlers.NewAPIHandlers(deps.Usecases),
	}
}
