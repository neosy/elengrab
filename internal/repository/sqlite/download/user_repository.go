package sldownload

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/download/mappers"
	"github.com/neosy/elengrab/internal/repository/sqlite/lock"
	"github.com/neosy/elengrab/pkg/dbutils"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type UserRepository struct {
	mappers *mappers.Mappers
	db      *sql.DB
	lock    lock.WriteLocker

	// options
	retryOptions retryOptions
}

// NewUserRepository returns a new object for the repository
func NewUserRepository(db *sql.DB, lock lock.WriteLocker) *UserRepository {
	return &UserRepository{
		mappers: mappers.NewMappers(),
		db:      db,
		lock:    lock,

		// options
		retryOptions: retryOptions{
			maxRetries: maxRetriesDefault,
			delay:      retryDelayDefault,
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
	err = execContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save user: %v", err)
	}

	return nil
}

func (r *UserRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*dauth.User, error) {
	var eUser edownload.User

	query, args, err := squirrel.Select(eUser.FieldsAll()...).
		From(eUser.TableName()).
		Where(squirrel.Eq{eUser.FieldName(&eUser.UserID): userID}).
		PlaceholderFormat(squirrel.Dollar).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	db := dbOrTx(ctx, r.db)

	for i := range r.retryOptions.maxRetries {
		// Scan result into entity
		err := db.QueryRowContext(ctx, query, args...).Scan(eUser.FieldPointers()...)
		if err == nil {
			break
		}
		if err == sql.ErrNoRows {
			return nil, nil
		}

		if sqlError, ok := err.(*sqlite.Error); ok && sqlError.Code() == sqlite3.SQLITE_BUSY {
			if i+1 == r.retryOptions.maxRetries {
				return nil, fmt.Errorf("failed to scan row: %w", err)
			}

			timer := time.NewTimer(r.retryOptions.delay)

			select {
			case <-ctx.Done():
				timer.Stop()
				return nil, ctx.Err()
			case <-timer.C:
				// Let's continue
			}

			continue
		}

		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	// Map entity to domain model
	user, err := r.mappers.MapUserEntityToDomain(&eUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error) {
	var eUser edownload.User

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
	db := dbOrTx(ctx, r.db)

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

	if err := fn(ctxWithTx(ctx, tx)); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
