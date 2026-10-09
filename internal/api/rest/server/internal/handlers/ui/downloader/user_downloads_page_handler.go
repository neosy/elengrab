package downloader

import (
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) UserDownloadsPageHandler(ctx *fasthttp.RequestCtx) {
	authCtx := policy.ResolveUserOrAnonym(ctx)

	login, ok := ctx.UserValue(qkeys.UserLoginKey.String()).(string)
	if !ok || login == "" {
		fasthttpx.WriteErrorx(ctx, apierrors.ErrUserNameIsRequired)
		return
	}

	user, err := h.authWeb.GetByLogin(ctx, login)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, apierrors.ErrUserNameIsIncorrect.Wrap(err))
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

	searchValues.Parameters.Filters.Add(qkeys.UserIDKey, idcodec.EncodeUUIDBase64URL(user.UserID))

	query, err := h.mappers.MapSearchValuesToUsecaseQuery(searchValues)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
		return
	}

	h.renderUserDownloadsPage(ctx, authCtx, *user, query)
}
