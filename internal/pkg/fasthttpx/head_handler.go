package fasthttpx

import (
	"github.com/valyala/fasthttp"
)

// MakeHeadHandler
func MakeHeadHandler(contentType string) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		WriteHeadResponse(ctx, contentType)
	}
}

// HTMLHeadHandler
func HTMLHeadHandler(ctx *fasthttp.RequestCtx) {
	WriteHTMLHeadResponse(ctx)
}

// CSSHeadHandler
func CSSHeadHandler(ctx *fasthttp.RequestCtx) {
	WriteHeadResponse(ctx, "text/css; charset=utf-8")
}

// JSHeadHandler
func JSHeadHandler(ctx *fasthttp.RequestCtx) {
	WriteHeadResponse(ctx, "text/javascript; charset=utf-8")
}
