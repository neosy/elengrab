package auth

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/pkg/dbutils"
	"github.com/neosy/elengrab/internal/ports/persistence"
	eauth "github.com/neosy/elengrab/internal/repository/sqlite/auth/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/auth/mappers"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
)

type UserRoleRepository struct {
	mappers *mappers.Mappers
	dbEntry persistence.DBEntry

	// options
	retryOptions dbexec.RetryOptions
}

// NewUserRoleRepository returns a new object for the repository
func NewUserRoleRepository(dbEntry persistence.DBEntry) persistence.UserRoleRepositoryFactory {
	return func() persistence.UserRoleRepository {
		return &UserRoleRepository{
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

func (r *UserRoleRepository) Insert(ctx context.Context, userRole *dauth.UserRole) error {
	return r.Save(ctx, userRole)
}

func (r *UserRoleRepository) Update(ctx context.Context, userRole *dauth.UserRole) error {
	return r.Save(ctx, userRole)
}

func (r *UserRoleRepository) Save(ctx context.Context, userRole *dauth.UserRole) error {
	if userRole == nil {
		return ierrors.ErrFuncParamNullPointer
	}

	// Convert the domain model to a database entity
	eUserRole, err := r.mappers.MapUserRoleDomainToEntity(userRole)
	if err != nil {
		return err
	}

	// Get the list of fields and values for insertion
	fields := eUserRole.Fields()
	values := eUserRole.Values()

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Insert(eUserRole.TableName()).
		Columns(fields...).
		Values(values...).
		Suffix(dbutils.UpsertSuffix(
			fields,
			eUserRole.FieldName(&eUserRole.UserID),
			eUserRole.FieldName(&eUserRole.RoleID),
		)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.dbEntry, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save userRole: %v", err)
	}

	return nil
}

func (r *UserRoleRepository) Delete(ctx context.Context, userID uuid.UUID, roleID string) error {
	var eUserRole eauth.UserRole

	sqlWhere := squirrel.And{
		squirrel.Eq{eUserRole.FieldName(&eUserRole.UserID): userID.String()},
		squirrel.Eq{eUserRole.FieldName(&eUserRole.RoleID): roleID},
	}

	// Build DELETE query
	sqlBuilder := squirrel.Delete(eUserRole.TableName()).
		Where(sqlWhere).
		PlaceholderFormat(squirrel.Dollar)

	// Generate SQL and args
	sqlStr, args, err := sqlBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.dbEntry, sqlStr, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to delete user role: %v", err)
	}

	return nil
}

func (r *UserRoleRepository) Find(
	ctx context.Context,
	userID uuid.UUID,
	roleID string,
) (*dauth.UserRole, error) {
	var eUserRole eauth.UserRole

	sqlQuery, args, err := squirrel.Select(eUserRole.FieldsAll()...).
		From(eUserRole.TableName()).
		Where(
			squirrel.And{
				squirrel.Eq{eUserRole.FieldName(&eUserRole.UserID): userID},
				squirrel.Eq{eUserRole.FieldName(&eUserRole.RoleID): roleID},
			}).
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
		err := row.Scan(eUserRole.FieldPointers()...)
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
	userRole, err := r.mappers.MapUserRoleEntityToDomain(&eUserRole)
	if err != nil {
		return nil, err
	}

	return userRole, nil
}

func (r *UserRoleRepository) Exists(
	ctx context.Context,
	userID uuid.UUID,
	roleID string,
) (bool, error) {
	var eUserRole eauth.UserRole

	// Build SQL query: SELECT 1 FROM table WHERE <id> = $1 LIMIT 1
	query, args, err := squirrel.Select("1").
		From(eUserRole.TableName()).
		Where(
			squirrel.And{
				squirrel.Eq{eUserRole.FieldName(&eUserRole.UserID): userID},
				squirrel.Eq{eUserRole.FieldName(&eUserRole.RoleID): roleID},
			}).
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

func (r *UserRoleRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.Tx(ctx, r.dbEntry, fn)
}

func (r *UserRoleRepository) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.TxIndependent(ctx, r.dbEntry, fn)
}
