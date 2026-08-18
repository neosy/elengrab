package static

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	"github.com/neosy/elengrab/internal/app/usecases/downloader"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail"
	appenv "github.com/neosy/elengrab/internal/pkg/config/app_env"
	"github.com/valyala/fasthttp"
)

type StaticHandlers struct {
	assets *assets.Assets

	cssHandler   fasthttp.RequestHandler
	fontsHandler fasthttp.RequestHandler
	imgHandler   fasthttp.RequestHandler
	iconHandler  fasthttp.RequestHandler
	jsHandler    fasthttp.RequestHandler
	pwaHandler   fasthttp.RequestHandler

	// usecases
	thumbnail  *thumbnail.Thumbnail
	downloader *downloader.Downloader

	// options
	appEnv appenv.AppEnv
}

func NewStaticHandlers(
	// Assets
	assets *assets.Assets,

	// usecases
	thumbnail *thumbnail.Thumbnail,
	downloader *downloader.Downloader,

	// options
	appEnv appenv.AppEnv,
) *StaticHandlers {

	h := &StaticHandlers{
		assets: assets,

		// usecases
		thumbnail:  thumbnail,
		downloader: downloader,

		// options
		appEnv: appEnv,
	}

	h.cssHandler = h.newFSHandler("css", "css")
	h.fontsHandler = h.newFSHandler("fonts", "fonts")
	h.imgHandler = h.newFSHandler("images", "images")
	h.iconHandler = h.newFSHandler("icons", "icons")
	h.jsHandler = h.newFSHandler("js", "js")
	h.pwaHandler = h.newFSHandler("pwa", "pwa")

	return h
}

func (h *StaticHandlers) newFSHandler(name string, httpPath string) fasthttp.RequestHandler {
	fs := &fasthttp.FS{
		Root:               filepath.Join(h.assets.DirPath(), "static", name),
		GenerateIndexPages: false,
	}

	var handler fasthttp.RequestHandler

	switch name {
	case "js":
		handler = h.newAssetHandler(fs.Root)
	case "css":
		handler = h.newAssetHandler(fs.Root)
	default:
		handler = fs.NewRequestHandler()
	}

	re := regexp.MustCompile(`^(.*?)(\.[0-9a-f]{6,})?(\.css|\.js|\.json)$`)

	return func(ctx *fasthttp.RequestCtx) {
		path := httppaths.StaticGroup + "/" + httpPath

		pathSuffix := string(ctx.Path()[len(path):])
		fileName, _ := strings.CutPrefix(pathSuffix, "/")

		matches := re.FindStringSubmatch(fileName)
		if len(matches) == 4 {
			fileName = matches[1] + matches[3]
		}

		fileURI := "/" + fileName

		ctx.Request.SetRequestURIBytes([]byte(fileURI))

		handler(ctx)
	}
}
