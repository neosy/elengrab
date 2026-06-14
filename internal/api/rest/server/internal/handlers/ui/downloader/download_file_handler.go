package downloader

import (
	"os"

	"github.com/google/uuid"
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) DownloadFileHandler(ctx *fasthttp.RequestCtx) {
	ctxUser := policy.ResolveUserOrAnonym(ctx)

	// Get the file name from the query parameter
	downloadIDStr := string(ctx.QueryArgs().Peek(downloadIDKey))
	if downloadIDStr == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrDownloadIDIsRequired)
		return
	}

	downloadID, err := uuid.Parse(downloadIDStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrDownloadIDIsIncorrect.Wrap(err))
		return
	}

	// Build the full path to the download
	downloadInfo, err := h.downloader.GetDownloadInfo(ctx, ctxUser, downloadID)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	var filePath string
	if downloadInfo.FileFullName != "" {
		filePath = h.downloadsStorage.Path(downloadInfo.FileFullName)
	}

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrFileNotFound.Wrap(err))
		return
	}

	// Detect content type by extension
	contentType := httpx.ContentTypeByExt(downloadInfo.FileExt)

	// Send file via streaming
	nfasthttp.SendFileDirect(ctx, filePath, downloadInfo.SafeReadableFullName, contentType)
}
