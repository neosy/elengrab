package edownload

import (
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/dbentity"
	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/download/table_names"
)

type MediaDownload struct {
	dbentity.BaseEntity[MediaDownload]
	DownloadID               uuid.UUID  `db:"file_id"`
	UserID                   *uuid.UUID `db:"user_id" pfield:"userID"`
	Status                   string     `db:"file_status"`
	MediaURL                 string     `db:"media_url"`
	MediaTitleOriginal       string     `db:"media_title_original"`
	MediaTitle               string     `db:"media_title"`
	MediaTitleLower          string     `db:"media_title_lower" pfield:"title"`
	MediaDescriptionOriginal *string    `db:"media_description_original"`
	MediaDescription         *string    `db:"media_description"`
	ChannelID                *string    `db:"channel_id"`
	FileName                 string     `db:"file_name"`
	Ext                      string     `db:"ext"`
	FileFullName             string     `db:"full_name"`
	FileSize                 *int64     `db:"file_size"`
	PartialHash              *string    `db:"partial_hash"`
	SafeReadableFileFullName string     `db:"safe_readable_full_name"`
	MediaInfo                *string    `db:"media_info"`
	ErrorMessage             *string    `db:"error_message"`
	Visibility               string     `db:"visibility"`
	DownloadedAt             *string    `db:"downloaded_at"`
	CreatedAt                time.Time  `db:"created_at" insert:"false" pfield:"createdAt"`
	UpdatedAt                time.Time  `db:"updated_at" sqlexpr:"CURRENT_TIMESTAMP"`
	DeletedAt                *time.Time `db:"deleted_at" insert:"false"`
}

var mediaDownloadMeta dbentity.EntityMetadata

func init() {
	mediaDownloadMeta = dbentity.NewEntityMetadata(MediaDownload{}.BaseEntity)
}

// TableName returns the table name
func (e *MediaDownload) TableName() string {
	return tablenames.Files
}

// FieldName field name from sql tag by structure field name
// Example:
// var ent <TableEntity>
// ent.FieldName(&ent.SalesId)
func (e *MediaDownload) FieldName(fieldPtr any) string {
	return e.BaseEntity.FieldName(e, fieldPtr)
}

// FieldNameWithAlias field name with alieas from sql tag by structure field pointer
// Example:
// var ent <TableEntity>
// ent.FieldName(ent, &ent.SalesId, "alias")
func (e *MediaDownload) FieldNameWithAlias(fieldPtr any, alias string) string {
	return e.BaseEntity.FieldNameWithAlias(e, fieldPtr, alias)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *MediaDownload) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}

// FieldPointer returns a pointer to the field of the given struct specified by tag.
func (e *MediaDownload) FieldPointer(fieldName string) any {
	ptr, _ := e.BaseEntity.FieldPointer(e, fieldName)
	return ptr
}

// PaginateFieldName returns the field name with alias from the pagination tag.
// Example:
//
//	var ent <TableEntity>
//	fieldName := ent.PaginateFieldName(&view.SomeField)
func (e *MediaDownload) PaginateFieldName(fieldPtr any) string {
	return e.BaseEntity.PaginateFieldName(e, fieldPtr)
}

// InsertFields returns fields included in insert operations.
// Fields with the `insert:"false"` tag are excluded.
func (e *MediaDownload) InsertFields() []string {
	return mediaDownloadMeta.InsertFields()
}

// QueryFields returns a list of fields that will be used for queries
func (e *MediaDownload) QueryFields() []string {
	return mediaDownloadMeta.QueryFields()
}

// QueryFieldsWithAlias returns a list of fields with alias that will be used for queries
func (e *MediaDownload) QueryFieldsWithAlias(alias string) []string {
	return e.BaseEntity.FieldsWithAlias(e.QueryFields(), alias)
}

// InsertValues returns values for fields included in insert operations.
// Fields with the `insert:"false"` tag are excluded.
func (e *MediaDownload) InsertValues() []any {
	return e.BaseEntity.InsertValues(e)
}

// FieldValues returns a map of field names to their corresponding values
func (e *MediaDownload) InsertFieldValues() map[string]any {
	return e.FieldValues(e.InsertFields(), e.InsertValues())
}
