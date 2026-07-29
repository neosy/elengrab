package admin

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin/mappers"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin/validators"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	httptemplates "github.com/neosy/elengrab/internal/api/rest/server/templates"
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

	templates  *httptemplates.Templates
	assets     *assets.Assets
	assetPaths paths.AssetPaths

	// Usecases
	usecases appUsecases

	// Options
	appMode dtypes.AppMode
	baseURL string
}

func NewAdminHandlers(
	logger *slog.Logger,

	templates *httptemplates.Templates,
	assets *assets.Assets,
	assetPaths paths.AssetPaths,

	// Usecases
	usecases *usecases.Usecases,

	appMode dtypes.AppMode,
	baseURL string,
) *AdminHandlers {
	return &AdminHandlers{
		logger:     logger,
		mappers:    mappers.NewMappers(),
		validators: validators.NewValidators(),

		templates:  templates,
		assets:     assets,
		assetPaths: assetPaths,

		// Usecases
		usecases: appUsecases{
			admin: usecases.Admin,
		},

		// Options
		appMode: appMode,
		baseURL: baseURL,
	}
}
