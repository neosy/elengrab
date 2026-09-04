package eauth

import (
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/dbentity"
	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/auth/table_names"
)

type UserRole struct {
	dbentity.BaseEntity[UserRole]

	// Reference to the user (many-to-many relationship)
	UserID uuid.UUID `db:"user_id"`

	// Reference to the role assigned to the user
	RoleID string `db:"role_id"`

	// Timestamp when the record was created
	CreatedAt time.Time `db:"created_at" insert:"false"`
}

// TableName returns the table name
func (e *UserRole) TableName() string {
	return tablenames.UserRoles
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
func (e *UserRole) FieldName(field any, alias ...string) string {
	return e.BaseEntity.FieldName(e, field, alias...)
}

// InsertValues returns values for fields included in insert operations.
// Fields with the `insert:"false"` tag are excluded.
func (e *UserRole) InsertValues() []any {
	return e.BaseEntity.InsertValues(e)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *UserRole) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}
