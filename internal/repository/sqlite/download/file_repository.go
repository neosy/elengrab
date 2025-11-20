package sldownload

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/download/mappers"
	"github.com/neosy/elengrab/pkg/dbutils"
)

type FileRepository struct {
	mappers *mappers.Mappers
	db      *sql.DB
	mu      *sync.RWMutex

	// options
	retryOptions retryOptions
}

type fileByFields struct {
	status     *dtypes.FileStatus
	beforeTime *time.Time
	limit      *uint64
}

// NewFileRepository returns a new object for the repository
func NewFileRepository(db *sql.DB, mu *sync.RWMutex) *FileRepository {
	return &FileRepository{
		mappers: mappers.NewMappers(),
		db:      db,
		mu:      mu,

		// options
		retryOptions: retryOptions{
			maxRetries: maxRetriesDefault,
			delay:      retryDelayDefault,
		},
	}
}

func (r *FileRepository) Insert(ctx context.Context, file *ddownload.File) error {
	return r.save(ctx, file, false)
}

func (r *FileRepository) Update(ctx context.Context, file *ddownload.File) error {
	return r.save(ctx, file, true)
}

func (r *FileRepository) save(ctx context.Context, file *ddownload.File, isUpd bool) error {
	if file == nil {
		return errors.New("function parameter is a null pointer")
	}

	// Convert the domain model to a database entity
	eFile := r.mappers.MapFileDomainToEntity(file)

	// Get the list of fields and values for insertion
	fields := eFile.Fields()
	values := eFile.Values()

	// If this is an update — add the UpdatedAt field with the current time
	if isUpd {
		fields = append(fields, eFile.FieldName(&eFile.UpdatedAt))
		values = append(values, time.Now())
	}

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Insert(eFile.TableName()).
		Columns(fields...).
		Values(values...).
		Suffix(dbutils.UpsertSuffix(fields, eFile.FieldName(&eFile.FileId))).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Execute the query
	err = execContext(ctx, r.db, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	return nil
}

func (r *FileRepository) UpdateStatusToNew(ctx context.Context) error {
	var ent edownload.File

	sqlWhere := squirrel.Or{
		squirrel.Eq{ent.FieldName(&ent.Status): dtypes.FileStatusPending},
		squirrel.Eq{ent.FieldName(&ent.Status): dtypes.FileStatusWorking},
	}

	// Build query
	sqlBuilder := squirrel.Update(ent.TableName()).
		Set(ent.FieldName(&ent.Status), dtypes.FileStatusNew).
		Where(sqlWhere).
		PlaceholderFormat(squirrel.Dollar)

	// Generate SQL and args
	sqlQuery, args, err := sqlBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error generating SQL: %v", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Execute the query
	err = execContext(ctx, r.db, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to update file: %v", err)
	}

	return nil
}

func (r *FileRepository) Delete(ctx context.Context, fileId uuid.UUID) error {
	var ent edownload.File

	// Build DELETE query
	sqlBuilder := squirrel.Delete(ent.TableName()).
		Where(squirrel.Eq{ent.FieldName(&ent.FileId): fileId.String()}).
		PlaceholderFormat(squirrel.Dollar)

	// Generate SQL and args
	sqlStr, args, err := sqlBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error generating SQL: %v", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Execute the query
	err = execContext(ctx, r.db, sqlStr, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

func (r *FileRepository) FindByFileId(ctx context.Context, fileId uuid.UUID) (*ddownload.File, error) {
	var (
		eFile edownload.File
		eTask edownload.DownloadTask

		aliasFiles = "f"
		aliasTasks = "t"
	)

	selectFields := append(eFile.FieldsAllWithAlias(aliasFiles), eTask.FieldsAllWithAlias(aliasTasks)...)

	sqlQuery, args, err := squirrel.Select(selectFields...).
		From(eFile.TableName() + " AS " + aliasFiles).
		LeftJoin(
			eTask.TableName() + " AS " + aliasTasks +
				" ON " + aliasTasks + "." + eTask.FieldName(&eTask.FileId) +
				" = " + aliasFiles + "." + eFile.FieldName(&eFile.FileId),
		).
		Where(squirrel.Eq{eFile.FieldNameWithAlias(&eFile.FileId, aliasFiles): fileId.String()}).
		PlaceholderFormat(squirrel.Dollar).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	row := r.db.QueryRowContext(ctx, sqlQuery, args...)

	// Scan result into entity
	if err := row.Scan(append(eFile.FieldPointers(), eTask.FieldPointers()...)...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // запись не найдена
		}
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	// Map entity to domain model
	file, err := r.mappers.MapFileEntityToDomain(&eFile, &eTask)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (r *FileRepository) getAll(ctx context.Context, byFields fileByFields, sortOrder string) ([]*ddownload.File, error) {
	var (
		eFile edownload.File
		eTask edownload.DownloadTask
		files []*ddownload.File

		aliasFiles = "f"
		aliasTasks = "t"
	)

	selectFields := append(eFile.FieldsAllWithAlias(aliasFiles), eTask.FieldsAllWithAlias(aliasTasks)...)

	var conditions = squirrel.And{}
	if byFields.status != nil {
		conditions = append(conditions, squirrel.Eq{eFile.FieldNameWithAlias(&eFile.Status, aliasFiles): byFields.status.String()})
	}
	if byFields.beforeTime != nil && !byFields.beforeTime.IsZero() {
		t := byFields.beforeTime.Add(-1 * time.Nanosecond)
		conditions = append(conditions, squirrel.Lt{eFile.FieldNameWithAlias(&eFile.CreatedAt, aliasFiles): t})
	}

	sqlWhere := squirrel.Expr("TRUE")
	if len(conditions) > 0 {
		sqlWhere = conditions
	}

	qb := squirrel.Select(selectFields...).
		From(eFile.TableName() + " AS " + aliasFiles).
		Where(sqlWhere).
		OrderBy(dbutils.OrderBy(dbutils.Flds{eFile.FieldNameWithAlias(&eFile.CreatedAt, aliasFiles): sortOrder})).
		LeftJoin(
			eTask.TableName() + " AS " + aliasTasks +
				" ON " + aliasTasks + "." + eTask.FieldName(&eTask.FileId) +
				" = " + aliasFiles + "." + eFile.FieldName(&eFile.FileId),
		).
		PlaceholderFormat(squirrel.Dollar)

	if byFields.limit != nil && *byFields.limit > 0 {
		qb = qb.Limit(*byFields.limit)
	}

	sqlQuery, args, err := qb.ToSql()

	if err != nil {
		return nil, fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	rows, err := r.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return files, nil // ничего не найдено
		}
		return nil, err
	}
	defer rows.Close()

	if rows != nil {
		files, err = r.mappers.MapRowsToFiles(rows)
		if err != nil {
			return nil, err
		}
	}

	return files, nil
}

func (r *FileRepository) GetAll(ctx context.Context) ([]*ddownload.File, error) {
	return r.getAll(ctx, fileByFields{}, dbutils.OrderDesc)
}

func (r *FileRepository) GetBeforeTime(ctx context.Context, before time.Time, limit uint64) ([]*ddownload.File, error) {
	return r.getAll(
		ctx,
		fileByFields{
			beforeTime: &before,
			limit:      &limit,
		},
		dbutils.OrderDesc,
	)
}

func (r *FileRepository) GetByStatus(ctx context.Context, status dtypes.FileStatus) ([]*ddownload.File, error) {
	return r.getAll(ctx, fileByFields{status: &status}, dbutils.OrderAsc)
}

func (r *FileRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(ctxWithTx(ctx, tx)); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
