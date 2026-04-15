package handlers

import (
	"os"
	"path/filepath"

	"github.com/google/uuid"
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/policy"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) FileStreamHandler(ctx *fasthttp.RequestCtx) {
	// Get user ID from context
	ctxUser := policy.ResolveUserOrAnonym(ctx)

	fileIdStr := ctx.UserValue(fileIdKey).(string)
	if fileIdStr == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrFileIdIsRequired)
		return
	}

	fileID, err := uuid.Parse(fileIdStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrFileIdIsIncorrect.Wrap(err))
		return
	}

	h.stream(ctx, *ctxUser, fileID, false)
}

func (h *DownloaderHandlers) StreamShortCodeHandler(ctx *fasthttp.RequestCtx) {
	shortCode := ctx.UserValue(shortCodeKey).(string)
	if shortCode == "" {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTPMessage("shortCode is required", fasthttp.StatusBadRequest))
		return
	}

	link, err := h.linkWeb.GetLastByShortCode(ctx, shortCode)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	if link == nil {
		nfasthttp.WriteErrorx(
			ctx,
			errorx.New(
				"link not found",
				exceptionx.NOT_FOUND,
				exceptionx.NOT_FOUND.ErrorMessage(),
			))
		return
	}

	fileID := stripUUIDFromPath(link.OriginalURL)
	if fileID == uuid.Nil {
		nfasthttp.WriteErrorx(
			ctx,
			errorx.New(
				"fileId is incorrect",
				exceptionx.WRONG_DATA,
				exceptionx.WRONG_DATA.ErrorMessage(),
			))
	}

	h.stream(ctx, dauth.UserContext{}, fileID, true)
}

func (h *DownloaderHandlers) stream(
	ctx *fasthttp.RequestCtx,
	ctxUser dauth.UserContext,
	fileID uuid.UUID,
	unrestricted bool,
) {
	var (
		fileInfo *dto.GetFileInfoResponse
		err      error
	)
	// Retrieve file info
	if unrestricted {
		fileInfo, err = h.downloader.GetFileInfoUnrestricted(ctx, fileID)
	} else {
		fileInfo, err = h.downloader.GetFileInfo(ctx, ctxUser, fileID)
	}

	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	// Build the full path to the file
	var filePath string
	if fileInfo.FileFullName != "" {
		filePath = filepath.Join(h.downloadsDir, fileInfo.FileFullName)
	}

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		nfasthttp.WriteErrorx(
			ctx,
			apierrors.ErrFileNotFound,
		)
		return
	}

	// Detect content type by extension
	contentType := httpx.ContentTypeByExt(fileInfo.FileExt)

	// Set headers for streaming in browser
	ctx.Response.Header.Set("Content-Type", contentType)
	ctx.Response.Header.Set("Content-Disposition", "inline") // Play in browser
	ctx.Response.Header.Set("Accept-Ranges", "bytes")        // Allow seeking
	ctx.Response.Header.Set("Cache-Control", "public, max-age=3600")

	// Stream the file directly (memory-efficient, supports Range)
	nfasthttp.SendFileDirect(ctx, filePath, fileInfo.SafeReadableFullName, contentType)
}
