package ewatchevent

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/dbentity"
	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/watch_event/table_names"
)

var (
	chunkConflictFields [3]string
)

func init() {
	var eChunk MediaWatchChunk

	chunkConflictFields = [3]string{
		eChunk.FieldName(&eChunk.DownloadID),
		eChunk.FieldName(&eChunk.UserID),
		eChunk.FieldName(&eChunk.ChunkIndex),
	}
}

type MediaWatchChunk struct {
	dbentity.BaseEntity[MediaWatchChunk]

	// Identifier of the watched media (UUID)
	DownloadID uuid.UUID `db:"download_id"`

	// Associated user identifier (UUID)
	UserID uuid.UUID `db:"user_id"`

	// Zero-based index of the 1000ms media chunk
	ChunkIndex int `db:"chunk_index"`

	// How many times this chunk was watched
	Qty int `db:"qty"`

	// Record creation timestamp, set automatically
	CreatedAt time.Time `db:"created_at" insert:"false"`
}

// TableName returns the table name
func (e *MediaWatchChunk) TableName() string {
	return tablenames.MediaWatchChunks
}

// FieldName field name from sql tag by structure field name
// Example:
// var ent <TableEntity>
// ent.FieldName(&ent.SalesId)
func (e *MediaWatchChunk) FieldName(fieldPtr any) string {
	return e.BaseEntity.FieldName(e, fieldPtr)
}

// FieldNameWithAlias field name with alieas from sql tag by structure field pointer
// Example:
// var ent <TableEntity>
// ent.FieldName(ent, &ent.SalesId, "alias")
func (e *MediaWatchChunk) FieldNameWithAlias(fieldPtr any, alias string) string {
	return e.BaseEntity.FieldNameWithAlias(e, fieldPtr, alias)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *MediaWatchChunk) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}

// FieldPointer returns a pointer to the field of the given struct specified by tag.
func (e *MediaWatchChunk) FieldPointer(fieldName string) any {
	ptr, _ := e.BaseEntity.FieldPointer(e, fieldName)
	return ptr
}

// Values returns a list of values for fields that will be used for updates
func (e *MediaWatchChunk) Values() []any {
	return e.BaseEntity.Values(e)
}

// FieldsMap returns a map of field names to their corresponding values
// using the entity's Fields() and Values() methods, ready for UPDATE statements.
func (e *MediaWatchChunk) FieldsMap() map[string]any {
	return e.BaseEntity.FieldsMap(e)
}

func (e *MediaWatchChunk) ConflictColumnsSQL() string {
	return strings.Join(chunkConflictFields[:], " ,")
}
