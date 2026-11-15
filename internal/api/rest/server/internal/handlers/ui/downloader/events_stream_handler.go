package downloader

import (
	"bufio"
	"bytes"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) EventsStreamHandler(ctx *fasthttp.RequestCtx) {
	authCtx := policy.ResolveUserOrAnonym(ctx)

	// SSE headers
	h.sse.SetupHeaders(ctx)

	ctx.SetBodyStreamWriter(func(w *bufio.Writer) {
		h.sse.StreamEvents(ctx, w, authCtx)
	})
}

// renderMediaItemRowForSSE renders a media item row for an SSE event.
func (h *DownloaderHandlers) renderMediaItemRowForSSE(
	ctx *fasthttp.RequestCtx,
	authCtx dauth.AuthContext,
	downloadInfo *ucdto.MediaDownloadInfo,
) ([]byte, error) {
	row := h.renderMediaItemRow(
		ctx,
		renderMediaItemRowParams{
			downloadInfo:    downloadInfo,
			isDownloadEvent: true,
			AuthCtx:         authCtx,
		},
	)
	if row.err != nil {
		return nil, row.err
	}
	if row.httpStatus == fasthttp.StatusNoContent {
		return nil, nil
	}

	var buf bytes.Buffer
	err := h.templates.Base.ExecuteTemplate(&buf, components.ResultRowStatusKey, row.data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
