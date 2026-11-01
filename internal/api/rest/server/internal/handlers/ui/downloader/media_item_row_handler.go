package downloader

import (
	"bytes"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) MediaItemRowHandler(ctx *fasthttp.RequestCtx) {
	authCtx := policy.ResolveUserOrAnonym(ctx)

	downloadID, err := h.extractDownloadID(ctx)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	downloadInfo, err := h.downloader.GetDownloadInfo(ctx, authCtx, downloadID)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	if ctx.IsHead() {
		fasthttpx.WriteHTMLHeadResponse(ctx)
		return
	}

	row := h.renderMediaItemRow(
		ctx,
		renderMediaItemRowParams{
			downloadInfo: downloadInfo,
			AuthCtx:      authCtx,
		},
	)
	if row.err != nil {
		fasthttpx.WriteErrorx(ctx, row.err)
		return
	}

	if row.httpStatus == fasthttp.StatusNoContent {
		ctx.SetStatusCode(row.httpStatus)
		ctx.Response.Header.Set("HX-Trigger", "no-op")
		return
	}

	var buf bytes.Buffer
	err = h.templates.Base.ExecuteTemplate(&buf, components.ResultRowStatusKey, row.data)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, errInternal(err))
		return
	}

	ctx.SetBody(buf.Bytes())
	ctx.SetStatusCode(fasthttp.StatusOK)
}
