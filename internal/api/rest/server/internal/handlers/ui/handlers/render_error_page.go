package handlers

import (
	"fmt"
	"html/template"
	"mime"
	"strings"

	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttp"
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
	cssStyleRaw, _ := uivalues.CssErrorFileName.Raw(h.assetFolders.Css())

	pageData := uivalues.ErrorPageData{
		BaseValues: uivalues.NewBaseValues(),
		BasePaths:  uivalues.NewBasePaths(),
		Values: uivalues.ErrorPageValues{
			Title:          fmt.Sprintf("Error %v (%s)!!!", req.statusCode, req.statusText),
			Header:         uivalues.Header,
			BaseURL:        fnx.Ternary(h.baseURL != "", h.baseURL, "/"),
			CssStyle:       template.HTML("<style>" + string(cssStyleRaw) + "</style>"),
			ErrorCode:      req.statusCode,
			ErrorTitle:     req.statusText,
			ErrorText:      req.errorText,
			DebugErrorText: req.errorxResp.Errors,
			DebugData:      template.HTML(strings.ReplaceAll(string(req.debugErrorData), "\n", "<br>")),
		},
	}

	// Load template
	tmpl, err := h.loadPageTemplate(uivalues.PageError.FileName())
	if err != nil {
		h.logger.Error("Failed to load template", "error", err)
		return
	}

	body := ctx.Response.Body()
	ctx.Response.SetBody(nil)

	// Execute template with PageTitle
	if err := tmpl.ExecuteTemplate(ctx, uivalues.PageError.Key(), pageData); err != nil {
		ctx.Response.SetBody(body)
		h.logger.Error("Failed to execute template", "error", err)
		return
	}

	// Set content type so browser renders HTML properly
	ctx.SetContentType(mime.TypeByExtension(".html"))
}
