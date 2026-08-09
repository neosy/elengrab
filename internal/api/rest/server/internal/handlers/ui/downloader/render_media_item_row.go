package downloader

import (
	"context"
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
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

type renderMediaItemRowResponse struct {
	data       pages.RowFragmentData
	httpStatus int
	err        error
}

type renderMediaItemRowParams struct {
	downloadInfo    *dto.MediaDownloadInfo
	isDownloadEvent bool
	lazyLoadImages  bool
}

func (h *DownloaderHandlers) renderMediaItemRow(
	ctx context.Context,
	params renderMediaItemRowParams,
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

	if params.downloadInfo == nil {
		return renderMediaItemRowResponse{
			httpStatus: fasthttp.StatusInternalServerError,
			err:        errorx.New("the request returned an empty"),
		}
	}

	var youtubeChannelID string
	if params.downloadInfo.ChannelID != nil && params.downloadInfo.IsYouTube() {
		youtubeChannelID = *params.downloadInfo.ChannelID
	}

	var thumbnailID string
	if params.downloadInfo.MediaInfo != nil {
		if params.downloadInfo.MediaInfo.PreferredThumbnailID() != uuid.Nil {
			thumbnailID = params.downloadInfo.MediaInfo.PreferredThumbnailID().String()
		}
	}

	downloadItemImageURL := httppaths.BuildPathMediaItemImage(
		params.downloadInfo.DownloadID,
		params.downloadInfo.ImageMetaHash(),
		[]dtypes.ImageSource{
			dtypes.ImageSourceThumbnail,
			dtypes.ImageSourceAvatar,
			dtypes.ImageSourceSite,
		},
	)

	downloadItemImageAvatarURL := httppaths.BuildPathMediaItemImage(
		params.downloadInfo.DownloadID,
		params.downloadInfo.ImageMetaHash(),
		[]dtypes.ImageSource{
			dtypes.ImageSourceAvatar,
			dtypes.ImageSourceSite,
		},
	)

	downloadItemImageSiteURL := httppaths.BuildPathMediaItemImage(
		params.downloadInfo.DownloadID,
		params.downloadInfo.ImageMetaHash(),
		[]dtypes.ImageSource{
			dtypes.ImageSourceSite,
		},
	)

	var isGrabResultItemHTMXOptionRepeat = false
	switch params.downloadInfo.Status {
	case dtypes.MediaDownloadStatusNew, dtypes.MediaDownloadStatusPending, dtypes.MediaDownloadStatusWorking:
		isGrabResultItemHTMXOptionRepeat = true
	}

	var watchURL, downloadURL, streamURL string
	if params.downloadInfo.Status == dtypes.MediaDownloadStatusDone {
		watchURL = httppaths.BuildPathMediaItemWatch(params.downloadInfo.DownloadID)
		downloadURL = httppaths.BuildPathMediaItemDownload(params.downloadInfo.DownloadID)
		streamURL = httppaths.BuildPathMediaItemStream(params.downloadInfo.DownloadID)
	}

	var (
		audioIcon     template.HTML
		mediaDuration = humanize.DurationClock(0)

		fileSize   = "-"
		format     = "-"
		dataFormat = "-"

		thumbnailWidth       template.CSS
		thumbnailHeight      template.CSS
		thumbnailAspectRatio template.CSS
	)
	if mediaInfo := params.downloadInfo.MediaInfo; mediaInfo != nil {
		if mediaInfo.FormatType == dtypes.FormatTypeAudioOnly {
			audioIcon = icons.MediaAudioIcon.FileRaw()
		}

		if params.downloadInfo.FileSize != nil && *params.downloadInfo.FileSize > 0 {
			fileSize = humanize.Bytes(*params.downloadInfo.FileSize)
		}

		if params.downloadInfo.FileExt != "" {
			format = params.downloadInfo.FileExt
			dataFormat = params.downloadInfo.FileExt
		}

		mediaDuration = mediaInfo.FormatDuration()

		var width, height int
		if mediaInfo.VideoInfo != nil {
			width = mediaInfo.VideoInfo.Width
			height = mediaInfo.VideoInfo.Height
		}

		if width == 0 || height == 0 || width > height {
			thumbnailWidth = "720px"
			thumbnailHeight = "auto"
			thumbnailAspectRatio = "16/9"
		} else {
			thumbnailHeight = "720px"
			thumbnailWidth = "auto"
			thumbnailAspectRatio = "9/16"
		}
	}

	visibility := h.getVisibilityResponse(params.downloadInfo)

	data := pages.RowFragmentValues{
		DownloadID:      idcodec.EncodeUUIDBase64URL(params.downloadInfo.DownloadID),
		DownloadStatus:  params.downloadInfo.Status.String(),
		WorkingStatus:   dltypes.MapUsecaseWorkingStatusToUI(params.downloadInfo.WorkingStatus).String(),
		Visibility:      params.downloadInfo.Visibility.String(),
		VisibilityLabel: params.downloadInfo.Visibility.Label(),
		IsReady:         params.downloadInfo.Status.IsReady(),
		HasShareLink:    h.hasShareLink(ctx, params.downloadInfo.DownloadID),

		YoutubeChannelID: youtubeChannelID,
		AvatarTitle:      params.downloadInfo.AvatarTitle,

		LazyLoadImages: params.lazyLoadImages,

		ThumbnailID:         thumbnailID,
		ThumbnailIsPortrait: params.downloadInfo.ThumbnalIsPortrait,
		ThumbnailURL:        h.thumbnailURLWithFallback(params.downloadInfo.MediaInfo),

		ThumbnailWidth:       thumbnailWidth,
		ThumbnailHeight:      thumbnailHeight,
		ThumbnailAspectRatio: thumbnailAspectRatio,

		MediaTitle: params.downloadInfo.MediaTitle,
		MediaURL:   params.downloadInfo.MediaURL,

		ContentTimeAgo:   params.downloadInfo.CreatedTimeAgo,
		ContentViewCount: humanize.CompactNumber(params.downloadInfo.ViewCount),

		WatchIndicatorEnabled: params.downloadInfo.WatchIndicatorEnabled(),
		WatchPercent:          params.downloadInfo.UserWatchDisplayPercent(),
		Watched:               params.downloadInfo.UserWatched,

		ImageURL:       downloadItemImageURL,
		ImageAvatarURL: downloadItemImageAvatarURL,
		ImageSiteURL:   downloadItemImageSiteURL,

		DownloadRowPath:    httppaths.BuildPathMediaItemRow(params.downloadInfo.DownloadID),
		DownloadRepeatPath: httppaths.BuildPathMediaItemDownloadRepeat(params.downloadInfo.DownloadID),

		FileSize:   fileSize,
		Format:     format,
		DataFormat: dataFormat,
		Duration:   mediaDuration,

		FormatTitle:   params.downloadInfo.MediaInfoText,
		FormatTooltip: params.downloadInfo.MediaInfoTooltip,

		DownloadURL: downloadURL,
		WatchURL:    watchURL,
		StreamURL:   streamURL,
		DeleteURL:   httppaths.BuildMediaItemPath(params.downloadInfo.DownloadID),

		ItemID:     idcodec.EncodeUUIDBase64URL(params.downloadInfo.DownloadID),
		RowID:      "row-" + idcodec.EncodeUUIDBase64URL(params.downloadInfo.DownloadID),
		ProgressID: "progress-" + idcodec.EncodeUUIDBase64URL(params.downloadInfo.DownloadID),

		IsItemHTMXOptionRepeat: isGrabResultItemHTMXOptionRepeat,
		IsDownloadEvent:        params.isDownloadEvent,
		ResultRowStatusTitle:   params.downloadInfo.StatusText,

		UserName: params.downloadInfo.UserLogin,

		RefreshingIcon:            icons.DownloadRefreshingIcon.FileRaw(),
		MetaUserNameSeparatorIcon: icons.DownloadMetaUserNameSeparatorIcon.FileRaw(),

		VisibilityIcon: visibility.Icon,

		AudioIcon:     audioIcon,
		ShareLinkIcon: icons.MediaShareLinkIcon.FileRaw(),
		WatchedIcon:   icons.MediaWatchedIcon.FileRaw(),

		DownloaderResultItemSourceLinkIcon: icons.DownloadSourceLinkIcon.FileRaw(),
		DownloaderResultItemStatusIcon:     icons.DownloaderIconByStatus(params.downloadInfo.Status).FileRaw(),
		DownloaderResultItemDeleteIcon:     icons.DownloadDeleteIcon.FileRaw(),

		ResultMediaUrlFade: "",
		ResultSizeFade:     "",
		ResultFormatFade:   "",
		IsItemSpiner:       params.downloadInfo.Status == dtypes.MediaDownloadStatusWorking,
	}

	if params.downloadInfo.Status == dtypes.MediaDownloadStatusFailed {
		data.DownloaderResultItemStatusFailedIcon = template.HTML(
			icons.DownloadRepeatIcon.FileRaw(),
		)
	}

	if params.downloadInfo.MediaInfo != nil {
		data.IsPortrait = params.downloadInfo.MediaInfo.IsPortrait()
		data.IsAudioOnly = params.downloadInfo.MediaInfo.IsAudioOnly()
		data.IsShorts = params.downloadInfo.MediaInfo.IsShorts()
		data.LoopPlayback = data.IsShorts || data.IsAudioOnly
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

	if params.downloadInfo.Progress != nil {
		extraData[items.DownloadingProgressPercentKey] = int(params.downloadInfo.Progress.Percent())
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
