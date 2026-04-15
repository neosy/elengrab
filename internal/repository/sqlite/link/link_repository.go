package link

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	dlink "github.com/neosy/elengrab/internal/domain/link"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/pkg/dbutils"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	elink "github.com/neosy/elengrab/internal/repository/sqlite/link/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/link/mappers"
)

type LinkRepository struct {
	mappers *mappers.Mappers
	db      *sql.DB
	lock    dbexec.WriteLocker

	// options
	retryOptions dbexec.RetryOptions
}

// NewLinkRepository returns a new object for the repository
func NewLinkRepository(db *sql.DB, lock dbexec.WriteLocker) *LinkRepository {
	return &LinkRepository{
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

func (r *LinkRepository) Insert(ctx context.Context, link *dlink.Link) error {
	return r.save(ctx, link)
}

func (r *LinkRepository) Update(ctx context.Context, link *dlink.Link) error {
	return r.save(ctx, link)
}

func (r *LinkRepository) save(ctx context.Context, link *dlink.Link) error {
	if link == nil {
		return ierrors.ErrFuncParamNullPointer
	}

	// Convert the domain model to a database entity
	eLink, err := r.mappers.MapLinkDomainToEntity(link)
	if err != nil {
		return err
	}

	// Get the list of fields and values for insertion
	fields := eLink.Fields()
	values := eLink.Values()

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Insert(eLink.TableName()).
		Columns(fields...).
		Values(values...).
		Suffix(dbutils.UpsertSuffix(fields, eLink.FieldName(&eLink.LinkID))).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save link: %v", err)
	}

	return nil
}

func (r *LinkRepository) SoftDelete(ctx context.Context, linkID uuid.UUID) error {
	var eLink elink.Link

	fieldsToUpdate := map[string]interface{}{
		eLink.FieldName(&eLink.DeletedAt): squirrel.Expr("CURRENT_TIMESTAMP"),
	}

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Update(eLink.TableName()).
		SetMap(fieldsToUpdate).
		Where(squirrel.Eq{eLink.FieldName(&eLink.LinkID): linkID}).
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

func (r *LinkRepository) HardDelete(ctx context.Context, linkID uuid.UUID) error {
	var eLink elink.Link

	// Build DELETE query
	sqlBuilder := squirrel.Delete(eLink.TableName()).
		Where(squirrel.Eq{eLink.FieldName(&eLink.LinkID): linkID.String()}).
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

func (r *LinkRepository) Find(ctx context.Context, linkID uuid.UUID) (*dlink.Link, error) {
	var eLink elink.Link

	sqlQuery, args, err := squirrel.Select(eLink.FieldsAll()...).
		From(eLink.TableName()).
		Where(squirrel.Eq{eLink.FieldName(&eLink.LinkID): linkID}).
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
		row := db.QueryRowContext(ctx, sqlQuery, args...).Scan(eLink.FieldPointers()...)
		// Scan result into entity
		err := row
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
	link, err := r.mappers.MapLinkEntityToDomain(&eLink)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (r *LinkRepository) Exists(ctx context.Context, linkID uuid.UUID) (bool, error) {
	var eLink elink.Link

	// Build SQL query: SELECT 1 FROM table WHERE <id> = $1 LIMIT 1
	query, args, err := squirrel.Select("1").
		From(eLink.TableName()).
		Where(squirrel.Eq{eLink.FieldName(&eLink.LinkID): linkID}).
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

func (r *LinkRepository) FindLastByShortCode(ctx context.Context, shortCode string) (*dlink.Link, error) {
	var eLink elink.Link

	sqlQuery, args, err := squirrel.Select(eLink.FieldsAll()...).
		From(eLink.TableName()).
		Where(squirrel.Eq{eLink.FieldName(&eLink.ShortCode): shortCode}).
		OrderBy(dbutils.OrderBy(dbutils.Flds{eLink.FieldName(&eLink.CreatedAt): dbutils.OrderDesc})).
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
		err := row.Scan(eLink.FieldPointers()...)
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
	link, err := r.mappers.MapLinkEntityToDomain(&eLink)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (r *LinkRepository) ExistsActiveShortCode(ctx context.Context, shortCode string) (bool, error) {
	var eLink elink.Link

	sqlWhere := squirrel.And{
		squirrel.Eq{eLink.FieldName(&eLink.ShortCode): shortCode},
		squirrel.Eq{eLink.FieldName(&eLink.DeletedAt): nil},
		squirrel.Or{
			squirrel.Eq{eLink.FieldName(&eLink.ExpiresAt): nil},
			squirrel.Gt{eLink.FieldName(&eLink.ExpiresAt): time.Now().UTC()},
		},
	}

	query, args, err := squirrel.Select("1").
		From(eLink.TableName()).
		Where(sqlWhere).
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

func (r *LinkRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.Tx(ctx, r.db, r.lock, fn)
}
