package dlstate

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/ports/persistence"
)

type DownloadState struct {
	logger *slog.Logger

	// repositories
	stateRep persistence.DownloadStateRepository
}

func NewDownloadState(
	logger *slog.Logger,
	stateRep persistence.DownloadStateRepository,
) *DownloadState {
	return &DownloadState{
		logger:   logger,
		stateRep: stateRep,
	}
}
