package fasthttpx

import (
	"github.com/valyala/fasthttp"
)

// WriteHeadResponse
func WriteHeadResponse(ctx *fasthttp.RequestCtx, contentType string) {
	ctx.SetContentType(contentType)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

// WriteHTMLHeadResponse writes a 200 OK response with an HTML Content-Type for HEAD requests.
func WriteHTMLHeadResponse(ctx *fasthttp.RequestCtx) {
	WriteHeadResponse(ctx, "text/html; charset=utf-8")
}

// WriteJSONHeadResponse writes a 200 OK response with an JSON Content-Type for HEAD requests.
func WriteJSONHeadResponse(ctx *fasthttp.RequestCtx) {
	WriteHeadResponse(ctx, "application/json")
}
