package elink

import (
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/dbentity"
	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/link/table_names"
)

type LinkClick struct {
	dbentity.BaseEntity[LinkClick]

	// Unique identifier for the click event
	LinkClickID uuid.UUID `db:"link_click_id"`

	// ID of the link that was clicked
	LinkID uuid.UUID `db:"link_id"`

	// The IP address from which the link was accessed
	IPAddress string `db:"ip_address"`

	// Full short URL, including domain, e.g., "https://s.nhub.ru/abc123"
	ShortURL string `db:"short_url"`

	// User ID who clicked the link (nullable, if not logged in or unknown)
	ClickedBy *string `db:"clicked_by"`

	// Timestamp of the click event
	ClickedAt time.Time `db:"clicked_at"`

	// User agent or device info (optional for tracking purposes)
	UserAgent *string `db:"user_agent"`

	// Referrer URL (optional, if available)
	Referrer *string `db:"referrer"`

	// Timestamp when the event was created
	CreatedAt time.Time `db:"created_at" insert:"false"`
}

// TableName returns the table name
func (e *LinkClick) TableName() string {
	return tablenames.LinkClicks
}

// FieldName field name from sql tag by structure field name
// Example:
// var ent <TableEntity>
// ent.FieldName(&ent.SalesId)
func (e *LinkClick) FieldName(fieldPtr any) string {
	return e.BaseEntity.FieldName(e, fieldPtr)
}

// FieldNameWithAlias field name with alieas from sql tag by structure field pointer
// Example:
// var ent <TableEntity>
// ent.FieldName(ent, &ent.SalesId, "alias")
func (e *LinkClick) FieldNameWithAlias(fieldPtr any, alias string) string {
	return e.BaseEntity.FieldNameWithAlias(e, fieldPtr, alias)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *LinkClick) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}

// FieldPointer returns a pointer to the field of the given struct specified by tag.
func (e *LinkClick) FieldPointer(fieldName string) any {
	ptr, _ := e.BaseEntity.FieldPointer(e, fieldName)
	return ptr
}

// Values returns a list of values for fields that will be used for updates
func (e *LinkClick) Values() []any {
	return e.BaseEntity.Values(e)
}

// FieldsMap returns a map of field names to their corresponding values
// using the entity's Fields() and Values() methods, ready for UPDATE statements.
func (e *LinkClick) FieldsMap() map[string]any {
	return e.BaseEntity.FieldsMap(e)
}
