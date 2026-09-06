package watchevent

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/pkg/dbutils"
	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	"github.com/neosy/elengrab/internal/repository/sqlite/types"
	ewatchevent "github.com/neosy/elengrab/internal/repository/sqlite/watch_event/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/watch_event/mappers"
)

type MediaWatchEventRepository struct {
	mappers *mappers.Mappers
	dbEntry persistence.DBEntry

	filtersByName types.FiltersByName
	queryOptions  dtypes.QueryOptions

	// options
	retryOptions dbexec.RetryOptions
}

// NewMediaWatchEventRepository returns a new object for the repository
func NewMediaWatchEventRepository(dbEntry persistence.DBEntry) persistence.MediaWatchEventRepositoryFactory {
	return func() persistence.MediaWatchEventRepository {
		return &MediaWatchEventRepository{
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
	fields := eEvent.InsertFields()
	values := eEvent.InsertValues()

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
	err = dbexec.ExecContext(ctx, r.dbEntry, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save record: %v", err)
	}

	return nil
}

func (r *MediaWatchEventRepository) DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error {
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
	err = dbexec.ExecContext(ctx, r.dbEntry, sqlStr, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to delete record: %v", err)
	}

	return nil
}

func (r *MediaWatchEventRepository) IterateAll(ctx context.Context, fn func(*ddownload.MediaWatchEvent) error) error {
	return r.iterateGetAll(ctx, dbutils.OrderAsc, fn)
}

func (r *MediaWatchEventRepository) iterateGetAll(
	ctx context.Context,
	sortOrderBy string,
	fn func(*ddownload.MediaWatchEvent) error,
) error {
	var eEvent ewatchevent.MediaWatchEvent

	var sqlWhere = squirrel.And{}
	for name, value := range r.filtersByName {
		if name != "" {
			sqlWhere = append(sqlWhere, squirrel.Eq{eEvent.FieldName(eEvent.FieldPointer(name)): value})
		}
	}

	// Create an ORDER BY clause based on fieldы with the specified sort order.
	orderBy := dbutils.OrderBy(
		dbutils.Flds{
			eEvent.FieldName(&eEvent.CreatedAt): sortOrderBy,
		})

	qb := squirrel.
		Select(eEvent.QueryFields()...).
		From(eEvent.TableName()).
		Where(sqlWhere).
		OrderBy(orderBy).
		PlaceholderFormat(squirrel.Dollar)

	if r.queryOptions.Limit != nil && *r.queryOptions.Limit > 0 {
		qb = qb.Limit(*r.queryOptions.Limit)
	}

	sqlQuery, args, err := qb.ToSql()
	if err != nil {
		return fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	db := dbexec.Resolve(ctx, r.dbEntry)
	rows, err := db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows != nil {
		var eEvent ewatchevent.MediaWatchEvent

		for rows.Next() {
			err := rows.Scan(eEvent.FieldPointers()...)
			if err != nil {
				return err
			}

			event, err := r.mappers.MapMediaWatchEventEntityToDomain(&eEvent)
			if err != nil {
				return err
			}

			err = fn(event)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *MediaWatchEventRepository) WithUserID() persistence.MediaWatchEventRepository {
	var eEvent ewatchevent.MediaWatchEvent

	r.filtersByName[eEvent.FieldName(&eEvent.UserID)] = nil

	return r
}

func (r *MediaWatchEventRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.Tx(ctx, r.dbEntry, fn)
}

func (r *MediaWatchEventRepository) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.TxIndependent(ctx, r.dbEntry, fn)
}
