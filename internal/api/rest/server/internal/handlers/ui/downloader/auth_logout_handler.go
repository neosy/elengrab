package downloader

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	authmw "github.com/neosy/elengrab/internal/api/rest/server/internal/middleware/auth"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) AuthLogoutHandler(ctx *fasthttp.RequestCtx) {
	if err := nfasthttp.EnforceHTTPS(ctx); err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	ctxUser := policy.ResolveUser(ctx)
	if ctxUser != nil {
		authmw.CookieSessionTokenKey.DeleteWithSecure(ctx)
	}
	ctx.Redirect(httppaths.IndexPath, fasthttp.StatusFound)
}
