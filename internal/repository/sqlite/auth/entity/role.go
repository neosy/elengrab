package eauth

import (
	"time"

	"github.com/neosy/elengrab/internal/pkg/dbentity"
	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/auth/table_names"
)

type Role struct {
	dbentity.BaseEntity[Role]

	// Unique role identifier (can be a readable key like 'admin', 'guest')
	RoleID string `db:"role_id"`

	// Human-readable role name, must be unique across the system
	Name string `db:"name"`

	// Optional description of role
	Description *string `db:"name"`

	// Timestamp when the record was created
	CreatedAt time.Time `db:"created_at" insert:"false"`

	// Timestamp when the record was last updated
	UpdatedAt time.Time `db:"updated_at" sqlexpr:"CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (e *Role) TableName() string {
	return tablenames.Roles
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
func (e *Role) FieldName(field any, alias ...string) string {
	return e.BaseEntity.FieldName(e, field, alias...)
}

// InsertValues returns values for fields included in insert operations.
// Fields with the `insert:"false"` tag are excluded.
func (e *Role) InsertValues() []any {
	return e.BaseEntity.InsertValues(e)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *Role) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}

// FieldPointer returns a pointer to the field of the given struct specified by tag.
func (e *Role) FieldPointer(fieldName string) any {
	ptr, _ := e.BaseEntity.FieldPointer(e, fieldName)
	return ptr
}
