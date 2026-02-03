package mappers

import (
	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
)

func (m *Mappers) MapSiteLogoDomainToEntity(logo *dmedia.SiteLogo) (*edownload.SiteLogo, error) {
	return &edownload.SiteLogo{
		LogoID:      logo.LogoID.String(),
		SiteURL:     logo.SiteURL,
		ImageURL:    logo.ImageURL,
		ImageRaw:    logo.ImageRaw,
		ImageFormat: logo.ImageFormat,
	}, nil
}

func (m *Mappers) MapSiteLogoEntityToDomain(eLogo *edownload.SiteLogo) (*dmedia.SiteLogo, error) {
	logoID, err := uuid.Parse(eLogo.LogoID)
	if err != nil {
		return nil, err
	}

	return &dmedia.SiteLogo{
		LogoID:      logoID,
		SiteURL:     eLogo.SiteURL,
		ImageURL:    eLogo.ImageURL,
		ImageRaw:    eLogo.ImageRaw,
		ImageFormat: eLogo.ImageFormat,
		CreatedAt:   eLogo.CreatedAt,
		UpdatedAt:   eLogo.UpdatedAt,
	}, nil
}
