package grabh

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/neosy/elengrab/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

var (
	mapContentTypeByExt = map[string]string{
		".mp3":  "audio/mpeg",
		".m4a":  "audio/mp4",
		".aac":  "audio/aac",
		".ogg":  "audio/ogg",
		".opus": "audio/opus",

		".mp4":  "video/mp4",
		".webm": "video/webm",
		".mkv":  "video/x-matroska",
		".mov":  "video/quicktime",
	}
)

func (h *GrabHandlers) DownloadHandler(ctx *fasthttp.RequestCtx) {
	// Get the file name from the query parameter
	fileName := string(ctx.QueryArgs().Peek("file"))
	if fileName == "" {
		nfasthttp.WriteError(ctx, errors.New("file is required"), fasthttp.StatusBadRequest)
		return
	}

	// Build the full path to the file
	filePath := h.usecases.Downloader.GetFilePath(fileName)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		nfasthttp.WriteError(ctx, errors.New("file not found"), fasthttp.StatusBadRequest)
		return
	}

	// Detect content type by extension
	contentType := mapContentTypeByExt[filepath.Ext(fileName)]
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Set headers for file download
	ctx.SetContentType(contentType)
	ctx.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))

	// Serve the file
	// fasthttp.ServeFile(ctx, filePath)
	fasthttp.ServeFileUncompressed(ctx, filePath)
}
