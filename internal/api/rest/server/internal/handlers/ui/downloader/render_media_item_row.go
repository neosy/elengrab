package downloader

import (
	"fmt"
	"html/template"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/items"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	dltypes "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/types"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/humanize"
	"github.com/valyala/fasthttp"
)

type renderMediaItemRowResponse struct {
	data       pages.RowFragmentData
	httpStatus int
	err        error
}

func (h *DownloaderHandlers) renderMediaItemRow(
	downloadInfo *dto.GetMediaDownloadInfoResponse,
	isDownloadEvent bool,
) renderMediaItemRowResponse {
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
		return renderMediaItemRowResponse{
			httpStatus: fasthttp.StatusInternalServerError,
			err:        errorx.New("the request returned an empty"),
		}
	}

	var youtubeChannelID string
	if downloadInfo.ChannelID != nil && downloadInfo.IsYouTube() {
		youtubeChannelID = *downloadInfo.ChannelID
	}

	var thumbnailID string
	if downloadInfo.MediaInfo != nil && downloadInfo.MediaInfo.PreferredThumbnailID() != uuid.Nil {
		thumbnailID = downloadInfo.MediaInfo.PreferredThumbnailID().String()
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

	data := pages.RowFragmentValues{
		DownloadID:     downloadInfo.DownloadID.String(),
		DownloadStatus: downloadInfo.Status.String(),
		WorkingStatus:  dltypes.MapUsecaseWorkingStatusToUI(downloadInfo.WorkingStatus).String(),
		IsReady:        downloadInfo.Status.IsReady(),

		YoutubeChannelID: youtubeChannelID,
		AvatarTitle:      downloadInfo.AvatarTitle,

		ThumbnailID:         thumbnailID,
		ThumbnailIsPortrait: downloadInfo.ThumbnalIsPortrait,
		ThumbnailURL:        h.thumbnailURLWithFallback(downloadInfo.MediaInfo),

		MediaTitle: downloadInfo.MediaTitle,
		MediaURL:   downloadInfo.MediaURL,

		ContentTimeAgo: downloadInfo.CreatedTimeAgo,

		ImageURL:       downloadItemImageURL,
		ImageAvatarURL: downloadItemImageAvatarURL,
		ImageSiteURL:   downloadItemImageSiteURL,

		DownloadRowPath:    httppaths.BuildPathMediaItemRow(downloadInfo.DownloadID),
		DownloadRepeatPath: httppaths.BuildPathMediaItemDownloadRepeat(downloadInfo.DownloadID),

		FileSize:   "-",
		Format:     "-",
		DataFormat: "-",

		FormatTitle:   downloadInfo.MediaInfoText,
		FormatTooltip: downloadInfo.MediaInfoTooltip,

		DownloadURL: downloadURL,
		WatchURL:    watchURL,
		StreamURL:   streamURL,
		DeleteURL:   httppaths.BuildMediaItemPath(downloadInfo.DownloadID),

		RowID:      "row-" + downloadInfo.DownloadID.String(),
		ProgressID: "progress-" + downloadInfo.DownloadID.String(),

		IsItemHTMXOptionRepeat: isGrabResultItemHTMXOptionRepeat,
		IsDownloadEvent:        isDownloadEvent,
		ResultRowStatusTitle:   downloadInfo.StatusText,

		UserName: downloadInfo.UserLogin,

		RefreshingIcon:            icons.DownloadRefreshingIcon.FileRaw(),
		MetaUserNameSeparatorIcon: icons.DownloadMetaUserNameSeparatorIcon.FileRaw(),

		DownloaderResultItemSourceLinkIcon: icons.DownloadSourceLinkIcon.FileRaw(),
		DownloaderResultItemStatusIcon:     icons.DownloaderIconByStatus(downloadInfo.Status).FileRaw(),
		DownloaderResultItemDeleteIcon:     icons.DownloadDeleteIcon.FileRaw(),

		ResultMediaUrlFade: "",
		ResultSizeFade:     "",
		ResultFormatFade:   "",
		IsItemSpiner:       downloadInfo.Status == dtypes.MediaDownloadStatusWorking,
	}

	if downloadInfo.Status == dtypes.MediaDownloadStatusFailed {
		data.DownloaderResultItemStatusFailedIcon = template.HTML(
			icons.DownloadRepeatIcon.FileRaw(),
		)
	}

	if downloadInfo.FileSize != nil && *downloadInfo.FileSize > 0 {
		data.FileSize = humanize.Bytes(*downloadInfo.FileSize)
	}
	if downloadInfo.FileExt != "" {
		data.Format = downloadInfo.FileExt
		data.DataFormat = downloadInfo.FileExt
	}
	if downloadInfo.MediaInfo != nil {
		data.Duration = downloadInfo.MediaInfo.FormatDuration()
		data.VideoIsShort = downloadInfo.MediaInfo.IsPortrait()
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
		extraData[items.DownloadingProgressPercentKey] = int(downloadInfo.Progress.Percent())
	}

	pageData := pages.RowFragmentData{
		BasePaths:     paths.NewHttpPaths(),
		Values:        &data,
		IconFileNames: icons.FileNamesByKey(),
		Extra:         extraData,
	}

	return renderMediaItemRowResponse{
		data:       pageData,
		httpStatus: fasthttp.StatusOK,
		err:        nil,
	}
}
