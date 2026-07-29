package admin

import (
	"net/http"
	"net/url"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *AdminHandlers) redirectAuth(ctx *fasthttp.RequestCtx) bool {
	ctxUser := policy.ResolveUser(ctx)
	if ctxUser == nil || ctxUser.UserType() < dtypes.UserTypeUser {
		loginPath := httppaths.GroupAccount + httppaths.PathLogin

		loginURL := &url.URL{
			Path: loginPath,
		}

		q := loginURL.Query()
		q.Set(redirectKey, string(ctx.URI().RequestURI()))
		loginURL.RawQuery = q.Encode()

		ctx.Redirect(loginURL.String(), fasthttp.StatusFound)
		return true
	}
	if ctxUser.UserType() != dtypes.UserTypeAdmin {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTPMessage("user without admin role", http.StatusForbidden))
		return true
	}
	return false
}

func errTemplateInternal(err error) error {
	return errorx.Errorf(
		"template execution error: %v", err,
		errorx.NewFromDomainException(exceptionx.ERROR),
	)
}
