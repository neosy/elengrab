package edownload

import (
	"time"

	"github.com/neosy/elengrab/internal/pkg/dbentity"
	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/download/table_names"
)

type DataMigration struct {
	dbentity.BaseEntity[DataMigration]

	// Unique identifier of the migration (e.g. "backfill_user_status")
	MigrationID string `db:"migration_id"`

	// Optional human-readable description of what this migration does
	Description *string `db:"description"`

	// Timestamp when this migration record was created (i.e. when migration was applied)
	CreatedAt time.Time `db:"created_at" insert:"false"`
}

// TableName returns the table name
func (e *DataMigration) TableName() string {
	return tablenames.DataMigrations
}

// FieldName field name from sql tag by structure field name
// Example:
// var ent <TableEntity>
// ent.FieldName(&ent.SalesId)
func (e *DataMigration) FieldName(field any) string {
	return e.BaseEntity.FieldName(e, field)
}

// FieldNameWithAlias field name with alieas from sql tag by structure field pointer
// Example:
// var ent <TableEntity>
// ent.FieldName(ent, &ent.SalesId, "alias")
func (e *DataMigration) FieldNameWithAlias(fieldPtr any, alias string) string {
	return e.BaseEntity.FieldNameWithAlias(e, fieldPtr, alias)
}

// Values returns a list of values for fields that will be used for updates
func (e *DataMigration) Values() []any {
	return e.BaseEntity.Values(e)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *DataMigration) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}
