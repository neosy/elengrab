package download

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"uuid"

	"github.com/Masterminds/squirrel"
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/pkg/dbutils"
	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/download/mappers"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type DownloadTaskRepository struct {
	mappers *mappers.Mappers
	dbEntry persistence.DBEntry

	// options
	retryOptions dbexec.RetryOptions
}

// NewTaskRepository returns a new object for the repository
func NewDownloadTaskRepository(dbEntry persistence.DBEntry) persistence.DownloadTaskRepositoryFactory {
	return func() persistence.DownloadTaskRepository {
		return &DownloadTaskRepository{
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

func (r *DownloadTaskRepository) Insert(ctx context.Context, task *ddownload.DownloadTask) error {
	return r.save(ctx, task)
}

func (r *DownloadTaskRepository) Update(ctx context.Context, task *ddownload.DownloadTask) error {
	return r.save(ctx, task)
}

func (r *DownloadTaskRepository) save(ctx context.Context, task *ddownload.DownloadTask) error {
	if task == nil {
		return ierrors.ErrFuncParamNullPointer
	}

	// Convert the domain model to a database entity
	eTask, err := r.mappers.MapDownloadTaskDomainToEntity(task)
	if err != nil {
		return err
	}

	// Get the list of fields and values for insertion
	fields := eTask.InsertFields()
	values := eTask.InsertValues()

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Insert(eTask.TableName()).
		Columns(fields...).
		Values(values...).
		Suffix(dbutils.UpsertSuffix(fields, eTask.FieldName(&eTask.TaskID))).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) &&
			sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return apperrors.ErrIDConflict
		}

		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.dbEntry, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save task: %v", err)
	}

	return nil
}

func (r *DownloadTaskRepository) UpdateStatusToNew(ctx context.Context) error {
	var ent edownload.DownloadTask

	sqlWhere := squirrel.Or{
		squirrel.Eq{ent.FieldName(&ent.Status): dtypes.MediaTaskStatusPending.String()},
		squirrel.Eq{ent.FieldName(&ent.Status): dtypes.MediaTaskStatusWorking.String()},
	}

	// Build query
	sqlBuilder := squirrel.Update(ent.TableName()).
		Set(ent.FieldName(&ent.Status), dtypes.MediaTaskStatusNew.String()).
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
		return fmt.Errorf("failed to update file: %v", err)
	}

	return nil
}

func (r *DownloadTaskRepository) findBy(ctx context.Context, filterFieldName string, filterValue any) (*ddownload.DownloadTask, error) {
	var eTask edownload.DownloadTask

	sqlWhere := squirrel.And{
		squirrel.Eq{filterFieldName: filterValue},
	}

	sqlQuery, args, err := squirrel.Select(eTask.QueryFields()...).
		From(eTask.TableName()).
		Where(sqlWhere).
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
		err := row.Scan(eTask.FieldPointers()...)
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
	task, err := r.mappers.MapDownloadTaskEntityToDomain(&eTask)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *DownloadTaskRepository) FindByTaskID(ctx context.Context, taskID uuid.UUID) (*ddownload.DownloadTask, error) {
	var eTask edownload.DownloadTask
	return r.findBy(ctx, eTask.FieldName(&eTask.TaskID), taskID.String())
}

func (r *DownloadTaskRepository) FindByDownloadID(ctx context.Context, downloadID uuid.UUID) (*ddownload.DownloadTask, error) {
	var eTask edownload.DownloadTask
	return r.findBy(ctx, eTask.FieldName(&eTask.DownloadID), downloadID.String())
}

func (r *DownloadTaskRepository) deleteBy(ctx context.Context, filterFieldName string, filterValue any) error {
	var eTask edownload.DownloadTask

	if filterFieldName == "" {
		return nil
	}

	sqlWhere := squirrel.And{
		squirrel.Eq{filterFieldName: filterValue},
	}

	// Build DELETE query
	sqlBuilder := squirrel.Delete(eTask.TableName()).
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
		return fmt.Errorf("failed to delete task: %v", err)
	}

	return nil
}

func (r *DownloadTaskRepository) Delete(ctx context.Context, taskID uuid.UUID) error {
	var eTask edownload.DownloadTask
	return r.deleteBy(ctx, eTask.FieldName(&eTask.TaskID), taskID.String())
}

func (r *DownloadTaskRepository) DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error {
	var eTask edownload.DownloadTask
	return r.deleteBy(ctx, eTask.FieldName(&eTask.DownloadID), downloadID.String())
}

func (r *DownloadTaskRepository) DeleteByStatus(ctx context.Context, status dtypes.MediaTaskStatus) error {
	var eTask edownload.DownloadTask
	return r.deleteBy(ctx, eTask.FieldName(&eTask.Status), status.String())
}

func (r *DownloadTaskRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.Tx(ctx, r.dbEntry, fn)
}

func (r *DownloadTaskRepository) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.TxIndependent(ctx, r.dbEntry, fn)
}
