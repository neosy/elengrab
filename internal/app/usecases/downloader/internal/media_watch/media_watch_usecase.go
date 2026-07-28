package mediawatch

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch/mappers"
	uwatchchunk "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch/media_user_watch_chunk"
	uwatchposition "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch/media_user_watch_position"
	uwatchstat "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch/media_user_watch_stat"
	watchevent "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch/media_watch_event"
	watchstat "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch/media_watch_stat"
	nworkerpool "github.com/neosy/elengrab/internal/pkg/workerpool"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type MediaWatch struct {
	logger  *slog.Logger
	mappers *mappers.Mappers

	// dispetchers
	watchEventDispatcher nworkerpool.JobDispatcher

	// State
	pendingStats statsUpdateQueue

	// internal usecase
	event        *watchevent.MediaWatchEvent
	userChunk    *uwatchchunk.MediaUserWatchChunk
	userStat     *uwatchstat.MediaUserWatchStat
	stat         *watchstat.MediaWatchStat
	userPosition *uwatchposition.MediaUserWatchPosition
}

func NewMediaWatch(
	logger *slog.Logger,

	// repositories
	eventRep persistence.MediaWatchEventRepository,
	chunkRep persistence.MediaUserWatchChunkRepository,
	userStatRep persistence.MediaUserWatchStatRepository,
	statRep persistence.MediaWatchStatRepository,
	positionRep persistence.MediaUserWatchPositionRepository,

	// in memory
	statCacheRep persistence.MediaWatchStatCacheRepository,
	positionCacheRep persistence.MediaUserWatchPositionCacheRepository,

	// dispetchers
	watchEventDispatcher nworkerpool.JobDispatcher,
) *MediaWatch {
	return &MediaWatch{
		logger:  logger,
		mappers: mappers.NewMappers(),

		// dispetchers
		watchEventDispatcher: watchEventDispatcher,

		// State
		pendingStats: newStatsUpdateQueue(),

		event:        watchevent.NewMediaWatchEvent(logger, eventRep),
		userChunk:    uwatchchunk.NewMediaUserWatchChunk(logger, chunkRep),
		userStat:     uwatchstat.NewMediaUserWatchStat(logger, userStatRep),
		stat:         watchstat.NewMediaWatchStat(logger, statRep, statCacheRep),
		userPosition: uwatchposition.NewMediaUserWatchPosition(logger, positionRep, positionCacheRep),
	}
}
