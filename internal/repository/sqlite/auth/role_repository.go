package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/pkg/dbutils"
	eauth "github.com/neosy/elengrab/internal/repository/sqlite/auth/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/auth/mappers"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
)

type RoleRepository struct {
	mappers *mappers.Mappers
	db      *sql.DB
	lock    dbexec.WriteLocker

	// options
	retryOptions dbexec.RetryOptions
}

// NewRoleRepository returns a new object for the repository
func NewRoleRepository(db *sql.DB, lock dbexec.WriteLocker) *RoleRepository {
	return &RoleRepository{
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

func (r *RoleRepository) Insert(ctx context.Context, role *dauth.Role) error {
	return r.Save(ctx, role)
}

func (r *RoleRepository) Update(ctx context.Context, role *dauth.Role) error {
	return r.Save(ctx, role)
}

func (r *RoleRepository) Save(ctx context.Context, role *dauth.Role) error {
	if role == nil {
		return errors.New("function parameter is a null pointer")
	}

	// Convert the domain model to a database entity
	eRole, err := r.mappers.MapRoleDomainToEntity(role)
	if err != nil {
		return err
	}

	// Get the list of fields and values for insertion
	fields := eRole.Fields()
	values := eRole.Values()

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Insert(eRole.TableName()).
		Columns(fields...).
		Values(values...).
		Suffix(dbutils.UpsertSuffix(fields, eRole.FieldName(&eRole.RoleID))).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save role: %v", err)
	}

	return nil
}

func (r *RoleRepository) Find(ctx context.Context, roleID string) (*dauth.Role, error) {
	var eRole eauth.Role

	sqlQuery, args, err := squirrel.Select(eRole.FieldsAll()...).
		From(eRole.TableName()).
		Where(squirrel.Eq{eRole.FieldName(&eRole.RoleID): roleID}).
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
		err := row.Scan(eRole.FieldPointers()...)
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
	role, err := r.mappers.MapRoleEntityToDomain(&eRole)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (r *RoleRepository) Exists(ctx context.Context, roleID string) (bool, error) {
	var eRole eauth.Role

	// Build SQL query: SELECT 1 FROM table WHERE <id> = $1 LIMIT 1
	query, args, err := squirrel.Select("1").
		From(eRole.TableName()).
		Where(squirrel.Eq{eRole.FieldName(&eRole.RoleID): roleID}).
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

func (r *RoleRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.Tx(ctx, r.db, r.lock, fn)
}
