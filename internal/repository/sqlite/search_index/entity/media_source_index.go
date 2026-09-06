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
	UserID *uuid.UUID `db:"user_id" pfield:"userID"`

	// Media title
	Title string `db:"title"`

	// Media title in lowercase for efficient case-insensitive searches
	TitleLower string `db:"title_lower" pfield:"title"`

	// Description media
	Description *string `db:"description"`

	// Description media in lowercase for efficient case-insensitive searches
	DescriptionLower string `db:"description_lower"`

	// Visibility access level for media (public or private)
	Visibility string `db:"visibility"`

	// Number of completed views
	Views int `db:"views" pfield:"views"`

	// Media source creation timestamp
	SourceCreatedAt time.Time `db:"source_created_at" pfield:"createdAt"`

	// Timestamp when the record was soft deleted
	DeletedAt *time.Time `db:"deleted_at" insert:"false"`
}

var mediaSourceIndexMeta dbentity.EntityMetadata

func init() {
	mediaSourceIndexMeta = dbentity.NewEntityMetadata(MediaSourceIndex{}.BaseEntity)
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

// PaginateFieldName returns the field name with alias from the pagination tag.
// Example:
//
//	var ent <TableEntity>
//	fieldName := ent.PaginateFieldName(&view.SomeField)
func (e *MediaSourceIndex) PaginateFieldName(fieldPtr any) string {
	return e.BaseEntity.PaginateFieldName(e, fieldPtr)
}

// InsertFields returns fields included in insert operations.
// Fields with the `insert:"false"` tag are excluded.
func (e *MediaSourceIndex) InsertFields() []string {
	return mediaSourceIndexMeta.InsertFields()
}

// QueryFields returns a list of fields that will be used for queries
func (e *MediaSourceIndex) QueryFields() []string {
	return mediaSourceIndexMeta.QueryFields()
}

// QueryFieldsWithAlias returns a list of fields with alias that will be used for queries
func (e *MediaSourceIndex) QueryFieldsWithAlias(alias string) []string {
	return e.BaseEntity.FieldsWithAlias(e.QueryFields(), alias)
}

// PaginationFieldNames returns a map from pagination field names to their corresponding database column names.
// It uses the `pfield` tag to determine the pagination field names.
func (e *MediaSourceIndex) PaginationFieldNames() map[string]string {
	return mediaSourceIndexMeta.PaginationFieldNames()
}

// InsertValues returns values for fields included in insert operations.
// Fields with the `insert:"false"` tag are excluded.
func (e *MediaSourceIndex) InsertValues() []any {
	return e.BaseEntity.InsertValues(e)
}

// FieldValues returns a map of field names to their corresponding values
func (e *MediaSourceIndex) InsertFieldValues() map[string]any {
	return e.FieldValues(e.InsertFields(), e.InsertValues())
}
