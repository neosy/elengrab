package auth

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/pkg/dbutils"
	"github.com/neosy/elengrab/internal/ports/persistence"
	eauth "github.com/neosy/elengrab/internal/repository/sqlite/auth/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/auth/mappers"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	"github.com/neosy/elengrab/internal/repository/sqlite/types"
)

type RoleRepository struct {
	mappers *mappers.Mappers
	dbEntry persistence.DBEntry

	filtersByName types.FiltersByName
	queryOptions  queryOptions

	// options
	retryOptions dbexec.RetryOptions
}

// NewRoleRepository returns a new object for the repository
func NewRoleRepository(dbEntry persistence.DBEntry) persistence.RoleRepositoryFactory {
	return func() persistence.RoleRepository {
		return &RoleRepository{
			mappers: mappers.NewMappers(),
			dbEntry: dbEntry,

			filtersByName: make(map[string]any),

			// options
			retryOptions: dbexec.RetryOptions{
				MaxRetries: maxRetriesDefault,
				Delay:      retryDelayDefault,
			},
		}
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
		return ierrors.ErrFuncParamNullPointer
	}

	// Convert the domain model to a database entity
	eRole, err := r.mappers.MapRoleDomainToEntity(role)
	if err != nil {
		return err
	}

	// Get the list of fields and values for insertion
	fields := eRole.InsertFields()
	values := eRole.InsertValues()

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
	err = dbexec.ExecContext(ctx, r.dbEntry, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save role: %v", err)
	}

	return nil
}

func (r *RoleRepository) Find(ctx context.Context, roleID string) (*dauth.Role, error) {
	var eRole eauth.Role

	sqlQuery, args, err := squirrel.Select(eRole.QueryFields()...).
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
	db := dbexec.Resolve(ctx, r.dbEntry)
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
	db := dbexec.Resolve(ctx, r.dbEntry)
	var exists bool
	execQuery := func() error {
		row := db.QueryRowContext(ctx, query, args...)
		var dummy int
		err := row.Scan(&dummy)
		if err == sql.ErrNoRows {
			exists = false
			return nil
		} else if err != nil {
			return err
		}
		exists = true
		return nil
	}

	err = dbexec.ExecRetry(ctx, r.retryOptions, execQuery)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *RoleRepository) GetAll(ctx context.Context) ([]*dauth.Role, error) {
	roles := make([]*dauth.Role, 0)

	err := r.iterateGetAll(
		ctx, dbutils.OrderAsc,
		func(role *dauth.Role) error {
			roles = append(roles, role)
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *RoleRepository) IterateGetAll(ctx context.Context, fn func(*dauth.Role) error) error {
	return r.iterateGetAll(ctx, dbutils.OrderAsc, fn)
}

func (r *RoleRepository) iterateGetAll(
	ctx context.Context,
	sortOrderBy string,
	fn func(*dauth.Role) error,
) error {
	var eRole eauth.Role

	var sqlWhere = squirrel.And{}
	for name, value := range r.filtersByName {
		if name != "" {
			sqlWhere = append(sqlWhere, squirrel.Eq{eRole.FieldName(eRole.FieldPointer(name)): value})
		}
	}

	if r.queryOptions.withoutGuest != nil && *r.queryOptions.withoutGuest {
		sqlWhere = append(sqlWhere, squirrel.NotEq{eRole.FieldName(&eRole.RoleID): dtypes.UserRoleGuest.String()})
	}

	// Create an ORDER BY clause based on fieldы with the specified sort order.
	orderBy := dbutils.OrderBy(
		dbutils.Flds{
			eRole.FieldName(&eRole.RoleID): sortOrderBy,
		})

	qb := squirrel.Select(eRole.QueryFields()...).
		From(eRole.TableName()).
		Where(sqlWhere).
		OrderBy(orderBy).
		PlaceholderFormat(squirrel.Dollar)

	if r.queryOptions.Limit != nil && *r.queryOptions.Limit > 0 {
		qb = qb.Limit(*r.queryOptions.Limit)
	}

	sqlQuery, args, err := qb.ToSql()
	if err != nil {
		return fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	db := dbexec.Resolve(ctx, r.dbEntry)
	rows, err := db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows != nil {
		for rows.Next() {
			err := rows.Scan(eRole.FieldPointers()...)
			if err != nil {
				return err
			}

			role, err := r.mappers.MapRoleEntityToDomain(&eRole)
			if err != nil {
				return err
			}

			err = fn(role)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *RoleRepository) WithFilters(filters map[string]any) persistence.RoleRepository {
	if len(filters) == 0 {
		return r
	}

	var (
		eRole eauth.Role

		fieldNameByAllowedFilter = map[string]string{
			"roleId": eRole.FieldName(&eRole.RoleID),
		}
	)

	for name, value := range filters {
		fieldName, exists := fieldNameByAllowedFilter[name]
		if exists {
			r.filtersByName[fieldName] = value
		}

	}

	return r
}

func (r *RoleRepository) WithoutGuest() persistence.RoleRepository {
	r.queryOptions.withoutGuest = new(true)
	return r
}

func (r *RoleRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.Tx(ctx, r.dbEntry, fn)
}

func (r *RoleRepository) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.TxIndependent(ctx, r.dbEntry, fn)
}
