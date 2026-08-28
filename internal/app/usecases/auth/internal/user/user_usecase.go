package authuser

import (
	"log/slog"

	authuserrole "github.com/neosy/elengrab/internal/app/usecases/auth/internal/user_role"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type User struct {
	logger *slog.Logger

	// repositories
	userRepo persistence.UserRepositoryFactory

	// internal
	userRole *authuserrole.UserRole
}

func NewUser(
	logger *slog.Logger,
	userRepo persistence.UserRepositoryFactory,
	userRole *authuserrole.UserRole,
) *User {
	return &User{
		logger:  logger,
		userRepo: userRepo,

		// internal
		userRole: userRole,
	}
}
