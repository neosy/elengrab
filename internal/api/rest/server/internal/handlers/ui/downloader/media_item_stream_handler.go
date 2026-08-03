package downloader

import (
	"os"

	"github.com/google/uuid"
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) MediaItemStreamHandler(ctx *fasthttp.RequestCtx) {
	// Get user ID from context
	authCtx := policy.ResolveUserOrAnonym(ctx)

	downloadIDStr, ok := ctx.UserValue(downloadIDKey).(string)
	if !ok || downloadIDStr == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrDownloadIDIsRequired)
		return
	}

	downloadID, err := idcodec.DecodeUUIDBase64URL(downloadIDStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrDownloadIDIsIncorrect.Wrap(err))
		return
	}

	h.stream(ctx, authCtx, downloadID, false)
}

func (h *DownloaderHandlers) StreamShortCodeHandler(ctx *fasthttp.RequestCtx) {
	shortCode, ok := ctx.UserValue(shortCodeKey).(string)
	if !ok || shortCode == "" {
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

	downloadID := stripUUIDFromIDPath(link.OriginalURL)
	if downloadID == uuid.Nil {
		nfasthttp.WriteErrorx(
			ctx,
			errorx.New(
				"downloadID is incorrect",
				exceptionx.WRONG_DATA,
				exceptionx.WRONG_DATA.ErrorMessage(),
			))
	}

	h.stream(ctx, dauth.AuthContext{}, downloadID, true)
}

func (h *DownloaderHandlers) stream(
	ctx *fasthttp.RequestCtx,
	authCtx dauth.AuthContext,
	downloadID uuid.UUID,
	unrestricted bool,
) {
	var (
		downloadInfo *dto.MediaDownloadInfo
		err          error
	)
	// Retrieve file info
	if unrestricted {
		downloadInfo, err = h.downloader.GetDownloadInfoUnrestricted(ctx, downloadID)
	} else {
		downloadInfo, err = h.downloader.GetDownloadInfo(ctx, authCtx, downloadID)
	}

	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	// Build the full path to the file
	var filePath string
	if downloadInfo.FileFullName != "" {
		filePath = h.downloadsStorage.Path(downloadInfo.FileFullName)
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
	contentType := httpx.ContentTypeByExt(downloadInfo.FileExt, httpx.ContentTypeOptionWithIsAudio(downloadInfo.IsAudioOnly()))

	// Set headers for streaming in browser
	ctx.Response.Header.Set("Content-Type", contentType)
	ctx.Response.Header.Set("Content-Disposition", "inline") // Play in browser
	ctx.Response.Header.Set("Accept-Ranges", "bytes")        // Allow seeking
	ctx.Response.Header.Set("Cache-Control", "public, max-age=3600")

	// Stream the file directly (memory-efficient, supports Range)
	nfasthttp.StreamFileDirect(ctx, filePath, contentType)
}
