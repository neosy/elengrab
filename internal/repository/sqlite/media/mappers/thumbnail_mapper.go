package mappers

import (
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	emedia "github.com/neosy/elengrab/internal/repository/sqlite/media/entity"
)

func (m *Mappers) MapThumbnailDomainToEntity(thumbnail *dmedia.Thumbnail) (*emedia.Thumbnail, error) {
	return &emedia.Thumbnail{
		ThumbID:    thumbnail.ThumbID,
		MediaID:    thumbnail.MediaID,
		Variant:    thumbnail.Variant.String(),
		Version:    thumbnail.Version,
		Width:      thumbnail.Width,
		Height:     thumbnail.Height,
		Format:     thumbnail.Format.String(),
		SourceType: thumbnail.SourceType.String(),
		SourceID:   thumbnail.SourceID,
		SourceURL:  thumbnail.SourceURL,
		StorageKey: thumbnail.StorageKey,
		IsPrimary:  thumbnail.IsPrimary,
	}, nil
}

func (m *Mappers) MapThumbnailEntityToDomain(eThumbnail *emedia.Thumbnail) (*dmedia.Thumbnail, error) {
	variant, err := dtypes.ParseThumbnailVariant(eThumbnail.Variant)
	if err != nil {
		return nil, err
	}

	format, err := dtypes.ParseImageFormat(eThumbnail.Format)
	if err != nil {
		return nil, err
	}

	sourceType, err := dtypes.ParseThumbnailSourceType(eThumbnail.SourceType)
	if err != nil {
		return nil, err
	}

	return &dmedia.Thumbnail{
		ThumbID:    eThumbnail.ThumbID,
		MediaID:    eThumbnail.MediaID,
		Variant:    variant,
		Version:    eThumbnail.Version,
		Width:      eThumbnail.Width,
		Height:     eThumbnail.Height,
		Format:     format,
		SourceType: sourceType,
		SourceID:   eThumbnail.SourceID,
		SourceURL:  eThumbnail.SourceURL,
		StorageKey: eThumbnail.StorageKey,
		IsPrimary:  eThumbnail.IsPrimary,
		CreatedAt:  eThumbnail.CreatedAt,
		UpdatedAt:  eThumbnail.UpdatedAt,
	}, nil
}
