package authmw

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/app/usecases/auth"
)

type AuthMiddleware struct {
	logger *slog.Logger

	// usecases
	auth *auth.Auth
}

func NewAuthMiddleware(logger *slog.Logger, auth *auth.Auth) *AuthMiddleware {
	return &AuthMiddleware{
		logger: logger,

		// usecases
		auth: auth,
	}
}
