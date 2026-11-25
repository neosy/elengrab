package searchindex

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"uuid"

	"github.com/Masterminds/squirrel"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/pkg/dbutils"
	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	esearchindex "github.com/neosy/elengrab/internal/repository/sqlite/search_index/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/search_index/mappers"
)

type MediaSourceIndexRepository struct {
	mappers *mappers.Mappers
	dbEntry persistence.DBEntry

	queryOptions queryOptions

	// options
	retryOptions dbexec.RetryOptions
}

// NewMediaSourceIndexRepository returns a new object for the repository
func NewMediaSourceIndexRepository(dbEntry persistence.DBEntry) persistence.MediaSourceIndexRepositoryFactory {
	return func() persistence.MediaSourceIndexRepository {
		return &MediaSourceIndexRepository{
			mappers: mappers.NewMappers(),
			dbEntry: dbEntry,

			queryOptions: newQueryOptions(),

			// options
			retryOptions: dbexec.RetryOptions{
				MaxRetries: maxRetriesDefault,
				Delay:      retryDelayDefault,
			},
		}
	}
}

func (r *MediaSourceIndexRepository) Insert(ctx context.Context, index *ddownload.MediaSourceIndex) error {
	return r.Save(ctx, index)
}

func (r *MediaSourceIndexRepository) Update(ctx context.Context, index *ddownload.MediaSourceIndex) error {
	return r.Save(ctx, index)
}

func (r *MediaSourceIndexRepository) Save(ctx context.Context, index *ddownload.MediaSourceIndex) error {
	if index == nil {
		return ierrors.ErrFuncParamNullPointer
	}

	// Convert the domain model to a database entity
	eIndex, err := r.mappers.MapMediaSourceIndexDomainToEntity(index)
	if err != nil {
		return err
	}

	// Get the list of fields and values for insertion
	fields := eIndex.InsertFields()
	values := eIndex.InsertValues()

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Insert(eIndex.TableName()).
		Columns(fields...).
		Values(values...).
		Suffix(dbutils.UpsertSuffix(fields, eIndex.FieldName(&eIndex.DownloadID))).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.dbEntry, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save media download: %v", err)
	}

	return nil
}

func (r *MediaSourceIndexRepository) UpdateOwner(ctx context.Context, fromID, toID uuid.UUID) error {
	var eIndex esearchindex.MediaSourceIndex

	sqlWhere := squirrel.Eq{eIndex.FieldName(&eIndex.UserID): fromID.String()}

	// Build query
	sqlBuilder := squirrel.Update(eIndex.TableName()).
		Set(eIndex.FieldName(&eIndex.UserID), toID.String()).
		Where(sqlWhere).
		PlaceholderFormat(squirrel.Dollar)

	// Generate SQL and args
	sqlQuery, args, err := sqlBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.dbEntry, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to update media download: %v", err)
	}

	return nil
}

func (r *MediaSourceIndexRepository) SoftDelete(ctx context.Context, downloadID uuid.UUID) error {
	var eIndex esearchindex.MediaSourceIndex

	fieldsToUpdate := map[string]any{
		eIndex.FieldName(&eIndex.DeletedAt): squirrel.Expr("CURRENT_TIMESTAMP"),
	}

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Update(eIndex.TableName()).
		SetMap(fieldsToUpdate).
		Where(squirrel.Eq{eIndex.FieldName(&eIndex.DownloadID): downloadID.String()}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.dbEntry, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save: %v", err)
	}

	return nil
}

func (r *MediaSourceIndexRepository) HardDelete(ctx context.Context, downloadID uuid.UUID) error {
	var ent esearchindex.MediaSourceIndex

	// Build DELETE query
	sqlBuilder := squirrel.
		Delete(ent.TableName()).
		Where(squirrel.Eq{ent.FieldName(&ent.DownloadID): downloadID.String()}).
		PlaceholderFormat(squirrel.Dollar)

	// Generate SQL and args
	sqlStr, args, err := sqlBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.dbEntry, sqlStr, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to delete: %v", err)
	}

	return nil
}

func (r *MediaSourceIndexRepository) Restore(ctx context.Context, downloadID uuid.UUID) error {
	var eIndex esearchindex.MediaSourceIndex

	fieldsToUpdate := map[string]any{
		eIndex.FieldName(&eIndex.DeletedAt): nil,
	}

	sqlWhere := squirrel.And{
		squirrel.Eq{eIndex.FieldName(&eIndex.DownloadID): downloadID.String()},
		squirrel.NotEq{eIndex.FieldName(&eIndex.DeletedAt): nil},
	}

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Update(eIndex.TableName()).
		SetMap(fieldsToUpdate).
		Where(sqlWhere).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.dbEntry, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save: %v", err)
	}

	return nil
}

func (r *MediaSourceIndexRepository) FindByDownloadID(ctx context.Context, downloadID uuid.UUID) (*ddownload.MediaSourceIndex, error) {
	var eIndex esearchindex.MediaSourceIndex

	sqlWhere := squirrel.And{}

	sqlWhere = append(sqlWhere,
		squirrel.Eq{
			eIndex.FieldName(&eIndex.DownloadID): downloadID.String(),
			eIndex.FieldName(&eIndex.DeletedAt):  nil,
		},
	)

	sqlWhere = r.applySearchFilter(sqlWhere)

	sqlQuery, args, err := squirrel.Select(eIndex.QueryFields()...).
		From(eIndex.TableName()).
		Where(sqlWhere).
		PlaceholderFormat(squirrel.Dollar).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	var notFound bool
	db := dbexec.Resolve(ctx, r.dbEntry)
	execQuery := func() error {
		row := db.QueryRowContext(ctx, sqlQuery, args...)
		// Scan result into entity
		err := row.Scan(eIndex.FieldPointers()...)
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
	index, err := r.mappers.MapSourceIndexEntityToDomain(&eIndex)
	if err != nil {
		return nil, err
	}

	return index, nil
}

func (r *MediaSourceIndexRepository) iterateAll(
	ctx context.Context,
	fn func(*ddownload.MediaSourceIndex) error,
) error {
	var (
		eIndex     esearchindex.MediaSourceIndex
		conditions = squirrel.And{}
	)

	conditions, err := r.applyFilters(conditions)
	if err != nil {
		return err
	}

	conditions = r.applyVisibilityFilter(conditions)
	conditions = r.applySearchFilter(conditions)
	conditions = r.applyContentTypeFilter(conditions)
	conditions = r.applyInProgressOwnerFilter(conditions)
	conditions = r.applyDeletedFilter(conditions)

	var orderBys dbutils.OrderByList

	switch r.queryOptions.ViewMode {
	case dtypes.QueryMediaViewModeNew:
		orderBys = dbutils.SortBy(eIndex.FieldName(&eIndex.SourceCreatedAt), dbutils.OrderDescending).List()
		orderBys.Add(eIndex.FieldName(&eIndex.RowID), dbutils.OrderDescending)
		if !r.queryOptions.LastRecord.CreatedAt.IsZero() {
			lastCreatedAt := r.queryOptions.LastRecord.CreatedAt.UTC()
			conditions = append(conditions, squirrel.Lt{eIndex.FieldName(&eIndex.SourceCreatedAt): lastCreatedAt})
		}
	case dtypes.QueryMediaViewModeOld:
		orderBys = dbutils.SortBy(eIndex.FieldName(&eIndex.SourceCreatedAt), dbutils.OrderAscending).List()
		orderBys.Add(eIndex.FieldName(&eIndex.RowID), dbutils.OrderAscending)
		if !r.queryOptions.LastRecord.CreatedAt.IsZero() {
			lastCreatedAt := r.queryOptions.LastRecord.CreatedAt.UTC()
			conditions = append(conditions, squirrel.Gt{eIndex.FieldName(&eIndex.SourceCreatedAt): lastCreatedAt})
		}
	case dtypes.QueryMediaViewModePopular:
		orderBys = dbutils.SortBy(eIndex.FieldName(&eIndex.Views), dbutils.OrderDescending).List()
		orderBys.Add(eIndex.FieldName(&eIndex.SourceCreatedAt), dbutils.OrderDescending)
		orderBys.Add(eIndex.FieldName(&eIndex.DownloadID), dbutils.OrderDescending)
		if !r.queryOptions.LastRecord.CreatedAt.IsZero() {
			lastViews := r.queryOptions.LastRecord.Views
			lastCreatedAt := r.queryOptions.LastRecord.CreatedAt.UTC()
			lastDownloadID := r.queryOptions.LastRecord.ID

			condition := squirrel.Or{
				squirrel.Lt{eIndex.FieldName(&eIndex.Views): lastViews},
				squirrel.And{
					squirrel.Eq{eIndex.FieldName(&eIndex.Views): lastViews},
					squirrel.Lt{eIndex.FieldName(&eIndex.SourceCreatedAt): lastCreatedAt},
				},
				squirrel.And{
					squirrel.Eq{eIndex.FieldName(&eIndex.Views): lastViews},
					squirrel.Eq{eIndex.FieldName(&eIndex.SourceCreatedAt): lastCreatedAt},
					squirrel.Lt{eIndex.FieldName(&eIndex.DownloadID): lastDownloadID.String()},
				},
			}

			conditions = append(conditions, condition)
		}
	}

	sqlWhere := conditions

	if len(r.queryOptions.OrderBys) > 0 {
		orderBys = append(orderBys, r.queryOptions.OrderBys...)
	}

	if len(orderBys) == 0 {
		orderBys = dbutils.SortBy(eIndex.FieldName(&eIndex.SourceCreatedAt), dbutils.OrderDescending).List()
	}

	qb := squirrel.Select(eIndex.QueryFields()...).
		From(eIndex.TableName()).
		Where(sqlWhere).
		OrderBy(orderBys.Query()).
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
		for rows.Next() {
			err := rows.Scan(eIndex.FieldPointers()...)
			if err != nil {
				return err
			}

			index, err := r.mappers.MapSourceIndexEntityToDomain(&eIndex)
			if err != nil {
				return err
			}

			err = fn(index)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *MediaSourceIndexRepository) IterateAll(ctx context.Context, fn func(*ddownload.MediaSourceIndex) error) error {
	return r.iterateAll(ctx, fn)
}

func (r *MediaSourceIndexRepository) FillEmptyMediaTitleLower(ctx context.Context) error {
	var eIndex esearchindex.MediaSourceIndex

	sqlWhere := squirrel.And{
		squirrel.Eq{eIndex.FieldName(&eIndex.TitleLower): ""},
	}

	sqlQuery, args, err := squirrel.Select(eIndex.QueryFields()...).
		From(eIndex.TableName()).
		Where(sqlWhere).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

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

	var indexes []*ddownload.MediaSourceIndex
	if rows != nil {
		indexes, err = r.mappers.MapRowsToSourceIndexes(rows)
		if err != nil {
			return err
		}
	}

	for _, index := range indexes {
		sqlQuery, args, err := squirrel.
			Update(eIndex.TableName()).
			SetMap(map[string]any{
				eIndex.FieldName(&eIndex.TitleLower): strings.ToLower(index.Title),
			}).
			Where(squirrel.Eq{eIndex.FieldName(&eIndex.DownloadID): index.DownloadID.String()}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if err != nil {
			return fmt.Errorf("error generating SQL: %v", err)
		}

		err = dbexec.ExecContext(ctx, r.dbEntry, sqlQuery, args, r.retryOptions)
		if err != nil {
			return fmt.Errorf("failed to save mediaDownload: %v", err)
		}
	}

	return nil
}

func (r *MediaSourceIndexRepository) Count(ctx context.Context) (int, error) {
	var eIndex esearchindex.MediaSourceIndex

	sqlWhere := squirrel.And{}

	sqlWhere, err := r.applyFilters(sqlWhere)
	if err != nil {
		return 0, err
	}

	sqlWhere = r.applyVisibilityFilter(sqlWhere)
	sqlWhere = r.applySearchFilter(sqlWhere)
	sqlWhere = r.applyContentTypeFilter(sqlWhere)
	sqlWhere = r.applyDeletedFilter(sqlWhere)

	sqlQuery, args, err := squirrel.Select("COUNT(*)").
		From(eIndex.TableName()).
		Where(sqlWhere).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	var count int
	db := dbexec.Resolve(ctx, r.dbEntry)
	execQuery := func() error {
		err := db.QueryRowContext(ctx, sqlQuery, args...).Scan(&count)
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	err = dbexec.ExecRetry(ctx, r.retryOptions, execQuery)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *MediaSourceIndexRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.Tx(ctx, r.dbEntry, fn)
}

func (r *MediaSourceIndexRepository) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.TxIndependent(ctx, r.dbEntry, fn)
}
