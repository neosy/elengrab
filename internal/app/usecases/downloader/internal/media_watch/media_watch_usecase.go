package mediawatch

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch/mappers"
	uwatchchunk "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch/media_user_watch_chunk"
	uwatchposition "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch/media_user_watch_position"
	uwatchstat "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch/media_user_watch_stat"
	watchevent "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch/media_watch_event"
	watchstat "github.com/neosy/elengrab/internal/app/usecases/downloader/internal/media_watch/media_watch_stat"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/pkg/workerpool"
	"github.com/neosy/elengrab/internal/ports/persistence"
)

type MediaWatch struct {
	logger  *slog.Logger
	mappers *mappers.Mappers

	// Dispatchers
	watchEventDispatcher workerpool.JobDispatcher

	// State
	pendingStats statsUpdateQueue

	// Internal usecase
	event        *watchevent.MediaWatchEvent
	userChunk    *uwatchchunk.MediaUserWatchChunk
	userStat     *uwatchstat.MediaUserWatchStat
	stat         *watchstat.MediaWatchStat
	userPosition *uwatchposition.MediaUserWatchPosition

	// Callbacks
	onWatchStatsUpdated func(
		ctx context.Context,
		authCtx dauth.AuthContext,
		downloadID uuid.UUID,
	)
	onWatchPositionUpdated func(
		ctx context.Context,
		authCtx dauth.AuthContext,
		downloadID uuid.UUID,
	)
}

func NewMediaWatch(
	logger *slog.Logger,

	// Repositories
	eventRepo persistence.MediaWatchEventRepositoryFactory,
	chunkRepo persistence.MediaUserWatchChunkRepositoryFactory,
	userStatRepo persistence.MediaUserWatchStatRepositoryFactory,
	statRepo persistence.MediaWatchStatRepositoryFactory,
	positionRepo persistence.MediaUserWatchPositionRepositoryFactory,

	// In memory
	userStatCacheRep persistence.MediaUserWatchStatCacheRepository,
	statCacheRep persistence.MediaWatchStatCacheRepository,
	positionCacheRep persistence.MediaUserWatchPositionCacheRepository,

	// Dispetchers
	watchEventDispatcher workerpool.JobDispatcher,

	// Callbacks
	onWatchStatsUpdated func(
		ctx context.Context,
		authCtx dauth.AuthContext,
		downloadID uuid.UUID,
	),
	onWatchPositionUpdated func(
		ctx context.Context,
		authCtx dauth.AuthContext,
		downloadID uuid.UUID,
	),
) *MediaWatch {
	return &MediaWatch{
		logger:  logger,
		mappers: mappers.NewMappers(),

		// dispetchers
		watchEventDispatcher: watchEventDispatcher,

		// State
		pendingStats: newStatsUpdateQueue(),

		event:        watchevent.NewMediaWatchEvent(logger, eventRepo),
		userChunk:    uwatchchunk.NewMediaUserWatchChunk(logger, chunkRepo),
		userStat:     uwatchstat.NewMediaUserWatchStat(logger, userStatRepo, userStatCacheRep),
		stat:         watchstat.NewMediaWatchStat(logger, statRepo, statCacheRep),
		userPosition: uwatchposition.NewMediaUserWatchPosition(logger, positionRepo, positionCacheRep),

		// Callbacks
		onWatchStatsUpdated:    onWatchStatsUpdated,
		onWatchPositionUpdated: onWatchPositionUpdated,
	}
}
