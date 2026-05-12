package download

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/pkg/dbutils"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/download/mappers"
	"github.com/neosy/elengrab/internal/repository/sqlite/sqlutil"
)

type MediaDownloadRepository struct {
	mappers *mappers.Mappers
	db      *sql.DB
	lock    dbexec.WriteLocker

	// filters
	filters downloadRepositoryFilters

	// options
	retryOptions dbexec.RetryOptions
}

type downloadByFields struct {
	statuses    []dtypes.MediaDownloadStatus
	beforeTime  *time.Time
	limit       *uint64
	partialHash **string
}

// NewMediaDownloadRepository returns a new object for the repository
func NewMediaDownloadRepository(db *sql.DB, lock dbexec.WriteLocker) *MediaDownloadRepository {
	return &MediaDownloadRepository{
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

func (r *MediaDownloadRepository) Copy() *MediaDownloadRepository {
	rep := uptr.Copy(r)

	rep.mappers = r.mappers
	rep.db = r.db
	rep.lock = r.lock

	rep.filters = rep.filters.copy()

	return rep
}

func (r *MediaDownloadRepository) downloadStatusesToStrings(statuses []dtypes.MediaDownloadStatus) []string {
	var statuseStrings = make([]string, 0, len(statuses))
	for _, status := range statuses {
		statuseStrings = append(statuseStrings, status.String())
	}
	return statuseStrings
}

func (r *MediaDownloadRepository) WithUser(userID uuid.UUID) persistence.MediaDownloadRepository {
	rep := r.Copy()
	rep.filters.userID = &userID
	return rep
}

func (r *MediaDownloadRepository) WithFilters(filters map[string]any) persistence.MediaDownloadRepository {
	rep := r.Copy()
	for key, value := range filters {
		switch key {
		case "userID":
			v, ok := value.(uuid.UUID)
			if ok {
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

func (r *MediaDownloadRepository) Insert(ctx context.Context, download *ddownload.MediaDownload) error {
	return r.save(ctx, download)
}

func (r *MediaDownloadRepository) Update(ctx context.Context, download *ddownload.MediaDownload) error {
	return r.save(ctx, download)
}

func (r *MediaDownloadRepository) save(ctx context.Context, download *ddownload.MediaDownload) error {
	if download == nil {
		return ierrors.ErrFuncParamNullPointer
	}

	// Convert the domain model to a database entity
	eDownload, err := r.mappers.MapDownloadDomainToEntity(download)
	if err != nil {
		return err
	}

	// Get the list of fields and values for insertion
	fields := eDownload.Fields()
	values := eDownload.Values()

	// If this is an update — add the UpdatedAt field with the current time
	// if isUpd {
	// 	fields = append(fields, eDownload.FieldName(&eDownload.UpdatedAt))
	// 	values = append(values, squirrel.Expr("CURRENT_TIMESTAMP"))
	// }

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Insert(eDownload.TableName()).
		Columns(fields...).
		Values(values...).
		Suffix(dbutils.UpsertSuffix(fields, eDownload.FieldName(&eDownload.DownloadID))).
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

func (r *MediaDownloadRepository) UpdateStatusToNew(ctx context.Context, statuses []dtypes.MediaDownloadStatus) error {
	if len(statuses) == 0 {
		return nil
	}

	statusStrings := r.downloadStatusesToStrings(statuses)

	var ent edownload.MediaDownload

	sqlWhere := squirrel.And{
		squirrel.Eq{ent.FieldName(&ent.Status): statusStrings},
		squirrel.Eq{ent.FieldName(&ent.DeletedAt): nil},
	}

	// Build query
	sqlBuilder := squirrel.Update(ent.TableName()).
		Set(ent.FieldName(&ent.Status), dtypes.MediaDownloadStatusNew.String()).
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

func (r *MediaDownloadRepository) UpdateOwner(ctx context.Context, fromID, toID uuid.UUID) error {
	var eDownload edownload.MediaDownload

	sqlWhere := squirrel.Eq{eDownload.FieldName(&eDownload.UserID): fromID}

	// Build query
	sqlBuilder := squirrel.Update(eDownload.TableName()).
		Set(eDownload.FieldName(&eDownload.UserID), toID).
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

func (r *MediaDownloadRepository) Delete(ctx context.Context, downloadID uuid.UUID, soft bool) error {
	if soft {
		return r.softDelete(ctx, downloadID)
	} else {
		return r.hardDelete(ctx, downloadID)
	}
}

func (r *MediaDownloadRepository) softDelete(ctx context.Context, downloadID uuid.UUID) error {
	var eDownload edownload.MediaDownload

	fieldsToUpdate := map[string]interface{}{
		eDownload.FieldName(&eDownload.UpdatedAt): squirrel.Expr("CURRENT_TIMESTAMP"),
		eDownload.FieldName(&eDownload.DeletedAt): squirrel.Expr("CURRENT_TIMESTAMP"),
	}

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Update(eDownload.TableName()).
		SetMap(fieldsToUpdate).
		Where(squirrel.Eq{eDownload.FieldName(&eDownload.DownloadID): downloadID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save mediaDownlaod: %v", err)
	}

	return nil
}

func (r *MediaDownloadRepository) hardDelete(ctx context.Context, downloadID uuid.UUID) error {
	var ent edownload.MediaDownload

	// Build DELETE query
	sqlBuilder := squirrel.Delete(ent.TableName()).
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
		return fmt.Errorf("failed to delete mediaDownload: %v", err)
	}

	return nil
}

func (r *MediaDownloadRepository) Restore(ctx context.Context, downloadID uuid.UUID) error {
	var eDownload edownload.MediaDownload

	fieldsToUpdate := map[string]any{
		eDownload.FieldName(&eDownload.UpdatedAt): squirrel.Expr("CURRENT_TIMESTAMP"),
		eDownload.FieldName(&eDownload.DeletedAt): nil,
	}

	sqlWhere := squirrel.And{
		squirrel.Eq{eDownload.FieldName(&eDownload.DownloadID): downloadID},
		squirrel.NotEq{eDownload.FieldName(&eDownload.DeletedAt): nil},
	}

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Update(eDownload.TableName()).
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
		return fmt.Errorf("failed to save mediaDownload: %v", err)
	}

	return nil
}

func (r *MediaDownloadRepository) FindByDownloadID(ctx context.Context, downloadID uuid.UUID) (*ddownload.MediaDownload, error) {
	var (
		eDownload edownload.MediaDownload
		eTask     edownload.DownloadTask

		aliasDownloads = "f"
		aliasTasks     = "t"
	)

	selectFields := append(eDownload.FieldsAllWithAlias(aliasDownloads), eTask.FieldsAllWithAlias(aliasTasks)...)

	sqlWhere := squirrel.Eq{
		eDownload.FieldNameWithAlias(&eDownload.DownloadID, aliasDownloads): downloadID.String(),
		eDownload.FieldNameWithAlias(&eDownload.DeletedAt, aliasDownloads):  nil,
	}

	if r.filters.userID != nil {
		sqlWhere[eDownload.FieldNameWithAlias(&eDownload.UserID, aliasDownloads)] = r.filters.userID
	}

	sqlQuery, args, err := squirrel.Select(selectFields...).
		From(eDownload.TableName() + " AS " + aliasDownloads).
		LeftJoin(
			eTask.TableName() + " AS " + aliasTasks +
				" ON " + aliasTasks + "." + eTask.FieldName(&eTask.DownloadID) +
				" = " + aliasDownloads + "." + eDownload.FieldName(&eDownload.DownloadID),
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
		err := row.Scan(append(eDownload.FieldPointers(), eTask.FieldPointers()...)...)
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
	download, err := r.mappers.MapDownloadEntityToDomain(&eDownload, &eTask)
	if err != nil {
		return nil, err
	}

	return download, nil
}

func (r *MediaDownloadRepository) iterateGetAll(
	ctx context.Context,
	byFields downloadByFields,
	sortOrderBy string,
	includeDeleted bool,
	fn func(*ddownload.MediaDownload) error,
) error {
	var (
		eDownload edownload.MediaDownload
		eTask     edownload.DownloadTask

		aliasDownloads = "f"
		aliasTasks     = "t"
	)

	selectFields := append(eDownload.FieldsAllWithAlias(aliasDownloads), eTask.FieldsAllWithAlias(aliasTasks)...)

	var conditions = squirrel.And{}
	if len(byFields.statuses) > 0 {
		statusStrings := r.downloadStatusesToStrings(byFields.statuses)
		conditions = append(conditions, squirrel.Eq{eDownload.FieldNameWithAlias(&eDownload.Status, aliasDownloads): statusStrings})
	}
	if r.filters.userID != nil {
		conditions = append(conditions, squirrel.Eq{eDownload.FieldNameWithAlias(&eDownload.UserID, aliasDownloads): *r.filters.userID})
	}
	if r.filters.title != nil && *r.filters.title != "" {
		conditions = append(conditions, sqlutil.Like(eDownload.FieldNameWithAlias(&eDownload.MediaTitleLower, aliasDownloads), *r.filters.title))
	}
	if byFields.beforeTime != nil && !byFields.beforeTime.IsZero() {
		t := byFields.beforeTime.Add(-1 * time.Nanosecond)
		conditions = append(conditions, squirrel.Lt{eDownload.FieldNameWithAlias(&eDownload.CreatedAt, aliasDownloads): t})
	}
	if byFields.partialHash != nil {
		if *byFields.partialHash == nil {
			conditions = append(conditions, squirrel.Expr(eDownload.FieldNameWithAlias(&eDownload.PartialHash, aliasDownloads)+" IS NULL"))
		} else {
			conditions = append(conditions, squirrel.Eq{eDownload.FieldNameWithAlias(&eDownload.PartialHash, aliasDownloads): **byFields.partialHash})
		}
	}
	if !includeDeleted {
		conditions = append(conditions, squirrel.Eq{eDownload.FieldNameWithAlias(&eDownload.DeletedAt, aliasDownloads): nil})
	}

	sqlWhere := conditions

	// Create an ORDER BY clause based on fieldы with the specified sort order.
	orderBy := dbutils.OrderBy(
		dbutils.Flds{
			eDownload.FieldNameWithAlias(&eDownload.CreatedAt, aliasDownloads): sortOrderBy,
		})

	qb := squirrel.Select(selectFields...).
		From(eDownload.TableName() + " AS " + aliasDownloads).
		Where(sqlWhere).
		OrderBy(orderBy).
		LeftJoin(
			eTask.TableName() + " AS " + aliasTasks +
				" ON " + eTask.FieldNameWithAlias(&eTask.DownloadID, aliasTasks) +
				" = " + eDownload.FieldNameWithAlias(&eDownload.DownloadID, aliasDownloads),
		).
		PlaceholderFormat(squirrel.Dollar)

	if byFields.limit != nil && *byFields.limit > 0 {
		qb = qb.Limit(*byFields.limit)
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
		err = r.mappers.MapRowsToDownloadsTask(rows, func(f *ddownload.MediaDownload) error {
			if err := fn(f); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *MediaDownloadRepository) IterateGetAll(ctx context.Context, includeDeleted bool, fn func(*ddownload.MediaDownload) error) error {
	return r.iterateGetAll(ctx, downloadByFields{}, dbutils.OrderDesc, includeDeleted, fn)
}

func (r *MediaDownloadRepository) GetAllFullNames(ctx context.Context, includeDeleted bool) (map[string]struct{}, error) {
	names := make(map[string]struct{})
	r.IterateFullNames(ctx, includeDeleted, func(name string) error {
		names[name] = struct{}{}
		return nil
	})
	return names, nil
}

func (r *MediaDownloadRepository) IterateFullNames(ctx context.Context, includeDeleted bool, fn func(string) error) error {
	var eDownload edownload.MediaDownload

	sqlWhere := squirrel.And{}

	if !includeDeleted {
		sqlWhere = append(sqlWhere, squirrel.Eq{eDownload.FieldName(&eDownload.DeletedAt): nil})
	}

	if r.filters.userID != nil {
		sqlWhere = append(sqlWhere, squirrel.Eq{eDownload.FieldName(&eDownload.UserID): *r.filters.userID})
	}

	sqlQuery, args, err := squirrel.Select(eDownload.FieldName(&eDownload.FileFullName)).
		From(eDownload.TableName()).
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

	for rows.Next() {
		var fullName string
		err = rows.Scan(&fullName)
		if err != nil {
			return err
		}

		if fullName == "" {
			continue
		}

		if err := fn(fullName); err != nil {
			return err
		}
	}

	return nil
}

func (r *MediaDownloadRepository) GetBeforeTime(ctx context.Context, before time.Time, limit uint64) ([]*ddownload.MediaDownload, error) {
	downloads := make([]*ddownload.MediaDownload, 0)

	err := r.iterateGetAll(
		ctx,
		downloadByFields{
			beforeTime: &before,
			limit:      &limit,
		},
		dbutils.OrderDesc,
		false,
		func(f *ddownload.MediaDownload) error {
			downloads = append(downloads, f)
			return nil
		})
	if err != nil {
		return nil, err
	}

	return downloads, nil
}

func (r *MediaDownloadRepository) GetByStatus(ctx context.Context, status dtypes.MediaDownloadStatus) ([]*ddownload.MediaDownload, error) {
	return r.GetByStatuses(ctx, []dtypes.MediaDownloadStatus{status})
}

func (r *MediaDownloadRepository) GetByStatuses(ctx context.Context, statuses []dtypes.MediaDownloadStatus) ([]*ddownload.MediaDownload, error) {
	var downloads = make([]*ddownload.MediaDownload, 0)

	err := r.iterateGetAll(
		ctx,
		downloadByFields{statuses: statuses},
		dbutils.OrderAsc, false,
		func(f *ddownload.MediaDownload) error {
			downloads = append(downloads, f)
			return nil
		})
	if err != nil {
		return nil, err
	}

	return downloads, nil
}

func (r *MediaDownloadRepository) GetByPartialHash(ctx context.Context, hash string) ([]*ddownload.MediaDownload, error) {
	var (
		h         = &hash
		downloads = make([]*ddownload.MediaDownload, 0)
	)

	err := r.iterateGetAll(
		ctx, downloadByFields{
			statuses:    []dtypes.MediaDownloadStatus{dtypes.MediaDownloadStatusDone},
			partialHash: &h,
		},
		dbutils.OrderDesc,
		false,
		func(f *ddownload.MediaDownload) error {
			downloads = append(downloads, f)
			return nil
		})
	if err != nil {
		return nil, err
	}

	return downloads, nil
}

func (r *MediaDownloadRepository) GetWithoutPartialHash(ctx context.Context) ([]*ddownload.MediaDownload, error) {
	var (
		h         *string
		downloads = make([]*ddownload.MediaDownload, 0)
	)

	err := r.iterateGetAll(
		ctx,
		downloadByFields{partialHash: &h},
		dbutils.OrderAsc, false,
		func(f *ddownload.MediaDownload) error {
			downloads = append(downloads, f)
			return nil
		})
	if err != nil {
		return nil, err
	}

	return downloads, nil
}

func (r *MediaDownloadRepository) GetDuplicateHashes(ctx context.Context, scope dtypes.UniquenessScope) ([]ddownload.DuplicateHashRow, error) {
	var eDownload edownload.MediaDownload

	sqlWhere := squirrel.And{
		squirrel.Expr(eDownload.FieldName(&eDownload.PartialHash) + " IS NOT NULL"),
		squirrel.NotEq{eDownload.FieldName(&eDownload.FileFullName): ""},
		squirrel.Eq{
			eDownload.FieldName(&eDownload.Status):    dtypes.MediaDownloadStatusDone.String(),
			eDownload.FieldName(&eDownload.DeletedAt): nil,
		},
	}

	if r.filters.userID != nil {
		sqlWhere = append(sqlWhere, squirrel.Eq{eDownload.FieldName(&eDownload.UserID): *r.filters.userID})
	}

	fields := make([]string, 0, 2)
	groupByFields := make([]string, 0, 2)

	fields = append(fields, eDownload.FieldName(&eDownload.PartialHash))
	groupByFields = append(groupByFields, eDownload.FieldName(&eDownload.PartialHash))

	if scope == dtypes.UniquenessScopePerUser {
		fields = append(fields, eDownload.FieldName(&eDownload.UserID))
		groupByFields = append(groupByFields, eDownload.FieldName(&eDownload.UserID))
	} else {
		fields = append(fields, "NULL AS "+eDownload.FieldName(&eDownload.UserID))
	}

	// SELECT partial_hash
	// FROM media_downloads
	// WHERE partial_hash IS NOT NULL
	// GROUP BY partial_hash
	// HAVING COUNT(*) > 1;
	sqlQuery, args, err := squirrel.Select(fields...).
		From(eDownload.TableName()).
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

func (r *MediaDownloadRepository) GetDeleted(ctx context.Context, from, to *time.Time) ([]*ddownload.MediaDownload, error) {
	var eDownload edownload.MediaDownload

	sqlWhere := squirrel.And{
		squirrel.NotEq{eDownload.FieldName(&eDownload.DeletedAt): nil},
	}

	if r.filters.userID != nil {
		sqlWhere = append(sqlWhere, squirrel.Eq{eDownload.FieldName(&eDownload.UserID): *r.filters.userID})
	}

	if from != nil {
		sqlWhere = append(sqlWhere, squirrel.GtOrEq{
			eDownload.FieldName(&eDownload.DeletedAt): *from,
		})
	}

	if to != nil {
		sqlWhere = append(sqlWhere, squirrel.LtOrEq{
			eDownload.FieldName(&eDownload.DeletedAt): *to,
		})
	}

	orderBy := dbutils.OrderBy(dbutils.Flds{eDownload.FieldName(&eDownload.DeletedAt): dbutils.OrderAsc})

	sqlQuery, args, err := squirrel.Select(eDownload.FieldsAll()...).
		From(eDownload.TableName()).
		Where(sqlWhere).
		OrderBy(orderBy).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error generating SQL: %v", err)
	}

	var downloads []*ddownload.MediaDownload

	// Execute the query
	db := dbexec.Resolve(ctx, r.db)
	rows, err := db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return downloads, nil
		}
		return nil, err
	}
	defer rows.Close()

	if rows != nil {
		downloads, err = r.mappers.MapRowsToDownloads(rows)
		if err != nil {
			return nil, err
		}
	}

	return downloads, nil
}

func (r *MediaDownloadRepository) FillEmptyMediaTitleLower(ctx context.Context) error {
	var eDownload edownload.MediaDownload

	sqlWhere := squirrel.And{
		squirrel.Eq{eDownload.FieldName(&eDownload.MediaTitleLower): ""},
	}

	sqlQuery, args, err := squirrel.Select(eDownload.FieldsAll()...).
		From(eDownload.TableName()).
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

	var downloads []*ddownload.MediaDownload
	if rows != nil {
		downloads, err = r.mappers.MapRowsToDownloads(rows)
		if err != nil {
			return err
		}
	}

	for _, download := range downloads {
		sqlQuery, args, err := squirrel.
			Update(eDownload.TableName()).
			SetMap(map[string]any{
				eDownload.FieldName(&eDownload.MediaTitleLower): strings.ToLower(download.MediaTitle),
			}).
			Where(squirrel.Eq{eDownload.FieldName(&eDownload.DownloadID): download.DownloadID}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if err != nil {
			return fmt.Errorf("error generating SQL: %v", err)
		}

		err = dbexec.ExecContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
		if err != nil {
			return fmt.Errorf("failed to save mediaDownload: %v", err)
		}
	}

	return nil
}

func (r *MediaDownloadRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.Tx(ctx, r.db, r.lock, fn)
}
