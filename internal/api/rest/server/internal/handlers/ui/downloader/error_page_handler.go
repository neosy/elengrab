package downloader

import (
	"encoding/json"
	"fmt"
	"net/http"

	errormw "github.com/neosy/elengrab/internal/api/rest/server/internal/middleware/error"
	appenv "github.com/neosy/elengrab/internal/pkg/config/app_env"
	nfasthttp "github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) ErrorPageHandler(ctx *fasthttp.RequestCtx) {
	var (
		isDebugMode    bool
		debugErrorData []byte
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
				debugErrorData = data
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

	h.renderErrorPage(ctx,
		renderErrorPageRequest{
			statusCode:     statusCode,
			errorText:      errorText,
			statusText:     statusText,
			debugErrorData: debugErrorData,
			errorxResp:     errorxResp,
		},
	)
}
