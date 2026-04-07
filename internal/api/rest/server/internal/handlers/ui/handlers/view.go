package handlers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	"github.com/neosy/elengrab/internal/pkg/nfasthttp"
	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) view(ctx *fasthttp.RequestCtx, streamPath string, fileID uuid.UUID) {
	fileInfo, err := h.downloader.GetFileInfoUnrestricted(ctx, fileID)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	baseValues := uivalues.BaseValues.Copy()
	baseValues.ShowHeader = true
	baseValues.ShowFooter = false

	cssPaths, err := uivalues.CssViewPaths(filepath.Join(h.assetsDir, dirStaticName, dirCssName))
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	jsScripts, err := uivalues.JsViewPaths(filepath.Join(h.assetsDir, dirStaticName, dirJsName))
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
	)
	if fileInfo.MediaInfo != nil {
		ext := fileInfo.MediaInfo.Format.Ext()

		isVideoPlayer = fileInfo.MediaInfo.Format.IsVideo()
		format = strings.ToUpper(ext)
		contentType = mapContentTypeByExt[ext]

		if fileInfo.MediaInfo.VideoInfo != nil {
			videoQuality = fmt.Sprintf(
				"%v • %v",
				strings.ToUpper(fileInfo.MediaInfo.VideoInfo.Codec.String()),
				fileInfo.MediaInfo.VideoInfo.Resolution,
			)
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

	viewerValues := uivalues.ViewerValues{
		ShowBackButton: false,
		IsVideoPlayer:  isVideoPlayer,
		MediaTitle:     fileInfo.MediaTitle,
		MediaParametes: mediaParametes,
		ContentType:    contentType,
	}

	dataMap := uivalues.MergeMaps(
		uivalues.PathValues,
		baseValues.ToMap(),
		viewerValues.ToMap(),
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
