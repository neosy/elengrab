package sqliterep

import (
	"database/sql"
	"sync"

	"github.com/neosy/elengrab/internal/ports/persistence"
	sldownload "github.com/neosy/elengrab/internal/repository/sqlite/download"
)

// Repositories groups all database repositories.
type Repositories struct {
	mu *sync.RWMutex

	File           persistence.FileRepository
	DownloadTask   persistence.DownloadTaskRepository
	YoutubeChannel persistence.YoutubeChannelRepository
}

// New returns a new Repositories struct with database connections.
func New(db *sql.DB) *Repositories {
	var mu sync.RWMutex

	return &Repositories{
		mu:             &mu,
		File:           sldownload.NewFileRepository(db, &mu),
		DownloadTask:   sldownload.NewDownloadTaskRepository(db, &mu),
		YoutubeChannel: sldownload.NewYoutubeChannelRepository(db, &mu),
	}
}
