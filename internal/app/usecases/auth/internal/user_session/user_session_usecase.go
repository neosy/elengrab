package usersession

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type UserSession struct {
	logger *slog.Logger

	// repositories
	userSessionRep persistence.UserSessionRepository
}

func NewUserSession(
	logger *slog.Logger,
	userSessionRep persistence.UserSessionRepository,
) *UserSession {
	return &UserSession{
		logger:         logger,
		userSessionRep: userSessionRep,
	}
}
