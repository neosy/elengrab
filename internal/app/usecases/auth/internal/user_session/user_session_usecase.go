package authsession

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type UserSession struct {
	logger *slog.Logger

	// repositories
	userSessionRepo persistence.UserSessionRepositoryFactory
}

func NewUserSession(
	logger *slog.Logger,
	userSessionRepo persistence.UserSessionRepositoryFactory,
) *UserSession {
	return &UserSession{
		logger:         logger,
		userSessionRepo: userSessionRepo,
	}
}
