package emedia

import (
	"time"

	"github.com/neosy/elengrab/internal/pkg/dbentity"
	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/media/table_names"
)

type Channel struct {
	dbentity.BaseEntity[Channel]

	// Unique ID for the channel
	ChannelID string `db:"channel_id"`

	// External channel identifier
	ExternalID string `db:"external_id"`

	// Platform identifier
	Platform string `db:"platform"`

	// Site URL
	ChannelURL string `db:"channel_url"`

	// Host from which the platform was detected
	Host string `db:"host"`

	// Channel username
	Username string `db:"username"`

	// Channel URL based on the username
	UsernameURL string `db:"username_url"`

	// Title of the channel
	Title string `db:"channel_title"`

	// URL of the channel avatar
	ImageURL *string `db:"image_url"`

	// Raw image data (binary)
	ImageRaw []byte `db:"image_raw"`

	// Format of the image (jpg, png, webp)
	ImageFormat *string `db:"image_format"`

	// Timestamp when the record was created
	CreatedAt time.Time `db:"created_at" insert:"false"`

	// Timestamp when the record was last updated
	UpdatedAt time.Time `db:"updated_at" sqlexpr:"CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (e *Channel) TableName() string {
	return tablenames.Channels
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
func (e *Channel) FieldName(field any, alias ...string) string {
	return e.BaseEntity.FieldName(e, field, alias...)
}

// InsertValues returns values for fields included in insert operations.
// Fields with the `insert:"false"` tag are excluded.
func (e *Channel) InsertValues() []any {
	return e.BaseEntity.InsertValues(e)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *Channel) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}
