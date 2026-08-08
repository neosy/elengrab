package mappers

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
)

func (m *Mappers) MapMediaDownloadChangedToEventResponse(
	req *ucdto.MediaDownloadChanged,
	getVisibility func(*ucdto.MediaDownloadInfo) *dto.VisibilityResponse,
	getHasShareLink func(downloadID uuid.UUID) bool,
) *dto.EventDownloadRowPatchResponse {
	var title string
	if req.TitleChanged {
		title = req.Info.MediaTitle
	}

	var visibility *dto.VisibilityResponse
	if req.VisibilityChanged {
		visibility = getVisibility(req.Info)
	}

	var hasShareLink *bool
	if req.ShareLinkChanged {
		hasShareLink = new(getHasShareLink(req.Info.DownloadID))
	}

	var watchPercent string
	if req.WatchPositionChanged && req.Info.WatchIndicatorEnabled() {
		watchPercent = strconv.FormatFloat(req.Info.UserWatchDisplayPercent(), 'f', 2, 64)
	}

	var watched *bool
	if req.WatchedChanged {
		watched = &req.Info.UserWatched
	}

	var viewCount *dto.ViewCountResponse
	if req.ViewCountChanged {
		viewCount = &dto.ViewCountResponse{
			Count: req.Info.ViewCount,
			Text:  req.Info.ViewCountText(),
		}
	}

	return &dto.EventDownloadRowPatchResponse{
		DownloadID: idcodec.EncodeUUIDBase64URL(req.DownloadID),

		Title:        title,
		Visibility:   visibility,
		HasShareLink: hasShareLink,

		WatchPercent: watchPercent,
		Watched:      watched,
		ViewCount:    viewCount,
	}
}
