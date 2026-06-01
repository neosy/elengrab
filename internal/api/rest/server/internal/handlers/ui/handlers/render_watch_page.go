package handlers

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/images"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/paths"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	iconfig "github.com/neosy/elengrab/internal/config"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/neosy/elengrab/internal/pkg/fnx"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/neosy/elengrab/internal/pkg/humanize"
	"github.com/valyala/fasthttp"
)

type renderWatchPageRequest struct {
	pageURL        string
	streamPath     string
	downloadID     uuid.UUID
	showBackButton bool
}

func (h *DownloaderHandlers) renderWatchPage(
	ctx *fasthttp.RequestCtx,
	req renderWatchPageRequest,
) {
	if ctx.IsHead() {
		ctx.SetContentType("text/html; charset=utf-8")
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}

	downloadInfo, err := h.downloader.GetDownloadInfoUnrestricted(ctx, req.downloadID)
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
		videoWidth    = 0
		videoHeight   = 0
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

			videoWidth = videoInfo.Width
			videoHeight = videoInfo.Height
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

	mediaURL := h.baseURL + req.streamPath
	prefixType := fnx.Ternary(isVideoPlayer, "video", "audio")
	description := downloadInfo.MediaTitle + fmt.Sprintf(" [%s]", downloadInfo.MediaInfoText)

	var imageData *dtypes.ImageData
	if downloadInfo.MediaInfo != nil && downloadInfo.MediaInfo.PreferredThumbnailID() != uuid.Nil {
		thumbnail, _ := h.thumbnail.GetByThumbID(ctx, downloadInfo.MediaInfo.PreferredThumbnailID())
		if thumbnail != nil {
			imageData = thumbnail.ImageData()
		}
		if imageData != nil {
			imageData.URL = h.baseURL + paths.ThumbnailPath(thumbnail.ThumbID.String())
		}
	}

	if imageData == nil && downloadInfo.ChannelID != nil && downloadInfo.IsYouTube() {
		channel, _ := h.downloader.FindYoutubeChannelInfo(ctx, *downloadInfo.ChannelID)
		if channel != nil && len(channel.ImageRaw) > 0 {
			imageData = channel.ImageData()
		}
		if imageData != nil {
			imageData.URL = h.baseURL + paths.YoutubeChannelPath(channel.ChannelID)
		}
	}

	if imageData == nil || imageData.Width < 120 {
		imageData = &dtypes.ImageData{
			URL:    h.baseURL + paths.ImagePath(images.Elengrab1280ImageJpgFileName),
			Width:  1280,
			Height: 720,
		}
	}

	metaOgItems := make(pages.MetaOgItems, 0, 20)
	metaOgItems.Add("site_name", iconfig.AppName)
	metaOgItems.Add("url", req.pageURL)
	metaOgItems.Add("title", downloadInfo.MediaTitle)
	metaOgItems.Add("description", description)
	metaOgItems.Add("image", imageData.URL)
	metaOgItems.Add("image:secure_url", imageData.URL)
	metaOgItems.Add("image:type", httpx.ContentTypeByExt(imageData.Format.String()))
	metaOgItems.Add("image:alt", "Elengrab logo")
	if imageData.Width != 0 && imageData.Height != 0 {
		metaOgItems.Add("image:width", strconv.Itoa(imageData.Width))
		metaOgItems.Add("image:height", strconv.Itoa(imageData.Height))
	}
	metaOgItems.Add("type", fnx.Ternary(isVideoPlayer, "video.other", "music.song"))

	metaOgItems.Add(fmt.Sprintf("%s:url", prefixType), mediaURL)
	metaOgItems.Add(fmt.Sprintf("%s:secure_url", prefixType), mediaURL)
	if contentType != "" {
		metaOgItems.Add(fmt.Sprintf("%s:type", prefixType), contentType)
	}
	if isVideoPlayer && videoWidth != 0 {
		metaOgItems.Add("video:width", strconv.Itoa(videoWidth))
		metaOgItems.Add("video:height", strconv.Itoa(videoHeight))
	}

	metaNameItems := make(pages.MetaNameItems, 0, 4)
	metaNameItems.Add("title", downloadInfo.MediaTitle)
	metaNameItems.Add("description", description)
	metaNameItems.Add("twitter:card", "summary_large_image")
	metaNameItems.Add("twitter:url", req.pageURL)
	metaNameItems.Add("twitter:title", downloadInfo.MediaTitle)
	metaNameItems.Add("twitter:description", description)
	metaNameItems.Add("twitter:image", imageData.URL)
	metaNameItems.Add("application-title", iconfig.AppName)

	baseValues := pages.NewBaseValues()
	baseValues.Title = downloadInfo.MediaTitle
	baseValues.ShowFooter = false
	baseValues.MetaOgItems = metaOgItems
	baseValues.MetaNameItems = metaNameItems

	cssPaths, err := paths.CssWatchPaths(h.assetFolders.Css())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	jsScripts, err := paths.WatchPageJsPaths(h.assetFolders.Js())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	jsImportJSON, err := paths.WatchPageJsImportJSON(h.assetFolders.Js())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	pwaManifestPath, err := paths.PwaManifestPath(h.assetFolders.Pwa())
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
				Name:  "Original Source",
				Value: mediaSourceFromURL(downloadInfo.MediaURL),
				URL:   downloadInfo.MediaURL,
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
		BasePaths:  paths.NewPaths(),
		Paths: pages.PagePaths{
			Css:             cssPaths,
			JsScripts:       jsScripts,
			JsImportMapJSON: template.HTML(jsImportJSON),
			PwaManifest:     pwaManifestPath,
			Stream:          req.streamPath,
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
	ctx.SetContentType("text/html; charset=utf-8")

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
