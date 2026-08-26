package mappers

import (
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/download_executor/types"
	dservices "github.com/neosy/elengrab/internal/domain/services"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (m *Mappers) MapMediaInfoDomain(
	srvMediaInfo *dservices.MediaInfo,
	thumbnailIDs types.ThumbnailIDs,
) *dtypes.MediaInfo {
	mediaInfo := srvMediaInfo.NewMediaInfoDomain()

	mediaInfo.ThumbnailID = thumbnailIDs.ThumbnailID
	mediaInfo.FrameThumbnailID = thumbnailIDs.FrameThumbnailID

	return mediaInfo
}
