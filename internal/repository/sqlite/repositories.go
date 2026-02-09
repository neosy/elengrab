package sqliterep

import (
	"database/sql"

	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/internal/repository/sqlite/auth"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	"github.com/neosy/elengrab/internal/repository/sqlite/download"
	"github.com/neosy/elengrab/internal/repository/sqlite/media"
)

// Repositories groups all database repositories.
type Repositories struct {
	lock dbexec.WriteLocker

	dbByName map[persistence.DBName]*sql.DB
	dbNames  []persistence.DBName

	File           persistence.FileRepository
	DownloadTask   persistence.DownloadTaskRepository
	YoutubeChannel persistence.YoutubeChannelRepository
	SiteLogo       persistence.SiteLogoRepository

	User        persistence.UserRepository
	UserSession persistence.UserSessionRepository
}

// New returns a new Repositories struct with database connections.
func New(dbByName map[persistence.DBName]*sql.DB) *Repositories {
	lock := dbexec.NewSQLiteLock()

	authDB := dbByName[persistence.DBAuthName]
	mainDB := dbByName[persistence.DBMainName]
	mediaDB := dbByName[persistence.DBMediaName]

	var dbNames []persistence.DBName
	for name := range dbByName {
		dbNames = append(dbNames, name)
	}

	return &Repositories{
		lock:     lock,
		dbByName: dbByName,
		dbNames:  dbNames,

		User:        auth.NewUserRepository(authDB, lock),
		UserSession: auth.NewUserSessionRepository(authDB, lock),

		File:         download.NewFileRepository(mainDB, lock),
		DownloadTask: download.NewDownloadTaskRepository(mainDB, lock),

		YoutubeChannel: media.NewYoutubeChannelRepository(mediaDB, lock),
		SiteLogo:       media.NewSiteLogoRepository(mediaDB, lock),
	}
}

func (r *Repositories) GetDBNames() []persistence.DBName {
	return r.dbNames
}
