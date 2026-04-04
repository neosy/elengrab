package linkclick

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type LinkClick struct {
	logger *slog.Logger

	// repositories
	linkClickRep persistence.LinkClickRepository
}

func NewLinkClick(
	logger *slog.Logger,
	linkClickRep persistence.LinkClickRepository,
) *LinkClick {
	return &LinkClick{
		logger:       logger,
		linkClickRep: linkClickRep,
	}
}
