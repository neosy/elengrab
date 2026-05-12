package handlers

import (
	"bytes"
	"encoding/json"

	"github.com/google/uuid"
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/dto"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/policy"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) RepeatDownloadHandler(ctx *fasthttp.RequestCtx) {
	ctxUser, err := policy.ResolveUserOrFallback(ctx, h.appMode)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	downloadIDStr, ok := ctx.UserValue(downloadIDKey).(string)
	if !ok || downloadIDStr == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrDownloadIDIsRequired)
		return
	}

	downloadID, err := uuid.Parse(downloadIDStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrDownloadIDIsIncorrect.Wrap(err))
		return
	}

	downloadInfo, err := h.downloader.RepeatDownload(ctx, *ctxUser, downloadID)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	row := h.genRow(downloadInfo, false)
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
	err = tmpl.ExecuteTemplate(&buf, uivalues.ComponentResultRowStatusKey, row.data)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errInternal(err))
		return
	}

	dataResp := dto.RepeatDownloadResponse{
		GuestCreated: ctxUser.GuestCreated,
	}

	dataJSON, _ := json.Marshal(dataResp)
	if dataJSON != nil {
		ctx.Response.Header.Set("HX-Trigger", string(dataJSON))
	}

	ctx.SetBody(buf.Bytes())
	ctx.SetStatusCode(fasthttp.StatusOK)
}
