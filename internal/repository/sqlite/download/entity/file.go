package edownload

import (
	"time"

	"github.com/google/uuid"
	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/download/table_names"
	"github.com/neosy/elengrab/pkg/dbentity"
)

type File struct {
	dbentity.BaseEntity[File]
	FileId               uuid.UUID `db:"file_id"`
	Title                string    `db:"title"`
	FileName             string    `db:"file_name"`
	Ext                  string    `db:"ext"`
	FullName             string    `db:"full_name"`
	SafeReadableFullName string    `db:"safe_readable_full_name"`
	CreatedAt            time.Time `db:"created_at" insert:"false"`
	UpdatedAt            time.Time `db:"updated_at" insert:"false"`
}

// TableName returns the table name
func (e *File) TableName() string {
	return tablenames.Files
}

// FieldName field name from sql tag by structure field name
// Example:
// var ent <TableEntity>
// ent.FieldName(&ent.SalesId)
func (e *File) FieldName(field any) string {
	return e.BaseEntity.FieldName(e, field)
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
