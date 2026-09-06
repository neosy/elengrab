package watchevent

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	"github.com/neosy/elengrab/internal/repository/sqlite/types"
	ewatchevent "github.com/neosy/elengrab/internal/repository/sqlite/watch_event/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/watch_event/mappers"
)

type MediaUserWatchChunkRepository struct {
	mappers *mappers.Mappers
	dbEntry persistence.DBEntry

	filtersByName types.FiltersByName
	queryOptions  dtypes.QueryOptions

	// options
	retryOptions dbexec.RetryOptions
}

// NewMediaUserWatchChunkRepository returns a new object for the repository
func NewMediaUserWatchChunkRepository(dbEntry persistence.DBEntry) persistence.MediaUserWatchChunkRepositoryFactory {
	return func() persistence.MediaUserWatchChunkRepository {
		return &MediaUserWatchChunkRepository{
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

func (r *MediaUserWatchChunkRepository) AddChunkQty(ctx context.Context, chunk *ddownload.MediaUserWatchChunk) error {
	if chunk == nil {
		return ierrors.ErrFuncParamNullPointer
	}

	// Convert the domain model to a database entity
	eChunk, err := r.mappers.MapMediaUserWatchChunkDomainToEntity(chunk)
	if err != nil {
		return err
	}

	// Get the list of fields and values for insertion
	fields := eChunk.InsertFields()
	values := eChunk.InsertValues()

	qtyFieldName := eChunk.FieldName(&eChunk.Qty)

	// Build INSERT query
	sqlBuilder := squirrel.
		Insert(eChunk.TableName()).
		Columns(fields...).
		Values(values...).
		Suffix("ON CONFLICT (" + eChunk.ConflictColumnsSQL() + ") DO UPDATE SET " + qtyFieldName + " = " + qtyFieldName + " + EXCLUDED." + qtyFieldName).
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
		return fmt.Errorf("failed to save record: %w", err)
	}

	return nil
}

func (r *MediaUserWatchChunkRepository) AddChunkQtyBatch(ctx context.Context, chunks []*ddownload.MediaUserWatchChunk) error {
	if chunks == nil {
		return ierrors.ErrFuncParamNullPointer
	}

	if len(chunks) == 0 {
		return nil
	}

	var eChunkTmpl ewatchevent.MediaUserWatchChunk

	fields := eChunkTmpl.InsertFields()

	// Build INSERT query
	sqlBuilder := squirrel.
		Insert(eChunkTmpl.TableName()).
		Columns(fields...)

	for _, chunk := range chunks {
		if chunk == nil {
			return ierrors.ErrFuncParamNullPointer
		}

		// Convert the domain model to a database entity
		eChunk, err := r.mappers.MapMediaUserWatchChunkDomainToEntity(chunk)
		if err != nil {
			return err
		}

		sqlBuilder = sqlBuilder.Values(eChunk.InsertValues()...)
	}

	qtyFieldName := eChunkTmpl.FieldName(&eChunkTmpl.Qty)

	sqlBuilder = sqlBuilder.
		Suffix("ON CONFLICT (" + eChunkTmpl.ConflictColumnsSQL() + ") DO UPDATE SET " + qtyFieldName + " = " + qtyFieldName + " + EXCLUDED." + qtyFieldName).
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
		return fmt.Errorf("failed to save records: %w", err)
	}

	return nil
}

func (r *MediaUserWatchChunkRepository) DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error {
	var eChunk ewatchevent.MediaUserWatchChunk

	// Build DELETE query
	sqlBuilder := squirrel.
		Delete(eChunk.TableName()).
		Where(squirrel.Eq{eChunk.FieldName(&eChunk.DownloadID): downloadID}).
		PlaceholderFormat(squirrel.Dollar)

	// Generate SQL and args
	sqlStr, args, err := sqlBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error generating SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.dbEntry, sqlStr, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}

	return nil
}

func (r *MediaUserWatchChunkRepository) DeleteAll(ctx context.Context) error {
	var eChunk ewatchevent.MediaUserWatchChunk

	// Build DELETE query
	sqlBuilder := squirrel.
		Delete(eChunk.TableName()).
		PlaceholderFormat(squirrel.Dollar)

	// Generate SQL and args
	sqlStr, args, err := sqlBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error generating SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.dbEntry, sqlStr, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to delete all records: %w", err)
	}

	return nil
}

func (r *MediaUserWatchChunkRepository) IterateDownloadUsers(
	ctx context.Context,
	fn func(downloadID, userID uuid.UUID,
	) error) error {
	var eChunk ewatchevent.MediaUserWatchChunk

	var sqlWhere = squirrel.And{}
	for name, value := range r.filtersByName {
		if name != "" {
			sqlWhere = append(sqlWhere, squirrel.Eq{
				eChunk.FieldName(eChunk.FieldPointer(name)): value,
			})
		}
	}

	selectFields := []string{
		eChunk.FieldName(&eChunk.DownloadID),
		eChunk.FieldName(&eChunk.UserID),
	}

	qb := squirrel.
		Select(selectFields...).
		From(eChunk.TableName()).
		Where(sqlWhere).
		GroupBy(selectFields...).
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
		var dID, uID uuid.UUID

		for rows.Next() {
			err := rows.Scan(&dID, &uID)
			if err != nil {
				return err
			}

			err = fn(dID, uID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

const countViewsSQL = `
WITH ranked AS (
    SELECT
        user_id,
        qty,
        ROW_NUMBER() OVER (
            PARTITION BY user_id
            ORDER BY qty DESC
        ) AS rn
    FROM media_user_watch_chunks
    WHERE download_id = ?
)
SELECT COALESCE(SUM(views), 0)
FROM (
    SELECT MIN(qty) AS views
    FROM ranked
    WHERE rn <= ?
    GROUP BY user_id
    HAVING COUNT(*) = ?
);
`

func (r *MediaUserWatchChunkRepository) CountViews(
	ctx context.Context,
	downloadID uuid.UUID,
	requiredChunks uint32,
) (uint32, error) {
	var (
		views    uint32
		notFound bool
	)

	// Execute the query
	db := dbexec.Resolve(ctx, r.dbEntry)
	execQuery := func() error {
		row := db.QueryRowContext(
			ctx, countViewsSQL,
			downloadID,
			requiredChunks, requiredChunks,
		)
		// Scan result into entity
		err := row.Scan(&views)
		if err == sql.ErrNoRows {
			notFound = true
			return nil
		}
		return err
	}
	err := dbexec.ExecRetry(ctx, r.retryOptions, execQuery)
	if err != nil {
		return 0, err
	}
	if notFound {
		return 0, nil
	}

	return views, nil
}

const countUserViewsSQL = `
WITH ranked AS (
    SELECT
        user_id,
        qty,
        ROW_NUMBER() OVER (
            PARTITION BY user_id
            ORDER BY qty DESC
        ) AS rn
    FROM media_user_watch_chunks
    WHERE download_id = ?
		AND user_id = ?
)
SELECT COALESCE(SUM(views), 0)
FROM (
    SELECT MIN(qty) AS views
    FROM ranked
    WHERE rn <= ?
    HAVING COUNT(*) = ?
);
`

func (r *MediaUserWatchChunkRepository) CountUserViews(
	ctx context.Context,
	downloadID uuid.UUID, userID uuid.UUID,
	requiredChunks uint32,
) (uint32, error) {
	var (
		views    uint32
		notFound bool
	)

	// Execute the query
	db := dbexec.Resolve(ctx, r.dbEntry)
	execQuery := func() error {
		row := db.QueryRowContext(
			ctx, countUserViewsSQL,
			downloadID, userID,
			requiredChunks, requiredChunks,
		)
		// Scan result into entity
		err := row.Scan(&views)
		if err == sql.ErrNoRows {
			notFound = true
			return nil
		}
		return err
	}
	err := dbexec.ExecRetry(ctx, r.retryOptions, execQuery)
	if err != nil {
		return 0, err
	}
	if notFound {
		return 0, nil
	}

	return views, nil
}

func (r *MediaUserWatchChunkRepository) WithUserID() persistence.MediaUserWatchChunkRepository {
	var eChunk ewatchevent.MediaUserWatchChunk

	r.filtersByName[eChunk.FieldName(&eChunk.UserID)] = nil

	return r
}

func (r *MediaUserWatchChunkRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.Tx(ctx, r.dbEntry, fn)
}

func (r *MediaUserWatchChunkRepository) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.TxIndependent(ctx, r.dbEntry, fn)
}
