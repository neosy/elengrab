package auth

import (
	"log/slog"
	"time"

	authrole "github.com/neosy/elengrab/internal/app/usecases/auth/internal/role"
	authuser "github.com/neosy/elengrab/internal/app/usecases/auth/internal/user"
	authuserrole "github.com/neosy/elengrab/internal/app/usecases/auth/internal/user_role"
	authsession "github.com/neosy/elengrab/internal/app/usecases/auth/internal/user_session"
	"github.com/neosy/elengrab/internal/app/usecases/auth/mappers"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type (
	Auth struct {
		logger  *slog.Logger
		mappers *mappers.Mappers

		// internal
		user        *authuser.User
		role        *authrole.Role
		userRole    *authuserrole.UserRole
		userSession *authsession.UserSession

		// options
		sessionTTL time.Duration
		// Refresh interval: how often we extend the session expiration on activity
		sessionRefreshInterval time.Duration
	}
)

func NewAuth(
	logger *slog.Logger,

	// repositories
	userRep persistence.UserRepository,
	roleRep persistence.RoleRepository,
	userRoleRep persistence.UserRoleRepository,
	userSessionRep persistence.UserSessionRepository,

	// optons
	opts ...AuthOption,
) *Auth {
	options := NewAuthOptions(opts...)
	userRole := authuserrole.NewUserRole(logger, userRoleRep)

	return &Auth{
		logger:  logger,
		mappers: mappers.NewMappers(),

		// internal
		user:        authuser.NewUser(logger, userRep, userRole),
		role:        authrole.NewRole(logger, roleRep),
		userRole:    userRole,
		userSession: authsession.NewUserSession(logger, userSessionRep),

		// options
		sessionTTL:             options.SessionTTL,
		sessionRefreshInterval: options.SessionRefreshInterval,
	}
}
