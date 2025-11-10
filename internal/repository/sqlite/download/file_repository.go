package sldownload

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/download/mappers"
	"github.com/neosy/elengrab/pkg/dbutils"
)

type FileRepository struct {
	db      *sql.DB
	mappers *mappers.Mappers
}

// NewFilesRepository returns a new object for the repository
func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{
		db:      db,
		mappers: mappers.NewMappers(),
	}
}

func (r *FileRepository) Insert(ctx context.Context, file *ddownload.File) error {
	return r.save(ctx, file, false)
}

func (r *FileRepository) Update(ctx context.Context, file *ddownload.File) error {
	return r.save(ctx, file, true)
}

func (r *FileRepository) save(ctx context.Context, file *ddownload.File, isUpd bool) error {
	if file == nil {
		return errors.New("function parameter is a null pointer")
	}

	// Convert the domain model to a database entity
	eFile := r.mappers.MapFileDomainToEntity(file)

	// Get the list of fields and values for insertion
	fields := eFile.Fields()
	values := eFile.Values()

	// If this is an update — add the UpdatedAt field with the current time
	if isUpd {
		fields = append(fields, eFile.FieldName(&eFile.UpdatedAt))
		values = append(values, time.Now())
	}

	// Generate SQL query with upsert logic
	sql, args, err := squirrel.
		Insert(eFile.TableName()).
		Columns(fields...).
		Values(values...).
		Suffix(dbutils.UpsertSuffix(fields, eFile.FieldName(&eFile.FileId))).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	// If SQL generation failed — return an error
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	// Execute the SQL query
	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to insert file: %w", err)
	}

	return nil
}

func (r *FileRepository) FindByFileId(ctx context.Context, fileId uuid.UUID) (*ddownload.File, error) {
	var ent edownload.File

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
	file := r.mappers.MapFileEntityToDomain(&ent)

	return file, nil
}
