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

type UserSessionRepository struct {
	mappers *mappers.Mappers
	dbEntry persistence.DBEntry

	// options
	retryOptions dbexec.RetryOptions
}

// NewUserSessionRepository returns a new object for the repository
func NewUserSessionRepository(dbEntry persistence.DBEntry) persistence.UserSessionRepositoryFactory {
	return func() persistence.UserSessionRepository {
		return &UserSessionRepository{
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

func (r *UserSessionRepository) Insert(ctx context.Context, session *dauth.UserSession) error {
	return r.Save(ctx, session)
}

func (r *UserSessionRepository) Update(ctx context.Context, session *dauth.UserSession) error {
	return r.Save(ctx, session)
}

func (r *UserSessionRepository) Save(ctx context.Context, session *dauth.UserSession) error {
	if session == nil {
		return ierrors.ErrFuncParamNullPointer
	}

	// Convert the domain model to a database entity
	eSession, err := r.mappers.MapUserSessionDomainToEntity(session)
	if err != nil {
		return err
	}

	// Get the list of fields and values for insertion
	fields := eSession.InsertFields()
	values := eSession.InsertValues()

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
	err = dbexec.ExecContext(ctx, r.dbEntry, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save user session: %v", err)
	}

	return nil
}

func (r *UserSessionRepository) FindBySessionID(ctx context.Context, sessionID uuid.UUID) (*dauth.UserSession, error) {
	var eSession eauth.UserSession

	sqlQuery, args, err := squirrel.Select(eSession.QueryFields()...).
		From(eSession.TableName()).
		Where(squirrel.Eq{eSession.FieldName(&eSession.SessionID): sessionID}).
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
		err := row.Scan(eSession.FieldPointers()...)
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
	session, err := r.mappers.MapUserSessionEntityToDomain(&eSession)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *UserSessionRepository) FindByToken(ctx context.Context, token string) (*dauth.UserSession, error) {
	var eSession eauth.UserSession

	sqlQuery, args, err := squirrel.Select(eSession.QueryFields()...).
		From(eSession.TableName()).
		Where(squirrel.Eq{eSession.FieldName(&eSession.SessionToken): token}).
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
		err := row.Scan(eSession.FieldPointers()...)
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
	session, err := r.mappers.MapUserSessionEntityToDomain(&eSession)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *UserSessionRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.Tx(ctx, r.dbEntry, fn)
}

func (r *UserSessionRepository) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.TxIndependent(ctx, r.dbEntry, fn)
}
