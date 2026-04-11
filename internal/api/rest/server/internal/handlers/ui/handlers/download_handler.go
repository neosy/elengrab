package handlers

import (
	"os"
	"path/filepath"

	"github.com/google/uuid"
	authmw "github.com/neosy/elengrab/internal/api/rest/server/internal/auth_middleware"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/neosy/elengrab/internal/pkg/httpx"
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
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTPMessage("file is required", fasthttp.StatusBadRequest))
		return
	}

	fileId, err := uuid.Parse(fileIdStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTPMessage("file is incorrect", fasthttp.StatusBadRequest))
		return
	}

	// Build the full path to the file
	fileInfo, err := h.downloader.GetFileInfo(ctx, *ctxUser, fileId)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	var filePath string
	if fileInfo.FileFullName != "" {
		filePath = filepath.Join(h.downloadsDir, fileInfo.FileFullName)
	}

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTP("file not found", fasthttp.StatusInternalServerError, err))
		return
	}

	// Detect content type by extension
	contentType := httpx.ContentTypeByExt(fileInfo.FileExt)

	// Send file via streaming
	nfasthttp.SendFileDirect(ctx, filePath, fileInfo.SafeReadableFullName, contentType)
}
