package sldownload

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/download/mappers"
	"github.com/neosy/elengrab/pkg/dbutils"
)

type DownloadTaskRepository struct {
	mappers *mappers.Mappers
	db      *sql.DB
	mu      *sync.RWMutex

	// options
	retryOptions retryOptions
}

type taskByFields struct {
	fileId *uuid.UUID
	status *dtypes.DownloadTaskStatus
}

// NewTaskRepository returns a new object for the repository
func NewDownloadTaskRepository(db *sql.DB, mu *sync.RWMutex) *DownloadTaskRepository {
	return &DownloadTaskRepository{
		mappers: mappers.NewMappers(),
		db:      db,
		mu:      mu,

		// options
		retryOptions: retryOptions{
			maxRetries: maxRetriesDefault,
			delay:      retryDelayDefault,
		},
	}
}

func (r *DownloadTaskRepository) Insert(ctx context.Context, task *ddownload.DownloadTask) error {
	return r.save(ctx, task, false)
}

func (r *DownloadTaskRepository) Update(ctx context.Context, task *ddownload.DownloadTask) error {
	return r.save(ctx, task, true)
}

func (r *DownloadTaskRepository) save(ctx context.Context, task *ddownload.DownloadTask, isUpd bool) error {
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

	// If this is an update — add the UpdatedAt field with the current time
	if isUpd {
		fields = append(fields, eTask.FieldName(&eTask.UpdatedAt))
		values = append(values, squirrel.Expr("CURRENT_TIMESTAMP"))
	}

	// Generate SQL query with upsert logic
	sqlQuery, args, err := squirrel.
		Insert(eTask.TableName()).
		Columns(fields...).
		Values(values...).
		Suffix(dbutils.UpsertSuffix(fields, eTask.FieldName(&eTask.TaskId))).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Execute the query
	err = execContext(ctx, r.db, sqlQuery, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	return nil
}

func (r *DownloadTaskRepository) UpdateStatusToNew(ctx context.Context) error {
	var ent edownload.DownloadTask

	sqlWhere := squirrel.Or{
		squirrel.Eq{ent.FieldName(&ent.Status): dtypes.DownloadTaskStatusPending},
		squirrel.Eq{ent.FieldName(&ent.Status): dtypes.DownloadTaskStatusWorking},
	}

	// Build query
	sqlBuilder := squirrel.Update(ent.TableName()).
		Set(ent.FieldName(&ent.Status), dtypes.DownloadTaskStatusNew).
		Where(sqlWhere).
		PlaceholderFormat(squirrel.Dollar)

	// Generate SQL and args
	sqlStr, args, err := sqlBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error generating SQL: %v", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Execute the query
	err = execContext(ctx, r.db, sqlStr, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to update file: %v", err)
	}

	return nil
}

func (r *DownloadTaskRepository) FindByTaskId(ctx context.Context, taskId uuid.UUID) (*ddownload.DownloadTask, error) {
	var ent edownload.DownloadTask

	sqlBuilder, args, err := squirrel.Select(ent.FieldsAll()...).
		From(ent.TableName()).
		Where(squirrel.Eq{ent.FieldName(&ent.TaskId): taskId.String()}).
		PlaceholderFormat(squirrel.Dollar).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	row := r.db.QueryRowContext(ctx, sqlBuilder, args...)

	// Scan result into entity
	if err := row.Scan(ent.FieldPointers()...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // запись не найдена
		}
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	// Map entity to domain model
	task, err := r.mappers.MapDownloadTaskEntityToDomain(&ent)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *DownloadTaskRepository) FindByFileId(ctx context.Context, fileId uuid.UUID) (*ddownload.DownloadTask, error) {
	var ent edownload.DownloadTask

	sqlBuilder, args, err := squirrel.Select(ent.FieldsAll()...).
		From(ent.TableName()).
		Where(squirrel.Eq{ent.FieldName(&ent.FileId): fileId.String()}).
		PlaceholderFormat(squirrel.Dollar).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("error generating SQL: %v", err)
	}

	// Execute the query
	row := r.db.QueryRowContext(ctx, sqlBuilder, args...)

	// Scan result into entity
	if err := row.Scan(ent.FieldPointers()...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // запись не найдена
		}
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	// Map entity to domain model
	task, err := r.mappers.MapDownloadTaskEntityToDomain(&ent)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *DownloadTaskRepository) Delete(ctx context.Context, taskId uuid.UUID) error {
	var ent edownload.DownloadTask

	// Build DELETE query
	sqlBuilder := squirrel.Delete(ent.TableName()).
		Where(squirrel.Eq{ent.FieldName(&ent.TaskId): taskId.String()}).
		PlaceholderFormat(squirrel.Dollar)

	// Generate SQL and args
	sqlStr, args, err := sqlBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error generating SQL: %v", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Execute the query
	err = execContext(ctx, r.db, sqlStr, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}

	return nil
}

func (r *DownloadTaskRepository) deleteBy(ctx context.Context, byFields taskByFields) error {
	var ent edownload.DownloadTask

	sqlWhere := squirrel.Expr("TRUE")
	if byFields.fileId != nil {
		sqlWhere = squirrel.Eq{ent.FieldName(&ent.FileId): byFields.fileId.String()}
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

	r.mu.Lock()
	defer r.mu.Unlock()

	// Execute the query
	err = execContext(ctx, r.db, sqlStr, args, r.retryOptions)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}

	return nil
}

func (r *DownloadTaskRepository) DeleteByFileId(ctx context.Context, fileId uuid.UUID) error {
	return r.deleteBy(ctx, taskByFields{fileId: &fileId})
}

func (r *DownloadTaskRepository) DeleteByStatus(ctx context.Context, status dtypes.DownloadTaskStatus) error {
	return r.deleteBy(ctx, taskByFields{status: &status})
}
