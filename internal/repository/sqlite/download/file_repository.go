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
	statuses    []dtypes.FileStatus
	beforeTime  *time.Time
	limit       *uint64
	partialHash **string
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
	eFile, err := r.mappers.MapFileDomainToEntity(file)
	if err != nil {
		return err
	}

	// Get the list of fields and values for insertion
	fields := eFile.Fields()
	values := eFile.Values()

	// If this is an update — add the UpdatedAt field with the current time
	if isUpd {
		fields = append(fields, eFile.FieldName(&eFile.UpdatedAt))
		values = append(values, squirrel.Expr("CURRENT_TIMESTAMP"))
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

	// Execute the query
	err = execContext(ctx, r.db, r.mu, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	return nil
}

func (r *FileRepository) UpdateStatusToNew(ctx context.Context, statuses []dtypes.FileStatus) error {
	if len(statuses) == 0 {
		return nil
	}

	var ent edownload.File

	sqlWhere := squirrel.And{
		squirrel.Eq{ent.FieldName(&ent.Status): statuses},
		squirrel.Eq{ent.FieldName(&ent.DeletedAt): nil},
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

	// Execute the query
	err = execContext(ctx, r.db, r.mu, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to update file: %v", err)
	}

	return nil
}

func (r *FileRepository) Delete(ctx context.Context, fileId uuid.UUID, soft bool) error {
	if soft {
		return r.softDelete(ctx, fileId)
	} else {
		return r.hardDelete(ctx, fileId)
	}
}

func (r *FileRepository) softDelete(ctx context.Context, fileId uuid.UUID) error {
	var eFile edownload.File

	fieldsToUpdate := map[string]interface{}{
		eFile.FieldName(&eFile.UpdatedAt): squirrel.Expr("CURRENT_TIMESTAMP"),
		eFile.FieldName(&eFile.DeletedAt): squirrel.Expr("CURRENT_TIMESTAMP"),
	}

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Update(eFile.TableName()).
		SetMap(fieldsToUpdate).
		Where(squirrel.Eq{eFile.FieldName(&eFile.FileId): fileId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = execContext(ctx, r.db, r.mu, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	return nil
}

func (r *FileRepository) hardDelete(ctx context.Context, fileId uuid.UUID) error {
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

	// Execute the query
	err = execContext(ctx, r.db, r.mu, sqlStr, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

func (r *FileRepository) Restore(ctx context.Context, fileId uuid.UUID) error {
	var eFile edownload.File

	fieldsToUpdate := map[string]any{
		eFile.FieldName(&eFile.UpdatedAt): squirrel.Expr("CURRENT_TIMESTAMP"),
		eFile.FieldName(&eFile.DeletedAt): nil,
	}

	sqlWhere := squirrel.And{
		squirrel.Eq{eFile.FieldName(&eFile.FileId): fileId},
		squirrel.NotEq{eFile.FieldName(&eFile.DeletedAt): nil},
	}

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Update(eFile.TableName()).
		SetMap(fieldsToUpdate).
		Where(sqlWhere).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = execContext(ctx, r.db, r.mu, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
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

	sqlWhere := squirrel.Eq{
		eFile.FieldNameWithAlias(&eFile.FileId, aliasFiles):    fileId.String(),
		eFile.FieldNameWithAlias(&eFile.DeletedAt, aliasFiles): nil,
	}

	sqlQuery, args, err := squirrel.Select(selectFields...).
		From(eFile.TableName() + " AS " + aliasFiles).
		LeftJoin(
			eTask.TableName() + " AS " + aliasTasks +
				" ON " + aliasTasks + "." + eTask.FieldName(&eTask.FileId) +
				" = " + aliasFiles + "." + eFile.FieldName(&eFile.FileId),
		).
		Where(sqlWhere).
		PlaceholderFormat(squirrel.Dollar).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	db := dbOrTx(ctx, r.db)
	row := db.QueryRowContext(ctx, sqlQuery, args...)

	// Scan result into entity
	if err := row.Scan(append(eFile.FieldPointers(), eTask.FieldPointers()...)...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
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

func (r *FileRepository) getAll(
	ctx context.Context,
	byFields fileByFields,
	sortOrderByCreatedAt string,
	includeDeleted bool,
) ([]*ddownload.File, error) {
	var (
		eFile edownload.File
		eTask edownload.DownloadTask
		files []*ddownload.File

		aliasFiles = "f"
		aliasTasks = "t"
	)

	selectFields := append(eFile.FieldsAllWithAlias(aliasFiles), eTask.FieldsAllWithAlias(aliasTasks)...)

	var conditions = squirrel.And{}
	if len(byFields.statuses) > 0 {
		conditions = append(conditions, squirrel.Eq{eFile.FieldNameWithAlias(&eFile.Status, aliasFiles): byFields.statuses})
	}
	if byFields.beforeTime != nil && !byFields.beforeTime.IsZero() {
		t := byFields.beforeTime.Add(-1 * time.Nanosecond)
		conditions = append(conditions, squirrel.Lt{eFile.FieldNameWithAlias(&eFile.CreatedAt, aliasFiles): t})
	}
	if byFields.partialHash != nil {
		if *byFields.partialHash == nil {
			conditions = append(conditions, squirrel.Expr(eFile.FieldNameWithAlias(&eFile.PartialHash, aliasFiles)+" IS NULL"))
		} else {
			conditions = append(conditions, squirrel.Eq{eFile.FieldNameWithAlias(&eFile.PartialHash, aliasFiles): **byFields.partialHash})
		}
	}
	if !includeDeleted {
		conditions = append(conditions, squirrel.Eq{eFile.FieldNameWithAlias(&eFile.DeletedAt, aliasFiles): nil})
	}

	sqlWhere := conditions

	qb := squirrel.Select(selectFields...).
		From(eFile.TableName() + " AS " + aliasFiles).
		Where(sqlWhere).
		OrderBy(dbutils.OrderBy(dbutils.Flds{eFile.FieldNameWithAlias(&eFile.CreatedAt, aliasFiles): sortOrderByCreatedAt})).
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
	db := dbOrTx(ctx, r.db)
	rows, err := db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows != nil {
		files, err = r.mappers.MapRowsToFilesTask(rows)
		if err != nil {
			return nil, err
		}
	}

	return files, nil
}

func (r *FileRepository) GetAll(ctx context.Context, includeDeleted bool) ([]*ddownload.File, error) {
	return r.getAll(ctx, fileByFields{}, dbutils.OrderDesc, includeDeleted)
}

func (r *FileRepository) GetAllFullNames(ctx context.Context, includeDeleted bool) ([]string, error) {
	var eFile edownload.File

	sqlWhere := squirrel.And{}

	if !includeDeleted {
		sqlWhere = append(sqlWhere, squirrel.Eq{eFile.FieldName(&eFile.DeletedAt): nil})
	}

	sqlQuery, args, err := squirrel.Select(eFile.FieldName(&eFile.FullName)).
		From(eFile.TableName()).
		Where(sqlWhere).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	db := dbOrTx(ctx, r.db)
	rows, err := db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fullNames []string
	for rows.Next() {
		var fullName string
		err = rows.Scan(&fullName)
		if err != nil {
			return nil, err
		}

		if fullName != "" {
			fullNames = append(fullNames, fullName)
		}
	}

	return fullNames, nil
}

func (r *FileRepository) GetBeforeTime(ctx context.Context, before time.Time, limit uint64) ([]*ddownload.File, error) {
	return r.getAll(
		ctx,
		fileByFields{
			beforeTime: &before,
			limit:      &limit,
		},
		dbutils.OrderDesc,
		false,
	)
}

func (r *FileRepository) GetByStatus(ctx context.Context, status dtypes.FileStatus) ([]*ddownload.File, error) {
	return r.GetByStatuses(ctx, []dtypes.FileStatus{status})
}

func (r *FileRepository) GetByStatuses(ctx context.Context, statuses []dtypes.FileStatus) ([]*ddownload.File, error) {
	return r.getAll(ctx, fileByFields{statuses: statuses}, dbutils.OrderAsc, false)
}

func (r *FileRepository) GetByPartialHash(ctx context.Context, hash string) ([]*ddownload.File, error) {
	var h = &hash
	return r.getAll(ctx, fileByFields{partialHash: &h, statuses: []dtypes.FileStatus{dtypes.FileStatusDone}}, dbutils.OrderDesc, false)
}

func (r *FileRepository) GetWithoutPartialHash(ctx context.Context) ([]*ddownload.File, error) {
	var h *string
	return r.getAll(ctx, fileByFields{partialHash: &h}, dbutils.OrderAsc, false)
}

func (r *FileRepository) GetDuplicateHashes(ctx context.Context) ([]string, error) {
	var eFile edownload.File

	sqlWhere := squirrel.And{
		squirrel.Expr(eFile.FieldName(&eFile.PartialHash) + " IS NOT NULL"),
		squirrel.NotEq{eFile.FieldName(&eFile.FullName): ""},
		squirrel.Eq{
			eFile.FieldName(&eFile.Status):    dtypes.FileStatusDone,
			eFile.FieldName(&eFile.DeletedAt): nil,
		},
	}

	// SELECT partial_hash
	// FROM files
	// WHERE partial_hash IS NOT NULL
	// GROUP BY partial_hash
	// HAVING COUNT(*) > 1;
	sqlQuery, args, err := squirrel.Select(eFile.FieldName(&eFile.PartialHash)).
		From(eFile.TableName()).
		Where(sqlWhere).
		GroupBy(eFile.FieldName(&eFile.PartialHash)).
		Having("COUNT(*) > 1").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error generating SQL: %v", err)
	}

	var hashes []string

	// Execute the query
	db := dbOrTx(ctx, r.db)
	rows, err := db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return hashes, nil
		}
		return nil, err
	}
	defer rows.Close()

	if rows != nil {
		for rows.Next() {
			var hash string
			err := rows.Scan(&hash)
			if err != nil {
				continue
			}
			hashes = append(hashes, hash)
		}
	}

	return hashes, nil
}

func (r *FileRepository) GetDeleted(ctx context.Context, from, to *time.Time) ([]*ddownload.File, error) {
	var eFile edownload.File

	sqlWhere := squirrel.And{
		squirrel.NotEq{eFile.FieldName(&eFile.DeletedAt): nil},
	}

	if from != nil {
		sqlWhere = append(sqlWhere, squirrel.GtOrEq{
			eFile.FieldName(&eFile.DeletedAt): *from,
		})
	}

	if to != nil {
		sqlWhere = append(sqlWhere, squirrel.LtOrEq{
			eFile.FieldName(&eFile.DeletedAt): *to,
		})
	}

	sqlQuery, args, err := squirrel.Select(eFile.FieldsAll()...).
		From(eFile.TableName()).
		Where(sqlWhere).
		OrderBy(dbutils.OrderBy(dbutils.Flds{eFile.FieldName(&eFile.DeletedAt): dbutils.OrderAsc})).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error generating SQL: %v", err)
	}

	var files []*ddownload.File

	// Execute the query
	db := dbOrTx(ctx, r.db)
	rows, err := db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return files, nil
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
