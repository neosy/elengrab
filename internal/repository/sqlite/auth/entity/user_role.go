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

// FieldName field name from sql tag by structure field name
// Example:
// var ent <TableEntity>
// ent.FieldName(&ent.SalesId)
func (e *UserRole) FieldName(fieldPtr any) string {
	return e.BaseEntity.FieldName(e, fieldPtr)
}

// FieldNameWithAlias field name with alieas from sql tag by structure field pointer
// Example:
// var ent <TableEntity>
// ent.FieldName(ent, &ent.SalesId, "alias")
func (e *UserRole) FieldNameWithAlias(fieldPtr any, alias string) string {
	return e.BaseEntity.FieldNameWithAlias(e, fieldPtr, alias)
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
