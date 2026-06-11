package dmedia

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

// Stores thumbnail metadata for media files
type Thumbnail struct {
	// Unique identifier for the thumbnail record
	ThumbID uuid.UUID

	// Reference to the parent media entity
	MediaID uuid.UUID

	// Thumbnail variant type (e.g. small, medium, large, original)
	Variant dtypes.ThumbnailVariant

	// Thumbnail generation version
	Version uint8

	// Thumbnail width in pixels
	Width *uint16

	// Thumbnail height in pixels
	Height *uint16

	// Image format (e.g. jpg, png, webp)
	Format dtypes.ImageFormat

	// Source of the thumbnail (youtube, vimeo, external, video_frame, generated, upload)
	SourceType dtypes.ThumbnailSourceType

	// Optional external identifier of the source
	SourceID *string

	// Optional external source URL used to derive the thumbnail
	SourceURL *string

	// Stable object key used to resolve file location in storage backend (FS/S3/etc.)
	StorageKey uuid.UUID

	// Flag indicating whether this is the primary thumbnail
	IsPrimary bool

	// Record creation timestamp
	CreatedAt time.Time

	// Record last update timestamp
	UpdatedAt time.Time

	// Raw image data for the thumbnail (not stored in DB, used for processing)
	ImageRaw []byte
}

var thumbnailNamespace = uuid.MustParse("d344e7d9-47f3-4b77-9f61-e822f3e0bca6")

// NewThumbnail creates a new Thumbnail instance with a unique ThumbID and default version and variant.
func NewThumbnail() *Thumbnail {
	return &Thumbnail{
		ThumbID: uuid.New(),
		Version: 1,
		Variant: dtypes.ThumbnailVariantOriginal,
	}
}

// MakeThumbnailStorageKey generates a (analog of UUID v5) based on mediaID, variant, and format using a predefined namespace
func MakeThumbnailStorageKey(thumbID uuid.UUID, variant string, format string) uuid.UUID {
	name := fmt.Sprintf("%s_%s_%s", thumbID.String(), variant, format)
	return uuid.NewSHA1(thumbnailNamespace, []byte(name))
}

// InitStorageKey generates a stable StorageKey (analog of UUID v5) based on MediaID, Variant, and Version
func (t *Thumbnail) InitStorageKey() error {
	if t.MediaID == uuid.Nil || !t.Variant.Exists() || t.Version == 0 {
		return fmt.Errorf("invalid media thumbnail data for storage key generation")
	}

	//t.StorageKey = MakeThumbnailStorageKey(t.MediaID, t.Variant.String(), int(t.Version))
	t.StorageKey = t.ThumbID

	return nil
}

// Validate checks the integrity of the Thumbnail data
func (t *Thumbnail) Validate() error {
	if t.ThumbID == uuid.Nil {
		return fmt.Errorf("thumbnail ID cannot be nil")
	}
	if t.MediaID == uuid.Nil {
		return fmt.Errorf("media ID cannot be nil")
	}
	if !t.Variant.Exists() {
		return fmt.Errorf("invalid thumbnail variant")
	}
	if t.Version == 0 {
		return fmt.Errorf("version must be greater than 0")
	}
	if !t.Format.Exists() {
		return fmt.Errorf("invalid thumbnail format")
	}
	if !t.SourceType.Exists() {
		return fmt.Errorf("invalid thumbnail source type")
	}
	if t.StorageKey == uuid.Nil {
		return fmt.Errorf("storage key cannot be nil")
	}
	return nil
}

// ImageData returns the thumbnail's image data as a structured object, or nil if no image data is available
func (t *Thumbnail) ImageData(url string) *dtypes.ImageData {
	var width, height int
	if t.Width != nil {
		width = int(*t.Width)
	}
	if t.Height != nil {
		height = int(*t.Height)
	}

	return &dtypes.ImageData{
		URL:    url,
		Format: t.Format,
		Width:  width,
		Height: height,
		Raw:    t.ImageRaw,
	}
}

// ImageDataWithSourceURL returns the thumbnail's image data as a structured object, or nil if no image data is available
func (t *Thumbnail) ImageDataWithSourceURL() *dtypes.ImageData {
	var url string
	if t.SourceURL != nil {
		url = *t.SourceURL
	}

	return t.ImageData(url)
}

func (t *Thumbnail) Copy() *Thumbnail {
	if t == nil {
		return nil
	}

	thumbnailCopy := *t

	thumbnailCopy.Width = uptr.Copy(t.Width)
	thumbnailCopy.Height = uptr.Copy(t.Height)
	thumbnailCopy.SourceID = uptr.Copy(t.SourceID)
	thumbnailCopy.SourceURL = uptr.Copy(t.SourceURL)

	imageRaw := make([]byte, len(thumbnailCopy.ImageRaw))
	copy(imageRaw, thumbnailCopy.ImageRaw)

	return &thumbnailCopy
}

func (t *Thumbnail) IsPortrait() bool {
	if t == nil {
		return false
	}

	if t.Width == nil || t.Height == nil {
		return false
	}

	return *t.Width < *t.Height
}
