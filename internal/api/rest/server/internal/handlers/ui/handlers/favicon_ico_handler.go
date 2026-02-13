package handlers

import "github.com/valyala/fasthttp"

func (h *DownloaderHandlers) FaviconICOHandler(ctx *fasthttp.RequestCtx) {
	fs := &fasthttp.FS{
		Root:               h.assetsDir + "/static/img",
		GenerateIndexPages: false,
		Compress:           false,
	}
	fs.NewRequestHandler()(ctx)
}
