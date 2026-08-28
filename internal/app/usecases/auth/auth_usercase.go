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
	userRepo persistence.UserRepositoryFactory,
	roleRepo persistence.RoleRepositoryFactory,
	userRoleRepo persistence.UserRoleRepositoryFactory,
	userSessionRepo persistence.UserSessionRepositoryFactory,

	// optons
	opts ...AuthOption,
) *Auth {
	options := NewAuthOptions(opts...)
	userRole := authuserrole.NewUserRole(logger, userRoleRepo)

	return &Auth{
		logger:  logger,
		mappers: mappers.NewMappers(),

		// internal
		user:        authuser.NewUser(logger, userRepo, userRole),
		role:        authrole.NewRole(logger, roleRepo),
		userRole:    userRole,
		userSession: authsession.NewUserSession(logger, userSessionRepo),

		// options
		sessionTTL:             options.SessionTTL,
		sessionRefreshInterval: options.SessionRefreshInterval,
	}
}
