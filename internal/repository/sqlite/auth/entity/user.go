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

	// User password hash, optional
	PasswordHash *string `db:"password_hash"`

	// Timestamp when the user's password was last updated
	PasswordUpdatedAt *time.Time `db:"password_updated_at"`

	// Active status
	IsActive int `db:"is_active"`

	// Timestamp when the record was created
	CreatedAt time.Time `db:"created_at" insert:"false"`

	// Timestamp when the record was last updated
	UpdatedAt time.Time `db:"updated_at" sqlexpr:"CURRENT_TIMESTAMP"`

	// Timestamp when the record was soft deleted
	DeletedAt *time.Time `db:"deleted_at" insert:"false"`
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

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *User) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}

// FieldPointer returns a pointer to the field of the given struct specified by tag.
func (e *User) FieldPointer(fieldName string) any {
	ptr, _ := e.BaseEntity.FieldPointer(e, fieldName)
	return ptr
}

// InsertValues returns values for fields included in insert operations.
// Fields with the `insert:"false"` tag are excluded.
func (e *User) InsertValues() []any {
	return e.BaseEntity.InsertValues(e)
}

// InsertFieldValues returns a map of field names to their corresponding values
// using the entity's Fields() and Values() methods, ready for UPDATE statements.
func (e *User) InsertFieldValues() map[string]any {
	return e.BaseEntity.InsertFieldValues(e)
}
