package middleware

import (
	"log/slog"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers"
	authmw "github.com/neosy/elengrab/internal/api/rest/server/internal/middleware/auth"
	errormw "github.com/neosy/elengrab/internal/api/rest/server/internal/middleware/error"
	"github.com/neosy/elengrab/internal/app/usecases"
	iconfig "github.com/neosy/elengrab/internal/config"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type Middlewares struct {
	Auth  *authmw.AuthMiddleware
	Error *errormw.ErrorMiddleware
}

func NewMiddlewares(
	logger *slog.Logger,
	usecases *usecases.Usecases,
	handlers *handlers.Handlers,

	// Options
	appMode dtypes.AppMode,
	reverseProxy iconfig.HTTPServerReverseProxyConfig,
) *Middlewares {
	return &Middlewares{
		Auth:  authmw.NewAuthMiddleware(logger, usecases.Auth, appMode, reverseProxy),
		Error: errormw.NewErrorMiddleware(logger, handlers.UI.Downloader.ErrorPageHandler),
	}
}
