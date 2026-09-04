package ewatchevent

import (
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/dbentity"
	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/watch_event/table_names"
)

type MediaWatchStat struct {
	dbentity.BaseEntity[MediaWatchStat]

	// Identifier of the watched media (UUID)
	DownloadID uuid.UUID `db:"download_id"`

	// Number of completed views
	Views int `db:"views"`

	// Record update timestamp, set automatically
	UpdatedAt time.Time `db:"updated_at" sqlexpr:"CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (e *MediaWatchStat) TableName() string {
	return tablenames.MediaWatchStats
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
func (e *MediaWatchStat) FieldName(field any, alias ...string) string {
	return e.BaseEntity.FieldName(e, field, alias...)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *MediaWatchStat) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}

// FieldPointer returns a pointer to the field of the given struct specified by tag.
func (e *MediaWatchStat) FieldPointer(fieldName string) any {
	ptr, _ := e.BaseEntity.FieldPointer(e, fieldName)
	return ptr
}

// InsertValues returns values for fields included in insert operations.
// Fields with the `insert:"false"` tag are excluded.
func (e *MediaWatchStat) InsertValues() []any {
	return e.BaseEntity.InsertValues(e)
}

// InsertFieldValues returns a map of field names to their corresponding values
// using the entity's Fields() and Values() methods, ready for UPDATE statements.
func (e *MediaWatchStat) InsertFieldValues() map[string]any {
	return e.BaseEntity.InsertFieldValues(e)
}
