package dto

import (
	"uuid"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type CreateThumbnailRequest struct {
	MediaID    uuid.UUID
	Variant    dtypes.ThumbnailVariant
	SourceType dtypes.ThumbnailSourceType
	SourceURL  *string
	SourceID   *string
	ImageData  *dtypes.ImageData
}
