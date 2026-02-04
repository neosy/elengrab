package edownload

import (
	"database/sql"

	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/download/table_names"
	"github.com/neosy/elengrab/pkg/dbentity"
)

type DownloadTask struct {
	dbentity.BaseEntity[DownloadTask]
	TaskId    sql.NullString `db:"task_id"`
	FileId    sql.NullString `db:"file_id"`
	Status    sql.NullString `db:"task_status"`
	MediaUrl  sql.NullString `db:"youtube_url"`
	Options   sql.NullString `db:"options"`
	WorkerId  sql.NullInt64  `db:"worker_id"`
	JobID     sql.NullString `db:"job_id"`
	CreatedAt sql.NullTime   `db:"created_at" insert:"false"`
	UpdatedAt sql.NullTime   `db:"updated_at" sqlexpr:"CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (e *DownloadTask) TableName() string {
	return tablenames.DownloadTasks
}

// FieldName field name from sql tag by structure field name
// Example:
// var ent <TableEntity>
// ent.FieldName(&ent.SalesId)
func (e *DownloadTask) FieldName(field any) string {
	return e.BaseEntity.FieldName(e, field)
}

// FieldNameWithAlias field name with alieas from sql tag by structure field pointer
// Example:
// var ent <TableEntity>
// ent.FieldName(ent, &ent.SalesId, "alias")
func (e *DownloadTask) FieldNameWithAlias(fieldPtr any, alias string) string {
	return e.BaseEntity.FieldNameWithAlias(e, fieldPtr, alias)
}

// Values returns a list of values for fields that will be used for updates
func (e *DownloadTask) Values() []any {
	return e.BaseEntity.Values(e)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *DownloadTask) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}
