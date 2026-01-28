package downloaderh

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) DeleteFileRowHanedler(ctx *fasthttp.RequestCtx) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		nfasthttp.WriteError(ctx, fmt.Errorf("authorization error: %v", err), fasthttp.StatusUnauthorized)
		return
	}

	fileIdStr := ctx.UserValue(fileIdKey).(string)
	if fileIdStr == "" {
		nfasthttp.WriteError(ctx, fmt.Errorf("fileId is required"), fasthttp.StatusBadRequest)
		return
	}

	fileId, err := uuid.Parse(fileIdStr)
	if err != nil {
		nfasthttp.WriteError(ctx, fmt.Errorf("fileId is incorrect: %v", err), fasthttp.StatusBadRequest)
		return
	}

	err = h.usecases.Downloader.DeleteDownload(ctx, userID, fileId)
	if err != nil {
		nfasthttp.WriteError(ctx, err, fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}
