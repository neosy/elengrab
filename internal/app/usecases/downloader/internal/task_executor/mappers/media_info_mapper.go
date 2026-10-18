package mappers

import (
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/task_executor/types"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

func (m *Mappers) MapMediaInfoDomain(
	srvMediaInfo *pservices.MediaInfo,
	thumbnailIDs types.ThumbnailIDs,
) *dtypes.MediaInfo {
	mediaInfo := srvMediaInfo.ToDomain()

	mediaInfo.ThumbnailID = thumbnailIDs.ThumbnailID
	mediaInfo.FrameThumbnailID = thumbnailIDs.FrameThumbnailID

	return mediaInfo
}
