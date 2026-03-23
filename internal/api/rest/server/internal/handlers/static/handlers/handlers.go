package handlers

import (
	"path/filepath"
	"regexp"
	"strings"

	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	"github.com/valyala/fasthttp"
)

type StaticHandlers struct {
	assetsDir string

	cssHandler  fasthttp.RequestHandler
	imgHandler  fasthttp.RequestHandler
	iconHandler fasthttp.RequestHandler
	jsHandler   fasthttp.RequestHandler
	pwaHandler  fasthttp.RequestHandler
}

func NewStaticHandlers(assetsDir string) *StaticHandlers {

	h := &StaticHandlers{
		assetsDir: assetsDir,
	}

	h.cssHandler = h.newFSHandler("css", "css")
	h.imgHandler = h.newFSHandler("img", "img")
	h.iconHandler = h.newFSHandler("img/icons", "icon")
	h.jsHandler = h.newFSHandler("js", "js")
	h.pwaHandler = h.newFSHandler("pwa", "pwa")

	return h
}

func (h *StaticHandlers) newFSHandler(name string, htmlPathName string) fasthttp.RequestHandler {
	fs := &fasthttp.FS{
		Root:               filepath.Join(h.assetsDir, "static", name),
		GenerateIndexPages: false,
	}

	handler := fs.NewRequestHandler()

	re := regexp.MustCompile(`^(.*?)(\.[0-9a-f]{6,})?(\.css|\.js)$`)

	return func(ctx *fasthttp.RequestCtx) {
		path := httppaths.GroupStatic + "/" + htmlPathName

		pathSuffix := string(ctx.Path()[len(path):])
		fileName, _ := strings.CutPrefix(pathSuffix, "/")

		matches := re.FindStringSubmatch(fileName)
		if len(matches) == 4 {
			fileName = matches[1] + matches[3]
		}

		pathSuffix = "/" + fileName

		ctx.Request.SetRequestURIBytes([]byte(pathSuffix))
		handler(ctx)
	}
}
