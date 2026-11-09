package errormw

import (
	"encoding/json"
	"fmt"
	"mime"

	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

// ErrorNotFoundHandler is a handler that handles 404 Not Found errors.
func (m *ErrorMiddleware) ErrorNotFoundHandler(ctx *fasthttp.RequestCtx) {
	var uri string
	if ctx.Request.URI() != nil && ctx.Request.URI().RequestURI() != nil {
		uri = string(ctx.Request.URI().RequestURI())
	}

	errorText := fmt.Sprintf("The requested URL %s was not found on this server.", uri)

	if ctx.IsGet() {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		m.renderErrorPage(ctx, errorText)
		return
	}

	fasthttpx.WriteErrorx(ctx,
		errorx.Errorf("the requested page could not be found: %s", uri,
			errorx.WithHttpStatus(fasthttp.StatusNotFound),
			errorx.WithErrorMessage("The requested page could not be found."),
		),
	)
}

// ErrorMethodNotAllowedHandler is a handler that handles 405 Method Not Allowed errors.
func (m *ErrorMiddleware) ErrorMethodNotAllowedHandler(ctx *fasthttp.RequestCtx) {
	var uri string
	if ctx.Request.URI() != nil && ctx.Request.URI().RequestURI() != nil {
		uri = string(ctx.Request.URI().RequestURI())
	}

	ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	ctx.SetUserValue(ErrorxResponseKey, fasthttpx.WriteErrorxResponse{
		HTTPStatus: fasthttp.StatusMethodNotAllowed,
		Message:    fasthttp.StatusMessage(fasthttp.StatusMethodNotAllowed),
	})
	m.renderErrorPage(ctx, fmt.Sprintf("The requested method %s is not allowed for URL %s.", ctx.Method(), uri))
}

// WrapErrorHandler wraps the given handler to render an HTML error page for failed GET requests.
// It runs after the next handler, so it sees the final status code and body.
// Errors of other methods (POST, PATCH, DELETE, HEAD) are left as is (JSON).
func (m *ErrorMiddleware) WrapErrorHandler(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		next(ctx)
		m.handleErrorResponse(ctx)
	}
}

// WrapWithErrorHandler wraps the given middleware with the error handler middleware.
func (m *ErrorMiddleware) WrapWithErrorHandler(next fasthttpx.Middleware) fasthttpx.Middleware {
	return fasthttpx.MiddlewareChain(m.WrapErrorHandler, next)
}

// handleErrorResponse replaces a failed GET response (status >= 400) with an HTML error page.
// For JSON errors it passes the message from the body to the page.
// Other methods and non-error responses are not touched.
func (m *ErrorMiddleware) handleErrorResponse(ctx *fasthttp.RequestCtx) {
	if !ctx.IsGet() || ctx.Response.StatusCode() < 400 {
		return
	}

	// htmx requests load fragments into an open page: keep the JSON error as is.
	if len(ctx.Request.Header.Peek("HX-Request")) > 0 {
		return
	}

	contentType := ctx.Response.Header.ContentType()
	if contentType == nil {
		return
	}

	if string(contentType) != mime.TypeByExtension(".json") {
		m.renderErrorPage(ctx, "")
		return
	}

	body := ctx.Response.Body()
	if len(body) == 0 {
		return
	}

	resp := fasthttpx.WriteErrorxResponse{}

	err := json.Unmarshal(body, &resp)
	if err != nil {
		return
	}

	ctx.SetUserValue(ErrorxResponseKey, resp)
	m.renderErrorPage(ctx, resp.Message)
}

func (m *ErrorMiddleware) renderErrorPage(ctx *fasthttp.RequestCtx, errorText string) {
	ctx.SetUserValue(ErrorTextKey, errorText)
	m.errorPageHandler(ctx)
}
