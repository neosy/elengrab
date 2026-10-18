package downloader

import (
	"encoding/json"
	"net/http"

	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	"github.com/neosy/elengrab/internal/exceptions"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) QueueUploadedFileHandler(ctx *fasthttp.RequestCtx) {
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

	upload, _, _ := h.uploadedFileCacheRepo.Find(ctx, fileID)
	if upload == nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTPMessage("Uploaded file not found", fasthttp.StatusNotFound))
		return
	}

	if upload.OwnerID() != authCtx.UserID {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTPMessage("You don't have permission", fasthttp.StatusForbidden))
		return
	}

	var receivedReq dto.QueueUploadedFileRequest
	if err := json.Unmarshal(ctx.PostBody(), &receivedReq); err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTPMessage("Invalid JSON request body", http.StatusBadRequest, err))
		return
	}

	scheduleUploadRequest, err := h.mappers.MapQueueUploadedFileRequestToScheduleUpload(receivedReq)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewFromError(err, exceptions.INVALID_REQUEST))
		return
	}

	uploadResponse, err := h.downloader.ScheduleUpload(ctx, authCtx, upload, scheduleUploadRequest)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	if uploadResponse == nil {
		h.logger.Error("Upload scheduling returned an empty response", "fileId", fileID)
		fasthttpx.WriteErrorx(ctx, apierrors.ErrEmptyResponse)
		return
	}

	err = h.uploadedFileCacheRepo.Forget(ctx, fileID)
	if err != nil {
		h.logger.Error("Failed to remove uploaded file from cache", "fileID", fileID, "error", err)
	}

	dataResponse := dto.QueueUploadedFileResponse{
		DownloadID: idcodec.EncodeUUIDBase64URL(uploadResponse.DownloadID),
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	fasthttpx.WriteResponse(ctx, dataResponse)
}
