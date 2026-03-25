package authuser

import (
	"log/slog"

	authuserrole "github.com/neosy/elengrab/internal/app/usecases/auth/internal/user_role"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type User struct {
	logger *slog.Logger

	// repositories
	userRep persistence.UserRepository

	// internal
	userRole *authuserrole.UserRole
}

func NewUser(
	logger *slog.Logger,
	userRep persistence.UserRepository,
	userRole *authuserrole.UserRole,
) *User {
	return &User{
		logger:  logger,
		userRep: userRep,

		// internal
		userRole: userRole,
	}
}
