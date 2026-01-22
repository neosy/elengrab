package downloaderh

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

var (
	mapContentTypeByExt = map[string]string{
		"mp3":  "audio/mpeg",
		"m4a":  "audio/mp4",
		"aac":  "audio/aac",
		"ogg":  "audio/ogg",
		"opus": "audio/opus",
		"flac": "audio/flac",

		"mp4":  "video/mp4",
		"webm": "video/webm",
		"mkv":  "video/x-matroska",
		"mov":  "video/quicktime",
	}
)

func (h *DownloaderHandlers) DownloadHandler(ctx *fasthttp.RequestCtx) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		nfasthttp.WriteError(ctx, fmt.Errorf("authorization error: %v", err), fasthttp.StatusUnauthorized)
		return
	}

	// Get the file name from the query parameter
	fileIdStr := string(ctx.QueryArgs().Peek("file"))
	if fileIdStr == "" {
		nfasthttp.WriteError(ctx, errors.New("file is required"), fasthttp.StatusBadRequest)
		return
	}

	fileId := uuid.MustParse(fileIdStr)

	// Build the full path to the file
	fileInfo, err := h.usecases.Downloader.GetFileInfo(ctx, userID, fileId)
	if err != nil {
		nfasthttp.WriteError(ctx, err, fasthttp.StatusNotFound)
		return
	}

	var filePath string
	if fileInfo.FileFullName != "" {
		filePath = filepath.Join(h.downloadsDir, fileInfo.FileFullName)
	}

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		nfasthttp.WriteError(ctx, errors.New("file not found"), fasthttp.StatusBadRequest)
		return
	}

	// Detect content type by extension
	contentType := mapContentTypeByExt[fileInfo.FileExt]
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Send file via streaming
	nfasthttp.SendFileDirect(ctx, filePath, fileInfo.SafeReadableFullName, contentType)
}
