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
	lockByDBName := make(map[persistence.DBName]dbexec.WriteLocker, len(dbByName))
	for name := range dbByName {
		lockByDBName[name] = dbexec.NewSQLiteLock()
	}

	authDB := dbByName[persistence.DBAuthName]
	mainDB := dbByName[persistence.DBMainName]
	mediaDB := dbByName[persistence.DBMediaName]

	var dbNames []persistence.DBName
	for name := range dbByName {
		dbNames = append(dbNames, name)
	}

	return &Repositories{
		dbByName: dbByName,
		dbNames:  dbNames,

		User:        auth.NewUserRepository(authDB, lockByDBName[persistence.DBAuthName]),
		UserSession: auth.NewUserSessionRepository(authDB, lockByDBName[persistence.DBAuthName]),

		File:         download.NewFileRepository(mainDB, lockByDBName[persistence.DBMainName]),
		DownloadTask: download.NewDownloadTaskRepository(mainDB, lockByDBName[persistence.DBMainName]),

		YoutubeChannel: media.NewYoutubeChannelRepository(mediaDB, lockByDBName[persistence.DBMediaName]),
		SiteLogo:       media.NewSiteLogoRepository(mediaDB, lockByDBName[persistence.DBMediaName]),
	}
}

func (r *Repositories) GetDBNames() []persistence.DBName {
	return r.dbNames
}
