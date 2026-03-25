package authuserrole

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type UserRole struct {
	logger *slog.Logger

	// repositories
	userRoleRep persistence.UserRoleRepository
}

func NewUserRole(
	logger *slog.Logger,
	userRoleRep persistence.UserRoleRepository,
) *UserRole {
	return &UserRole{
		logger:      logger,
		userRoleRep: userRoleRep,
	}
}
