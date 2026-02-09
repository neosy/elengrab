package emedia

import (
	"time"

	tablenames "github.com/neosy/elengrab/internal/repository/sqlite/media/table_names"
	"github.com/neosy/elengrab/pkg/dbentity"
)

type SiteLogo struct {
	dbentity.BaseEntity[SiteLogo]

	// Unique ID for the logo
	LogoID string `db:"logo_id"`

	// Site URL
	SiteURL string `db:"site_url"`

	// URL of the logo image
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
func (e *SiteLogo) TableName() string {
	return tablenames.SiteLogos
}

// FieldName field name from sql tag by structure field name
// Example:
// var ent <TableEntity>
// ent.FieldName(&ent.SalesId)
func (e *SiteLogo) FieldName(fieldPtr any) string {
	return e.BaseEntity.FieldName(e, fieldPtr)
}

// FieldNameWithAlias field name with alieas from sql tag by structure field pointer
// Example:
// var ent <TableEntity>
// ent.FieldName(ent, &ent.SalesId, "alias")
func (e *SiteLogo) FieldNameWithAlias(fieldPtr any, alias string) string {
	return e.BaseEntity.FieldNameWithAlias(e, fieldPtr, alias)
}

// Values returns a list of values for fields that will be used for updates
func (e *SiteLogo) Values() []any {
	return e.BaseEntity.Values(e)
}

// FieldPointers returns a slice of pointers to all exported fields of the given struct.
func (e *SiteLogo) FieldPointers() []any {
	ptrs, _ := e.BaseEntity.FieldPointers(e)
	return ptrs
}
