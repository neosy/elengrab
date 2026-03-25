package download

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/dbutils"
	"github.com/neosy/elengrab/internal/repository/sqlite/dbexec"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/download/mappers"
)

type DownloadTaskRepository struct {
	mappers *mappers.Mappers
	db      *sql.DB
	lock    dbexec.WriteLocker

	// options
	retryOptions dbexec.RetryOptions
}

type taskByFields struct {
	fileID *uuid.UUID
	status *dtypes.DownloadTaskStatus
}

// NewTaskRepository returns a new object for the repository
func NewDownloadTaskRepository(db *sql.DB, lock dbexec.WriteLocker) *DownloadTaskRepository {
	return &DownloadTaskRepository{
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

func (r *DownloadTaskRepository) Insert(ctx context.Context, task *ddownload.DownloadTask) error {
	return r.save(ctx, task)
}

func (r *DownloadTaskRepository) Update(ctx context.Context, task *ddownload.DownloadTask) error {
	return r.save(ctx, task)
}

func (r *DownloadTaskRepository) save(ctx context.Context, task *ddownload.DownloadTask) error {
	if task == nil {
		return errors.New("function parameter is a null pointer")
	}

	// Convert the domain model to a database entity
	eTask, err := r.mappers.MapDownloadTaskDomainToEntity(task)
	if err != nil {
		return err
	}

	// Get the list of fields and values for insertion
	fields := eTask.Fields()
	values := eTask.Values()

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
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save task: %v", err)
	}

	return nil
}

func (r *DownloadTaskRepository) UpdateStatusToNew(ctx context.Context) error {
	var ent edownload.DownloadTask

	sqlWhere := squirrel.Or{
		squirrel.Eq{ent.FieldName(&ent.Status): dtypes.DownloadTaskStatusPending.String()},
		squirrel.Eq{ent.FieldName(&ent.Status): dtypes.DownloadTaskStatusWorking.String()},
	}

	// Build query
	sqlBuilder := squirrel.Update(ent.TableName()).
		Set(ent.FieldName(&ent.Status), dtypes.DownloadTaskStatusNew.String()).
		Where(sqlWhere).
		PlaceholderFormat(squirrel.Dollar)

	// Generate SQL and args
	sqlStr, args, err := sqlBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlStr, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to update file: %v", err)
	}

	return nil
}

func (r *DownloadTaskRepository) FindByTaskID(ctx context.Context, taskID uuid.UUID) (*ddownload.DownloadTask, error) {
	var ent edownload.DownloadTask

	sqlQuery, args, err := squirrel.Select(ent.FieldsAll()...).
		From(ent.TableName()).
		Where(squirrel.Eq{ent.FieldName(&ent.TaskID): taskID.String()}).
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
		err := row.Scan(ent.FieldPointers()...)
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
	task, err := r.mappers.MapDownloadTaskEntityToDomain(&ent)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *DownloadTaskRepository) FindByFileID(ctx context.Context, fileID uuid.UUID) (*ddownload.DownloadTask, error) {
	var ent edownload.DownloadTask

	sqlQuery, args, err := squirrel.Select(ent.FieldsAll()...).
		From(ent.TableName()).
		Where(squirrel.Eq{ent.FieldName(&ent.FileID): fileID.String()}).
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
		err := row.Scan(ent.FieldPointers()...)
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
	task, err := r.mappers.MapDownloadTaskEntityToDomain(&ent)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *DownloadTaskRepository) Delete(ctx context.Context, taskID uuid.UUID) error {
	var ent edownload.DownloadTask

	// Build DELETE query
	sqlBuilder := squirrel.Delete(ent.TableName()).
		Where(squirrel.Eq{ent.FieldName(&ent.TaskID): taskID.String()}).
		PlaceholderFormat(squirrel.Dollar)

	// Generate SQL and args
	sqlStr, args, err := sqlBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlStr, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}

	return nil
}

func (r *DownloadTaskRepository) deleteBy(ctx context.Context, byFields taskByFields) error {
	var ent edownload.DownloadTask

	sqlWhere := squirrel.Expr("TRUE")
	if byFields.fileID != nil {
		sqlWhere = squirrel.Eq{ent.FieldName(&ent.FileID): byFields.fileID.String()}
	} else if byFields.status != nil {
		sqlWhere = squirrel.Eq{ent.FieldName(&ent.Status): byFields.status.String()}
	}

	// Build DELETE query
	sqlBuilder := squirrel.Delete(ent.TableName()).
		Where(sqlWhere).
		PlaceholderFormat(squirrel.Dollar)

	// Generate SQL and args
	sqlStr, args, err := sqlBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	err = dbexec.ExecContext(ctx, r.db, r.lock, sqlStr, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}

	return nil
}

func (r *DownloadTaskRepository) DeleteByFileID(ctx context.Context, fileID uuid.UUID) error {
	return r.deleteBy(ctx, taskByFields{fileID: &fileID})
}

func (r *DownloadTaskRepository) DeleteByStatus(ctx context.Context, status dtypes.DownloadTaskStatus) error {
	return r.deleteBy(ctx, taskByFields{status: &status})
}
