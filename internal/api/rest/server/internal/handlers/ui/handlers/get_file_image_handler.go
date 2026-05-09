package handlers

import (
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/policy"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
	"github.com/neosy/elengrab/internal/pkg/httpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) GetFileImageHandler(ctx *fasthttp.RequestCtx) {
	ctxUser := policy.ResolveUserOrAnonym(ctx)

	fileIdStr, ok := ctx.UserValue(fileIdKey).(string)
	if !ok || fileIdStr == "" {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrFileIdIsRequired)
		return
	}

	fileID, err := uuid.Parse(fileIdStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, apierrors.ErrFileIdIsIncorrect.Wrap(err))
		return
	}

	args := ctx.QueryArgs()
	argImageSource := string(args.Peek(sourceKey))
	argImageSources := strings.Split(argImageSource, ",")

	var imageSources []dtypes.ImageSource

	for _, src := range argImageSources {
		if err := h.validators.Validate.Var(src, "required,imageSource"); err != nil {
			nfasthttp.WriteErrorx(ctx, errorx.NewFromError(err, exceptionx.VALIDATE))
			return
		}

		source, err := dtypes.ParseImageSource(src)
		if err == nil {
			imageSources = append(imageSources, source)
		}
	}

	imageData, err := h.downloader.GetFileImage(ctx, *ctxUser, fileID, imageSources)
	if err == nil && imageData != nil && len(imageData.Raw) > 0 {
		ctx.SetContentType(httpx.ContentTypeByExt(imageData.Format.String()))
		ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
		ctx.SetBody(imageData.Raw)
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}

	iconsDir := filepath.Join(h.assetsDir, "static/img/icons")
	defaultAvatarSVG := uivalues.IconFileRawByKey(uivalues.MediaDefaultIconNameKey, iconsDir)

	ctx.SetContentType("image/svg+xml")
	ctx.Response.Header.Set("Cache-Control", "public, max-age=86400")
	ctx.SetBody([]byte(defaultAvatarSVG))
	ctx.SetStatusCode(fasthttp.StatusOK)
}
