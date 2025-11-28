package handlers

import (
	"html/template"

	htmxh "github.com/neosy/elengrab/internal/api/rest/server/handlers/htmx"
	"github.com/neosy/elengrab/internal/app/usecases"
)

type Dependencies struct {
	Usecases  *usecases.Usecases
	Templates *template.Template

	// Options
	AssetsDir string
}

type handlers struct {
	HTMX *htmxh.Handlers
}

func New(deps *Dependencies) *handlers {
	return &handlers{
		HTMX: htmxh.NewHTMXHandlers(deps.Usecases, deps.Templates, deps.AssetsDir),
	}
}
