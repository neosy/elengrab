package errormw

import (
	"log/slog"

	"github.com/valyala/fasthttp"
)

type ErrorMiddleware struct {
	logger *slog.Logger

	writeErrorHandler func(ctx *fasthttp.RequestCtx)
}

func NewErrorMiddleware(logger *slog.Logger, writeErrorHandler func(ctx *fasthttp.RequestCtx)) *ErrorMiddleware {
	return &ErrorMiddleware{
		logger:            logger,
		writeErrorHandler: writeErrorHandler,
	}
}
