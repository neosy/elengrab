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
	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	"github.com/neosy/elengrab/internal/repository/sqlite/types"
	ewatchevent "github.com/neosy/elengrab/internal/repository/sqlite/watch_event/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/watch_event/mappers"
)

type MediaWatchStatRepository struct {
	mappers *mappers.Mappers
	dbEntry persistence.DBEntry

	filtersByName types.FiltersByName

	// options
	retryOptions dbexec.RetryOptions
}

// NewMediaWatchStatRepository returns a new object for the repository
func NewMediaWatchStatRepository(dbEntry persistence.DBEntry) persistence.MediaWatchStatRepositoryFactory {
	return func() persistence.MediaWatchStatRepository {
		return &MediaWatchStatRepository{
			mappers: mappers.NewMappers(),
			dbEntry: dbEntry,

			filtersByName: make(map[string]any),

			// options
			retryOptions: dbexec.RetryOptions{
				MaxRetries: maxRetriesDefault,
				Delay:      retryDelayDefault,
			},
		}
	}
}

func (r *MediaWatchStatRepository) Insert(ctx context.Context, stat *ddownload.MediaWatchStat) error {
	return r.save(ctx, stat)
}

func (r *MediaWatchStatRepository) Update(ctx context.Context, stat *ddownload.MediaWatchStat) error {
	return r.save(ctx, stat)
}

func (r *MediaWatchStatRepository) Write(ctx context.Context, stat *ddownload.MediaWatchStat) error {
	return r.save(ctx, stat)
}

func (r *MediaWatchStatRepository) save(ctx context.Context, stat *ddownload.MediaWatchStat) error {
	if stat == nil {
		return ierrors.ErrFuncParamNullPointer
	}

	// Convert the domain model to a database entity
	eStat, err := r.mappers.MapMediaWatchStatDomainToEntity(stat)
	if err != nil {
		return err
	}

	// Get the list of fields and values for insertion
	fields := eStat.Fields()
	values := eStat.Values()

	// Build INSERT query
	sqlBuilder := squirrel.
		Insert(eStat.TableName()).
		Columns(fields...).
		Values(values...).
		Suffix(dbutils.UpsertSuffix(fields, eStat.FieldName(&eStat.DownloadID))).
		PlaceholderFormat(squirrel.Dollar)

	// Generate SQL query with upsert logic
	sqlQuery, args, err := sqlBuilder.ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.dbEntry, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save record: %v", err)
	}

	return nil
}

func (r *MediaWatchStatRepository) DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error {
	var eStat ewatchevent.MediaWatchStat

	// Build DELETE query
	sqlBuilder := squirrel.
		Delete(eStat.TableName()).
		Where(squirrel.Eq{eStat.FieldName(&eStat.DownloadID): downloadID}).
		PlaceholderFormat(squirrel.Dollar)

	// Generate SQL and args
	sqlQuery, args, err := sqlBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.dbEntry, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to delete record: %v", err)
	}

	return nil
}

func (r *MediaWatchStatRepository) DeleteAll(ctx context.Context) error {
	var eStat ewatchevent.MediaWatchStat

	// Build DELETE query
	sqlBuilder := squirrel.
		Delete(eStat.TableName()).
		PlaceholderFormat(squirrel.Dollar)

	// Generate SQL and args
	sqlQuery, args, err := sqlBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.dbEntry, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to delete record: %v", err)
	}

	return nil
}

func (r *MediaWatchStatRepository) Find(ctx context.Context, downloadID uuid.UUID) (*ddownload.MediaWatchStat, error) {
	var eStat ewatchevent.MediaWatchStat

	// Build SELECT query
	sqlBuilder := squirrel.
		Select(eStat.FieldsAll()...).
		From(eStat.TableName()).
		Where(squirrel.Eq{eStat.FieldName(&eStat.DownloadID): downloadID}).
		PlaceholderFormat(squirrel.Dollar).
		Limit(1)

	// Generate SQL and args
	sqlQuery, args, err := sqlBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	var notFound bool
	db := dbexec.Resolve(ctx, r.dbEntry)
	execQuery := func() error {
		row := db.QueryRowContext(ctx, sqlQuery, args...)
		// Scan result into entity
		err := row.Scan(eStat.FieldPointers()...)
		if err == sql.ErrNoRows {
			notFound = true
			return nil
		}
		return err
	}
	err = dbexec.ExecRetry(ctx, r.retryOptions, execQuery)
	if err != nil {
		return nil, err
	}
	if notFound {
		return nil, nil
	}

	// Map entity to domain model
	stat, err := r.mappers.MapMediaWatchStatEntityToDomain(&eStat)
	if err != nil {
		return nil, err
	}

	return stat, nil
}

func (r *MediaWatchStatRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.Tx(ctx, r.dbEntry, fn)
}

func (r *MediaWatchStatRepository) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.TxIndependent(ctx, r.dbEntry, fn)
}
