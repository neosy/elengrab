package handlers

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/neosy/elengrab/internal/api/rest/server/assets"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	"github.com/neosy/elengrab/internal/app/usecases/downloader"
	"github.com/neosy/elengrab/internal/app/usecases/thumbnail"
	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/valyala/fasthttp"
)

type StaticHandlers struct {
	assetsDir string

	assetFolders assets.FolderPaths

	// caches
	assetFileCacheRep persistence.AssetFileCacheRepository

	cssHandler   fasthttp.RequestHandler
	fontsHandler fasthttp.RequestHandler
	imgHandler   fasthttp.RequestHandler
	iconHandler  fasthttp.RequestHandler
	jsHandler    fasthttp.RequestHandler
	pwaHandler   fasthttp.RequestHandler

	// usecases
	thumbnail  *thumbnail.Thumbnail
	downloader *downloader.Downloader
}

func NewStaticHandlers(
	assetsDir string,

	// caches
	assetFileCacheRep persistence.AssetFileCacheRepository,

	// usecases
	thumbnail *thumbnail.Thumbnail,
	downloader *downloader.Downloader,
) *StaticHandlers {

	h := &StaticHandlers{
		assetsDir:    assetsDir,
		assetFolders: assets.NewFolderPaths(assetsDir),

		// caches
		assetFileCacheRep: assetFileCacheRep,

		// usecases
		thumbnail:  thumbnail,
		downloader: downloader,
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
		Root:               filepath.Join(h.assetsDir, "static", name),
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
		path := httppaths.GroupStatic + "/" + httpPath

		pathSuffix := string(ctx.Path()[len(path):])
		fileName, _ := strings.CutPrefix(pathSuffix, "/")

		matches := re.FindStringSubmatch(fileName)
		var hash string
		if len(matches) == 4 {
			fileName = matches[1] + matches[3]
			hash = strings.TrimPrefix(matches[2], ".")
		}

		filePath := "/" + fileName

		ctx.Request.SetRequestURIBytes([]byte(filePath))
		if hash != "" {
			ctx.SetUserValue(hashKey, hash)
		}

		handler(ctx)
	}
}
