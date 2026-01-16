package sqliterep

import (
	"database/sql"

	"github.com/neosy/elengrab/internal/ports/persistence"
	sldownload "github.com/neosy/elengrab/internal/repository/sqlite/download"
	"github.com/neosy/elengrab/internal/repository/sqlite/lock"
)

// Repositories groups all database repositories.
type Repositories struct {
	lock lock.WriteLocker

	db *sql.DB

	File           persistence.FileRepository
	DownloadTask   persistence.DownloadTaskRepository
	YoutubeChannel persistence.YoutubeChannelRepository

	User        persistence.UserRepository
	UserSession persistence.UserSessionRepository
}

// New returns a new Repositories struct with database connections.
func New(db *sql.DB) *Repositories {
	lock := lock.NewSQLiteLock()

	return &Repositories{
		lock:           lock,
		db:             db,
		File:           sldownload.NewFileRepository(db, lock),
		DownloadTask:   sldownload.NewDownloadTaskRepository(db, lock),
		YoutubeChannel: sldownload.NewYoutubeChannelRepository(db, lock),
		User:           sldownload.NewUserRepository(db, lock),
		UserSession:    sldownload.NewUserSessionRepository(db, lock),
	}
}
