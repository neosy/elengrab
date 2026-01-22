package handlers

import (
	"html/template"

	apihandlers "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/api"
	statich "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/static"
	uih "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui"
	"github.com/neosy/elengrab/internal/app/usecases"
)

type Dependencies struct {
	Usecases  *usecases.Usecases
	Templates *template.Template

	// Options
	AssetsDir    string
	DownloadsDir string
}

type handlers struct {
	Static *statich.StaticHandlers
	UI     *uih.UIHandlers
	API    *apihandlers.APIHandlers
}

func New(deps *Dependencies) *handlers {
	return &handlers{
		Static: statich.NewStaticHandlers(deps.AssetsDir),
		UI:     uih.NewUIHandlers(deps.Usecases, deps.Templates, deps.AssetsDir, deps.DownloadsDir),
		API:    apihandlers.NewAPIHandlers(deps.Usecases),
	}
}
