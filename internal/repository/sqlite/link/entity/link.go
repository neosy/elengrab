package elink

import (
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/dbentity"
	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/link/table_names"
)

type Link struct {
	dbentity.BaseEntity[Link]

	// Unique identifier of the link
	LinkID uuid.UUID `db:"link_id"`

	// Original (long) URL to be shortened
	OriginalURL string `db:"original_url"`

	// Generated short code used in the shortened URL, e.g., "abc123
	ShortCode string `db:"short_code"`

	// Full short URL, including domain, e.g., "https://s.nhub.ru/abc123"
	ShortURL string `db:"short_url"`

	// Indicates if the full short URL should be used for exact match
	IsMatchShortURL int `db:"is_match_short_url"`

	// Maximum number of allowed clicks; nil means unlimited
	MaxClicks *int `db:"max_clicks"`

	// JSON array of user IDs allowed to access the link; nil means no restrictions
	AllowedUserIDs *string `db:"allowed_user_ids"`

	// JSON array of IP addresses allowed to access the link; nil means no restrictions
	AllowedIPs *string `db:"allowed_ips"`

	// Expiration date and time for the link; nil means no expiration
	ExpiresAt *time.Time `db:"expires_at"`

	// Timestamp when the link was created
	CreatedAt time.Time `db:"created_at" insert:"false"`

	// Timestamp when the link was last updated
	UpdatedAt time.Time `db:"updated_at" sqlexpr:"CURRENT_TIMESTAMP"`

	// Timestamp when the link was soft-deleted
	DeletedAt *time.Time `db:"deleted_at" insert:"false"`
}

// TableName returns the table name
func (e *Link) TableName() string {
	return tablenames.Links
}

// FieldName field name from sql tag by structure field name
// Example:
// var ent <TableEntity>
// ent.FieldName(&ent.SalesId)
func (e *Link) FieldName(fieldPtr any) string {
	return e.BaseEntity.FieldName(e, fieldPtr)
}

// FieldNameWithAlias field name with alieas from sql tag by structure field pointer
// Example:
// var ent <TableEntity>
// ent.FieldName(ent, &ent.SalesId, "alias")
func (e *Link) FieldNameWithAlias(fieldPtr any, alias string) string {
	return e.BaseEntity.FieldNameWithAlias(e, fieldPtr, alias)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *Link) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}

// FieldPointer returns a pointer to the field of the given struct specified by tag.
func (e *Link) FieldPointer(fieldName string) any {
	ptr, _ := e.BaseEntity.FieldPointer(e, fieldName)
	return ptr
}

// InsertValues returns values for fields included in insert operations.
// Fields with the `insert:"false"` tag are excluded.
func (e *Link) InsertValues() []any {
	return e.BaseEntity.InsertValues(e)
}

// InsertFieldValues returns a map of field names to their corresponding values
// using the entity's Fields() and Values() methods, ready for UPDATE statements.
func (e *Link) InsertFieldValues() map[string]any {
	return e.BaseEntity.InsertFieldValues(e)
}
