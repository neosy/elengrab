package authmw

import (
	"log/slog"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

type AuthMiddleware struct {
	logger *slog.Logger

	// services
	auth pservices.AuthMiddleware

	// options
	appMode dtypes.AppMode
}

func NewAuthMiddleware(logger *slog.Logger, auth pservices.AuthMiddleware, appMode dtypes.AppMode) *AuthMiddleware {
	return &AuthMiddleware{
		logger: logger,

		// services
		auth: auth,

		// options
		appMode: appMode,
	}
}
