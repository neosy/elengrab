package mappers

import (
	"github.com/google/uuid"
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

func (m *Mappers) MapCreateThumbnailRequestToThumbnailDomain(
	req *dto.CreateThumbnailRequest,
) (*dmedia.Thumbnail, error) {
	if req == nil {
		return nil, apperrors.ErrFuncParamNullPointer
	}

	if req.ImageData == nil {
		return nil, apperrors.ErrFuncContainsEmptyFields
	}

	thumbnail := dmedia.NewThumbnail()

	thumbnail.MediaID = req.MediaID
	if thumbnail.MediaID == uuid.Nil {
		thumbnail.MediaID = uuid.New()
	}

	thumbnail.Format = req.ImageData.Format
	thumbnail.SourceType = req.SourceType
	thumbnail.SourceURL = req.SourceURL
	thumbnail.SourceID = req.SourceID

	if req.ImageData.URL != "" {
		thumbnail.SourceURL = &req.ImageData.URL
	}

	if req.ImageData.Width > 0 {
		thumbnail.Width = uptr.Any(uint16(req.ImageData.Width))
	}

	if req.ImageData.Height > 0 {
		thumbnail.Height = uptr.Any(uint16(req.ImageData.Height))
	}

	return thumbnail, nil
}
