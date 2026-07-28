package sqliterep

import (
	"database/sql"

	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/internal/repository/sqlite/auth"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	"github.com/neosy/elengrab/internal/repository/sqlite/download"
	"github.com/neosy/elengrab/internal/repository/sqlite/link"
	"github.com/neosy/elengrab/internal/repository/sqlite/media"
	sqlitetypes "github.com/neosy/elengrab/internal/repository/sqlite/types"
	watchevent "github.com/neosy/elengrab/internal/repository/sqlite/watch_event"
)

// Repositories groups all database repositories.
type Repositories struct {
	dbRegistry *sqlitetypes.DBRegistry

	User        persistence.UserRepository
	Role        persistence.RoleRepository
	UserRole    persistence.UserRoleRepository
	UserSession persistence.UserSessionRepository

	MediaDownload         persistence.MediaDownloadRepository
	DownloadTask          persistence.DownloadTaskRepository
	DownloadDataMigration persistence.DownloadDataMigrationRepository

	MediaWatchEvent        persistence.MediaWatchEventRepository
	MediaUserWatchChunk    persistence.MediaUserWatchChunkRepository
	MediaUserWatchStat     persistence.MediaUserWatchStatRepository
	MediaWatchStat         persistence.MediaWatchStatRepository
	MediaUserWatchPosition persistence.MediaUserWatchPositionRepository

	YoutubeChannel persistence.YoutubeChannelRepository
	SiteLogo       persistence.SiteLogoRepository
	Thumbnail      persistence.ThumbnailRepository

	Link      persistence.LinkRepository
	LickClick persistence.LinkClickRepository
}

// New returns a new Repositories struct with database connections.
func New(dbEntries []persistence.DBEntry) *Repositories {
	var entriesByName = make(map[string]persistence.DBEntry, len(dbEntries))

	for _, e := range dbEntries {
		entriesByName[e.DBName()] = e
	}

	type dbEntry struct {
		db     *sql.DB
		locker dbexec.WriteLocker
	}

	eAuth := dbEntry{
		db:     entriesByName[AuthSchema.DBName()].DB(),
		locker: dbexec.NewSQLiteLock(),
	}

	eMain := dbEntry{
		db:     entriesByName[MainSchema.DBName()].DB(),
		locker: dbexec.NewSQLiteLock(),
	}

	eMedia := dbEntry{
		db:     entriesByName[MediaSchema.DBName()].DB(),
		locker: dbexec.NewSQLiteLock(),
	}

	eLink := dbEntry{
		db:     entriesByName[LinkSchema.DBName()].DB(),
		locker: dbexec.NewSQLiteLock(),
	}

	eWatchEvent := dbEntry{
		db:     entriesByName[WatchEventSchema.DBName()].DB(),
		locker: dbexec.NewSQLiteLock(),
	}

	return &Repositories{
		dbRegistry: sqlitetypes.NewRegistry(entriesByName),

		User:        auth.NewUserRepository(eAuth.db, eAuth.locker),
		Role:        auth.NewRoleRepository(eAuth.db, eAuth.locker),
		UserRole:    auth.NewUserRoleRepository(eAuth.db, eAuth.locker),
		UserSession: auth.NewUserSessionRepository(eAuth.db, eAuth.locker),

		MediaDownload:         download.NewMediaDownloadRepository(eMain.db, eMain.locker),
		DownloadTask:          download.NewDownloadTaskRepository(eMain.db, eMain.locker),
		DownloadDataMigration: download.NewDataMigrationRepository(eMain.db, eMain.locker),

		MediaWatchEvent:        watchevent.NewMediaWatchEventRepository(eWatchEvent.db, eWatchEvent.locker),
		MediaUserWatchChunk:    watchevent.NewMediaUserWatchChunkRepository(eWatchEvent.db, eWatchEvent.locker),
		MediaUserWatchStat:     watchevent.NewMediaUserWatchStatRepository(eWatchEvent.db, eWatchEvent.locker),
		MediaWatchStat:         watchevent.NewMediaWatchStatRepository(eWatchEvent.db, eWatchEvent.locker),
		MediaUserWatchPosition: watchevent.NewMediaUserWatchPositionRepository(eWatchEvent.db, eWatchEvent.locker),

		YoutubeChannel: media.NewYoutubeChannelRepository(eMedia.db, eMedia.locker),
		SiteLogo:       media.NewSiteLogoRepository(eMedia.db, eMedia.locker),
		Thumbnail:      media.NewThumbnailRepository(eMedia.db, eMedia.locker),

		Link:      link.NewLinkRepository(eLink.db, eLink.locker),
		LickClick: link.NewLinkClickRepository(eLink.db, eLink.locker),
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
