package authweb

import (
	"log/slog"

	pservices "github.com/neosy/elengrab/internal/ports/services"
)

type AuthWeb struct {
	logger *slog.Logger

	// services
	auth pservices.Auth

	// options
	defaultAdminLogin    string
	defaultAdminPassword string
}

func NewAuthWeb(
	logger *slog.Logger,

	// services
	auth pservices.Auth,

	// options
	defaultAdminLogin string,
	defaultAdminPassword string,
) *AuthWeb {
	return &AuthWeb{
		logger: logger,

		// services
		auth: auth,

		// options
		defaultAdminLogin:    defaultAdminLogin,
		defaultAdminPassword: defaultAdminPassword,
	}
}
