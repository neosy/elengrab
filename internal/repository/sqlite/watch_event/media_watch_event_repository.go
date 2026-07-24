package watchevent

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/pkg/dbutils"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	ewatchevent "github.com/neosy/elengrab/internal/repository/sqlite/watch_event/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/watch_event/mappers"
)

type MediaWatchEventRepository struct {
	mappers *mappers.Mappers
	db      *sql.DB
	lock    dbexec.WriteLocker

	filtersByName filtersByName

	// options
	retryOptions dbexec.RetryOptions
}

// NewMediaWatchEventRepository returns a new object for the repository
func NewMediaWatchEventRepository(db *sql.DB, lock dbexec.WriteLocker) *MediaWatchEventRepository {
	return &MediaWatchEventRepository{
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

func (r *MediaWatchEventRepository) Copy() *MediaWatchEventRepository {
	rep := uptr.Copy(r)

	rep.mappers = r.mappers
	rep.db = r.db
	rep.lock = r.lock

	rep.filtersByName = rep.filtersByName.copy()

	return rep
}

func (r *MediaWatchEventRepository) Insert(ctx context.Context, event *ddownload.MediaWatchEvent) error {
	return r.save(ctx, event)
}

func (r *MediaWatchEventRepository) Update(ctx context.Context, event *ddownload.MediaWatchEvent) error {
	return r.save(ctx, event)
}

func (r *MediaWatchEventRepository) Write(ctx context.Context, event *ddownload.MediaWatchEvent) error {
	return r.save(ctx, event)
}

func (r *MediaWatchEventRepository) save(ctx context.Context, event *ddownload.MediaWatchEvent) error {
	if event == nil {
		return ierrors.ErrFuncParamNullPointer
	}

	// Convert the domain model to a database entity
	eEvent, err := r.mappers.MapMediaWatchEventDomainToEntity(event)
	if err != nil {
		return err
	}

	// Get the list of fields and values for insertion
	fields := eEvent.Fields()
	values := eEvent.Values()

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Insert(eEvent.TableName()).
		Columns(fields...).
		Values(values...).
		Suffix(dbutils.UpsertSuffix(fields, eEvent.FieldName(&eEvent.EventID))).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save record: %v", err)
	}

	return nil
}

func (r *MediaWatchEventRepository) Delete(ctx context.Context, downloadID uuid.UUID) error {
	var eEvent ewatchevent.MediaWatchEvent

	// Build DELETE query
	sqlBuilder := squirrel.
		Delete(eEvent.TableName()).
		Where(squirrel.Eq{eEvent.FieldName(&eEvent.DownloadID): downloadID}).
		PlaceholderFormat(squirrel.Dollar)

	// Generate SQL and args
	sqlStr, args, err := sqlBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlStr, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to delete record: %v", err)
	}

	return nil
}

func (r *MediaWatchEventRepository) WithUserID() persistence.MediaWatchEventRepository {
	rep := r.Copy()

	var eEvent ewatchevent.MediaWatchEvent

	rep.filtersByName[eEvent.FieldName(&eEvent.UserID)] = nil

	return rep
}

func (r *MediaWatchEventRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.Tx(ctx, r.db, r.lock, fn)
}

func (r *MediaWatchEventRepository) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.TxIndependent(ctx, r.db, r.lock, fn)
}
