package ewatchevent

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/dbentity"
	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/watch_event/table_names"
)

var (
	userChunkConflictFields [3]string
)

func init() {
	var eChunk MediaUserWatchChunk

	userChunkConflictFields = [3]string{
		eChunk.FieldName(&eChunk.DownloadID),
		eChunk.FieldName(&eChunk.UserID),
		eChunk.FieldName(&eChunk.ChunkIndex),
	}
}

type MediaUserWatchChunk struct {
	dbentity.BaseEntity[MediaUserWatchChunk]

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
func (e *MediaUserWatchChunk) TableName() string {
	return tablenames.MediaUserWatchChunks
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
func (e *MediaUserWatchChunk) FieldName(field any, alias ...string) string {
	return e.BaseEntity.FieldName(e, field, alias...)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *MediaUserWatchChunk) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}

// FieldPointer returns a pointer to the field of the given struct specified by tag.
func (e *MediaUserWatchChunk) FieldPointer(fieldName string) any {
	ptr, _ := e.BaseEntity.FieldPointer(e, fieldName)
	return ptr
}

// InsertValues returns values for fields included in insert operations.
// Fields with the `insert:"false"` tag are excluded.
func (e *MediaUserWatchChunk) InsertValues() []any {
	return e.BaseEntity.InsertValues(e)
}

// InsertFieldValues returns a map of field names to their corresponding values
// using the entity's Fields() and Values() methods, ready for UPDATE statements.
func (e *MediaUserWatchChunk) InsertFieldValues() map[string]any {
	return e.BaseEntity.InsertFieldValues(e)
}

func (e *MediaUserWatchChunk) ConflictColumnsSQL() string {
	return strings.Join(userChunkConflictFields[:], " ,")
}
