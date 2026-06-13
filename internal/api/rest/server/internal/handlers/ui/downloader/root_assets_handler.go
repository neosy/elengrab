package downloader

import (
	"path/filepath"

	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) AssetRobotsHandler(ctx *fasthttp.RequestCtx) {
	fasthttp.ServeFileUncompressed(ctx, filepath.Join(h.assets.FolderPaths().Static(), "robots.txt"))
}

func (h *DownloaderHandlers) AssetFabiconHandler(ctx *fasthttp.RequestCtx) {
	fasthttp.ServeFileUncompressed(ctx, filepath.Join(h.assets.FolderPaths().Img(), "favicon.ico"))
}
