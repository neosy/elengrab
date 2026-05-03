package handlers

import (
	"os"

	"github.com/google/uuid"
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/policy"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) DownloadHandler(ctx *fasthttp.RequestCtx) {
	ctxUser := policy.ResolveUserOrAnonym(ctx)

	// Get the file name from the query parameter
	fileIdStr := string(ctx.QueryArgs().Peek("file"))
	if fileIdStr == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrFileIdIsRequired)
		return
	}

	fileId, err := uuid.Parse(fileIdStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrFileIdIsIncorrect.Wrap(err))
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
		filePath = h.downloadsStorage.Path(fileInfo.FileFullName)
	}

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrFileNotFound.Wrap(err))
		return
	}

	// Detect content type by extension
	contentType := httpx.ContentTypeByExt(fileInfo.FileExt)

	// Send file via streaming
	nfasthttp.SendFileDirect(ctx, filePath, fileInfo.SafeReadableFullName, contentType)
}
