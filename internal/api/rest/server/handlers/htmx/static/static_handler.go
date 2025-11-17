package statich

import (
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/paths"
	"github.com/valyala/fasthttp"
)

func (h *StaticHandlers) staticHandler(ctx *fasthttp.RequestCtx, name string) {
	// Remove /static prefix
	ctx.Request.SetRequestURIBytes(ctx.Path()[len(httppaths.GroupStatic+"/"+name):])

	fs := &fasthttp.FS{
		Root:               h.assetsDir + "/static/" + name,
		GenerateIndexPages: false,
	}

	fs.NewRequestHandler()(ctx)
}

func (h *StaticHandlers) StaticCssHandler(ctx *fasthttp.RequestCtx) {
	h.staticHandler(ctx, "css")
}

func (h *StaticHandlers) StaticImgHandler(ctx *fasthttp.RequestCtx) {
	h.staticHandler(ctx, "img")
}

func (h *StaticHandlers) StaticJsHandler(ctx *fasthttp.RequestCtx) {
	h.staticHandler(ctx, "js")
}

func (h *StaticHandlers) StaticPwaHandler(ctx *fasthttp.RequestCtx) {
	h.staticHandler(ctx, "pwa")
}
