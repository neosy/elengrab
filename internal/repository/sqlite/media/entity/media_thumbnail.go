package emedia

import (
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/pkg/dbentity"
	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/media/table_names"
)

// Stores thumbnail metadata for media files
type Thumbnail struct {
	dbentity.BaseEntity[Thumbnail]

	// Unique identifier for the thumbnail record
	ThumbID uuid.UUID `db:"thumb_id"`

	// Reference to the parent media entity
	MediaID uuid.UUID `db:"media_id"`

	// Thumbnail variant type (e.g. small, medium, large, original)
	Variant string `db:"variant"`

	// Thumbnail generation version
	Version uint8 `db:"version"`

	// Thumbnail width in pixels
	Width *uint16 `db:"width"`

	// Thumbnail height in pixels
	Height *uint16 `db:"height"`

	// Image format (e.g. jpg, png, webp)
	Format string `db:"format"`

	// Source of the thumbnail (youtube, vimeo, external, video_frame, generated, upload)
	SourceType string `db:"source_type"`

	// Optional external identifier of the source
	SourceID *string `db:"source_id"`

	// Optional external source URL used to derive the thumbnail
	SourceURL *string `db:"source_url"`

	// Stable object key used to resolve file location in storage backend (FS/S3/etc.)
	StorageKey uuid.UUID `db:"storage_key"`

	// Flag indicating whether this is the primary thumbnail
	IsPrimary bool `db:"is_primary"`

	// Record creation timestamp
	CreatedAt time.Time `db:"created_at" insert:"false"`

	// Record last update timestamp
	UpdatedAt time.Time `db:"updated_at" sqlexpr:"CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (e *Thumbnail) TableName() string {
	return tablenames.Thumbnails
}

// FieldName field name from sql tag by structure field name
// Example:
// var ent <TableEntity>
// ent.FieldName(&ent.SalesId)
func (e *Thumbnail) FieldName(fieldPtr any) string {
	return e.BaseEntity.FieldName(e, fieldPtr)
}

// FieldNameWithAlias field name with alieas from sql tag by structure field pointer
// Example:
// var ent <TableEntity>
// ent.FieldName(ent, &ent.SalesId, "alias")
func (e *Thumbnail) FieldNameWithAlias(fieldPtr any, alias string) string {
	return e.BaseEntity.FieldNameWithAlias(e, fieldPtr, alias)
}

// InsertValues returns values for fields included in insert operations.
// Fields with the `insert:"false"` tag are excluded.
func (e *Thumbnail) InsertValues() []any {
	return e.BaseEntity.InsertValues(e)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *Thumbnail) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}
