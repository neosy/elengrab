package handlers

import (
	"github.com/google/uuid"
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) FileWatchHandler(ctx *fasthttp.RequestCtx) {
	fileIdStr, ok := ctx.UserValue(fileIdKey).(string)
	if !ok || fileIdStr == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrFileIdIsRequired)
		return
	}

	fileID, err := uuid.Parse(fileIdStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrFileIdIsIncorrect.Wrap(err))
		return
	}

	h.watch(
		ctx,
		httppaths.BuildPathFileWatch(fileID),
		httppaths.BuildPathFileStream(fileID),
		fileID,
		true,
	)
}
