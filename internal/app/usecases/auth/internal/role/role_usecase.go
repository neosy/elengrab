package authrole

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type Role struct {
	logger *slog.Logger

	// repositories
	roleRep persistence.RoleRepository
}

func NewRole(
	logger *slog.Logger,
	roleRep persistence.RoleRepository,
) *Role {
	return &Role{
		logger:  logger,
		roleRep: roleRep,
	}
}
