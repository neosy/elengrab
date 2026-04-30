package mappers

import (
	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	emedia "github.com/neosy/elengrab/internal/repository/sqlite/media/entity"
)

func (m *Mappers) MapSiteLogoDomainToEntity(logo *dmedia.SiteLogo) (*emedia.SiteLogo, error) {
	return &emedia.SiteLogo{
		LogoID:      logo.LogoID.String(),
		SiteURL:     logo.SiteURL,
		SiteTitle:   logo.SiteTitle,
		ImageURL:    logo.ImageURL,
		ImageRaw:    logo.ImageRaw,
		ImageFormat: logo.ImageFormat.String(),
	}, nil
}

func (m *Mappers) MapSiteLogoEntityToDomain(eLogo *emedia.SiteLogo) (*dmedia.SiteLogo, error) {
	logoID, err := uuid.Parse(eLogo.LogoID)
	if err != nil {
		return nil, err
	}

	imagesFormat, err := dtypes.ParseImageFormat(eLogo.ImageFormat)
	if err != nil {
		return nil, err
	}

	return &dmedia.SiteLogo{
		LogoID:      logoID,
		SiteURL:     eLogo.SiteURL,
		SiteTitle:   eLogo.SiteTitle,
		ImageURL:    eLogo.ImageURL,
		ImageRaw:    eLogo.ImageRaw,
		ImageFormat: imagesFormat,
		CreatedAt:   eLogo.CreatedAt,
		UpdatedAt:   eLogo.UpdatedAt,
	}, nil
}
