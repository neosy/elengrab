package downloader

import (
	"html/template"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/valyala/fasthttp"
)

type renderMediaItemImageRequest struct {
	downloadID uuid.UUID
	authCtx    dauth.AuthContext
	shortLink  *shortLink
}

func (h *DownloaderHandlers) MediaItemImageHandler(ctx *fasthttp.RequestCtx) {
	authCtx := policy.ResolveUserOrAnonym(ctx)

	downloadID, err := h.extractDownloadID(ctx)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	h.renderMediaItemImage(ctx,
		renderMediaItemImageRequest{
			authCtx:    authCtx,
			downloadID: downloadID,
		},
	)
}

func (h *DownloaderHandlers) ShortLinkImageHandler(ctx *fasthttp.RequestCtx) {
	_, downloadID, err := h.resolveShortLinkToDownloadID(ctx, false)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	h.renderMediaItemImage(ctx,
		renderMediaItemImageRequest{
			downloadID: downloadID,
		},
	)
}

func (h *DownloaderHandlers) renderMediaItemImage(ctx *fasthttp.RequestCtx, req renderMediaItemImageRequest) {
	args := ctx.QueryArgs()
	argImageSource := string(args.Peek(qkeys.SourceKey.String()))
	argImageSource = strings.TrimSpace(argImageSource)

	var imageSources []dtypes.ImageSource

	if argImageSource != "" {
		for src := range strings.SplitSeq(argImageSource, ",") {
			if src == "" {
				continue
			}

			if err := h.validators.Validate.Var(src, "required,imageSource"); err != nil {
				nfasthttp.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
				return
			}

			source, err := dtypes.ParseImageSource(src)
			if err == nil {
				imageSources = append(imageSources, source)
			}
		}
	}

	imageData, err := h.downloader.GetDownloadImage(ctx, req.authCtx, req.downloadID, imageSources)
	if err == nil && imageData != nil && len(imageData.Raw) > 0 {
		ctx.SetContentType(httpx.ContentTypeByExt(imageData.Format.String()))
		ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
		ctx.SetBody(imageData.Raw)
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}

	hasChannelOrSite := func() bool {
		return slices.ContainsFunc(imageSources, func(v dtypes.ImageSource) bool {
			return v == dtypes.ImageSourceChannel || v == dtypes.ImageSourceSite
		})
	}

	var defaultImageSVG template.HTML
	if hasChannelOrSite() {
		downloadInfo, _ := h.downloader.GetDownloadInfo(ctx, req.authCtx, req.downloadID)
		if downloadInfo != nil && !downloadInfo.Status.IsReady() {
			defaultImageSVG = icons.DownloadAvatarPlaceholderIcon.FileRaw()
		}
	}

	if defaultImageSVG == "" {
		defaultImageSVG = icons.MediaDefaultIcon.FileRaw()
	}

	ctx.SetContentType("image/svg+xml")
	ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
	ctx.SetBody([]byte(defaultImageSVG))
	ctx.SetStatusCode(fasthttp.StatusOK)
}
