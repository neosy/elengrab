package downloader

import (
	"bytes"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/components"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) renderMediaItemRowPlaceholder(
	ctx *fasthttp.RequestCtx,
	downloadInfo *ucdto.MediaDownloadInfo,
) (*bytes.Buffer, error) {
	row := h.renderMediaItemRow(ctx, renderMediaItemRowParams{
		downloadInfo:   downloadInfo,
		lazyLoadImages: false,
		ResultRowFade:  "fade-in",
	})

	if row.err != nil {
		nfasthttp.WriteErrorx(ctx, row.err)
		return nil, row.err
	}
	if row.httpStatus == fasthttp.StatusNoContent {
		ctx.SetStatusCode(row.httpStatus)
		ctx.Response.Header.Set("HX-Trigger", "no-op")
		return nil, nil
	}

	var buf bytes.Buffer
	err := h.templates.Base.ExecuteTemplate(&buf, components.ResultNewRowKey, row.data)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return nil, err
	}

	return &buf, err
}
