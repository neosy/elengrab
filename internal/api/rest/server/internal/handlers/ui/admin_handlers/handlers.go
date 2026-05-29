package handlers

import (
	"html/template"
	"log/slog"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin_handlers/mappers"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin_handlers/validators"
	"github.com/neosy/elengrab/internal/app/usecases"
	adminuc "github.com/neosy/elengrab/internal/app/usecases/admin"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type appUsecases struct {
	admin *adminuc.Admin
}

type AdminHandlers struct {
	logger     *slog.Logger
	mappers    *mappers.Mappers
	validators *validators.Validators

	templates    *template.Template
	assetFolders assets.FolderPaths

	// Usecases
	usecases appUsecases

	// Options
	appMode dtypes.AppMode
	baseURL string
}

func NewAdminHandlers(
	logger *slog.Logger,

	templates *template.Template,
	assetFolders assets.FolderPaths,

	// Usecases
	usecases *usecases.Usecases,

	appMode dtypes.AppMode,
	baseURL string,
) *AdminHandlers {
	return &AdminHandlers{
		logger:     logger,
		mappers:    mappers.NewMappers(),
		validators: validators.NewValidators(),

		templates:    templates,
		assetFolders: assetFolders,

		// Usecases
		usecases: appUsecases{
			admin: usecases.Admin,
		},

		// Options
		appMode: appMode,
		baseURL: baseURL,
	}
}
