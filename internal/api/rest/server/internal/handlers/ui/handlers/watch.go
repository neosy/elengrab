package handlers

import (
	"fmt"
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
	fileID uuid.UUID,
	showBackButton bool,
) {
	if ctx.IsHead() {
		ctx.SetContentType("text/html; charset=utf-8")
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}

	fileInfo, err := h.downloader.GetFileInfoUnrestricted(ctx, fileID)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	var (
		format        = "-"
		videoQuality  = ""
		audioQuality  = ""
		fileSize      = "-"
		isVideoPlayer = true
		contentType   = ""
		videoWidth    = 0
		videoHeight   = 0
	)
	if fileInfo.MediaInfo != nil {
		ext := fileInfo.MediaInfo.Format.Ext()

		isVideoPlayer = fileInfo.MediaInfo.Format.IsVideo()
		format = strings.ToUpper(ext)
		contentType = httpx.ContentTypeByExt(ext)

		if videoInfo := fileInfo.MediaInfo.VideoInfo; videoInfo != nil {
			qValues := []string{strings.ToUpper(videoInfo.Codec.String())}
			if videoInfo.Resolution != dtypes.VideoResolutionNone {
				qValues = append(qValues, videoInfo.Resolution.String())
			}
			qValues = append(qValues, videoInfo.ResolutionString())
			videoQuality = strings.Join(qValues, " • ")
			videoWidth = videoInfo.Width
			videoHeight = videoInfo.Height
		}
		if fileInfo.MediaInfo.AudioInfo != nil {
			parts := make([]string, 0, 2)
			parts = append(parts, fmt.Sprintf("%v kbps ", fileInfo.MediaInfo.AudioInfo.Bitrate))
			if fileInfo.MediaInfo.AudioInfo.SampleRate != nil {
				parts = append(parts, fmt.Sprintf("%v Hz ", *fileInfo.MediaInfo.AudioInfo.SampleRate))
			}
			audioQuality = strings.Join(parts, " • ")
		}
	}

	if fileInfo.FileSize != nil {
		fileSize = uformat.BytesHuman(*fileInfo.FileSize)
	}

	const disableMediaType = true
	mediaURL := h.baseURL + streamPath
	prefixType := fnx.Ternary(isVideoPlayer, "video", "audio")
	description := fileInfo.MediaTitle + fmt.Sprintf(" [%s]", fileInfo.MediaInfoText)

	var imageData *dtypes.ImageData
	if fileInfo.MediaInfo != nil && fileInfo.MediaInfo.GetThumbnailID() != nil {
		thumbnail, _ := h.thumbnail.GetInfoByThumbID(ctx, *fileInfo.MediaInfo.GetThumbnailID())
		if thumbnail != nil {
			imageData = thumbnail.ImageData()
		}
		if imageData != nil {
			imageData.URL = h.baseURL + uivalues.ThumbnailHttpPath(thumbnail.ThumbID.String())
		}
	}

	if imageData == nil && fileInfo.YoutubeChannelID != nil {
		channel, _ := h.downloader.FindYoutubeChannelInfo(ctx, *fileInfo.YoutubeChannelID)
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
	if disableMediaType {
		metaOgItems.Add("type", "website")
	} else {
		metaOgItems.Add("type", fnx.Ternary(isVideoPlayer, "video.other", "music.song"))
	}
	metaOgItems.Add("title", fileInfo.MediaTitle)
	metaOgItems.Add("description", description)
	metaOgItems.Add("url", pageURL)
	metaOgItems.Add("image", imageData.URL)
	metaOgItems.Add("image:secure_url", imageData.URL)
	metaOgItems.Add("image:type", httpx.ContentTypeFromPath(imageData.URL))
	if imageData.Width != 0 && imageData.Height != 0 {
		metaOgItems.Add("image:width", strconv.Itoa(imageData.Width))
		metaOgItems.Add("image:height", strconv.Itoa(imageData.Height))
	}
	metaOgItems.Add("image:alt", "Elengrab logo")

	if !disableMediaType {
		metaOgItems.Add(prefixType, mediaURL)
		metaOgItems.Add(fmt.Sprintf("%s:secure_url", prefixType), mediaURL)
		if contentType != "" {
			metaOgItems.Add(fmt.Sprintf("%s:type", prefixType), contentType)
		}
		if isVideoPlayer && videoWidth != 0 {
			metaOgItems.Add("video:width", strconv.Itoa(videoWidth))
			metaOgItems.Add("video:height", strconv.Itoa(videoHeight))
		}
	}

	metaNameItems := make(uivalues.MetaNameItems, 0, 4)
	metaNameItems.Add("twitter:card", "summary_large_image")
	metaNameItems.Add("twitter:title", fileInfo.MediaTitle)
	metaNameItems.Add("twitter:description", description)
	metaNameItems.Add("twitter:image", imageData.URL)

	baseValues := uivalues.NewBaseValues()
	baseValues.Title = fileInfo.MediaTitle
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

	type mediaParameter struct {
		Name  string
		Value string
	}

	mediaParametes := make([]mediaParameter, 0, 4)
	mediaParametes = append(mediaParametes, mediaParameter{"Format", format})
	if videoQuality != "" {
		mediaParametes = append(mediaParametes, mediaParameter{"Video", videoQuality})
	}
	if audioQuality != "" {
		mediaParametes = append(mediaParametes, mediaParameter{"Audio", audioQuality})
	}
	mediaParametes = append(mediaParametes, mediaParameter{"File Size", fileSize})

	watcherValues := uivalues.WatcherValues{
		ShowBackButton: showBackButton,
		IsVideoPlayer:  isVideoPlayer,
		MediaTitle:     fileInfo.MediaTitle,
		MediaParametes: mediaParametes,
		ContentType:    contentType,
	}

	dataMap := uivalues.MergeMaps(
		uivalues.PathValues,
		baseValues.ToMap(),
		watcherValues.ToMap(),
	)
	dataMap[uivalues.CssPathsKey] = cssPaths
	dataMap[uivalues.JsScriptsKey] = jsScripts
	dataMap[uivalues.PathStreamKey] = streamPath

	// Set content type so browser renders HTML properly
	ctx.SetContentType("text/html; charset=utf-8")

	// Load template
	tmpl, err := h.loadPage(uivalues.PageWatch.FileName())
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}

	// Execute template with PageTitle
	if err := tmpl.ExecuteTemplate(ctx, uivalues.PageWatch.Key(), dataMap); err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}
}
