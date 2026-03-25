package eauth

import (
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/dbentity"
	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/auth/table_names"
)

type User struct {
	dbentity.BaseEntity[User]

	// Unique user identifier (UUID)
	UserID uuid.UUID `db:"user_id"`

	// Username/login, must be unique
	Login string `db:"login"`

	// User email address, optional
	Email *string `db:"email"`

	// Active status
	IsActive int `db:"is_active"`

	// Timestamp when the record was created
	CreatedAt time.Time `db:"created_at" insert:"false"`

	// Timestamp when the record was last updated
	UpdatedAt time.Time `db:"updated_at" sqlexpr:"CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (e *User) TableName() string {
	return tablenames.Users
}

// FieldName field name from sql tag by structure field name
// Example:
// var ent <TableEntity>
// ent.FieldName(&ent.SalesId)
func (e *User) FieldName(fieldPtr any) string {
	return e.BaseEntity.FieldName(e, fieldPtr)
}

// FieldNameWithAlias field name with alieas from sql tag by structure field pointer
// Example:
// var ent <TableEntity>
// ent.FieldName(ent, &ent.SalesId, "alias")
func (e *User) FieldNameWithAlias(fieldPtr any, alias string) string {
	return e.BaseEntity.FieldNameWithAlias(e, fieldPtr, alias)
}

// Values returns a list of values for fields that will be used for updates
func (e *User) Values() []any {
	return e.BaseEntity.Values(e)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *User) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}
