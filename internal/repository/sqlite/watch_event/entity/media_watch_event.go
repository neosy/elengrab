package ewatchevent

import (
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/dbentity"
	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/watch_event/table_names"
)

type MediaWatchEvent struct {
	dbentity.BaseEntity[MediaWatchEvent]

	// Unique event identifier (UUID)
	EventID uuid.UUID `db:"event_id"`

	// Identifier of the watched media (UUID)
	DownloadID uuid.UUID `db:"download_id"`

	// Associated user identifier (UUID)
	UserID *uuid.UUID `db:"user_id"`

	// User session identifier (UUID)
	SessionID *uuid.UUID `db:"session_id"`

	// Playback position in milliseconds
	PositionMs int `db:"position_ms"`

	// Playback duration since the previous event in milliseconds
	IntervalMs int `db:"interval_ms"`

	// Record creation timestamp, set automatically
	CreatedAt time.Time `db:"created_at" insert:"false"`
}

// TableName returns the table name
func (e *MediaWatchEvent) TableName() string {
	return tablenames.MediaWatchEvents
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
func (e *MediaWatchEvent) FieldName(field any, alias ...string) string {
	return e.BaseEntity.FieldName(e, field, alias...)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *MediaWatchEvent) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}

// FieldPointer returns a pointer to the field of the given struct specified by tag.
func (e *MediaWatchEvent) FieldPointer(fieldName string) any {
	ptr, _ := e.BaseEntity.FieldPointer(e, fieldName)
	return ptr
}

// InsertValues returns values for fields included in insert operations.
// Fields with the `insert:"false"` tag are excluded.
func (e *MediaWatchEvent) InsertValues() []any {
	return e.BaseEntity.InsertValues(e)
}

// InsertFieldValues returns a map of field names to their corresponding values
// using the entity's Fields() and Values() methods, ready for UPDATE statements.
func (e *MediaWatchEvent) InsertFieldValues() map[string]any {
	return e.BaseEntity.InsertFieldValues(e)
}
