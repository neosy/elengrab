package useruc

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type User struct {
	logger *slog.Logger

	// repositories
	userRep persistence.UserRepository
}

func NewUser(
	logger *slog.Logger,
	userRep persistence.UserRepository,
) *User {
	return &User{
		logger:  logger,
		userRep: userRep,
	}
}
