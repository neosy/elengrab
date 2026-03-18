package download

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/download/mappers"
	"github.com/neosy/elengrab/internal/repository/sqlite/sqlutil"
	"github.com/neosy/elengrab/pkg/dbutils"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

type FileRepository struct {
	mappers *mappers.Mappers
	db      *sql.DB
	lock    dbexec.WriteLocker

	// filters
	filters fileRepositoryFilters

	// options
	retryOptions dbexec.RetryOptions
}

type fileByFields struct {
	statuses    []dtypes.FileStatus
	beforeTime  *time.Time
	limit       *uint64
	partialHash **string
}

// NewFileRepository returns a new object for the repository
func NewFileRepository(db *sql.DB, lock dbexec.WriteLocker) *FileRepository {
	return &FileRepository{
		mappers: mappers.NewMappers(),
		db:      db,
		lock:    lock,

		// options
		retryOptions: dbexec.RetryOptions{
			MaxRetries: maxRetriesDefault,
			Delay:      retryDelayDefault,
		},
	}
}

func (r *FileRepository) Copy() *FileRepository {
	rep := uptr.Copy(r)

	rep.mappers = r.mappers
	rep.db = r.db
	rep.lock = r.lock

	rep.filters = rep.filters.copy()

	return rep
}

func (r *FileRepository) fileStatusesToStrings(statuses []dtypes.FileStatus) []string {
	var statuseStrings = make([]string, 0, len(statuses))
	for _, status := range statuses {
		statuseStrings = append(statuseStrings, status.String())
	}
	return statuseStrings
}

func (r *FileRepository) WithUser(userID uuid.UUID) persistence.FileRepository {
	rep := r.Copy()
	rep.filters.userID = &userID
	return rep
}

func (r *FileRepository) WithFilters(filters map[string]any) persistence.FileRepository {
	rep := r.Copy()
	for key, value := range filters {
		switch key {
		case "userID":
			v, ok := value.(uuid.UUID)
			if ok && v != uuid.Nil {
				rep.filters.userID = &v
			}
		case "title":
			v, ok := value.(string)
			if ok && v != "" {
				rep.filters.title = &v
			}
		}
	}
	return rep
}

func (r *FileRepository) Insert(ctx context.Context, file *ddownload.File) error {
	return r.save(ctx, file)
}

func (r *FileRepository) Update(ctx context.Context, file *ddownload.File) error {
	return r.save(ctx, file)
}

func (r *FileRepository) save(ctx context.Context, file *ddownload.File) error {
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
	// if isUpd {
	// 	fields = append(fields, eFile.FieldName(&eFile.UpdatedAt))
	// 	values = append(values, squirrel.Expr("CURRENT_TIMESTAMP"))
	// }

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Insert(eFile.TableName()).
		Columns(fields...).
		Values(values...).
		Suffix(dbutils.UpsertSuffix(fields, eFile.FieldName(&eFile.FileID))).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	return nil
}

func (r *FileRepository) UpdateStatusToNew(ctx context.Context, statuses []dtypes.FileStatus) error {
	if len(statuses) == 0 {
		return nil
	}

	statusStrings := r.fileStatusesToStrings(statuses)

	var ent edownload.File

	sqlWhere := squirrel.And{
		squirrel.Eq{ent.FieldName(&ent.Status): statusStrings},
		squirrel.Eq{ent.FieldName(&ent.DeletedAt): nil},
	}

	// Build query
	sqlBuilder := squirrel.Update(ent.TableName()).
		Set(ent.FieldName(&ent.Status), dtypes.FileStatusNew.String()).
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
		Where(squirrel.Eq{eFile.FieldName(&eFile.FileID): fileId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	return nil
}

func (r *FileRepository) hardDelete(ctx context.Context, fileId uuid.UUID) error {
	var ent edownload.File

	// Build DELETE query
	sqlBuilder := squirrel.Delete(ent.TableName()).
		Where(squirrel.Eq{ent.FieldName(&ent.FileID): fileId.String()}).
		PlaceholderFormat(squirrel.Dollar)

	// Generate SQL and args
	sqlStr, args, err := sqlBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlStr, args, r.retryOptions)
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
		squirrel.Eq{eFile.FieldName(&eFile.FileID): fileId},
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
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	return nil
}

func (r *FileRepository) FindByFileID(ctx context.Context, fileId uuid.UUID) (*ddownload.File, error) {
	var (
		eFile edownload.File
		eTask edownload.DownloadTask

		aliasFiles = "f"
		aliasTasks = "t"
	)

	selectFields := append(eFile.FieldsAllWithAlias(aliasFiles), eTask.FieldsAllWithAlias(aliasTasks)...)

	sqlWhere := squirrel.Eq{
		eFile.FieldNameWithAlias(&eFile.FileID, aliasFiles):    fileId.String(),
		eFile.FieldNameWithAlias(&eFile.DeletedAt, aliasFiles): nil,
	}

	if r.filters.userID != nil {
		sqlWhere[eFile.FieldNameWithAlias(&eFile.UserID, aliasFiles)] = r.filters.userID
	}

	sqlQuery, args, err := squirrel.Select(selectFields...).
		From(eFile.TableName() + " AS " + aliasFiles).
		LeftJoin(
			eTask.TableName() + " AS " + aliasTasks +
				" ON " + aliasTasks + "." + eTask.FieldName(&eTask.FileID) +
				" = " + aliasFiles + "." + eFile.FieldName(&eFile.FileID),
		).
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
		err := row.Scan(append(eFile.FieldPointers(), eTask.FieldPointers()...)...)
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
	file, err := r.mappers.MapFileEntityToDomain(&eFile, &eTask)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (r *FileRepository) getAll(
	ctx context.Context,
	byFields fileByFields,
	sortOrderBy string,
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
		statusStrings := r.fileStatusesToStrings(byFields.statuses)
		conditions = append(conditions, squirrel.Eq{eFile.FieldNameWithAlias(&eFile.Status, aliasFiles): statusStrings})
	}
	if r.filters.userID != nil {
		conditions = append(conditions, squirrel.Eq{eFile.FieldNameWithAlias(&eFile.UserID, aliasFiles): *r.filters.userID})
	}
	if r.filters.title != nil && *r.filters.title != "" {
		conditions = append(conditions, sqlutil.Like(eFile.FieldNameWithAlias(&eFile.MediaTitleLower, aliasFiles), *r.filters.title))
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

	// Create an ORDER BY clause based on fieldы with the specified sort order.
	orderBy := dbutils.OrderBy(
		dbutils.Flds{
			eFile.FieldNameWithAlias(&eFile.CreatedAt, aliasFiles): sortOrderBy,
		})

	qb := squirrel.Select(selectFields...).
		From(eFile.TableName() + " AS " + aliasFiles).
		Where(sqlWhere).
		OrderBy(orderBy).
		LeftJoin(
			eTask.TableName() + " AS " + aliasTasks +
				" ON " + eTask.FieldNameWithAlias(&eTask.FileID, aliasTasks) +
				" = " + eFile.FieldNameWithAlias(&eFile.FileID, aliasFiles),
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
	db := dbexec.Resolve(ctx, r.db)
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

	if r.filters.userID != nil {
		sqlWhere = append(sqlWhere, squirrel.Eq{eFile.FieldName(&eFile.UserID): *r.filters.userID})
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
	db := dbexec.Resolve(ctx, r.db)
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
	return r.getAll(
		ctx, fileByFields{
			statuses:    []dtypes.FileStatus{dtypes.FileStatusDone},
			partialHash: &h,
		},
		dbutils.OrderDesc,
		false,
	)
}

func (r *FileRepository) GetWithoutPartialHash(ctx context.Context) ([]*ddownload.File, error) {
	var h *string
	return r.getAll(ctx, fileByFields{partialHash: &h}, dbutils.OrderAsc, false)
}

func (r *FileRepository) GetDuplicateHashes(ctx context.Context, scope dtypes.UniquenessScope) ([]ddownload.DuplicateHashRow, error) {
	var eFile edownload.File

	sqlWhere := squirrel.And{
		squirrel.Expr(eFile.FieldName(&eFile.PartialHash) + " IS NOT NULL"),
		squirrel.NotEq{eFile.FieldName(&eFile.FullName): ""},
		squirrel.Eq{
			eFile.FieldName(&eFile.Status):    dtypes.FileStatusDone.String(),
			eFile.FieldName(&eFile.DeletedAt): nil,
		},
	}

	if r.filters.userID != nil {
		sqlWhere = append(sqlWhere, squirrel.Eq{eFile.FieldName(&eFile.UserID): *r.filters.userID})
	}

	fields := make([]string, 0, 2)
	groupByFields := make([]string, 0, 2)

	fields = append(fields, eFile.FieldName(&eFile.PartialHash))
	groupByFields = append(groupByFields, eFile.FieldName(&eFile.PartialHash))

	if scope == dtypes.UniquenessScopePerUser {
		fields = append(fields, eFile.FieldName(&eFile.UserID))
		groupByFields = append(groupByFields, eFile.FieldName(&eFile.UserID))
	} else {
		fields = append(fields, "NULL AS "+eFile.FieldName(&eFile.UserID))
	}

	// SELECT partial_hash
	// FROM files
	// WHERE partial_hash IS NOT NULL
	// GROUP BY partial_hash
	// HAVING COUNT(*) > 1;
	sqlQuery, args, err := squirrel.Select(fields...).
		From(eFile.TableName()).
		Where(sqlWhere).
		GroupBy(groupByFields...).
		Having("COUNT(*) > 1").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error generating SQL: %v", err)
	}

	var hashRows []ddownload.DuplicateHashRow

	// Execute the query
	db := dbexec.Resolve(ctx, r.db)
	rows, err := db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return hashRows, nil
		}
		return nil, err
	}
	defer rows.Close()

	if rows != nil {
		for rows.Next() {
			var (
				hashRow ddownload.DuplicateHashRow
				userID  sql.NullString
			)
			err := rows.Scan(&hashRow.Hash, &userID)
			if err != nil {
				continue
			}
			if userID.Valid {
				uid, err := uuid.Parse(userID.String)
				if err != nil {
					return nil, err
				}
				hashRow.UserID = &uid
			}
			hashRows = append(hashRows, hashRow)
		}
	}

	return hashRows, nil
}

func (r *FileRepository) GetDeleted(ctx context.Context, from, to *time.Time) ([]*ddownload.File, error) {
	var eFile edownload.File

	sqlWhere := squirrel.And{
		squirrel.NotEq{eFile.FieldName(&eFile.DeletedAt): nil},
	}

	if r.filters.userID != nil {
		sqlWhere = append(sqlWhere, squirrel.Eq{eFile.FieldName(&eFile.UserID): *r.filters.userID})
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

	orderBy := dbutils.OrderBy(dbutils.Flds{eFile.FieldName(&eFile.DeletedAt): dbutils.OrderAsc})

	sqlQuery, args, err := squirrel.Select(eFile.FieldsAll()...).
		From(eFile.TableName()).
		Where(sqlWhere).
		OrderBy(orderBy).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error generating SQL: %v", err)
	}

	var files []*ddownload.File

	// Execute the query
	db := dbexec.Resolve(ctx, r.db)
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
	r.lock.Lock()
	defer r.lock.Unlock()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(dbexec.CtxWithTx(ctx, tx)); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *FileRepository) FillEmptyMediaTitleLower(ctx context.Context) error {
	var eFile edownload.File

	sqlWhere := squirrel.And{
		squirrel.Eq{eFile.FieldName(&eFile.MediaTitleLower): ""},
	}

	sqlQuery, args, err := squirrel.Select(eFile.FieldsAll()...).
		From(eFile.TableName()).
		Where(sqlWhere).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

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

	var files []*ddownload.File
	if rows != nil {
		files, err = r.mappers.MapRowsToFiles(rows)
		if err != nil {
			return err
		}
	}

	for _, file := range files {
		sqlQuery, args, err := squirrel.
			Update(eFile.TableName()).
			SetMap(map[string]any{
				eFile.FieldName(&eFile.MediaTitleLower): strings.ToLower(file.MediaTitle),
			}).
			Where(squirrel.Eq{eFile.FieldName(&eFile.FileID): file.FileID}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if err != nil {
			return fmt.Errorf("error generating SQL: %v", err)
		}

		err = dbexec.ExecContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
		if err != nil {
			return fmt.Errorf("failed to save file: %v", err)
		}
	}

	return nil
}
