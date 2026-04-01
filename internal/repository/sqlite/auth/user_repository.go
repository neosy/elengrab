package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
	eauth "github.com/neosy/elengrab/internal/repository/sqlite/auth/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/auth/mappers"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
)

type UserRepository struct {
	mappers *mappers.Mappers
	db      *sql.DB
	lock    dbexec.WriteLocker

	// options
	retryOptions dbexec.RetryOptions
}

// NewUserRepository returns a new object for the repository
func NewUserRepository(db *sql.DB, lock dbexec.WriteLocker) *UserRepository {
	return &UserRepository{
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

func (r *UserRepository) Insert(ctx context.Context, user *dauth.User) error {
	if user == nil {
		return errors.New("function parameter is a null pointer")
	}

	// Convert the domain model to a database entity
	eUser, err := r.mappers.MapUserDomainToEntity(user)
	if err != nil {
		return err
	}

	// Update password hash and timestamp if it's not empty
	if eUser.PasswordHash != nil {
		eUser.PasswordUpdatedAt = uptr.Any(time.Now().UTC())
	}

	// Get the list of fields and values for insertion
	fields := eUser.Fields()
	values := eUser.Values()

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Insert(eUser.TableName()).
		Columns(fields...).
		Values(values...).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save user: %v", err)
	}

	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *dauth.User) error {
	if user == nil {
		return errors.New("function parameter is a null pointer")
	}

	// Convert the domain model to a database entity
	eUser, err := r.mappers.MapUserDomainToEntity(user)
	if err != nil {
		return err
	}

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Update(eUser.TableName()).
		SetMap(eUser.FieldsMap()).
		Where(squirrel.Eq{eUser.FieldName(&eUser.UserID): eUser.UserID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save user: %v", err)
	}

	return nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, newPasswHash string) error {
	if userID == uuid.Nil {
		return errors.New("user ID is nil")
	}

	var eUser eauth.User

	fieldsMap := map[string]any{
		eUser.FieldName(eUser.PasswordHash):      newPasswHash,
		eUser.FieldName(eUser.PasswordUpdatedAt): squirrel.Expr("CURRENT_TIMESTAMP"),
	}

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Update(eUser.TableName()).
		SetMap(fieldsMap).
		Where(squirrel.Eq{eUser.FieldName(&eUser.UserID): userID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save user: %v", err)
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, userID uuid.UUID, soft bool) error {
	if soft {
		return r.softDelete(ctx, userID)
	} else {
		return r.hardDelete(ctx, userID)
	}
}

func (r *UserRepository) softDelete(ctx context.Context, userID uuid.UUID) error {
	var eUser eauth.User

	fieldsToUpdate := map[string]interface{}{
		eUser.FieldName(&eUser.DeletedAt): squirrel.Expr("CURRENT_TIMESTAMP"),
	}

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Update(eUser.TableName()).
		SetMap(fieldsToUpdate).
		Where(squirrel.Eq{eUser.FieldName(&eUser.UserID): userID}).
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

func (r *UserRepository) hardDelete(ctx context.Context, userID uuid.UUID) error {
	var eUser eauth.User

	// Build DELETE query
	sqlBuilder := squirrel.Delete(eUser.TableName()).
		Where(squirrel.Eq{eUser.FieldName(&eUser.UserID): userID.String()}).
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

func (r *UserRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*dauth.User, error) {
	var eUser eauth.User
	return r.findByFieldName(ctx, eUser.FieldName(&eUser.UserID), userID)
}

func (r *UserRepository) FindByLogin(ctx context.Context, login dtypes.Login) (*dauth.User, error) {
	var eUser eauth.User
	return r.findByFieldName(ctx, eUser.FieldName(&eUser.Login), login.String())
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*dauth.User, error) {
	var eUser eauth.User
	return r.findByFieldName(ctx, eUser.FieldName(&eUser.Email), email)
}

func (r *UserRepository) findByFieldName(ctx context.Context, fieldName string, value any) (*dauth.User, error) {
	var (
		eUser     eauth.User
		eUserRole eauth.UserRole

		aliasUsers     = "u"
		aliasUserRoles = "r"
	)

	selectFields := append(
		eUser.FieldsAllWithAlias(aliasUsers),
		"GROUP_CONCAT("+eUserRole.FieldNameWithAlias(&eUserRole.RoleID, aliasUserRoles)+") AS roles",
	)

	var sqlWhere squirrel.Sqlizer

	if fieldName == eUser.FieldName(&eUser.Login) {
		sqlWhere = squirrel.Expr(
			eUser.FieldNameWithAlias(eUser.FieldPointer(fieldName), aliasUsers)+" = ? COLLATE NOCASE",
			value,
		)

	} else {
		sqlWhere = squirrel.Eq{
			eUser.FieldNameWithAlias(eUser.FieldPointer(fieldName), aliasUsers): value,
		}
	}

	sqlQuery, args, err := squirrel.Select(selectFields...).
		From(eUser.TableName() + " AS " + aliasUsers).
		LeftJoin(
			eUserRole.TableName() + " AS " + aliasUserRoles +
				" ON " + eUserRole.FieldNameWithAlias(&eUserRole.UserID, aliasUserRoles) +
				" = " + eUser.FieldNameWithAlias(&eUser.UserID, aliasUsers),
		).
		Where(sqlWhere).
		GroupBy(eUser.FieldNameWithAlias(&eUser.UserID, aliasUsers)).
		PlaceholderFormat(squirrel.Dollar).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	var (
		notFound bool
		roles    string
	)
	db := dbexec.Resolve(ctx, r.db)
	execQuery := func() error {
		row := db.QueryRowContext(ctx, sqlQuery, args...)
		// Scan result into entity
		err := row.Scan(append(eUser.FieldPointers(), &roles)...)
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
	user, err := r.mappers.MapUserEntityToDomain(&eUser, roles)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error) {
	var eUser eauth.User
	return r.existsByFieldName(ctx, eUser.FieldName(&eUser.UserID), userID)
}

func (r *UserRepository) ExistsByLogin(ctx context.Context, login dtypes.Login) (bool, error) {
	var eUser eauth.User
	return r.existsByFieldName(ctx, eUser.FieldName(&eUser.Login), login.String())
}

func (r *UserRepository) existsByFieldName(ctx context.Context, fieldName string, value any) (bool, error) {
	var (
		eUser    eauth.User
		sqlWhere squirrel.Sqlizer
	)

	if fieldName == eUser.FieldName(&eUser.Login) {
		sqlWhere = squirrel.Expr(
			eUser.FieldName(eUser.FieldPointer(fieldName))+" = ? COLLATE NOCASE",
			value,
		)

	} else {
		sqlWhere = squirrel.Eq{
			eUser.FieldName(eUser.FieldPointer(fieldName)): value,
		}
	}

	// Build SQL query: SELECT 1 FROM table WHERE <id> = $1 LIMIT 1
	query, args, err := squirrel.Select("1").
		From(eUser.TableName()).
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

func (r *UserRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.Tx(ctx, r.db, r.lock, fn)
}
