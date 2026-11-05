package downloader

import (
	"encoding/json"
	"net/http"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) ExtractByURLHandler(ctx *fasthttp.RequestCtx) {
	authCtx, err := policy.ResolveUserOrFallback(ctx, h.appMode)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	var receivedReq dto.ExtractByURLRequest
	if err := json.Unmarshal(ctx.PostBody(), &receivedReq); err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTPMessage("Invalid JSON request body", http.StatusBadRequest, err))
		return
	}

	extractResponse, err := h.downloader.ExtractByURL(ctx, authCtx, receivedReq.URL)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	dataResponse, err := h.mappers.MapExtractByURLResponseDomainToDto(extractResponse)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	fasthttpx.WriteResponse(ctx, dataResponse)
}
