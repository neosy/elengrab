package handlers

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/google/uuid"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	iconfig "github.com/neosy/elengrab/internal/config"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/neosy/elengrab/internal/pkg/fnx"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) watch(
	ctx *fasthttp.RequestCtx,
	pageURL string,
	streamPath string,
	downloadID uuid.UUID,
	showBackButton bool,
) {
	if ctx.IsHead() {
		ctx.SetContentType("text/html; charset=utf-8")
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}

	downloadInfo, err := h.downloader.GetDownloadInfoUnrestricted(ctx, downloadID)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	var (
		format         = "-"
		videoQuality   string
		audioQuality   string
		fileSize       = "-"
		isVideoPlayer  = true
		contentType    = "text/html"
		videoWidth     = 0
		videoHeight    = 0
		originalSource template.HTML
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

	originalSource = template.HTML(
		fmt.Sprintf(
			`<a href="%s" target="_blank">%s</a>`,
			downloadInfo.MediaURL,
			mediaSourceFromURL(downloadInfo.MediaURL),
		),
	)

	if downloadInfo.FileSize != nil {
		fileSize = uformat.BytesHuman(*downloadInfo.FileSize)
	}

	mediaURL := h.baseURL + streamPath
	prefixType := fnx.Ternary(isVideoPlayer, "video", "audio")
	description := downloadInfo.MediaTitle + fmt.Sprintf(" [%s]", downloadInfo.MediaInfoText)

	var imageData *dtypes.ImageData
	if downloadInfo.MediaInfo != nil && downloadInfo.MediaInfo.GetThumbnailID() != nil {
		thumbnail, _ := h.thumbnail.GetInfoByThumbID(ctx, *downloadInfo.MediaInfo.GetThumbnailID())
		if thumbnail != nil {
			imageData = thumbnail.ImageData()
		}
		if imageData != nil {
			imageData.URL = h.baseURL + uivalues.ThumbnailHttpPath(thumbnail.ThumbID.String())
		}
	}

	if imageData == nil && downloadInfo.ChannelID != nil && downloadInfo.IsYouTube() {
		channel, _ := h.downloader.FindYoutubeChannelInfo(ctx, *downloadInfo.ChannelID)
		if channel != nil && len(channel.ImageRaw) > 0 {
			imageData = channel.ImageData()
		}
		if imageData != nil {
			imageData.URL = h.baseURL + uivalues.YoutubeChannelHttpPath(channel.ChannelID)
		}
	}

	if imageData == nil || imageData.Width < 120 {
		imageData = &dtypes.ImageData{
			URL:    h.baseURL + uivalues.ImageHttpPath(uivalues.Elengrab1280ImageJpgFileName),
			Width:  1280,
			Height: 720,
		}
	}

	metaOgItems := make(uivalues.MetaOgItems, 0, 20)
	metaOgItems.Add("site_name", iconfig.AppName)
	metaOgItems.Add("type", fnx.Ternary(isVideoPlayer, "video.other", "music.song"))
	metaOgItems.Add("title", downloadInfo.MediaTitle)
	metaOgItems.Add("description", description)
	metaOgItems.Add("url", pageURL)
	metaOgItems.Add("image", imageData.URL)
	metaOgItems.Add("image:secure_url", imageData.URL)
	metaOgItems.Add("image:type", httpx.ContentTypeByExt(imageData.Format.String()))
	if imageData.Width != 0 && imageData.Height != 0 {
		metaOgItems.Add("image:width", strconv.Itoa(imageData.Width))
		metaOgItems.Add("image:height", strconv.Itoa(imageData.Height))
	}
	metaOgItems.Add("image:alt", "Elengrab logo")

	metaOgItems.Add(fmt.Sprintf("%s:url", prefixType), mediaURL)
	metaOgItems.Add(fmt.Sprintf("%s:secure_url", prefixType), mediaURL)
	if contentType != "" {
		metaOgItems.Add(fmt.Sprintf("%s:type", prefixType), contentType)
	}
	if isVideoPlayer && videoWidth != 0 {
		metaOgItems.Add("video:width", strconv.Itoa(videoWidth))
		metaOgItems.Add("video:height", strconv.Itoa(videoHeight))
	}

	metaNameItems := make(uivalues.MetaNameItems, 0, 4)
	metaNameItems.Add("twitter:card", "summary_large_image")
	metaNameItems.Add("twitter:title", downloadInfo.MediaTitle)
	metaNameItems.Add("twitter:description", description)
	metaNameItems.Add("twitter:image", imageData.URL)

	baseValues := uivalues.NewBaseValues()
	baseValues.Title = downloadInfo.MediaTitle
	baseValues.ShowFooter = false
	baseValues.MetaOgItems = metaOgItems
	baseValues.MetaNameItems = metaNameItems

	cssPaths, err := uivalues.CssWatchPaths(h.assetFolders.Css())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	jsScripts, err := uivalues.JsWatchPaths(h.assetFolders.Js())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	pwaManifestPath, err := uivalues.PwaManifestPath(h.assetFolders.Pwa())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	mediaParameters := make([]uivalues.MediaParameter, 0, 4)
	mediaParameters = append(mediaParameters, uivalues.MediaParameter{Name: "Format", Value: format})
	if videoQuality != "" {
		mediaParameters = append(mediaParameters, uivalues.MediaParameter{Name: "Video", Value: videoQuality})
	}
	if audioQuality != "" {
		mediaParameters = append(mediaParameters, uivalues.MediaParameter{Name: "Audio", Value: audioQuality})
	}
	mediaParameters = append(mediaParameters, uivalues.MediaParameter{Name: "File Size", Value: fileSize})
	mediaParameters = append(mediaParameters, uivalues.MediaParameter{Name: "Original Source", Value: originalSource})

	pageData := uivalues.WatchPageData{
		BaseValues: baseValues,
		BasePaths:  uivalues.NewBasePaths(),
		Paths: uivalues.PagePaths{
			Css:         cssPaths,
			JsScripts:   jsScripts,
			PwaManifest: pwaManifestPath,
			Stream:      streamPath,
		},
		Values: uivalues.WatchPageValues{
			ShowBackButton:   showBackButton,
			IsVideoPlayer:    isVideoPlayer,
			MediaTitle:       downloadInfo.MediaTitle,
			MediaDescription: downloadInfo.MediaDescription,
			MediaParameters:  mediaParameters,
			ContentType:      contentType,
		},
	}

	// Set content type so browser renders HTML properly
	ctx.SetContentType("text/html; charset=utf-8")

	// Load template
	tmpl, err := h.loadPage(uivalues.PageWatch.FileName())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}

	// Execute template with PageTitle
	if err := tmpl.ExecuteTemplate(ctx, uivalues.PageWatch.Key(), pageData); err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}
}
