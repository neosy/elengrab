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

// FieldName returns the field name from the SQL tag using a structure field name or pointer,
// optionally prefixed with a table alias.
//
// Examples:
//
//	var entity <TableEntity>
//	entity.FieldName("created_at", "")              // "created_at"
//	entity.FieldName(&entity.CreatedAt, "")         // "created_at"
//	entity.FieldName(&entity.CreatedAt, "users")    // "users.created_at"
func (e *DataMigration) FieldName(field any, alias ...string) string {
	return e.BaseEntity.FieldName(e, field, alias...)
}

// InsertValues returns values for fields included in insert operations.
// Fields with the `insert:"false"` tag are excluded.
func (e *DataMigration) InsertValues() []any {
	return e.BaseEntity.InsertValues(e)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *DataMigration) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}
