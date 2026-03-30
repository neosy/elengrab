package handlers

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	authmw "github.com/neosy/elengrab/internal/api/rest/server/internal/auth_middleware"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) DownloadHandler(ctx *fasthttp.RequestCtx) {
	ctxUser := authmw.UserFromContext(ctx)
	if ctxUser == nil {
		// anonymous
		ctxUser = dauth.UserContextAnonymous()
	}

	// Get the file name from the query parameter
	fileIdStr := string(ctx.QueryArgs().Peek("file"))
	if fileIdStr == "" {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTP("file is required", fasthttp.StatusBadRequest))
		return
	}

	fileId := uuid.MustParse(fileIdStr)

	// Build the full path to the file
	fileInfo, err := h.downloader.GetFileInfo(ctx, *ctxUser, fileId)
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
