package auth

import (
	"log/slog"

	useruc "github.com/neosy/elengrab/internal/app/usecases/auth/internal/user"
	usersession "github.com/neosy/elengrab/internal/app/usecases/auth/internal/user_session"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type Auth struct {
	logger *slog.Logger

	// internal
	user        *useruc.User
	userSession *usersession.UserSession
}

func NewAuth(
	logger *slog.Logger,

	// repositories
	userRep persistence.UserRepository,
	userSessionRep persistence.UserSessionRepository,
) *Auth {
	return &Auth{
		logger:      logger,
		user:        useruc.NewUser(logger, userRep),
		userSession: usersession.NewUserSession(logger, userSessionRep),
	}
}
