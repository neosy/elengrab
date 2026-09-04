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

// FieldName field name from sql tag by structure field name
// Example:
// var ent <TableEntity>
// ent.FieldName(&ent.SalesId)
func (e *MediaWatchStat) FieldName(fieldPtr any) string {
	return e.BaseEntity.FieldName(e, fieldPtr)
}

// FieldNameWithAlias field name with alieas from sql tag by structure field pointer
// Example:
// var ent <TableEntity>
// ent.FieldName(ent, &ent.SalesId, "alias")
func (e *MediaWatchStat) FieldNameWithAlias(fieldPtr any, alias string) string {
	return e.BaseEntity.FieldNameWithAlias(e, fieldPtr, alias)
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

// Values returns a list of values for fields that will be used for updates
func (e *MediaWatchStat) Values() []any {
	return e.BaseEntity.Values(e)
}

// FieldValues returns a map of field names to their corresponding values
// using the entity's Fields() and Values() methods, ready for UPDATE statements.
func (e *MediaWatchStat) FieldValues() map[string]any {
	return e.BaseEntity.FieldValues(e)
}
