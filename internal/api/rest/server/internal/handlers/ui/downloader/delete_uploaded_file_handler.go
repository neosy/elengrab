package downloader

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) DeleteUploadedFileHandler(ctx *fasthttp.RequestCtx) {
	authCtx, err := policy.ResolveUserOrFallback(ctx, h.appMode)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	if ctx.IsHead() {
		fasthttpx.WriteHTMLHeadResponse(ctx)
		return
	}

	fileID, err := h.extractFileID(ctx)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	upload, _, err := h.uploadedFileCacheRepo.Find(ctx, fileID)
	if err != nil {
		h.logger.Error("Failed to find the uploaded file", "fileID", fileID, "error", err)
		nfasthttp.WriteErrorx(ctx, err)
		return
	}
	if upload == nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTPMessage("Uploaded file not found", fasthttp.StatusNotFound))
		return
	}

	if upload.OwnerID() != authCtx.UserID {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTPMessage("You don't have permission", fasthttp.StatusForbidden))
		return
	}

	err = upload.Cleanup()
	if err != nil {
		// Cleanup has already logged the failure. Return before touching the cache:
		// the file is still in the temporary directory, so the entry must keep tracking it.
		nfasthttp.WriteErrorx(ctx, errorx.Errorf("failed to remove the uploaded file: %w", err,
			errorx.WithErrorMessage("Failed to delete the file."),
		))
		return
	}

	err = h.uploadedFileCacheRepo.Delete(ctx, fileID)
	if err != nil {
		h.logger.Error("Failed to delete the uploaded file", "fileID", fileID, "error", err)
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusNoContent)
}
