package dmedia

import (
	"time"

	"github.com/google/uuid"
	uptr "github.com/neosy/elengrab/pkg/utils/pointer"
)

type ImageData struct {
	// URL of the image
	URL string

	// Raw image data (binary)
	Raw []byte

	// Format of the image (jpg, png, webp)
	Format string
}

type SiteLogo struct {
	// Unique ID for the logo
	LogoID uuid.UUID

	// Site URL
	SiteURL string

	// URL of the logo image
	ImageURL string

	// Raw image data (binary)
	ImageRaw []byte

	// Format of the image (jpg, png, webp)
	ImageFormat string

	// Timestamp when the record was created
	CreatedAt time.Time

	// Timestamp when the record was last updated
	UpdatedAt time.Time
}

func (l *SiteLogo) ImageData() ImageData {
	return ImageData{
		URL:    l.ImageURL,
		Raw:    l.ImageRaw,
		Format: l.ImageFormat,
	}
}

func (src *SiteLogo) Copy() *SiteLogo {
	if src == nil {
		return nil
	}

	copy := uptr.Copy(src)

	if len(src.ImageRaw) > 1 {
		copy.ImageRaw = append([]byte{}, src.ImageRaw...)
	}

	return copy
}
