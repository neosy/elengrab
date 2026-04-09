package nfasthttp

import (
	"context"
	"fmt"
	"log/slog"

	appenv "github.com/neosy/elengrab/internal/pkg/nconfig/app_env"
	"github.com/valyala/fasthttp"
)

func NewHandler(baseCtx context.Context, logger *slog.Logger, env appenv.AppEnv, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return loggerMiddleware(logger, corsMiddleware(baseCtx, env, next))
}

func loggerMiddleware(logger *slog.Logger, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		next(ctx)

		status := ctx.Response.StatusCode()
		if status >= 400 {
			logger.DebugContext(
				ctx,
				fmt.Sprintf(
					"Request %s %s %s -> %d",
					ctx.Method(),
					ctx.Host(),
					ctx.RequestURI(),
					status,
				),
				slog.String("Client IP", GetClientIP(ctx)),
			)
		}
	}
}

func corsMiddleware(baseCtx context.Context, env appenv.AppEnv, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetUserValue(RequestCtxKey, baseCtx)
		ctx.SetUserValue(AppConfigCtxKey, env)
		next(ctx)
	}
}
