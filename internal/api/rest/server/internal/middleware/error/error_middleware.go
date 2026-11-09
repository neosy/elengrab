package errormw

import (
	"log/slog"

	"github.com/valyala/fasthttp"
)

type ErrorMiddleware struct {
	logger *slog.Logger

	errorPageHandler func(ctx *fasthttp.RequestCtx)
}

func NewErrorMiddleware(logger *slog.Logger, errorPageHandler func(ctx *fasthttp.RequestCtx)) *ErrorMiddleware {
	return &ErrorMiddleware{
		logger:           logger,
		errorPageHandler: errorPageHandler,
	}
}
