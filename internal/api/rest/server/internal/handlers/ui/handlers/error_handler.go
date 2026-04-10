package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"mime"
	"net/http"
	"strings"

	errormw "github.com/neosy/elengrab/internal/api/rest/server/internal/error_middleware"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	"github.com/neosy/elengrab/internal/pkg/fnx"
	appenv "github.com/neosy/elengrab/internal/pkg/nconfig/app_env"
	"github.com/neosy/elengrab/internal/pkg/nfasthttp"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) WriteErrorHandler(ctx *fasthttp.RequestCtx) {
	var (
		isDebugMode    bool
		debugErrorText string
		errorxResp     nfasthttp.WriteErrorxResponse
	)

	userValue, ok := ctx.UserValue(nfasthttp.AppConfigCtxKey).(appenv.AppEnv)
	if ok {
		switch userValue {
		case appenv.AppEnvDevelop, appenv.AppEnvLocal, appenv.AppEnvTest:
			isDebugMode = true
		}
	}

	errorText, _ := ctx.UserValue(errormw.ErrorTextKey).(string)
	if isDebugMode {
		errorxResp, ok = ctx.UserValue(errormw.ErrorxResponseKey).(nfasthttp.WriteErrorxResponse)
		if ok {
			data, err := json.MarshalIndent(errorxResp, "", " ")
			if err == nil {
				debugErrorText = string(data)
			}
		}
	}

	statusCode := ctx.Response.StatusCode()
	statusText := http.StatusText(statusCode)

	if errorText == "" {
		switch statusCode {
		case fasthttp.StatusNotFound:
			var uri string
			if ctx.Request.URI() != nil && ctx.Request.URI().RequestURI() != nil {
				uri = string(ctx.Request.URI().RequestURI())
			}
			statusText = fmt.Sprintf("The requested URL %s was not found on this server.", uri)
		}
	}

	cssStyle := uivalues.CssFileRaw(uivalues.CssErrorFileName, h.assetFolders.Css())

	errorValues := uivalues.NewErrorValues()
	errorValues.Title = fmt.Sprintf("Error %v (%s)!!!", statusCode, statusText)
	errorValues.BaseURL = fnx.Ternary(h.baseURL != "", h.baseURL, "/")
	errorValues.ErrorCode = statusCode
	errorValues.ErrorTitle = statusText
	errorValues.ErrorText = errorText
	errorValues.DebugErrorText = errorxResp.Errors

	dataMap := uivalues.MergeMaps(uivalues.PathValues, errorValues.ToMap())
	dataMap[uivalues.CssStyleKey] = template.HTML("<style>" + cssStyle + "</style>")
	dataMap[uivalues.DebugDataKey] = template.HTML(strings.ReplaceAll(debugErrorText, "\n", "<br>"))

	// Load template
	tmpl, err := h.loadPage(uivalues.PageError.FileName())
	if err != nil {
		h.logger.Error("Failed to load template", "error", err)
		return
	}

	body := ctx.Response.Body()
	ctx.Response.SetBody(nil)

	// Execute template with PageTitle
	if err := tmpl.ExecuteTemplate(ctx, uivalues.PageError.Key(), dataMap); err != nil {
		ctx.Response.SetBody(body)
		h.logger.Error("Failed to execute template", "error", err)
		return
	}

	// Set content type so browser renders HTML properly
	ctx.SetContentType(mime.TypeByExtension(".html"))
}
