package searchindex

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/pkg/dbutils"
	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	esearchindex "github.com/neosy/elengrab/internal/repository/sqlite/search_index/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/search_index/mappers"
	"github.com/neosy/elengrab/internal/repository/sqlite/sqlutil"
)

type MediaSourceIndexRepository struct {
	mappers *mappers.Mappers
	db      *sql.DB
	lock    dbexec.WriteLocker

	filtersByName filtersByName
	queryOptions  dtypes.QueryMediaOptions

	// options
	retryOptions dbexec.RetryOptions
}

// NewMediaSourceIndexRepository returns a new object for the repository
func NewMediaSourceIndexRepository(db *sql.DB, lock dbexec.WriteLocker) persistence.MediaSourceIndexRepositoryFactory {
	return func() persistence.MediaSourceIndexRepository {
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
}

func (r *MediaSourceIndexRepository) WithOptions(options dtypes.QueryMediaOptions) persistence.MediaSourceIndexRepository {
	if options.Before != nil {
		r.queryOptions.Before = options.Before
	}

	if options.Limit != nil {
		r.queryOptions.Limit = options.Limit
	}

	if options.Visibility != nil {
		r.queryOptions.Visibility = options.Visibility
	}

	return r
}

func (r *MediaSourceIndexRepository) WithUser(userID uuid.UUID) persistence.MediaSourceIndexRepository {
	var eIndex esearchindex.MediaSourceIndex

	if userID == uuid.Nil {
		r.queryOptions.Visibility = new(dtypes.QueryMediaVisibilityPublic)
	} else {
		r.filtersByName[eIndex.FieldName(&eIndex.UserID)] = userID
	}

	return r
}

func (r *MediaSourceIndexRepository) WithFilters(filters map[string]any) persistence.MediaSourceIndexRepository {
	var (
		eIndex esearchindex.MediaSourceIndex

		fieldNameByAllowedFilter = map[string]string{
			dtypes.QueryFilterNameUserID: eIndex.FieldName(&eIndex.UserID),
			dtypes.QueryFilterNameTitle:  eIndex.FieldName(&eIndex.TitleLower),
		}
	)

	for name, value := range filters {
		fieldName, exists := fieldNameByAllowedFilter[name]
		if exists {
			switch fieldName {
			case eIndex.FieldName(&eIndex.UserID):
				id, ok := value.(uuid.UUID)
				if !ok {
					continue
				}
				if id == uuid.Nil {
					r.queryOptions.Visibility = new(dtypes.QueryMediaVisibilityPublic)
					continue
				}
				r.filtersByName[fieldName] = value
			default:
				r.filtersByName[fieldName] = value
			}

		}
	}

	return r
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
	fields := eIndex.Fields()
	values := eIndex.Values()

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
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save media download: %v", err)
	}

	return nil
}

func (r *MediaSourceIndexRepository) UpdateOwner(ctx context.Context, fromID, toID uuid.UUID) error {
	var eIndex esearchindex.MediaSourceIndex

	sqlWhere := squirrel.Eq{eIndex.FieldName(&eIndex.UserID): fromID}

	// Build query
	sqlBuilder := squirrel.Update(eIndex.TableName()).
		Set(eIndex.FieldName(&eIndex.UserID), toID).
		Where(sqlWhere).
		PlaceholderFormat(squirrel.Dollar)

	// Generate SQL and args
	sqlQuery, args, err := sqlBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to update media download: %v", err)
	}

	return nil
}

func (r *MediaSourceIndexRepository) SoftDelete(ctx context.Context, downloadID uuid.UUID) error {
	var eIndex esearchindex.MediaSourceIndex

	fieldsToUpdate := map[string]interface{}{
		eIndex.FieldName(&eIndex.DeletedAt): squirrel.Expr("CURRENT_TIMESTAMP"),
	}

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Update(eIndex.TableName()).
		SetMap(fieldsToUpdate).
		Where(squirrel.Eq{eIndex.FieldName(&eIndex.DownloadID): downloadID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
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
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlStr, args, r.retryOptions)
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
		squirrel.Eq{eIndex.FieldName(&eIndex.DownloadID): downloadID},
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
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
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

	for name, value := range r.filtersByName {
		if name != "" {
			sqlWhere = append(sqlWhere, squirrel.Eq{eIndex.FieldName(eIndex.FieldPointer(name)): value})
		}
	}

	sqlQuery, args, err := squirrel.Select(eIndex.FieldsAll()...).
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
	db := dbexec.Resolve(ctx, r.db)
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

func (r *MediaSourceIndexRepository) iterateGetAll(
	ctx context.Context,
	sortOrderBy string,
	includeDeleted bool,
	fn func(*ddownload.MediaSourceIndex) error,
) error {
	var eIndex esearchindex.MediaSourceIndex

	var (
		filterUserID string
		conditions   = squirrel.And{}
	)
	for name, value := range r.filtersByName {
		switch name {
		case "":
			continue
		case eIndex.FieldName(&eIndex.TitleLower):
			filter, ok := value.(string)
			if !ok {
				return fmt.Errorf("expected string, got %T (%v)", value, value)
			}
			conditions = append(conditions, sqlutil.Like(eIndex.FieldName(&eIndex.TitleLower), filter))
		case eIndex.FieldName(&eIndex.UserID):
			userID, ok := value.(uuid.UUID)
			if !ok {
				return fmt.Errorf("expected uuid.UUID, got %T", value)
			}
			filterUserID = userID.String()
		default:
			conditions = append(conditions, squirrel.Eq{eIndex.FieldName(eIndex.FieldPointer(name)): value})
		}
	}

	if r.queryOptions.Visibility != nil {
		if *r.queryOptions.Visibility > dtypes.QueryMediaVisibilityAll {
			sqlOr := squirrel.Or{
				squirrel.Eq{eIndex.FieldName(&eIndex.UserID): nil},
				squirrel.Eq{eIndex.FieldName(&eIndex.UserID): uuid.Nil},
				squirrel.Eq{eIndex.FieldName(&eIndex.Visibility): dtypes.MediaVisibilityPublic.String()},
			}
			if filterUserID != "" && *r.queryOptions.Visibility == dtypes.QueryMediaVisibilityAuthenticated {
				sqlOr = append(sqlOr, squirrel.Eq{eIndex.FieldName(&eIndex.UserID): filterUserID})
				filterUserID = ""
			}
			conditions = append(conditions, sqlOr)
		}
	}

	if filterUserID != "" {
		conditions = append(conditions, squirrel.Eq{eIndex.FieldName(&eIndex.UserID): filterUserID})
	}

	if r.queryOptions.Before != nil && !r.queryOptions.Before.IsZero() {
		t := r.queryOptions.Before.Add(-1 * time.Nanosecond)
		conditions = append(conditions, squirrel.Lt{eIndex.FieldName(&eIndex.SourceCreatedAt): t})
	}
	if !includeDeleted {
		conditions = append(conditions, squirrel.Eq{eIndex.FieldName(&eIndex.DeletedAt): nil})
	}

	sqlWhere := conditions

	// Create an ORDER BY clause based on fieldы with the specified sort order.
	orderBy := dbutils.OrderBy(
		dbutils.Flds{
			eIndex.FieldName(&eIndex.SourceCreatedAt): sortOrderBy,
		})

	qb := squirrel.Select(eIndex.FieldsAll()...).
		From(eIndex.TableName()).
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
	db := dbexec.Resolve(ctx, r.db)
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

			download, err := r.mappers.MapSourceIndexEntityToDomain(&eIndex)
			if err != nil {
				return err
			}

			err = fn(download)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *MediaSourceIndexRepository) IterateGetAll(ctx context.Context, includeDeleted bool, fn func(*ddownload.MediaSourceIndex) error) error {
	return r.iterateGetAll(ctx, dbutils.OrderDesc, includeDeleted, fn)
}

func (r *MediaSourceIndexRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.Tx(ctx, r.db, r.lock, fn)
}

func (r *MediaSourceIndexRepository) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.TxIndependent(ctx, r.db, r.lock, fn)
}
