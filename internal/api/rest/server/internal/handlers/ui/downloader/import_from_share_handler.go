package downloader

import (
	"strings"

	apierrors "github.com/neosy/elengrab/internal/api/errors"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) ImportFromShareHandler(ctx *fasthttp.RequestCtx) {
	ctxUser, err := policy.ResolveUserOrFallback(ctx, h.appMode)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	args := ctx.QueryArgs()

	url := string(args.Peek(qkeys.UrlKey.String()))
	if url == "" {
		url = string(args.Peek(qkeys.TextKey.String()))
	}

	url = strings.TrimSpace(url)
	if url == "" {
		fasthttpx.WriteErrorx(ctx, apierrors.ErrURLIsRequired)
		return
	}

	if err := h.validators.Validate.Var(url, "url"); err != nil {
		fasthttpx.WriteErrorx(ctx, apierrors.ErrInvalidURL.Wrap(err))
		return
	}

	resp, err := h.downloader.ScheduleDownload(
		ctx,
		ctxUser,
		url,
		&ddownload.DownloadOptions{
			FormatType:      dtypes.FormatTypeVideoAudio,
			VideoCodec:      dtypes.VideoCodecH264,
			VideoResolution: dtypes.VideoResolution1080p,
			VideoFormat:     dtypes.VideoFormatAuto,
		},
	)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	if resp == nil {
		fasthttpx.WriteErrorx(ctx, apierrors.ErrEmptyResponse)
		return
	}

	ctx.Redirect(httppaths.IndexPath, fasthttp.StatusFound)
}
