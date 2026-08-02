package watchevent

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/pkg/dbutils"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	ewatchevent "github.com/neosy/elengrab/internal/repository/sqlite/watch_event/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/watch_event/mappers"
)

type MediaUserWatchPositionRepository struct {
	mappers *mappers.Mappers
	db      *sql.DB
	lock    dbexec.WriteLocker

	filtersByName filtersByName

	// options
	retryOptions dbexec.RetryOptions
}

// NewMediaUserWatchPositionRepository returns a new object for the repository
func NewMediaUserWatchPositionRepository(db *sql.DB, lock dbexec.WriteLocker) *MediaUserWatchPositionRepository {
	return &MediaUserWatchPositionRepository{
		mappers: mappers.NewMappers(),
		db:      db,
		lock:    lock,

		filtersByName: make(map[string]any),

		// options
		retryOptions: dbexec.RetryOptions{
			MaxRetries: maxRetriesDefault,
			Delay:      retryDelayDefault,
		},
	}
}

func (r *MediaUserWatchPositionRepository) Copy() *MediaUserWatchPositionRepository {
	rep := uptr.Copy(r)

	rep.mappers = r.mappers
	rep.db = r.db
	rep.lock = r.lock

	rep.filtersByName = rep.filtersByName.copy()

	return rep
}

func (r *MediaUserWatchPositionRepository) Insert(ctx context.Context, position *ddownload.MediaUserWatchPosition) error {
	return r.save(ctx, position)
}

func (r *MediaUserWatchPositionRepository) Update(ctx context.Context, position *ddownload.MediaUserWatchPosition) error {
	return r.save(ctx, position)
}

func (r *MediaUserWatchPositionRepository) Write(ctx context.Context, position *ddownload.MediaUserWatchPosition) error {
	return r.save(ctx, position)
}

func (r *MediaUserWatchPositionRepository) save(ctx context.Context, position *ddownload.MediaUserWatchPosition) error {
	if position == nil {
		return ierrors.ErrFuncParamNullPointer
	}

	// Convert the domain model to a database entity
	ePosition, err := r.mappers.MapMediaUserWatchPositionDomainToEntity(position)
	if err != nil {
		return err
	}

	// Get the list of fields and values for insertion
	fields := ePosition.Fields()
	values := ePosition.Values()

	// Build INSERT query
	sqlBuilder := squirrel.
		Insert(ePosition.TableName()).
		Columns(fields...).
		Values(values...).
		Suffix(dbutils.UpsertSuffix(fields, ePosition.ConflictFields()...)).
		PlaceholderFormat(squirrel.Dollar)

	// Generate SQL query with upsert logic
	sqlQuery, args, err := sqlBuilder.ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save record: %v", err)
	}

	return nil
}

func (r *MediaUserWatchPositionRepository) DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error {
	var ePosition ewatchevent.MediaUserWatchPosition

	// Build DELETE query
	sqlBuilder := squirrel.
		Delete(ePosition.TableName()).
		Where(squirrel.Eq{ePosition.FieldName(&ePosition.DownloadID): downloadID}).
		PlaceholderFormat(squirrel.Dollar)

	// Generate SQL and args
	sqlQuery, args, err := sqlBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to delete record: %v", err)
	}

	return nil
}

func (r *MediaUserWatchPositionRepository) Find(
	ctx context.Context,
	downloadID uuid.UUID,
	userID uuid.UUID,
	sessionID uuid.UUID,
) (*ddownload.MediaUserWatchPosition, error) {
	var ePosition ewatchevent.MediaUserWatchPosition

	var sessionIDFilter string
	if sessionID != uuid.Nil {
		sessionIDFilter = sessionID.String()
	}

	var sqlWhere = squirrel.And{}

	sqlWhere = append(sqlWhere,
		squirrel.Eq{
			ePosition.FieldName(&ePosition.DownloadID): downloadID,
			ePosition.FieldName(&ePosition.UserID):     userID,
			ePosition.FieldName(&ePosition.SessionID):  sessionIDFilter,
		},
	)

	// Build SELECT query
	sqlBuilder := squirrel.
		Select(ePosition.FieldsAll()...).
		From(ePosition.TableName()).
		Where(sqlWhere).
		PlaceholderFormat(squirrel.Dollar).
		Limit(1)

	// Generate SQL and args
	sqlQuery, args, err := sqlBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	var notFound bool
	db := dbexec.Resolve(ctx, r.db)
	execQuery := func() error {
		row := db.QueryRowContext(ctx, sqlQuery, args...)
		// Scan result into entity
		err := row.Scan(ePosition.FieldPointers()...)
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
	position, err := r.mappers.MapMediaUserWatchPositionEntityToDomain(&ePosition)
	if err != nil {
		return nil, err
	}

	return position, nil
}

func (r *MediaUserWatchPositionRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.Tx(ctx, r.db, r.lock, fn)
}

func (r *MediaUserWatchPositionRepository) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.TxIndependent(ctx, r.db, r.lock, fn)
}
