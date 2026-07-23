package downloader

import (
	"fmt"
	"mime"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/clientcap"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/images"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	iconfig "github.com/neosy/elengrab/internal/config"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/neosy/elengrab/internal/pkg/humanize"
	"github.com/valyala/fasthttp"
)

type renderWatchPageRequest struct {
	pageURL        string
	streamURLPath  string
	downloadID     uuid.UUID
	showBackButton bool

	allowAnonymous bool
	authCtx        dauth.UserContext
}

func (h *DownloaderHandlers) renderWatchPage(
	ctx *fasthttp.RequestCtx,
	req renderWatchPageRequest,
) {
	if ctx.IsHead() {
		ctx.SetContentType(mime.TypeByExtension(".html"))
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}

	var (
		downloadInfo *dto.GetMediaDownloadInfoResponse
		err          error
	)
	if req.allowAnonymous {
		downloadInfo, err = h.downloader.GetDownloadInfoUnrestricted(ctx, req.downloadID)
	} else {
		downloadInfo, err = h.downloader.GetDownloadInfo(ctx, req.authCtx, req.downloadID)
	}
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	var (
		format        = "-"
		videoQuality  string
		audioQuality  string
		fileSize      = "-"
		isVideoPlayer = true
		contentType   = "text/html"
	)
	if downloadInfo.MediaInfo != nil {
		ext := downloadInfo.MediaInfo.Format.Ext()

		isVideoPlayer = downloadInfo.MediaInfo.Format.IsVideo()
		format = strings.ToUpper(ext)

		if videoInfo := downloadInfo.MediaInfo.VideoInfo; videoInfo != nil {
			values := []string{strings.ToUpper(videoInfo.Codec.String())}
			if videoInfo.Resolution != dtypes.VideoResolutionNone {
				values = append(values, videoInfo.Resolution.String())
			}
			values = append(values, videoInfo.ResolutionString())

			videoQuality = strings.Join(values, " • ")

		}
		if audioInfo := downloadInfo.MediaInfo.AudioInfo; audioInfo != nil {
			values := make([]string, 0, 3)
			if audioInfo.Codec.Title() != "" {
				values = append(values, audioInfo.Codec.Title())
			}
			if audioInfo.Bitrate > 0 {
				values = append(values, fmt.Sprintf("%d kbps", audioInfo.Bitrate))
			}
			if audioInfo.SampleRate != nil && *audioInfo.SampleRate > 0 {
				values = append(values, fmt.Sprintf("%d Hz", *audioInfo.SampleRate))
			}
			audioQuality = strings.Join(values, " • ")
		}
	}

	if downloadInfo.FileSize != nil {
		fileSize = humanize.Bytes(*downloadInfo.FileSize)
	}

	description := downloadInfo.MediaDescriptionUI()

	imageData := h.thumbnailImageData(ctx, downloadInfo.MediaInfo)

	if imageData == nil && downloadInfo.ChannelID != nil && downloadInfo.IsYouTube() {
		channel, _ := h.downloader.FindYoutubeChannelInfo(ctx, *downloadInfo.ChannelID)
		if channel != nil && len(channel.ImageRaw) > 0 {
			imageData = channel.ImageData()
		}
		if imageData != nil {
			imageData.URL = paths.YoutubeChannelPath(channel.ChannelID)
		}
	}

	if imageData == nil || imageData.Width < 120 {
		imageData = h.thumbnailImageDataWithFallback(ctx, downloadInfo.MediaInfo)
	}

	if imageData == nil {
		imageData = &dtypes.ImageData{
			URL:    paths.ImagePath(images.Elengrab1280ImageJpgFileName),
			Format: dtypes.ImageFormatJPEG,
			Width:  1280,
			Height: 720,
		}
	}

	metaOgItems := make(pages.MetaOgItems, 0, 20)
	metaOgItems.Add("site_name", iconfig.AppName)
	metaOgItems.Add("url", req.pageURL)
	metaOgItems.Add("title", downloadInfo.MediaTitle)
	metaOgItems.Add("description", description)
	metaOgItems.Add("image", imageData.FullURL(h.baseURL))
	metaOgItems.Add("image:secure_url", imageData.FullURL(h.baseURL))
	metaOgItems.Add("image:type", httpx.ContentTypeByExt(imageData.Format.String()))
	metaOgItems.Add("image:alt", "Elengrab logo")
	if imageData.Width != 0 && imageData.Height != 0 {
		metaOgItems.Add("image:width", strconv.Itoa(imageData.Width))
		metaOgItems.Add("image:height", strconv.Itoa(imageData.Height))
	}
	metaOgItems.Add("type", "article")

	metaNameItems := make(pages.MetaNameItems, 0, 4)
	metaNameItems.Add("title", downloadInfo.MediaTitle)
	metaNameItems.Add("description", description)
	metaNameItems.Add("twitter:card", "summary_large_image")
	metaNameItems.Add("twitter:url", req.pageURL)
	metaNameItems.Add("twitter:title", downloadInfo.MediaTitle)
	metaNameItems.Add("twitter:description", description)
	metaNameItems.Add("twitter:image", imageData.FullURL(h.baseURL))
	metaNameItems.Add("application-title", iconfig.AppName)

	baseValues := pages.NewBaseValues()
	baseValues.Title = downloadInfo.MediaTitle
	baseValues.MetaOgItems = metaOgItems
	baseValues.MetaNameItems = metaNameItems

	cssPaths, err := h.assetPaths.WatchPageCssPaths()
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	caps := clientcap.Detect(string(ctx.UserAgent()))

	jsScripts, err := h.assetPaths.WatchPageJsPaths(caps.IsLegacyWebKit)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	pwaManifestPath, err := h.assetPaths.PwaManifestPath()
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	mediaParameters := make([]pages.MediaParameter, 0, 4)
	mediaParameters = append(mediaParameters, pages.MediaParameter{Name: "Format", Value: format})
	mediaParameters = append(mediaParameters, pages.MediaParameter{Name: "Duration", Value: downloadInfo.MediaInfo.FormatDuration()})
	if videoQuality != "" {
		mediaParameters = append(mediaParameters, pages.MediaParameter{Name: "Video", Value: videoQuality})
	}
	if audioQuality != "" {
		mediaParameters = append(mediaParameters, pages.MediaParameter{Name: "Audio", Value: audioQuality})
	}
	mediaParameters = append(mediaParameters, pages.MediaParameter{Name: "File Size", Value: fileSize})
	if downloadInfo.MediaURL != "" {
		mediaParameters = append(mediaParameters,
			pages.MediaParameter{
				Name:     "Original Source",
				Value:    mediaSourceFromURL(downloadInfo.MediaURL),
				URL:      downloadInfo.MediaURL,
				CopyIcon: icons.CopyURLIcon.FileRaw(),
			},
		)
	}

	shortURL, _ := h.linkWeb.FindLastShortURL(
		ctx,
		h.buildMediaWatchURL(downloadInfo.DownloadID),
	)
	if shortURL != "" {
		mediaParameters = append(mediaParameters,
			pages.MediaParameter{
				Type:     pages.MediaParameterTypeShareLink,
				Name:     "Share",
				Value:    "Short link",
				URL:      shortURL,
				CopyIcon: icons.CopyURLIcon.FileRaw(),
			},
		)
	}

	titleImageURL := httppaths.BuildPathMediaItemImage(
		downloadInfo.DownloadID,
		downloadInfo.ImageMetaHash(),
		[]dtypes.ImageSource{
			dtypes.ImageSourceAvatar,
			dtypes.ImageSourceSite,
		},
	)

	pageData := pages.WatchPageData{
		BaseValues: baseValues,
		BasePaths:  paths.NewHttpPaths(),
		Paths: pages.PagePaths{
			Css:         cssPaths,
			JsScripts:   jsScripts,
			PwaManifest: pwaManifestPath,
			StreamURL:   req.streamURLPath,
		},
		Values: pages.WatchPageValues{
			ShowBackButton:     req.showBackButton,
			IsVideoPlayer:      isVideoPlayer,
			MediaTitle:         downloadInfo.MediaTitle,
			MediaDescription:   downloadInfo.MediaDescription,
			MediaParameters:    mediaParameters,
			ContentType:        contentType,
			MediaTitleImageURL: titleImageURL,
		},
	}

	// Set content type so browser renders HTML properly
	ctx.SetContentType(mime.TypeByExtension(".html"))

	// Load template
	tmpl, err := h.loadPageTemplate(pages.WatchPage.FileName())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}

	// Execute template with PageTitle
	if err := tmpl.ExecuteTemplate(ctx, pages.WatchPage.Key(), pageData); err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}
}
