package ewatchevent

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/dbentity"
	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/watch_event/table_names"
)

var (
	userPositionConflictFields [3]string
)

func init() {
	var ePosition MediaUserWatchPosition

	userPositionConflictFields = [3]string{
		ePosition.FieldName(&ePosition.DownloadID),
		ePosition.FieldName(&ePosition.UserID),
		ePosition.FieldName(&ePosition.SessionID),
	}
}

type MediaUserWatchPosition struct {
	dbentity.BaseEntity[MediaUserWatchPosition]

	// Identifier of the watched media (UUID)
	DownloadID uuid.UUID `db:"download_id"`

	// Associated user identifier (UUID)
	UserID uuid.UUID `db:"user_id"`

	// User session identifier (UUID)
	SessionID string `db:"session_id"`

	// Last saved playback position in milliseconds
	PositionMs int `db:"position_ms"`

	// Record creation timestamp, set automatically
	CreatedAt time.Time `db:"created_at" insert:"false"`

	// Record last update timestamp, set automatically
	UpdatedAt time.Time `db:"updated_at" sqlexpr:"CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (e *MediaUserWatchPosition) TableName() string {
	return tablenames.MediaUserWatchPositions
}

// FieldName field name from sql tag by structure field name
// Example:
// var ent <TableEntity>
// ent.FieldName(&ent.SalesId)
func (e *MediaUserWatchPosition) FieldName(fieldPtr any) string {
	return e.BaseEntity.FieldName(e, fieldPtr)
}

// FieldNameWithAlias field name with alieas from sql tag by structure field pointer
// Example:
// var ent <TableEntity>
// ent.FieldName(ent, &ent.SalesId, "alias")
func (e *MediaUserWatchPosition) FieldNameWithAlias(fieldPtr any, alias string) string {
	return e.BaseEntity.FieldNameWithAlias(e, fieldPtr, alias)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *MediaUserWatchPosition) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}

// FieldPointer returns a pointer to the field of the given struct specified by tag.
func (e *MediaUserWatchPosition) FieldPointer(fieldName string) any {
	ptr, _ := e.BaseEntity.FieldPointer(e, fieldName)
	return ptr
}

// InsertValues returns values for fields included in insert operations.
// Fields with the `insert:"false"` tag are excluded.
func (e *MediaUserWatchPosition) InsertValues() []any {
	return e.BaseEntity.InsertValues(e)
}

// InsertFieldValues returns a map of field names to their corresponding values
// using the entity's Fields() and Values() methods, ready for UPDATE statements.
func (e *MediaUserWatchPosition) InsertFieldValues() map[string]any {
	return e.BaseEntity.InsertFieldValues(e)
}

func (e *MediaUserWatchPosition) ConflictFields() []string {
	return userPositionConflictFields[:]
}

func (e *MediaUserWatchPosition) ConflictColumnsSQL() string {
	return strings.Join(userPositionConflictFields[:], " ,")
}
