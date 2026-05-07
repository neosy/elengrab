package edownload

import (
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/dbentity"
	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/download/table_names"
)

type File struct {
	dbentity.BaseEntity[File]
	FileID               uuid.UUID  `db:"file_id"`
	UserID               *uuid.UUID `db:"user_id"`
	Status               string     `db:"file_status"`
	MediaUrl             string     `db:"media_url"`
	MediaTitle           string     `db:"media_title"`
	MediaTitleLower      string     `db:"media_title_lower"`
	MediaDescription     *string    `db:"media_description"`
	ChannelID            *string    `db:"channel_id"`
	FileName             string     `db:"file_name"`
	Ext                  string     `db:"ext"`
	FullName             string     `db:"full_name"`
	FileSize             *int64     `db:"file_size"`
	PartialHash          *string    `db:"partial_hash"`
	SafeReadableFullName string     `db:"safe_readable_full_name"`
	MediaInfo            *string    `db:"media_info"`
	ErrorMessage         *string    `db:"error_message"`
	DownloadedAt         *string    `db:"downloaded_at"`
	CreatedAt            time.Time  `db:"created_at" insert:"false"`
	UpdatedAt            time.Time  `db:"updated_at" sqlexpr:"CURRENT_TIMESTAMP"`
	DeletedAt            *time.Time `db:"deleted_at" insert:"false"`
}

// TableName returns the table name
func (e *File) TableName() string {
	return tablenames.Files
}

// FieldName field name from sql tag by structure field name
// Example:
// var ent <TableEntity>
// ent.FieldName(&ent.SalesId)
func (e *File) FieldName(fieldPtr any) string {
	return e.BaseEntity.FieldName(e, fieldPtr)
}

// FieldNameWithAlias field name with alieas from sql tag by structure field pointer
// Example:
// var ent <TableEntity>
// ent.FieldName(ent, &ent.SalesId, "alias")
func (e *File) FieldNameWithAlias(fieldPtr any, alias string) string {
	return e.BaseEntity.FieldNameWithAlias(e, fieldPtr, alias)
}

// Values returns a list of values for fields that will be used for updates
func (e *File) Values() []any {
	return e.BaseEntity.Values(e)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *File) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}
