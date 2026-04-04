package linklink

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type Link struct {
	logger *slog.Logger

	// repositories
	linkRep persistence.LinkRepository

	// options
	deduplicate bool
}

func NewLink(
	logger *slog.Logger,
	linkRep persistence.LinkRepository,
	deduplicate bool,
) *Link {
	return &Link{
		logger:      logger,
		linkRep:     linkRep,
		deduplicate: deduplicate,
	}
}
