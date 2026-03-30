package handlers

import (
	authmw "github.com/neosy/elengrab/internal/api/rest/server/internal/auth_middleware"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) AuthLogoutHandler(ctx *fasthttp.RequestCtx) {
	ctxUser := authmw.UserFromContext(ctx)
	if ctxUser != nil {
		authmw.CookieSessionTokenKey.Delete(ctx)
	}
	ctx.Redirect(httppaths.PathIndex, fasthttp.StatusFound)
}
