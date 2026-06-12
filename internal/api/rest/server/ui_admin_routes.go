package httpsrv

import (
	"github.com/fasthttp/router"
	handlers "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin_handlers"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
)

// setupUIDownloaderRoutes setup UI routes.
func (s *httpServer) setupUIAdminRoutes(r *router.Router, handlers *handlers.AdminHandlers) {
	middlewareError := s.errorMiddleware.ErrorHandler

	// Admin with middleware (error, auth or anonym)
	g := nfasthttp.NewRouterGroup(httppaths.AdminGroup, r)
	g.Use(middlewareError, s.authMiddleware.AuthOrAnonym)
	{
		g.HEAD("", handlers.AdminPageHandler)
		g.GET("", handlers.AdminPageHandler)

		g.HEAD(httppaths.AdminUsersPath, handlers.AdminPageHandler)
		g.GET(httppaths.AdminUsersPath, handlers.AdminPageHandler)
	}

	// With middleware (error, require admin auth)
	g = nfasthttp.NewRouterGroup(httppaths.AdminGroup, r)
	g.Use(middlewareError, s.authMiddleware.RequireAdmin)
	{
		g.HEAD(httppaths.AdminUserDetailPath, handlers.UserDetailHandler)
		g.GET(httppaths.AdminUserDetailPath, handlers.UserDetailHandler)

		g.HEAD(httppaths.AdminUserTableRowPath, handlers.UserTableRowHandler)
		g.GET(httppaths.AdminUserTableRowPath, handlers.UserTableRowHandler)

		g.POST(httppaths.AdminUserRolesPath, handlers.SetUserRolesHandler)
	}
}
