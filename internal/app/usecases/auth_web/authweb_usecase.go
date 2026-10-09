package authweb

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/app/usecases/auth_web/mappers"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

type authWeb struct {
	logger  *slog.Logger
	mappers *mappers.Mappers

	// services
	auth pservices.AuthService

	// options
	defaultAdminLogin    string
	defaultAdminPassword string
}

func NewAuthWeb(
	logger *slog.Logger,

	// services
	auth pservices.AuthService,

	// options
	defaultAdminLogin string,
	defaultAdminPassword string,
) AuthWeb {
	return &authWeb{
		logger:  logger,
		mappers: mappers.NewMappers(),

		// services
		auth: auth,

		// options
		defaultAdminLogin:    defaultAdminLogin,
		defaultAdminPassword: defaultAdminPassword,
	}
}
