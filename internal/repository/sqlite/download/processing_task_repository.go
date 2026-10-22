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
	"github.com/neosy/elengrab/internal/repository/sqlite/types"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type ProcessingTaskRepository struct {
	mappers *mappers.Mappers
	dbEntry persistence.DBEntry

	queryOptions types.QueryMediaOptions

	// options
	retryOptions dbexec.RetryOptions
}

func NewProcessingTaskRepository(dbEntry persistence.DBEntry) persistence.MediaProcessingTaskRepositoryFactory {
	return func() persistence.MediaProcessingTaskRepository {
		return &ProcessingTaskRepository{
			mappers: mappers.NewMappers(),
			dbEntry: dbEntry,

			queryOptions: types.NewQueryMediaOptions(),

			// options
			retryOptions: dbexec.RetryOptions{
				MaxRetries: maxRetriesDefault,
				Delay:      retryDelayDefault,
			},
		}
	}
}

func (r *ProcessingTaskRepository) Insert(ctx context.Context, task *ddownload.MediaProcessingTask) error {
	return r.Write(ctx, task)
}

func (r *ProcessingTaskRepository) Update(ctx context.Context, task *ddownload.MediaProcessingTask) error {
	return r.Write(ctx, task)
}

func (r *ProcessingTaskRepository) Write(ctx context.Context, task *ddownload.MediaProcessingTask) error {
	if task == nil {
		return ierrors.ErrFuncParamNullPointer
	}

	// Convert the domain model to a database entity
	eTask, err := r.mappers.MapMediaProcessingTaskToEntity(task)
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

func (r *ProcessingTaskRepository) UpdateStatusToNew(ctx context.Context) error {
	var eTask edownload.MediaProcessingTask

	sqlWhere := squirrel.Or{
		squirrel.Eq{eTask.FieldName(&eTask.Status): dtypes.MediaTaskStatusPending.String()},
		squirrel.Eq{eTask.FieldName(&eTask.Status): dtypes.MediaTaskStatusWorking.String()},
	}

	// Build query
	sqlBuilder := squirrel.Update(eTask.TableName()).
		Set(eTask.FieldName(&eTask.Status), dtypes.MediaTaskStatusNew.String()).
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

func (r *ProcessingTaskRepository) findBy(ctx context.Context, filterFieldName string, filterValue any) (*ddownload.MediaProcessingTask, error) {
	var eTask edownload.MediaProcessingTask

	sqlWhere := squirrel.And{
		squirrel.Eq{filterFieldName: filterValue},
	}

	for name, filter := range r.queryOptions.Filters {
		if name != "" {
			sqlWhere = append(sqlWhere, filter.SqlCondition())
		}
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
	task, err := r.mappers.MapMediaProcessingTaskEntityToDomain(&eTask)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *ProcessingTaskRepository) FindByTaskID(ctx context.Context, taskID uuid.UUID) (*ddownload.MediaProcessingTask, error) {
	var eTask edownload.MediaProcessingTask
	return r.findBy(ctx, eTask.FieldName(&eTask.TaskID), taskID.String())
}

func (r *ProcessingTaskRepository) FindByDownloadID(ctx context.Context, downloadID uuid.UUID) (*ddownload.MediaProcessingTask, error) {
	var eTask edownload.MediaProcessingTask
	return r.findBy(ctx, eTask.FieldName(&eTask.DownloadID), downloadID.String())
}

func (r *ProcessingTaskRepository) deleteBy(ctx context.Context, filterFieldName string, filterValue any) error {
	var eTask edownload.MediaProcessingTask

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

func (r *ProcessingTaskRepository) Delete(ctx context.Context, taskID uuid.UUID) error {
	var eTask edownload.MediaProcessingTask
	return r.deleteBy(ctx, eTask.FieldName(&eTask.TaskID), taskID.String())
}

func (r *ProcessingTaskRepository) DeleteByDownloadID(ctx context.Context, downloadID uuid.UUID) error {
	var eTask edownload.MediaProcessingTask
	return r.deleteBy(ctx, eTask.FieldName(&eTask.DownloadID), downloadID.String())
}

func (r *ProcessingTaskRepository) DeleteByStatus(ctx context.Context, status dtypes.MediaTaskStatus) error {
	var eTask edownload.MediaProcessingTask
	return r.deleteBy(ctx, eTask.FieldName(&eTask.Status), status.String())
}

func (r *ProcessingTaskRepository) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.Tx(ctx, r.dbEntry, fn)
}

func (r *ProcessingTaskRepository) TxIndependent(ctx context.Context, fn func(ctx context.Context) error) error {
	return dbexec.TxIndependent(ctx, r.dbEntry, fn)
}
