package grabberh

import (
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

type fileRowInfoData struct {
	PathFileRow  string
	YoutubeTitle string
	YoutubeURL   string
	FileSize     string
	Format       string
	DownloadURL  string
}

func (h *GrabberHandlers) GetFileRow(ctx *fasthttp.RequestCtx) {
	fileIdStr := ctx.UserValue(fileIdKey).(string)
	if fileIdStr == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("FileId is required")
		return
	}

	fileId, err := uuid.Parse(fileIdStr)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("FileId is incorrect")
		return
	}

	fileInfo, err := h.usecases.Downloader.GetFileInfo(ctx, fileId)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString(err.Error())
		return
	}

	buf, httpStatus, err := h.genFileRow(ctx, fileInfo, false)
	if err != nil {
		ctx.SetStatusCode(httpStatus)
		ctx.SetBodyString(err.Error())
		return
	}
	if httpStatus == fasthttp.StatusNoContent {
		ctx.SetStatusCode(httpStatus)
		ctx.Response.Header.Set("HX-Trigger", "no-op")
		return
	}

	ctx.SetBody(buf.Bytes())
	ctx.SetStatusCode(fasthttp.StatusOK)
}
