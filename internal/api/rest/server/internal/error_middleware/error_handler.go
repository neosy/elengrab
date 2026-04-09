package errormw

import (
	"encoding/json"
	"fmt"
	"mime"

	"github.com/neosy/elengrab/internal/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

// ErrorNotFoundHandler is a handler that handles 404 Not Found errors.
func (m *ErrorMiddleware) ErrorNotFoundHandler(ctx *fasthttp.RequestCtx) {
	var uri string
	if ctx.Request.URI() != nil && ctx.Request.URI().RequestURI() != nil {
		uri = string(ctx.Request.URI().RequestURI())
	}

	ctx.SetStatusCode(fasthttp.StatusNotFound)
	m.writeError(ctx, fmt.Sprintf("The requested URL %s was not found on this server.", uri))
}

// ErrorMethodNotAllowedHandler is a handler that handles 405 Method Not Allowed errors.
func (m *ErrorMiddleware) ErrorMethodNotAllowedHandler(ctx *fasthttp.RequestCtx) {
	var uri string
	if ctx.Request.URI() != nil && ctx.Request.URI().RequestURI() != nil {
		uri = string(ctx.Request.URI().RequestURI())
	}

	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.SetUserValue(ErrorxResponseKey, nfasthttp.WriteErrorxResponse{
		HTTPStatus: fasthttp.StatusMethodNotAllowed,
		Message:    fasthttp.StatusMessage(fasthttp.StatusMethodNotAllowed),
	})
	m.writeError(ctx, fmt.Sprintf("The requested URL %s was not found on this server.", uri))
}

// ErrorHandler is a middleware that handles errors and writes error responses.
func (m *ErrorMiddleware) ErrorHandler(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		next(ctx)
		m.errorHandler(ctx)
	}
}

// WrapWithErrorHandler wraps the given middleware with the error handler middleware.
func (m *ErrorMiddleware) WrapWithErrorHandler(next nfasthttp.Middleware) nfasthttp.Middleware {
	return nfasthttp.MiddlewareChain(m.ErrorHandler, next)
}

func (m *ErrorMiddleware) errorHandler(ctx *fasthttp.RequestCtx) {
	if !ctx.IsGet() || ctx.Response.StatusCode() < 400 {
		return
	}

	contentType := ctx.Response.Header.ContentType()
	if contentType == nil {
		return
	}

	if string(contentType) != mime.TypeByExtension(".json") {
		m.writeError(ctx, "")
		return
	}

	body := ctx.Response.Body()
	if len(body) == 0 {
		return
	}

	resp := nfasthttp.WriteErrorxResponse{}

	err := json.Unmarshal(body, &resp)
	if err != nil {
		return
	}

	ctx.SetUserValue(ErrorxResponseKey, resp)
	m.writeError(ctx, resp.Message)
}

func (m *ErrorMiddleware) writeError(ctx *fasthttp.RequestCtx, errorText string) {
	ctx.SetUserValue(ErrorTextKey, errorText)
	m.writeErrorHandler(ctx)
}
