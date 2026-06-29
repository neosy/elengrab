package downloader

import (
	"bytes"
	"encoding/json"

	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) RetryImportMediaHandler(ctx *fasthttp.RequestCtx) {
	authCtx, err := policy.ResolveUserOrFallback(ctx, h.appMode)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

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

	downloadInfo, err := h.downloader.RetryDownload(ctx, authCtx, downloadID)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	row := h.renderMediaItemRow(downloadInfo, false)
	if row.err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}
	if row.httpStatus == fasthttp.StatusNoContent {
		ctx.SetStatusCode(row.httpStatus)
		ctx.Response.Header.Set("HX-Trigger", "no-op")
		return
	}

	tmpl, err := h.templates.Clone()
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, components.ResultRowStatusKey, row.data)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}

	dataResp := dto.RepeatDownloadResponse{
		GuestCreated: authCtx.GuestCreated,
	}

	dataJSON, _ := json.Marshal(dataResp)
	if dataJSON != nil {
		ctx.Response.Header.Set("HX-Trigger", string(dataJSON))
	}

	ctx.SetBody(buf.Bytes())
	ctx.SetStatusCode(fasthttp.StatusOK)
}
