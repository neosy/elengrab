package static

import (
	"context"
	"path"
	"path/filepath"
	"strings"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
	appenv "github.com/neosy/elengrab/internal/pkg/config/app_env"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/valyala/fasthttp"
)

func (h *StaticHandlers) newAssetHandler(subdirPath string) func(*fasthttp.RequestCtx) {
	handler := func(ctx *fasthttp.RequestCtx) {
		uri := string(ctx.RequestURI())

		// remove query string if there is
		if idx := strings.IndexByte(uri, '?'); idx >= 0 {
			uri = uri[:idx]
		}

		// protection from path traversal
		cleanPath := path.Clean("/" + uri)
		filePath := filepath.Join(subdirPath, filepath.FromSlash(cleanPath))

		var (
			assetFile *dtypes.AssetFile
			err       error
		)

		if h.appEnv == appenv.AppEnvDevelop {
			assetFile, err = h.assets.ReadAssetFileNoCache(filePath)
		} else {
			assetFile, err = h.assets.ReadAssetFile(context.Background(), filePath)
		}
		if err != nil {
			nfasthttp.WriteErrorx(ctx, err)
			return
		}

		ctx.SetContentType(httpx.ContentTypeByExt(filepath.Ext(assetFile.FileName)))
		ctx.SetStatusCode(fasthttp.StatusOK)

		if ctx.IsHead() {
			return
		}

		ctx.Write(assetFile.Raw)
	}

	return handler
}
