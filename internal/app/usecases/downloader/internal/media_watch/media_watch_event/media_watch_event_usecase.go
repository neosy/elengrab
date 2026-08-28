package watchevent

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type MediaWatchEvent struct {
	logger *slog.Logger

	// repositories
	eventRepo persistence.MediaWatchEventRepositoryFactory
}

func NewMediaWatchEvent(
	logger *slog.Logger,

	// repositories
	eventRepo persistence.MediaWatchEventRepositoryFactory,
) *MediaWatchEvent {
	return &MediaWatchEvent{
		logger:    logger,
		eventRepo: eventRepo,
	}
}
