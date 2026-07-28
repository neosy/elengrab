package media

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/pkg/dbutils"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	emedia "github.com/neosy/elengrab/internal/repository/sqlite/media/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/media/mappers"
)

type ThumbnailRepository struct {
	mappers *mappers.Mappers
	db      *sql.DB
	lock    dbexec.WriteLocker

	// options
	retryOptions dbexec.RetryOptions
}

// NewThumbnailRepository returns a new object for the repository
func NewThumbnailRepository(db *sql.DB, lock dbexec.WriteLocker) *ThumbnailRepository {
	return &ThumbnailRepository{
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

func (r *ThumbnailRepository) Insert(ctx context.Context, thumbnail *dmedia.Thumbnail) error {
	return r.Save(ctx, thumbnail)
}

func (r *ThumbnailRepository) Update(ctx context.Context, thumbnail *dmedia.Thumbnail) error {
	return r.Save(ctx, thumbnail)
}

func (r *ThumbnailRepository) Save(ctx context.Context, thumbnail *dmedia.Thumbnail) error {
	if thumbnail == nil {
		return ierrors.ErrFuncParamNullPointer
	}

	// Convert the domain model to a database entity
	eThumbnail, err := r.mappers.MapThumbnailDomainToEntity(thumbnail)
	if err != nil {
		return err
	}

	// Get the list of fields and values for insertion
	fields := eThumbnail.Fields()
	values := eThumbnail.Values()

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Insert(eThumbnail.TableName()).
		Columns(fields...).
		Values(values...).
		Suffix(dbutils.UpsertSuffix(fields, eThumbnail.FieldName(&eThumbnail.ThumbID))).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save siteThumbnail: %v", err)
	}

	return nil
}

func (r *ThumbnailRepository) Delete(ctx context.Context, thumbID uuid.UUID) error {
	var eThumbnail emedia.Thumbnail

	// Build DELETE query
	sqlBuilder := squirrel.Delete(eThumbnail.TableName()).
		Where(squirrel.Eq{eThumbnail.FieldName(&eThumbnail.ThumbID): thumbID}).
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

func (r *ThumbnailRepository) FindByThumbID(
	ctx context.Context, thumbID uuid.UUID) (*dmedia.Thumbnail, error) {
	var ent emedia.Thumbnail

	// Build SQL query
	sqlQuery, args, err := squirrel.Select(ent.FieldsAll()...).
		From(ent.TableName()).
		Where(squirrel.Eq{ent.FieldName(&ent.ThumbID): thumbID.String()}).
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
		err := row.Scan(ent.FieldPointers()...)
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
	logo, err := r.mappers.MapThumbnailEntityToDomain(&ent)
	if err != nil {
		return nil, err
	}

	return logo, nil
}

func (r *ThumbnailRepository) ExistsByThumbID(ctx context.Context, thumbID uuid.UUID) (bool, error) {
	var ent emedia.Thumbnail

	query, args, err := squirrel.Select("1").
		From(ent.TableName()).
		Where(squirrel.Eq{ent.FieldName(&ent.ThumbID): thumbID.String()}).
		PlaceholderFormat(squirrel.Dollar).
		Limit(1).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	db := dbexec.Resolve(ctx, r.db)

	// Execute query and check if any row exists
	var exists int
	err = db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *ThumbnailRepository) FindByVersion(
	ctx context.Context,
	mediaID uuid.UUID,
	variant dtypes.ThumbnailVariant,
	version uint8,
) (*dmedia.Thumbnail, error) {
	var ent emedia.Thumbnail

	sqlWehere := squirrel.Eq{
		ent.FieldName(&ent.MediaID): mediaID.String(),
		ent.FieldName(&ent.Variant): variant.String(),
		ent.FieldName(&ent.Version): version,
	}

	// Build SQL query
	sqlQuery, args, err := squirrel.Select(ent.FieldsAll()...).
		From(ent.TableName()).
		Where(sqlWehere).
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
		err := row.Scan(ent.FieldPointers()...)
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
	logo, err := r.mappers.MapThumbnailEntityToDomain(&ent)
	if err != nil {
		return nil, err
	}

	return logo, nil
}

func (r *ThumbnailRepository) FindByMediaIDBest(
	ctx context.Context,
	mediaID uuid.UUID,
) (*dmedia.Thumbnail, error) {
	var ent emedia.Thumbnail

	sqlWehere := squirrel.Eq{
		ent.FieldName(&ent.MediaID): mediaID.String(),
	}

	sqlOrderVariant := `CASE ` + ent.FieldName(&ent.Variant) + `
			WHEN '` + dtypes.ThumbnailVariantOriginal.String() + `' THEN 1
			WHEN '` + dtypes.ThumbnailVariantLarge.String() + `' THEN 2
			WHEN '` + dtypes.ThumbnailVariantMedium.String() + `' THEN 3
			WHEN '` + dtypes.ThumbnailVariantSmall.String() + `' THEN 4
			ELSE 5
		END`

	// Build SQL query
	sqlQuery, args, err := squirrel.Select(ent.FieldsAll()...).
		From(ent.TableName()).
		Where(sqlWehere).
		OrderBy(
			sqlOrderVariant,
			dbutils.OrderBy(dbutils.Flds{ent.FieldName(&ent.Version): dbutils.OrderDesc}),
		).
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
		err := row.Scan(ent.FieldPointers()...)
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
	logo, err := r.mappers.MapThumbnailEntityToDomain(&ent)
	if err != nil {
		return nil, err
	}

	return logo, nil
}

func (r *ThumbnailRepository) GetByMediaID(
	ctx context.Context,
	mediaID uuid.UUID,
) ([]*dmedia.Thumbnail, error) {
	var ent emedia.Thumbnail

	sqlWehere := squirrel.Eq{
		ent.FieldName(&ent.MediaID): mediaID.String(),
	}

	// Build SQL query
	sqlQuery, args, err := squirrel.Select(ent.FieldsAll()...).
		From(ent.TableName()).
		Where(sqlWehere).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error generating SQL: %v", err)
	}

	var thumbnails []*dmedia.Thumbnail

	// Execute the query
	db := dbexec.Resolve(ctx, r.db)
	rows, err := db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return thumbnails, nil
		}
		return nil, err
	}
	defer rows.Close()

	if rows != nil {
		var (
			eThumbnail emedia.Thumbnail
			thumbnails []*dmedia.Thumbnail
		)

		for rows.Next() {
			err := rows.Scan(eThumbnail.FieldPointers()...)
			if err != nil {
				return nil, err
			}

			thumbnail, err := r.mappers.MapThumbnailEntityToDomain(&eThumbnail)
			if err != nil {
				return nil, err
			}

			thumbnails = append(thumbnails, thumbnail)
		}

	}

	return thumbnails, nil
}

func (r *ThumbnailRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.Tx(ctx, r.db, r.lock, fn)
}

func (r *ThumbnailRepository) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.TxIndependent(ctx, r.db, r.lock, fn)
}
