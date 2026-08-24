package esearchindex

import (
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/dbentity"
	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/search_index/table_names"
)

type MediaSourceIndex struct {
	dbentity.BaseEntity[MediaSourceIndex]

	// Identifier of the watched media (UUID)
	DownloadID uuid.UUID `db:"download_id"`

	// User identifier
	UserID *uuid.UUID `db:"user_id"`

	// Media title
	Title string `db:"media_title"`

	// Media title in lowercase for efficient case-insensitive searches
	TitleLower string `db:"title_lower"`

	// Description media
	Description *string `db:"description"`

	// Description media in lowercase for efficient case-insensitive searches
	DescriptionLower string `db:"description_lower"`

	// Visibility access level for media (public or private)
	Visibility string `db:"visibility"`

	// Number of completed views
	Views int `db:"views"`

	// Media source creation timestamp
	SourceCreatedAt time.Time `db:"source_created_at"`
}

// TableName returns the table name
func (e *MediaSourceIndex) TableName() string {
	return tablenames.MediaSourcesIndex
}

// FieldName field name from sql tag by structure field name
// Example:
// var ent <TableEntity>
// ent.FieldName(&ent.SalesId)
func (e *MediaSourceIndex) FieldName(fieldPtr any) string {
	return e.BaseEntity.FieldName(e, fieldPtr)
}

// FieldNameWithAlias field name with alieas from sql tag by structure field pointer
// Example:
// var ent <TableEntity>
// ent.FieldName(ent, &ent.SalesId, "alias")
func (e *MediaSourceIndex) FieldNameWithAlias(fieldPtr any, alias string) string {
	return e.BaseEntity.FieldNameWithAlias(e, fieldPtr, alias)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *MediaSourceIndex) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}

// FieldPointer returns a pointer to the field of the given struct specified by tag.
func (e *MediaSourceIndex) FieldPointer(fieldName string) any {
	ptr, _ := e.BaseEntity.FieldPointer(e, fieldName)
	return ptr
}

// Values returns a list of values for fields that will be used for updates
func (e *MediaSourceIndex) Values() []any {
	return e.BaseEntity.Values(e)
}

// FieldsMap returns a map of field names to their corresponding values
// using the entity's Fields() and Values() methods, ready for UPDATE statements.
func (e *MediaSourceIndex) FieldsMap() map[string]any {
	return e.BaseEntity.FieldsMap(e)
}
