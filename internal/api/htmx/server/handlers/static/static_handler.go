package statich

import (
	"github.com/valyala/fasthttp"
)

func (h *StaticHandlers) StaticHandler(ctx *fasthttp.RequestCtx) {
	// Remove /static prefix
	ctx.Request.SetRequestURIBytes(ctx.Path()[len("/static"):])

	fs := &fasthttp.FS{
		Root:               h.assetsDir + "/static",
		GenerateIndexPages: false,
	}

	fs.NewRequestHandler()(ctx)
}
