package handlers

import (
	"fmt"
	"html/template"
	"path/filepath"

	dltypes "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/types"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
	"github.com/valyala/fasthttp"
)

type genRowResult struct {
	data       uivalues.RowFragmentData
	httpStatus int
	err        error
}

func (h *DownloaderHandlers) genRow(
	downloadInfo *dto.GetMediaDownloadInfoResponse,
	isDownloadEvent bool,
) genRowResult {
	var (
		cacheChanged = struct {
			youtubeChannelID bool
			mediaTitle       bool
			FileSize         bool
			Format           bool
			Status           bool
			ProgressPercent  bool
		}{}
	)

	if downloadInfo == nil {
		return genRowResult{
			httpStatus: fasthttp.StatusInternalServerError,
			err:        errorx.New("the request returned an empty"),
		}
	}

	var youtubeChannelID string
	if downloadInfo.ChannelID != nil && downloadInfo.IsYouTube() {
		youtubeChannelID = *downloadInfo.ChannelID
	}

	var thumbnailID string
	if downloadInfo.MediaInfo != nil && downloadInfo.MediaInfo.GetThumbnailID() != nil {
		thumbnailID = downloadInfo.MediaInfo.GetThumbnailID().String()
	}

	downloadItemImageURL := httppaths.BuildPathMediaItemImage(
		downloadInfo.DownloadID,
		downloadInfo.ImageMetaHash(),
		[]dtypes.ImageSource{
			dtypes.ImageSourceThumbnail,
			dtypes.ImageSourceAvatar,
			dtypes.ImageSourceSite,
		},
	)

	downloadItemImageAvatarURL := httppaths.BuildPathMediaItemImage(
		downloadInfo.DownloadID,
		downloadInfo.ImageMetaHash(),
		[]dtypes.ImageSource{
			dtypes.ImageSourceAvatar,
			dtypes.ImageSourceSite,
		},
	)

	downloadItemImageSiteURL := httppaths.BuildPathMediaItemImage(
		downloadInfo.DownloadID,
		downloadInfo.ImageMetaHash(),
		[]dtypes.ImageSource{
			dtypes.ImageSourceSite,
		},
	)

	iconsDir := filepath.Join(h.assetsDir, "static/img/icons")

	var isGrabResultItemHTMXOptionRepeat = false
	switch downloadInfo.Status {
	case dtypes.MediaDownloadStatusNew, dtypes.MediaDownloadStatusPending, dtypes.MediaDownloadStatusWorking:
		isGrabResultItemHTMXOptionRepeat = true
	}

	var watchURL, downloadURL, streamURL string
	if downloadInfo.Status == dtypes.MediaDownloadStatusDone {
		watchURL = httppaths.BuildPathMediaItemWatch(downloadInfo.DownloadID)
		downloadURL = httppaths.BuildPathMediaItemDownload(downloadInfo.DownloadID)
		streamURL = httppaths.BuildPathMediaItemStream(downloadInfo.DownloadID)
	}

	data := uivalues.RowFragmentValues{
		DownloadID:     downloadInfo.DownloadID.String(),
		DownloadStatus: downloadInfo.Status.String(),
		WorkingStatus:  dltypes.MapUsecaseWorkingStatusToUI(downloadInfo.WorkingStatus).String(),

		YoutubeChannelID: youtubeChannelID,
		ThumbnailID:      thumbnailID,
		AvatarTitle:      downloadInfo.AvatarTitle,

		MediaTitle: downloadInfo.MediaTitle,
		MediaURL:   downloadInfo.MediaURL,

		ContentTimeAgo: downloadInfo.CreatedTimeAgo,

		ImageURL:       downloadItemImageURL,
		ImageAvatarURL: downloadItemImageAvatarURL,
		ImageSiteURL:   downloadItemImageSiteURL,

		DownloadRowPath:    httppaths.BuildPathMediaItemRow(downloadInfo.DownloadID),
		DownloadRepeatPath: httppaths.BuildPathMediaItemDownloadRepeat(downloadInfo.DownloadID),

		FileSize:      "-",
		Format:        "-",
		DataFormat:    "-",
		FormatTitle:   downloadInfo.MediaInfoText,
		FormatTooltip: downloadInfo.MediaInfoTooltip,

		DownloadURL: downloadURL,
		WatchURL:    watchURL,
		StreamURL:   streamURL,
		DeleteURL:   httppaths.BuildPathMediaItem(httppaths.PathMediaItem, downloadInfo.DownloadID),

		RowID:      "row-" + downloadInfo.DownloadID.String(),
		ProgressID: "progress-" + downloadInfo.DownloadID.String(),

		IsItemHTMXOptionRepeat:         isGrabResultItemHTMXOptionRepeat,
		IsDownloadEvent:                isDownloadEvent,
		ResultRowStatusTitle:           downloadInfo.StatusText,
		DownloaderResultItemStatusIcon: template.HTML(uivalues.DownloaderResultStatusIconSvgRaw(downloadInfo.Status, iconsDir)),
		DownloaderResultItemDeleteIcon: template.HTML(uivalues.IconFileRaw(uivalues.IconFileName(uivalues.DownloadDeleteIconNameKey), iconsDir)),
		ResultMediaUrlFade:             "",
		ResultSizeFade:                 "",
		ResultFormatFade:               "",
		IsItemSpiner:                   downloadInfo.Status == dtypes.MediaDownloadStatusWorking,
	}

	if downloadInfo.Status == dtypes.MediaDownloadStatusFailed {
		data.DownloaderResultItemStatusFailedIcon = template.HTML(
			uivalues.IconFileRaw(uivalues.IconFileName(uivalues.DownloadRepeatIconNameKey), iconsDir),
		)
	}

	if downloadInfo.FileSize != nil && *downloadInfo.FileSize > 0 {
		data.FileSize = uformat.BytesHuman(*downloadInfo.FileSize)
	}
	if downloadInfo.FileExt != "" {
		data.Format = downloadInfo.FileExt
		data.DataFormat = downloadInfo.FileExt
	}
	if downloadInfo.MediaInfo != nil {
		data.IsAudio = fmt.Sprint(downloadInfo.MediaInfo.FormatType == dtypes.FormatTypeAudioOnly)
	}

	if cacheChanged.mediaTitle {
		data.ResultMediaUrlFade = "fade-text"
	}
	if cacheChanged.FileSize {
		data.ResultSizeFade = "fade-text"
	}
	if cacheChanged.Format {
		data.ResultFormatFade = "fade-text"
	}

	extraData := make(map[string]any)

	if downloadInfo.Progress != nil {
		extraData[uivalues.DownloadingProgressPercentKey] = int(downloadInfo.Progress.Percent())
	}

	pageData := uivalues.RowFragmentData{
		BasePaths:     uivalues.NewBasePaths(),
		Values:        data,
		IconFileNames: uivalues.IconFileNames(),
		Extra:         extraData,
	}

	return genRowResult{
		data:       pageData,
		httpStatus: fasthttp.StatusOK,
		err:        nil,
	}
}
