package dmedia

import (
	"time"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

type SiteLogo struct {
	// Unique ID for the logo
	LogoID uuid.UUID

	// Site URL
	SiteURL string

	// Title of the site
	SiteTitle string

	// URL of the logo image
	ImageURL string

	// Raw image data (binary)
	ImageRaw []byte

	// Format of the image (jpg, png, webp)
	ImageFormat dtypes.ImageFormat

	// Timestamp when the record was created
	CreatedAt time.Time

	// Timestamp when the record was last updated
	UpdatedAt time.Time
}

// NewSiteLogo creates a new SiteLogo record with the required fields.
func NewSiteLogo(siteURL string, siteTitle string, imgData *dtypes.ImageData) *SiteLogo {
	siteLogo := &SiteLogo{}
	siteLogo.SetRequired(siteURL, siteTitle, imgData)
	return siteLogo
}

// SetRequired sets the required fields for a SiteLogo record.
func (l *SiteLogo) SetRequired(siteURL string, siteTitle string, imgData *dtypes.ImageData) {
	l.SiteURL = siteURL
	l.SiteTitle = siteTitle
	l.SetImage(imgData)
}

// SetImage sets the image data for a SiteLogo record.
func (l *SiteLogo) SetImage(imgData *dtypes.ImageData) {
	if imgData != nil {
		l.ImageURL = imgData.URL
		l.ImageRaw = imgData.Raw
		l.ImageFormat = imgData.Format
	} else {
		l.ImageURL = ""
		l.ImageRaw = nil
		l.ImageFormat = dtypes.ImageFormatUnknown
	}
}

func (l *SiteLogo) ImageData() *dtypes.ImageData {
	return &dtypes.ImageData{
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
