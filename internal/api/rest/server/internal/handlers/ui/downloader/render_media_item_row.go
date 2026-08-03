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

	link, _ := h.linkWeb.ResolveURL(ctx, h.buildMediaWatchURL(params.downloadInfo.DownloadID))

	var shareLinkIcon template.HTML
	if link != nil {
		shareLinkIcon = icons.MediaShareLinkIcon.FileRaw()
	}

	var watchURL, downloadURL, streamURL string
	if params.downloadInfo.Status == dtypes.MediaDownloadStatusDone {
		watchURL = httppaths.BuildPathMediaItemWatch(params.downloadInfo.DownloadID)
		downloadURL = httppaths.BuildPathMediaItemDownload(params.downloadInfo.DownloadID)
		streamURL = httppaths.BuildPathMediaItemStream(params.downloadInfo.DownloadID)
	}

	var audioIcon template.HTML
	if params.downloadInfo.MediaInfo != nil {
		if params.downloadInfo.MediaInfo.FormatType == dtypes.FormatTypeAudioOnly {
			audioIcon = icons.MediaAudioIcon.FileRaw()
		}
	}

	var watchedIcon template.HTML
	if params.downloadInfo.UserWatched {
		watchedIcon = icons.MediaWatchedIcon.FileRaw()
	}

	data := pages.RowFragmentValues{
		DownloadID:     idcodec.EncodeUUIDBase64URL(params.downloadInfo.DownloadID),
		DownloadStatus: params.downloadInfo.Status.String(),
		WorkingStatus:  dltypes.MapUsecaseWorkingStatusToUI(params.downloadInfo.WorkingStatus).String(),
		Visibility:     params.downloadInfo.Visibility.String(),
		IsReady:        params.downloadInfo.Status.IsReady(),

		YoutubeChannelID: youtubeChannelID,
		AvatarTitle:      params.downloadInfo.AvatarTitle,

		ThumbnailID:         thumbnailID,
		ThumbnailIsPortrait: params.downloadInfo.ThumbnalIsPortrait,
		ThumbnailURL:        h.thumbnailURLWithFallback(params.downloadInfo.MediaInfo),

		MediaTitle: params.downloadInfo.MediaTitle,
		MediaURL:   params.downloadInfo.MediaURL,

		ContentTimeAgo:   params.downloadInfo.CreatedTimeAgo,
		ContentViewCount: humanize.CompactNumber(params.downloadInfo.ViewCount),

		WatchPercent: params.downloadInfo.UserWatchDisplayPercent(),
		Watched:      params.downloadInfo.UserWatched,

		ImageURL:       downloadItemImageURL,
		ImageAvatarURL: downloadItemImageAvatarURL,
		ImageSiteURL:   downloadItemImageSiteURL,

		DownloadRowPath:    httppaths.BuildPathMediaItemRow(params.downloadInfo.DownloadID),
		DownloadRepeatPath: httppaths.BuildPathMediaItemDownloadRepeat(params.downloadInfo.DownloadID),

		FileSize:   "-",
		Format:     "-",
		DataFormat: "-",

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

		PublicIcon:  icons.MediaPublicIcon.FileRaw(),
		PrivateIcon: icons.MediaPrivateIcon.FileRaw(),

		AudioIcon:     audioIcon,
		ShareLinkIcon: shareLinkIcon,
		WatchedIcon:   watchedIcon,

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

	if params.downloadInfo.FileSize != nil && *params.downloadInfo.FileSize > 0 {
		data.FileSize = humanize.Bytes(*params.downloadInfo.FileSize)
	}

	if params.downloadInfo.FileExt != "" {
		data.Format = params.downloadInfo.FileExt
		data.DataFormat = params.downloadInfo.FileExt
	}

	if params.downloadInfo.MediaInfo != nil {
		data.Duration = params.downloadInfo.MediaInfo.FormatDuration()
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
