package linkclick

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type LinkClick struct {
	logger *slog.Logger

	// repositories
	linkClickRepo persistence.LinkClickRepositoryFactory
}

func NewLinkClick(
	logger *slog.Logger,
	linkClickRepo persistence.LinkClickRepositoryFactory,
) *LinkClick {
	return &LinkClick{
		logger:        logger,
		linkClickRepo: linkClickRepo,
	}
}
