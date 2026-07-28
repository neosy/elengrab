package mediawatch

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch/mappers"
	uwatchstat "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch/media_user_watch_stat"
	watchchunk "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch/media_watch_chunk"
	watchevent "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch/media_watch_event"
	watchposition "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch/media_watch_position"
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
	event    *watchevent.MediaWatchEvent
	chunk    *watchchunk.MediaWatchChunk
	userStat *uwatchstat.MediaUserWatchStat
	stat     *watchstat.MediaWatchStat
	position *watchposition.MediaWatchPosition
}

func NewMediaWatch(
	logger *slog.Logger,

	// repositories
	eventRep persistence.MediaWatchEventRepository,
	chunkRep persistence.MediaWatchChunkRepository,
	userStatRep persistence.MediaUserWatchStatRepository,
	statRep persistence.MediaWatchStatRepository,
	positionRep persistence.MediaWatchPositionRepository,

	// in memory
	statCacheRep persistence.MediaWatchStatCacheRepository,
	positionCacheRep persistence.MediaWatchPositionCacheRepository,

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

		event:    watchevent.NewMediaWatchEvent(logger, eventRep),
		chunk:    watchchunk.NewMediaWatchChunk(logger, chunkRep),
		userStat: uwatchstat.NewMediaUserWatchStat(logger, userStatRep),
		stat:     watchstat.NewMediaWatchStat(logger, statRep, statCacheRep),
		position: watchposition.NewMediaWatchPosition(logger, positionRep, positionCacheRep),
	}
}
