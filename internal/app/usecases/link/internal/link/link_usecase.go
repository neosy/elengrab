package link

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type Link struct {
	logger *slog.Logger

	// repositories
	linkRepo persistence.LinkRepositoryFactory
}

func NewLink(
	logger *slog.Logger,
	linkRepo persistence.LinkRepositoryFactory,
) *Link {
	return &Link{
		logger:   logger,
		linkRepo: linkRepo,
	}
}
