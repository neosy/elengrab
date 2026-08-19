package downloader

import (
	"fmt"
	"html/template"
	"mime"
	"strings"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/fnx"
	"github.com/valyala/fasthttp"
)

type renderErrorPageRequest struct {
	statusCode            int
	errorText, statusText string
	debugErrorData        []byte
	errorxResp            nfasthttp.WriteErrorxResponse
}

func (h *DownloaderHandlers) renderErrorPage(
	ctx *fasthttp.RequestCtx,
	req renderErrorPageRequest,
) {
	cssStyleRaw, _ := paths.ErrorCssFileName.Raw(ctx, h.assets)

	pageData := pages.ErrorPageData{
		BaseValues: pages.NewBaseValues(),
		BasePaths:  paths.NewHttpPaths(),
		Values: pages.ErrorPageValues{
			Title:          fmt.Sprintf("Error %v (%s)!!!", req.statusCode, req.statusText),
			Header:         pages.Header,
			BaseURL:        fnx.Ternary(h.baseURL != "", h.baseURL, "/"),
			CssStyle:       template.HTML("<style>" + string(cssStyleRaw) + "</style>"),
			ErrorCode:      req.statusCode,
			ErrorTitle:     req.statusText,
			ErrorText:      req.errorText,
			DebugErrorText: req.errorxResp.Errors,
			DebugData:      template.HTML(strings.ReplaceAll(string(req.debugErrorData), "\n", "<br>")),
		},
	}

	body := ctx.Response.Body()
	ctx.Response.SetBody(nil)

	// Execute template with PageTitle
	if err := h.templates.Pages[pages.ErrorPage.Key()].ExecuteTemplate(ctx, pages.ErrorPage.Key(), pageData); err != nil {
		ctx.Response.SetBody(body)
		h.logger.Error("Failed to execute template", "error", err)
		return
	}

	// Set content type so browser renders HTML properly
	ctx.SetContentType(mime.TypeByExtension(".html"))
}
