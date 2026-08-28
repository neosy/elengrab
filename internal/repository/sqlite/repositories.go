package sqliterep

import (
	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/internal/repository/sqlite/auth"
	"github.com/neosy/elengrab/internal/repository/sqlite/download"
	"github.com/neosy/elengrab/internal/repository/sqlite/link"
	"github.com/neosy/elengrab/internal/repository/sqlite/media"
	searchindex "github.com/neosy/elengrab/internal/repository/sqlite/search_index"
	"github.com/neosy/elengrab/internal/repository/sqlite/types"
	watchevent "github.com/neosy/elengrab/internal/repository/sqlite/watch_event"
)

// Repositories groups all database repositories.
type Repositories struct {
	dbRegistry *types.DBRegistry

	User        persistence.UserRepositoryFactory
	Role        persistence.RoleRepositoryFactory
	UserRole    persistence.UserRoleRepositoryFactory
	UserSession persistence.UserSessionRepositoryFactory

	DownloadDataMigration persistence.DownloadDataMigrationRepositoryFactory
	MediaDownload         persistence.MediaDownloadRepositoryFactory
	DownloadTask          persistence.DownloadTaskRepositoryFactory

	MediaSourceIndex persistence.MediaSourceIndexRepositoryFactory

	MediaWatchEvent        persistence.MediaWatchEventRepositoryFactory
	MediaUserWatchChunk    persistence.MediaUserWatchChunkRepositoryFactory
	MediaUserWatchStat     persistence.MediaUserWatchStatRepositoryFactory
	MediaWatchStat         persistence.MediaWatchStatRepositoryFactory
	MediaUserWatchPosition persistence.MediaUserWatchPositionRepositoryFactory

	YoutubeChannel persistence.YoutubeChannelRepositoryFactory
	SiteLogo       persistence.SiteLogoRepositoryFactory
	Thumbnail      persistence.ThumbnailRepositoryFactory

	Link      persistence.LinkRepositoryFactory
	LickClick persistence.LinkClickRepositoryFactory
}

// New returns a new Repositories struct with database connections.
func New(dbEntries []persistence.DBEntry) *Repositories {
	var entriesByName = make(map[string]persistence.DBEntry, len(dbEntries))

	for _, e := range dbEntries {
		entriesByName[e.DBName()] = e
	}

	authEntry := types.NewDBEntry(
		AuthSchema,
		entriesByName[AuthSchema.DBName()].DB(),
	)

	mainEntry := types.NewDBEntry(
		MainSchema,
		entriesByName[MainSchema.DBName()].DB(),
	)

	mediaEntry := types.NewDBEntry(
		MediaSchema,
		entriesByName[MediaSchema.DBName()].DB(),
	)

	linkEntry := types.NewDBEntry(
		LinkSchema,
		entriesByName[LinkSchema.DBName()].DB(),
	)

	watchEventEntry := types.NewDBEntry(
		WatchEventSchema,
		entriesByName[WatchEventSchema.DBName()].DB(),
	)

	searchIndexEntry := types.NewDBEntry(
		SearchIndexSchema,
		entriesByName[SearchIndexSchema.DBName()].DB(),
	)

	return &Repositories{
		dbRegistry: types.NewRegistry(entriesByName),

		User:        auth.NewUserRepository(authEntry),
		Role:        auth.NewRoleRepository(authEntry),
		UserRole:    auth.NewUserRoleRepository(authEntry),
		UserSession: auth.NewUserSessionRepository(authEntry),

		DownloadDataMigration: download.NewDataMigrationRepository(mainEntry),
		MediaDownload:         download.NewMediaDownloadRepository(mainEntry),
		DownloadTask:          download.NewDownloadTaskRepository(mainEntry),

		MediaWatchEvent:        watchevent.NewMediaWatchEventRepository(watchEventEntry),
		MediaUserWatchChunk:    watchevent.NewMediaUserWatchChunkRepository(watchEventEntry),
		MediaUserWatchStat:     watchevent.NewMediaUserWatchStatRepository(watchEventEntry),
		MediaWatchStat:         watchevent.NewMediaWatchStatRepository(watchEventEntry),
		MediaUserWatchPosition: watchevent.NewMediaUserWatchPositionRepository(watchEventEntry),

		MediaSourceIndex: searchindex.NewMediaSourceIndexRepository(searchIndexEntry),

		YoutubeChannel: media.NewYoutubeChannelRepository(mediaEntry),
		SiteLogo:       media.NewSiteLogoRepository(mediaEntry),
		Thumbnail:      media.NewThumbnailRepository(mediaEntry),

		Link:      link.NewLinkRepository(linkEntry),
		LickClick: link.NewLinkClickRepository(linkEntry),
	}
}

func (r *Repositories) Schemas() []persistence.DBSchema {
	return r.dbRegistry.Schemas()
}

func (r *Repositories) SchemasByName() map[string]persistence.DBSchema {
	return r.dbRegistry.SchemasByName()
}

func (r *Repositories) EntriesByName() map[string]persistence.DBEntry {
	return r.dbRegistry.EntriesByName()
}
