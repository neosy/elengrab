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

// FieldName field name from sql tag by structure field name
// Example:
// var ent <TableEntity>
// ent.FieldName(&ent.SalesId)
func (e *MediaWatchEvent) FieldName(fieldPtr any) string {
	return e.BaseEntity.FieldName(e, fieldPtr)
}

// FieldNameWithAlias field name with alieas from sql tag by structure field pointer
// Example:
// var ent <TableEntity>
// ent.FieldName(ent, &ent.SalesId, "alias")
func (e *MediaWatchEvent) FieldNameWithAlias(fieldPtr any, alias string) string {
	return e.BaseEntity.FieldNameWithAlias(e, fieldPtr, alias)
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

// Values returns a list of values for fields that will be used for updates
func (e *MediaWatchEvent) Values() []any {
	return e.BaseEntity.Values(e)
}

// FieldsMap returns a map of field names to their corresponding values
// using the entity's Fields() and Values() methods, ready for UPDATE statements.
func (e *MediaWatchEvent) FieldsMap() map[string]any {
	return e.BaseEntity.FieldsMap(e)
}
