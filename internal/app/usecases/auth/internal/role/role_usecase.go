package authrole

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type Role struct {
	logger *slog.Logger

	// repositories
	roleRepo persistence.RoleRepositoryFactory
}

func NewRole(
	logger *slog.Logger,
	roleRepo persistence.RoleRepositoryFactory,
) *Role {
	return &Role{
		logger:   logger,
		roleRepo: roleRepo,
	}
}
