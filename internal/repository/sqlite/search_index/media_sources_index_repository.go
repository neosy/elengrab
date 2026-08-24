package searchindex

import (
	"context"
	"database/sql"

	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	"github.com/neosy/elengrab/internal/repository/sqlite/search_index/mappers"
)

type MediaSourceIndexRepository struct {
	mappers *mappers.Mappers
	db      *sql.DB
	lock    dbexec.WriteLocker

	filtersByName filtersByName
	queryOptions  queryOptions

	// options
	retryOptions dbexec.RetryOptions
}

// NewMediaSourceIndexRepository returns a new object for the repository
func NewMediaSourceIndexRepository(db *sql.DB, lock dbexec.WriteLocker) *MediaSourceIndexRepository {
	return &MediaSourceIndexRepository{
		mappers: mappers.NewMappers(),
		db:      db,
		lock:    lock,

		filtersByName: make(map[string]any),

		// options
		retryOptions: dbexec.RetryOptions{
			MaxRetries: maxRetriesDefault,
			Delay:      retryDelayDefault,
		},
	}
}

func (r *MediaSourceIndexRepository) Copy() *MediaSourceIndexRepository {
	rep := uptr.Copy(r)

	rep.mappers = r.mappers
	rep.db = r.db
	rep.lock = r.lock

	rep.filtersByName = rep.filtersByName.copy()
	rep.queryOptions = rep.queryOptions.copy()

	return rep
}

func (r *MediaSourceIndexRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.Tx(ctx, r.db, r.lock, fn)
}

func (r *MediaSourceIndexRepository) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.TxIndependent(ctx, r.db, r.lock, fn)
}
