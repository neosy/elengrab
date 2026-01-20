package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	eauth "github.com/neosy/elengrab/internal/repository/sqlite/auth/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/auth/mappers"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	"github.com/neosy/elengrab/pkg/dbutils"
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
	return r.Save(ctx, user)
}

func (r *UserRepository) Update(ctx context.Context, user *dauth.User) error {
	return r.Save(ctx, user)
}

func (r *UserRepository) Save(ctx context.Context, user *dauth.User) error {
	if user == nil {
		return errors.New("function parameter is a null pointer")
	}

	// Convert the domain model to a database entity
	eUser, err := r.mappers.MapUserDomainToEntity(user)
	if err != nil {
		return err
	}

	// Get the list of fields and values for insertion
	fields := eUser.Fields()
	values := eUser.Values()

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Insert(eUser.TableName()).
		Columns(fields...).
		Values(values...).
		Suffix(dbutils.UpsertSuffix(fields, eUser.FieldName(&eUser.UserID))).
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

func (r *UserRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*dauth.User, error) {
	var eUser eauth.User

	sqlQuery, args, err := squirrel.Select(eUser.FieldsAll()...).
		From(eUser.TableName()).
		Where(squirrel.Eq{eUser.FieldName(&eUser.UserID): userID}).
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
		err := row.Scan(eUser.FieldPointers()...)
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
	user, err := r.mappers.MapUserEntityToDomain(&eUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error) {
	var eUser eauth.User

	// Build SQL query: SELECT 1 FROM table WHERE <id> = $1 LIMIT 1
	query, args, err := squirrel.Select("1").
		From(eUser.TableName()).
		Where(squirrel.Eq{eUser.FieldName(&eUser.UserID): userID}).
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
