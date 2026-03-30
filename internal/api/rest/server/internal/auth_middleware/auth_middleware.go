package authmw

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/app/usecases/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

type AuthMiddleware struct {
	logger *slog.Logger

	// services
	auth pservices.Auth

	// options
	appMode dtypes.AppMode
}

func NewAuthMiddleware(logger *slog.Logger, auth *auth.Auth, appMode dtypes.AppMode) *AuthMiddleware {
	return &AuthMiddleware{
		logger: logger,

		// services
		auth: auth,

		// options
		appMode: appMode,
	}
}
