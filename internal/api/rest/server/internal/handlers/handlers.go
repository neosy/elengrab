package handlers

import (
	"fmt"
	"html/template"
	"log/slog"
	"path/filepath"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/api"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/static"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	httptemplates "github.com/neosy/elengrab/internal/api/rest/server/templates"
	"github.com/neosy/elengrab/internal/app/usecases"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	appenv "github.com/neosy/elengrab/internal/pkg/config/app_env"
	pstorage "github.com/neosy/elengrab/internal/ports/storage"
)

type Dependencies struct {
	// Storages
	DownloadsStorage pstorage.DownloadsStorage

	// Assets
	Assets *assets.Assets

	Usecases *usecases.Usecases
	Template *template.Template

	// Options
	AppEnv          appenv.AppEnv
	AppMode         dtypes.AppMode
	BaseURL         string
	ShortLinkPrefix string
}

type Handlers struct {
	Static *static.StaticHandlers
	UI     *ui.Handlers
	API    *api.Handlers
}

func New(logger *slog.Logger, deps *Dependencies) (*Handlers, error) {
	templates, err := loadPageTemplates(deps.Template, deps.Assets)
	if err != nil {
		return nil, err
	}

	return &Handlers{
		Static: static.NewStaticHandlers(
			deps.Assets,
			deps.Usecases.Thumbnail,
			deps.Usecases.Downloader,
			deps.AppEnv,
		),
		UI: ui.NewHandlers(
			logger,
			deps.DownloadsStorage,
			deps.Assets,
			deps.Usecases,
			templates,
			deps.AppMode,
			deps.BaseURL,
			deps.ShortLinkPrefix,
		),
		API: api.NewHandlers(deps.Usecases),
	}, nil
}

func loadPageTemplates(source *template.Template, assets *assets.Assets) (*httptemplates.Templates, error) {
	templatePages := make(map[string]*template.Template)

	for _, page := range pages.AllPages() {
		key := page.Key()

		if _, exists := templatePages[key]; exists {
			return nil, fmt.Errorf("duplicate page template key: %s", key)
		}

		t, err := source.Clone()
		if err != nil {
			return nil, err
		}

		t, err = t.ParseFiles(filepath.Join(assets.FolderPaths().Pages(), page.FileName()))
		if err != nil {
			return nil, fmt.Errorf("parse page %s: %w", key, err)
		}

		templatePages[key] = t
	}

	return &httptemplates.Templates{
		Base:  source,
		Pages: templatePages,
	}, nil
}
