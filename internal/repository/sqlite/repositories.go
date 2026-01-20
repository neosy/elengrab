package sqliterep

import (
	"database/sql"
	"fmt"

	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/internal/repository/sqlite/auth"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	"github.com/neosy/elengrab/internal/repository/sqlite/download"
)

// Repositories groups all database repositories.
type Repositories struct {
	lock dbexec.WriteLocker

	dbByName map[persistence.DBName]*sql.DB
	dbNames  []persistence.DBName

	File           persistence.FileRepository
	DownloadTask   persistence.DownloadTaskRepository
	YoutubeChannel persistence.YoutubeChannelRepository

	User        persistence.UserRepository
	UserSession persistence.UserSessionRepository
}

// New returns a new Repositories struct with database connections.
func New(dbByName map[persistence.DBName]*sql.DB) *Repositories {
	lock := dbexec.NewSQLiteLock()

	dbDownload := dbByName[persistence.DBMainName]
	dbAuth := dbByName[persistence.DBAuthName]

	var dbNames []persistence.DBName
	for name := range dbByName {
		dbNames = append(dbNames, name)
	}

	return &Repositories{
		lock:           lock,
		dbByName:       dbByName,
		dbNames:        dbNames,
		File:           download.NewFileRepository(dbDownload, lock),
		DownloadTask:   download.NewDownloadTaskRepository(dbDownload, lock),
		YoutubeChannel: download.NewYoutubeChannelRepository(dbDownload, lock),
		User:           auth.NewUserRepository(dbAuth, lock),
		UserSession:    auth.NewUserSessionRepository(dbAuth, lock),
	}
}

func (r *Repositories) GetDBNames() []persistence.DBName {
	return r.dbNames
}

func (r *Repositories) DBFileName(dbName persistence.DBName) string {
	return DBFileName(dbName)
}

func DBFileName(dbName persistence.DBName) string {
	return fmt.Sprintf("%s.db", dbName)
}
