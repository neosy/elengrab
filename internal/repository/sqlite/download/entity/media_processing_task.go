package edownload

import (
	"database/sql"

	"github.com/neosy/elengrab/internal/pkg/dbentity"
	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/download/table_names"
)

type MediaProcessingTask struct {
	dbentity.BaseEntity[MediaProcessingTask]

	// Unique task identifier (UUID)
	TaskID sql.NullString `db:"task_id"`

	// Unique file identifier (UUID)
	DownloadID sql.NullString `db:"download_id"`

	// Status
	Status sql.NullString `db:"task_status"`

	// Path to the staged uploaded file
	SourcePath sql.NullString `db:"source_path"`

	// ID of the worker currently processing the task
	WorkerID sql.NullInt64 `db:"worker_id"`

	// ID of the job currently processing the task
	JobID sql.NullString `db:"job_id"`

	// Task creation timestamp
	CreatedAt sql.NullTime `db:"created_at" insert:"false"`

	// Last update timestamp
	UpdatedAt sql.NullTime `db:"updated_at" sqlexpr:"CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (e *MediaProcessingTask) TableName() string {
	return tablenames.ProcessingTasks
}

// FieldName returns the field name from the SQL tag using a structure field name or pointer,
// optionally prefixed with a table alias.
//
// Examples:
//
//	var entity <TableEntity>
//	entity.FieldName("created_at", "")              // "created_at"
//	entity.FieldName(&entity.CreatedAt, "")         // "created_at"
//	entity.FieldName(&entity.CreatedAt, "users")    // "users.created_at"
func (e *MediaProcessingTask) FieldName(field any, alias ...string) string {
	return e.BaseEntity.FieldName(e, field, alias...)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *MediaProcessingTask) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}

// FieldPointer returns a pointer to the field of the given struct specified by tag.
func (e *MediaProcessingTask) FieldPointer(fieldName string) any {
	ptr, _ := e.BaseEntity.FieldPointer(e, fieldName)
	return ptr
}

// InsertValues returns values for fields included in insert operations.
// Fields with the `insert:"false"` tag are excluded.
func (e *MediaProcessingTask) InsertValues() []any {
	return e.BaseEntity.InsertValues(e)
}

// InsertFieldValues returns a map of field names to their corresponding values
// using the entity's Fields() and Values() methods, ready for UPDATE statements.
func (e *MediaProcessingTask) InsertFieldValues() map[string]any {
	return e.BaseEntity.InsertFieldValues(e)
}
