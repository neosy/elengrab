package mediawatch

import (
	"log/slog"

	watchchunk "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch/media_watch_chunk"
	watchevent "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch/media_watch_event"
	watchstat "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch/media_watch_stat"
	nworkerpool "github.com/neosy/elengrab/internal/pkg/workerpool"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type MediaWatch struct {
	logger *slog.Logger

	// dispetchers
	watchEventDispatcher nworkerpool.JobDispatcher

	// State
	pendingStats statsUpdateQueue

	// internal usecase
	event *watchevent.MediaWatchEvent
	chunk *watchchunk.MediaWatchChunk
	stat  *watchstat.MediaWatchStat
}

func NewMediaWatch(
	logger *slog.Logger,

	// repositories
	eventRep persistence.MediaWatchEventRepository,
	chunkRep persistence.MediaWatchChunkRepository,
	statRep persistence.MediaWatchStatRepository,

	// in memory
	statCacheRep persistence.MediaWatchStatCacheRepository,

	// dispetchers
	watchEventDispatcher nworkerpool.JobDispatcher,
) *MediaWatch {
	return &MediaWatch{
		logger: logger,

		// dispetchers
		watchEventDispatcher: watchEventDispatcher,

		// State
		pendingStats: newStatsUpdateQueue(),

		event: watchevent.NewMediaWatchEvent(logger, eventRep),
		chunk: watchchunk.NewMediaWatchChunk(logger, chunkRep),
		stat:  watchstat.NewMediaWatchStat(logger, statRep, statCacheRep),
	}
}
