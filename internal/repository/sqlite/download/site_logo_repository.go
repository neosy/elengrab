package download

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/download/mappers"
	"github.com/neosy/elengrab/pkg/dbutils"
)

type SiteLogoRepository struct {
	mappers *mappers.Mappers
	db      *sql.DB
	lock    dbexec.WriteLocker

	// options
	retryOptions dbexec.RetryOptions
}

// NewSiteLogoRepository returns a new object for the repository
func NewSiteLogoRepository(db *sql.DB, lock dbexec.WriteLocker) *SiteLogoRepository {
	return &SiteLogoRepository{
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

func (r *SiteLogoRepository) Insert(ctx context.Context, logo *dmedia.SiteLogo) error {
	return r.Save(ctx, logo)
}

func (r *SiteLogoRepository) Update(ctx context.Context, logo *dmedia.SiteLogo) error {
	return r.Save(ctx, logo)
}

func (r *SiteLogoRepository) Save(ctx context.Context, logo *dmedia.SiteLogo) error {
	if logo == nil {
		return errors.New("function parameter is a null pointer")
	}

	// Convert the domain model to a database entity
	eLogo, err := r.mappers.MapSiteLogoDomainToEntity(logo)
	if err != nil {
		return err
	}

	// Get the list of fields and values for insertion
	fields := eLogo.Fields()
	values := eLogo.Values()

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Insert(eLogo.TableName()).
		Columns(fields...).
		Values(values...).
		Suffix(dbutils.UpsertSuffix(fields, eLogo.FieldName(&eLogo.LogoID))).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save siteLogo: %v", err)
	}

	return nil
}

func (r *SiteLogoRepository) FindByLogoID(ctx context.Context, logoID uuid.UUID) (*dmedia.SiteLogo, error) {
	var ent edownload.SiteLogo

	// Build SQL query
	sqlQuery, args, err := squirrel.Select(ent.FieldsAll()...).
		From(ent.TableName()).
		Where(squirrel.Eq{ent.FieldName(&ent.LogoID): logoID.String()}).
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
	logo, err := r.mappers.MapSiteLogoEntityToDomain(&ent)
	if err != nil {
		return nil, err
	}

	return logo, nil
}

func (r *SiteLogoRepository) ExistsByLogoID(ctx context.Context, logoID uuid.UUID) (bool, error) {
	var ent edownload.SiteLogo

	// Build SQL query: SELECT 1 FROM table WHERE logo_id = $1 LIMIT 1
	query, args, err := squirrel.Select("1").
		From(ent.TableName()).
		Where(squirrel.Eq{ent.FieldName(&ent.LogoID): logoID.String()}).
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

func (r *SiteLogoRepository) FindBySiteURL(ctx context.Context, siteURL string) (*dmedia.SiteLogo, error) {
	var ent edownload.SiteLogo

	// Build SQL query
	sqlQuery, args, err := squirrel.Select(ent.FieldsAll()...).
		From(ent.TableName()).
		Where(squirrel.Eq{ent.FieldName(&ent.SiteURL): siteURL}).
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
	logo, err := r.mappers.MapSiteLogoEntityToDomain(&ent)
	if err != nil {
		return nil, err
	}

	return logo, nil
}

func (r *SiteLogoRepository) ExistsBySiteURL(ctx context.Context, siteURL string) (bool, error) {
	var ent edownload.SiteLogo

	// Build SQL query: SELECT 1 FROM table WHERE site_url = $1 LIMIT 1
	query, args, err := squirrel.Select("1").
		From(ent.TableName()).
		Where(squirrel.Eq{ent.FieldName(&ent.SiteURL): siteURL}).
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
