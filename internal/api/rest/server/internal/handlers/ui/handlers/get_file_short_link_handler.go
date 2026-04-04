package handlers

import (
	"strings"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/dto"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) GetFileShortLinkHandler(ctx *fasthttp.RequestCtx) {
	fileIdStr := ctx.UserValue(fileIdKey).(string)
	if fileIdStr == "" {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTP("fileId is required", fasthttp.StatusBadRequest))
		return
	}

	fileId, err := uuid.Parse(fileIdStr)
	if err != nil {
		nfasthttp.WriteErrorx(ctx, errorx.NewHTTP("fileId is incorrect", fasthttp.StatusBadRequest))
		return
	}

	url, err := h.linkWeb.CreateShortLink(ctx, strings.TrimSuffix(h.baseURL, "/")+httppaths.BuildPathFile(fileId))
	if err != nil {
		nfasthttp.WriteErrorx(ctx, err)
		return
	}

	resp := dto.GetFileShareLinkResponse{
		FileID: fileId.String(),
		URL:    url,
	}

	nfasthttp.WriteResponse(ctx, resp)
}
