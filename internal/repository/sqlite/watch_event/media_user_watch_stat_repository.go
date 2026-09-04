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
	ewatchevent "github.com/neosy/elengrab/internal/repository/sqlite/watch_event/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/watch_event/mappers"
)

type MediaUserWatchStatRepository struct {
	mappers *mappers.Mappers
	dbEntry persistence.DBEntry

	// options
	retryOptions dbexec.RetryOptions
}

// NewMediaUserWatchStatRepository returns a new object for the repository
func NewMediaUserWatchStatRepository(dbEntry persistence.DBEntry) persistence.MediaUserWatchStatRepositoryFactory {
	return func() persistence.MediaUserWatchStatRepository {
		return &MediaUserWatchStatRepository{
			mappers: mappers.NewMappers(),
			dbEntry: dbEntry,

			// options
			retryOptions: dbexec.RetryOptions{
				MaxRetries: maxRetriesDefault,
				Delay:      retryDelayDefault,
			},
		}
	}
}

func (r *MediaUserWatchStatRepository) Insert(ctx context.Context, stat *ddownload.MediaUserWatchStat) error {
	return r.save(ctx, stat)
}

func (r *MediaUserWatchStatRepository) Update(ctx context.Context, stat *ddownload.MediaUserWatchStat) error {
	return r.save(ctx, stat)
}

func (r *MediaUserWatchStatRepository) Write(ctx context.Context, stat *ddownload.MediaUserWatchStat) error {
	return r.save(ctx, stat)
}

func (r *MediaUserWatchStatRepository) save(ctx context.Context, stat *ddownload.MediaUserWatchStat) error {
	if stat == nil {
		return ierrors.ErrFuncParamNullPointer
	}

	// Convert the domain model to a database entity
	eStat, err := r.mappers.MapMediaUserWatchStatDomainToEntity(stat)
	if err != nil {
		return err
	}

	// Get the list of fields and values for insertion
	fields := eStat.InsertFields()
	values := eStat.InsertValues()

	// Build INSERT query
	sqlBuilder := squirrel.
		Insert(eStat.TableName()).
		Columns(fields...).
		Values(values...).
		Suffix(dbutils.UpsertSuffix(fields, eStat.ConflictFields()...)).
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

func (r *MediaUserWatchStatRepository) Delete(ctx context.Context, downloadID uuid.UUID, userID uuid.UUID) error {
	var eStat ewatchevent.MediaUserWatchStat

	// Build DELETE query
	sqlBuilder := squirrel.
		Delete(eStat.TableName()).
		Where(squirrel.Eq{
			eStat.FieldName(&eStat.DownloadID): downloadID,
			eStat.FieldName(&eStat.UserID):     userID,
		}).
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

func (r *MediaUserWatchStatRepository) DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error {
	var eStat ewatchevent.MediaUserWatchStat

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

func (r *MediaUserWatchStatRepository) DeleteAll(ctx context.Context) error {
	var eStat ewatchevent.MediaUserWatchStat

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

func (r *MediaUserWatchStatRepository) Find(ctx context.Context, downloadID uuid.UUID, userID uuid.UUID) (*ddownload.MediaUserWatchStat, error) {
	var eStat ewatchevent.MediaUserWatchStat

	// Build SELECT query
	sqlBuilder := squirrel.
		Select(eStat.QueryFields()...).
		From(eStat.TableName()).
		Where(squirrel.Eq{
			eStat.FieldName(&eStat.DownloadID): downloadID,
			eStat.FieldName(&eStat.UserID):     userID,
		}).
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
	stat, err := r.mappers.MapMediaUserWatchStatEntityToDomain(&eStat)
	if err != nil {
		return nil, err
	}

	return stat, nil
}

func (r *MediaUserWatchStatRepository) Exists(
	ctx context.Context,
	downloadID uuid.UUID, userID uuid.UUID,
) (bool, error) {
	var eStat ewatchevent.MediaUserWatchStat

	// Build SELECT query
	sqlBuilder := squirrel.
		Select("1").
		From(eStat.TableName()).
		Where(squirrel.Eq{
			eStat.FieldName(&eStat.DownloadID): downloadID,
			eStat.FieldName(&eStat.UserID):     userID,
		}).
		PlaceholderFormat(squirrel.Dollar).
		Limit(1)

	// Generate SQL and args
	sqlQuery, args, err := sqlBuilder.ToSql()
	if err != nil {
		return false, fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	var exists bool
	db := dbexec.Resolve(ctx, r.dbEntry)
	execQuery := func() error {
		row := db.QueryRowContext(ctx, sqlQuery, args...)
		var dummy int
		err := row.Scan(&dummy)
		if err == sql.ErrNoRows {
			exists = false
			return nil
		}
		if err != nil {
			return err
		}
		exists = true
		return nil
	}
	err = dbexec.ExecRetry(ctx, r.retryOptions, execQuery)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *MediaUserWatchStatRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.Tx(ctx, r.dbEntry, fn)
}

func (r *MediaUserWatchStatRepository) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.TxIndependent(ctx, r.dbEntry, fn)
}
