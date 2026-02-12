package emedia

import (
	"time"

	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/media/table_names"
	"github.com/neosy/elengrab/pkg/dbentity"
)

type YoutubeChannel struct {
	dbentity.BaseEntity[YoutubeChannel]

	// Unique ID for the channel
	ChannelID string `db:"channel_id"`

	// Site URL
	ChannelURL string `db:"channel_url"`

	// Title of the channel
	ChannelTitle string `db:"channel_title"`

	// URL of the channel avatar
	ImageURL string `db:"image_url"`

	// Raw image data (binary)
	ImageRaw []byte `db:"image_raw"`

	// Format of the image (jpg, png, webp)
	ImageFormat string `db:"image_format"`

	// Timestamp when the record was created
	CreatedAt time.Time `db:"created_at" insert:"false"`

	// Timestamp when the record was last updated
	UpdatedAt time.Time `db:"updated_at" sqlexpr:"CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (e *YoutubeChannel) TableName() string {
	return tablenames.YoutubeChannels
}

// FieldName field name from sql tag by structure field name
// Example:
// var ent <TableEntity>
// ent.FieldName(&ent.SalesId)
func (e *YoutubeChannel) FieldName(fieldPtr any) string {
	return e.BaseEntity.FieldName(e, fieldPtr)
}

// FieldNameWithAlias field name with alieas from sql tag by structure field pointer
// Example:
// var ent <TableEntity>
// ent.FieldName(ent, &ent.SalesId, "alias")
func (e *YoutubeChannel) FieldNameWithAlias(fieldPtr any, alias string) string {
	return e.BaseEntity.FieldNameWithAlias(e, fieldPtr, alias)
}

// Values returns a list of values for fields that will be used for updates
func (e *YoutubeChannel) Values() []any {
	return e.BaseEntity.Values(e)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *YoutubeChannel) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}
