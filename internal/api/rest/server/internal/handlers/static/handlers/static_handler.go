package handlers

import (
	"github.com/valyala/fasthttp"
)

func (h *StaticHandlers) StaticCssHandler(ctx *fasthttp.RequestCtx) {
	h.cssHandler(ctx)
}

func (h *StaticHandlers) StaticImgHandler(ctx *fasthttp.RequestCtx) {
	h.imgHandler(ctx)
}

func (h *StaticHandlers) StaticIconHandler(ctx *fasthttp.RequestCtx) {
	h.iconHandler(ctx)
}

func (h *StaticHandlers) StaticJsHandler(ctx *fasthttp.RequestCtx) {
	h.jsHandler(ctx)
}

func (h *StaticHandlers) StaticPwaHandler(ctx *fasthttp.RequestCtx) {
	h.pwaHandler(ctx)
}
