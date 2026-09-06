package ewatchevent

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/dbentity"
	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/watch_event/table_names"
)

var (
	userWatchStatConflictFields [2]string
)

func init() {
	var eStat MediaUserWatchStat

	userWatchStatConflictFields = [2]string{
		eStat.FieldName(&eStat.DownloadID),
		eStat.FieldName(&eStat.UserID),
	}
}

type MediaUserWatchStat struct {
	dbentity.BaseEntity[MediaUserWatchStat]

	// Identifier of the watched media (UUID)
	DownloadID uuid.UUID `db:"download_id"`

	// Associated user identifier (UUID)
	UserID uuid.UUID `db:"user_id"`

	// Number of completed views
	Views int `db:"views"`

	// Record update timestamp, set automatically
	UpdatedAt time.Time `db:"updated_at" sqlexpr:"CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (e *MediaUserWatchStat) TableName() string {
	return tablenames.MediaUserWatchStats
}

// FieldName field name from sql tag by structure field name
// Example:
// var ent <TableEntity>
// ent.FieldName(&ent.SalesId)
func (e *MediaUserWatchStat) FieldName(fieldPtr any) string {
	return e.BaseEntity.FieldName(e, fieldPtr)
}

// FieldNameWithAlias field name with alieas from sql tag by structure field pointer
// Example:
// var ent <TableEntity>
// ent.FieldName(ent, &ent.SalesId, "alias")
func (e *MediaUserWatchStat) FieldNameWithAlias(fieldPtr any, alias string) string {
	return e.BaseEntity.FieldNameWithAlias(e, fieldPtr, alias)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *MediaUserWatchStat) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}

// FieldPointer returns a pointer to the field of the given struct specified by tag.
func (e *MediaUserWatchStat) FieldPointer(fieldName string) any {
	ptr, _ := e.BaseEntity.FieldPointer(e, fieldName)
	return ptr
}

// InsertValues returns values for fields included in insert operations.
// Fields with the `insert:"false"` tag are excluded.
func (e *MediaUserWatchStat) InsertValues() []any {
	return e.BaseEntity.InsertValues(e)
}

// FieldValues returns a map of field names to their corresponding values
// using the entity's Fields() and Values() methods, ready for UPDATE statements.
func (e *MediaUserWatchStat) FieldValues() map[string]any {
	return e.BaseEntity.FieldValues(e)
}

func (e *MediaUserWatchStat) ConflictFields() []string {
	return userWatchStatConflictFields[:]
}

func (e *MediaUserWatchStat) ConflictColumnsSQL() string {
	return strings.Join(userWatchStatConflictFields[:], " ,")
}
