package handlers

import (
	"strings"

	"github.com/google/uuid"
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/icons"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/policy"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) MediaItemImageHandler(ctx *fasthttp.RequestCtx) {
	ctxUser := policy.ResolveUserOrAnonym(ctx)

	downloadIDStr, ok := ctx.UserValue(downloadIDKey).(string)
	if !ok || downloadIDStr == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrDownloadIDIsRequired)
		return
	}

	downloadID, err := uuid.Parse(downloadIDStr)
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

	imageData, err := h.downloader.GetDownloadImage(ctx, *ctxUser, downloadID, imageSources)
	if err == nil && imageData != nil && len(imageData.Raw) > 0 {
		ctx.SetContentType(httpx.ContentTypeByExt(imageData.Format.String()))
		ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
		ctx.SetBody(imageData.Raw)
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}

	defaultAvatarSVG := icons.FileRawByKey(icons.MediaDefaultIconNameKey, h.assetFolders.Icons())

	ctx.SetContentType("image/svg+xml")
	ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
	ctx.SetBody([]byte(defaultAvatarSVG))
	ctx.SetStatusCode(fasthttp.StatusOK)
}
