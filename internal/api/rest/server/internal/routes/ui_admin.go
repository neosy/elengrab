package routes

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
)

// registerUIAdmin register UI admin routes.
func (r *routes) registerUIAdmin(handlers *admin.AdminHandlers) {
	middlewareError := r.middlewares.Error.ErrorHandler

	// Admin with middleware (error, auth or anonym)
	g := nfasthttp.NewRouterGroup(httppaths.AdminGroup, r.router)
	g.Use(middlewareError, r.middlewares.Auth.AuthOrAnonym)
	{
		g.HEAD("", handlers.AdminPageHandler)
		g.GET("", handlers.AdminPageHandler)

		g.HEAD(httppaths.AdminUsersPath, handlers.AdminPageHandler)
		g.GET(httppaths.AdminUsersPath, handlers.AdminPageHandler)
	}

	// With middleware (error, require admin auth)
	g = nfasthttp.NewRouterGroup(httppaths.AdminGroup, r.router)
	g.Use(middlewareError, r.middlewares.Auth.RequireAdmin)
	{
		g.HEAD(httppaths.AdminUserDetailPath, handlers.UserDetailHandler)
		g.GET(httppaths.AdminUserDetailPath, handlers.UserDetailHandler)

		g.HEAD(httppaths.AdminUserTableRowPath, handlers.UserTableRowHandler)
		g.GET(httppaths.AdminUserTableRowPath, handlers.UserTableRowHandler)

		g.POST(httppaths.AdminUserRolesPath, handlers.SetUserRolesHandler)
	}
}
