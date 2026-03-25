package eauth

import (
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/dbentity"
	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/auth/table_names"
)

type UserSession struct {
	dbentity.BaseEntity[UserSession]

	// Unique session identifier (UUID)
	SessionID uuid.UUID `db:"session_id"`

	// Associated user identifier (UUID)
	UserID uuid.UUID `db:"user_id"`

	// Random session token stored in cookie
	SessionToken string `db:"session_token"`

	// Timestamp when the record was created
	CreatedAt time.Time `db:"created_at" insert:"false"`

	// Session expiration timestamp
	ExpiresAt time.Time `db:"expires_at"`
}

// TableName returns the table name
func (e *UserSession) TableName() string {
	return tablenames.UserSessions
}

// FieldName field name from sql tag by structure field name
// Example:
// var ent <TableEntity>
// ent.FieldName(&ent.SalesId)
func (e *UserSession) FieldName(fieldPtr any) string {
	return e.BaseEntity.FieldName(e, fieldPtr)
}

// FieldNameWithAlias field name with alieas from sql tag by structure field pointer
// Example:
// var ent <TableEntity>
// ent.FieldName(ent, &ent.SalesId, "alias")
func (e *UserSession) FieldNameWithAlias(fieldPtr any, alias string) string {
	return e.BaseEntity.FieldNameWithAlias(e, fieldPtr, alias)
}

// Values returns a list of values for fields that will be used for updates
func (e *UserSession) Values() []any {
	return e.BaseEntity.Values(e)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *UserSession) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}
