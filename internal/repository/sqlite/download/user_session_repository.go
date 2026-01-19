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

type UserSessionRepository struct {
	mappers *mappers.Mappers
	db      *sql.DB
	lock    lock.WriteLocker

	// options
	retryOptions retryOptions
}

// NewUserSessionRepository returns a new object for the repository
func NewUserSessionRepository(db *sql.DB, lock lock.WriteLocker) *UserSessionRepository {
	return &UserSessionRepository{
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

func (r *UserSessionRepository) Insert(ctx context.Context, session *dauth.UserSession) error {
	return r.Save(ctx, session)
}

func (r *UserSessionRepository) Update(ctx context.Context, session *dauth.UserSession) error {
	return r.Save(ctx, session)
}

func (r *UserSessionRepository) Save(ctx context.Context, session *dauth.UserSession) error {
	if session == nil {
		return errors.New("function parameter is a null pointer")
	}

	// Convert the domain model to a database entity
	eSession, err := r.mappers.MapUserSessionDomainToEntity(session)
	if err != nil {
		return err
	}

	// Get the list of fields and values for insertion
	fields := eSession.Fields()
	values := eSession.Values()

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Insert(eSession.TableName()).
		Columns(fields...).
		Values(values...).
		Suffix(dbutils.UpsertSuffix(fields, eSession.FieldName(&eSession.SessionID))).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = execContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save user session: %v", err)
	}

	return nil
}

func (r *UserSessionRepository) FindBySessionID(ctx context.Context, sessionID uuid.UUID) (*dauth.UserSession, error) {
	var eSession edownload.UserSession

	query, args, err := squirrel.Select(eSession.FieldsAll()...).
		From(eSession.TableName()).
		Where(squirrel.Eq{eSession.FieldName(&eSession.SessionID): sessionID}).
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
		err := db.QueryRowContext(ctx, query, args...).Scan(eSession.FieldPointers()...)
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
	session, err := r.mappers.MapUserSessionEntityToDomain(&eSession)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *UserSessionRepository) FindByToken(ctx context.Context, token string) (*dauth.UserSession, error) {
	var eSession edownload.UserSession

	query, args, err := squirrel.Select(eSession.FieldsAll()...).
		From(eSession.TableName()).
		Where(squirrel.Eq{eSession.FieldName(&eSession.SessionToken): token}).
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
		err := db.QueryRowContext(ctx, query, args...).Scan(eSession.FieldPointers()...)
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
	session, err := r.mappers.MapUserSessionEntityToDomain(&eSession)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *UserSessionRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
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
