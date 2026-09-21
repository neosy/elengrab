package downloader

import (
	"bytes"
	"encoding/json"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/consts"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) SearchHandler(ctx *fasthttp.RequestCtx) {
	authCtx := policy.ResolveUserOrAnonym(ctx)

	query := udto.MediaDownloadQueryDefault(consts.LoadHistoryLimit)

	if string(ctx.Request.Header.ContentType()) == "application/json" {
		var req dto.SearchRequest

		err := json.Unmarshal(ctx.PostBody(), &req)
		if err != nil {
			fasthttpx.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
			return
		}

		err = h.validators.Validate.Struct(req)
		if err != nil {
			nfasthttp.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
			return
		}

		query, err = h.mappers.MapSearchRequestToUsecaseQuery(req)
		if err != nil {
			nfasthttp.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
			return
		}
	} else {
		params, err := parsePostSearchParameters(ctx)
		if err != nil {
			nfasthttp.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
			return
		}

		paramValues, err := params.ParseValues()
		if err != nil {
			nfasthttp.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
			return
		}

		query, err = h.mappers.MapSearchParameterValuesToUsecaseQuery(&paramValues)
		if err != nil {
			nfasthttp.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
			return
		}
	}

	var bodyBuffer bytes.Buffer
	err := h.renderDownloadRows(ctx, &bodyBuffer, authCtx, query)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(bodyBuffer.Bytes())
}
