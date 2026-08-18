package routes

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
)

func (r *routes) registerUIAuth(handlers *downloader.DownloaderHandlers) {
	middlewareError := r.middlewares.Error.ErrorHandler

	// Account with middleware (error, auth or anonym)
	g := nfasthttp.NewRouterGroup(httppaths.AuthGroup, r.router)
	g.Use(middlewareError, r.middlewares.Auth.AuthOrAnonym)
	{
		g.GET(httppaths.RegisterPath, handlers.AuthRegisterPageHandler)
		g.HEAD(httppaths.RegisterPath, handlers.AuthRegisterPageHandler)
		g.POST(httppaths.RegisterPath, handlers.AuthRegisterSubmitHandler)

		g.GET(httppaths.LoginPath, handlers.AuthLoginPageHandler)
		g.HEAD(httppaths.LoginPath, handlers.AuthLoginPageHandler)
		g.POST(httppaths.LoginPath, handlers.AuthLoginSubmitHandler)

		g.GET(httppaths.LogoutPath, handlers.AuthLogoutHandler)
	}
}
