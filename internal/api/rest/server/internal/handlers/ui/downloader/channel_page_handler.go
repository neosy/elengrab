package downloader

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) ChannelPageHandler(ctx *fasthttp.RequestCtx) {
	authCtx := policy.ResolveUserOrAnonym(ctx)

	channelID, err := h.extractChannelID(ctx)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	_, err = h.downloader.GetChannelByID(ctx, channelID)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	searchValues, err := h.parseSearchGetRequest(ctx)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
		return
	}

	if searchValues.Parameters.Filters == nil {
		searchValues.Parameters.Filters = types.NewQueryFilters()
	}

	searchValues.Parameters.Filters.Add(qkeys.ChannelIDKey, idcodec.EncodeUUIDBase64URL(channelID))

	query, err := h.mappers.MapSearchValuesToUsecaseQuery(searchValues)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
		return
	}

	h.renderChannelPage(ctx, authCtx, query)
}
