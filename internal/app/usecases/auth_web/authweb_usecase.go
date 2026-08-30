package authweb

import (
	"log/slog"

	pservices "github.com/neosy/elengrab/internal/ports/services"
)

type authWeb struct {
	logger *slog.Logger

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
		logger: logger,

		// services
		auth: auth,

		// options
		defaultAdminLogin:    defaultAdminLogin,
		defaultAdminPassword: defaultAdminPassword,
	}
}
