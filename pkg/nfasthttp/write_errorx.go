package nfasthttp

import (
	"encoding/json"
	"time"

	"github.com/neosy/elengrab/pkg/errorx"
	appenv "github.com/neosy/elengrab/pkg/nconfig/app_env"
	"github.com/valyala/fasthttp"
)

type writeErrorxResponseDto struct {
	HTTPStatus int    `json:"status"`
	Code       int    `json:"code,omitempty"`
	Message    string `json:"message"`
	Timestamp  string `json:"timestamp"`
}

// WriteErrorx
func WriteErrorx(ctx *fasthttp.RequestCtx, err error) {
	var message string
	var httpStatusCode int
	var errCode uint
	var errx *errorx.Errorx

	switch e := err.(type) {
	case *errorx.Errorx:
		errx = e
	}

	httpStatusCode = fasthttp.StatusInternalServerError
	if errx != nil {
		exType := errx.ExceptionType()
		exCode := errx.ExceptionCode()

		if exType == nil && exCode != nil {
			tmpType := exCode.Type()
			exType = tmpType
		}

		if exType != nil {
			httpStatusCode = exType.HttpStatusCode()
		}

		if exCode != nil {
			errCode = exCode.Num()
		}

		userValue := ctx.UserValue(AppConfigCtxKey)
		if userValue != nil {
			switch userValue {
			case appenv.AppEnvDevelop, appenv.AppEnvLocal:
				message = errorx.NewErrorTexts().AddUnwrapErr(errx).Join()
			}
		}

		if message == "" && errx.Message() != nil && errx.Message().String() != "" {
			message = errx.Message().String()
		}

		if message == "" && exCode != nil {
			message = exCode.String()
		}

		if message == "" && exType != nil {
			message = exType.String()
		}
	}

	if message == "" && err != nil {
		message = err.Error()
	}

	if message == "" {
		message = "Internal Server Error"
	}

	errorResponse := writeErrorxResponseDto{
		HTTPStatus: httpStatusCode,
		Code:       int(errCode),
		Message:    message,
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	ctx.SetStatusCode(httpStatusCode)
	ctx.SetContentType("application/json")

	jErrorResponse, err := json.Marshal(errorResponse)
	if err != nil {
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetBody(jErrorResponse)
}
