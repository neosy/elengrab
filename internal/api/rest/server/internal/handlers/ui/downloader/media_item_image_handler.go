package downloader

import (
	"strings"

	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/icons"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) MediaItemImageHandler(ctx *fasthttp.RequestCtx) {
	authCtx := policy.ResolveUserOrAnonym(ctx)

	downloadIDStr, ok := ctx.UserValue(downloadIDKey).(string)
	if !ok || downloadIDStr == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrDownloadIDIsRequired)
		return
	}

	downloadID, err := idcodec.DecodeUUIDBase64URL(downloadIDStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrDownloadIDIsIncorrect.Wrap(err))
		return
	}

	args := ctx.QueryArgs()
	argImageSource := string(args.Peek(sourceKey))
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

	imageData, err := h.downloader.GetDownloadImage(ctx, authCtx, downloadID, imageSources)
	if err == nil && imageData != nil && len(imageData.Raw) > 0 {
		ctx.SetContentType(httpx.ContentTypeByExt(imageData.Format.String()))
		ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
		ctx.SetBody(imageData.Raw)
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}

	defaultAvatarSVG := icons.MediaDefaultIcon.FileRaw()

	ctx.SetContentType("image/svg+xml")
	ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
	ctx.SetBody([]byte(defaultAvatarSVG))
	ctx.SetStatusCode(fasthttp.StatusOK)
}
