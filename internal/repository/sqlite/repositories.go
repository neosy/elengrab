package sqliterep

import (
	"database/sql"

	"github.com/neosy/elengrab/internal/ports/persistence"
	sldownload "github.com/neosy/elengrab/internal/repository/sqlite/download"
)

// Repositories groups all database repositories.
type Repositories struct {
	File persistence.FileRepository
}

// New returns a new Repositories struct with database connections.
func New(db *sql.DB) *Repositories {
	return &Repositories{
		File: sldownload.NewFileRepository(db),
	}
}
