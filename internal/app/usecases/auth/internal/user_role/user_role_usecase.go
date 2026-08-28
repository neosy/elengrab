package authuserrole

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type UserRole struct {
	logger *slog.Logger

	// repositories
	userRoleRepo persistence.UserRoleRepositoryFactory
}

func NewUserRole(
	logger *slog.Logger,
	userRoleRepo persistence.UserRoleRepositoryFactory,
) *UserRole {
	return &UserRole{
		logger:       logger,
		userRoleRepo: userRoleRepo,
	}
}
