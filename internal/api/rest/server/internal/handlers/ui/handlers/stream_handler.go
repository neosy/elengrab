package handlers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) StreamHandler(ctx *fasthttp.RequestCtx) {
	// Get user ID from context
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		nfasthttp.WriteError(ctx, fmt.Errorf("authorization error: %v", err), fasthttp.StatusUnauthorized)
		return
	}

	// Get the file ID from query parameter
	fileIdStr := string(ctx.QueryArgs().Peek("file"))
	if fileIdStr == "" {
		nfasthttp.WriteError(ctx, errors.New("file is required"), fasthttp.StatusBadRequest)
		return
	}

	fileId := uuid.MustParse(fileIdStr)

	// Retrieve file info
	fileInfo, err := h.usecases.Downloader.GetFileInfo(ctx, userID, fileId)
	if err != nil {
		nfasthttp.WriteError(ctx, err, fasthttp.StatusNotFound)
		return
	}

	// Build the full path to the file
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

	// Set headers for streaming in browser
	ctx.Response.Header.Set("Content-Type", contentType)
	ctx.Response.Header.Set("Content-Disposition", "inline") // Play in browser
	ctx.Response.Header.Set("Accept-Ranges", "bytes")        // Allow seeking
	ctx.Response.Header.Set("Cache-Control", "public, max-age=3600")

	// Stream the file directly (memory-efficient, supports Range)
	nfasthttp.SendFileDirect(ctx, filePath, fileInfo.SafeReadableFullName, contentType)
}
