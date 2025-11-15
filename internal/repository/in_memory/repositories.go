package inmemoryrep

import (
	"github.com/neosy/elengrab/internal/ports/persistence"
	imdownload "github.com/neosy/elengrab/internal/repository/in_memory/download"
)

// Repositories groups all database repositories.
type Repositories struct {
	DownloadState persistence.DownloadStateRepository
}

// New returns a new Repositories struct.
func New() *Repositories {
	return &Repositories{
		DownloadState: imdownload.NewDownloadStateRepository(),
	}
}
