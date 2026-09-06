package download

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/pkg/dbutils"
	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/download/mappers"
)

type DataMigrationRepository struct {
	mappers *mappers.Mappers
	dbEntry persistence.DBEntry

	// options
	retryOptions dbexec.RetryOptions
}

// NewDataMigrationRepository returns a new object for the repository
func NewDataMigrationRepository(dbEntry persistence.DBEntry) persistence.DownloadDataMigrationRepositoryFactory {
	return func() persistence.DownloadDataMigrationRepository {
		return &DataMigrationRepository{
			mappers: mappers.NewMappers(),
			dbEntry: dbEntry,

			// options
			retryOptions: dbexec.RetryOptions{
				MaxRetries: maxRetriesDefault,
				Delay:      retryDelayDefault,
			},
		}
	}
}

func (r *DataMigrationRepository) Insert(ctx context.Context, migration *ddownload.DataMigration) error {
	return r.Save(ctx, migration)
}

func (r *DataMigrationRepository) Update(ctx context.Context, migration *ddownload.DataMigration) error {
	return r.Save(ctx, migration)
}

func (r *DataMigrationRepository) Save(ctx context.Context, migration *ddownload.DataMigration) error {
	if migration == nil {
		return ierrors.ErrFuncParamNullPointer
	}

	// Convert the domain model to a database entity
	eMigration, err := r.mappers.MapDataMigrationDomainToEntity(migration)
	if err != nil {
		return err
	}

	// Get the list of fields and values for insertion
	fields := eMigration.InsertFields()
	values := eMigration.InsertValues()

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Insert(eMigration.TableName()).
		Columns(fields...).
		Values(values...).
		Suffix(dbutils.UpsertSuffix(fields, eMigration.FieldName(&eMigration.MigrationID))).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.dbEntry, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save siteMigration: %v", err)
	}

	return nil
}

func (r *DataMigrationRepository) Find(
	ctx context.Context, migrationID string) (*ddownload.DataMigration, error) {
	var ent edownload.DataMigration

	// Build SQL query
	sqlQuery, args, err := squirrel.Select(ent.QueryFields()...).
		From(ent.TableName()).
		Where(squirrel.Eq{ent.FieldName(&ent.MigrationID): migrationID}).
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
	migragion, err := r.mappers.MapDataMigrationEntityToDomain(&ent)
	if err != nil {
		return nil, err
	}

	return migragion, nil
}

func (r *DataMigrationRepository) Exists(ctx context.Context, migrationID string) (bool, error) {
	var ent edownload.DataMigration

	query, args, err := squirrel.Select("1").
		From(ent.TableName()).
		Where(squirrel.Eq{ent.FieldName(&ent.MigrationID): migrationID}).
		PlaceholderFormat(squirrel.Dollar).
		Limit(1).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	db := dbexec.Resolve(ctx, r.dbEntry)

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

func (r *DataMigrationRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.Tx(ctx, r.dbEntry, fn)
}
