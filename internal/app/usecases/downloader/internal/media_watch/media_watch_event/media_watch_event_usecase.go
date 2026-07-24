package watchevent

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type MediaWatchEvent struct {
	logger *slog.Logger

	// repositories
	eventRep persistence.MediaWatchEventRepository
}

func NewMediaWatchEvent(
	logger *slog.Logger,

	// repositories
	eventRep persistence.MediaWatchEventRepository,
) *MediaWatchEvent {
	return &MediaWatchEvent{
		logger:   logger,
		eventRep: eventRep,
	}
}
