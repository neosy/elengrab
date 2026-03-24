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
	Errors     string `json:"errors,omitempty"`
	Timestamp  string `json:"timestamp"`
}

// WriteErrorx
func WriteErrorx(ctx *fasthttp.RequestCtx, err error) {
	var (
		message             string
		errors              string
		httpStatusCode      int
		errCode             uint
		shouldCombineErrors bool
	)

	errx := errorx.NewFromError(err)

	userValue := ctx.UserValue(AppConfigCtxKey)
	if userValue != nil {
		switch userValue {
		case appenv.AppEnvDevelop, appenv.AppEnvLocal:
			shouldCombineErrors = true
		}
	}

	httpStatusCode = fasthttp.StatusInternalServerError
	if errx != nil {
		exception := errx.Exception()

		if httpStatus := errx.HttpStatusCode(); httpStatus != nil {
			httpStatusCode = *httpStatus
		}

		if exception != nil {
			errCode = exception.Num()
		}

		if message == "" && errx.Message() != nil {
			text := errx.Message()
			message = *text
		}

		if message == "" {
			message = errx.Error()
		}

		if message == "" && exception != nil {
			message = exception.String()
		}

		if shouldCombineErrors {
			errors = errx.Error()
		}
	} else {
		if err != nil {
			message = err.Error()

			if shouldCombineErrors {
				errors = err.Error()
			}
		}
	}

	if message == "" {
		message = "Internal Server Error"
	}

	errorResponse := writeErrorxResponseDto{
		HTTPStatus: httpStatusCode,
		Code:       int(errCode),
		Message:    message,
		Errors:     errors,
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	ctx.SetStatusCode(httpStatusCode)
	ctx.SetContentType("application/json")

	errorResponseJSON, err := json.Marshal(errorResponse)
	if err != nil {
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetBody(errorResponseJSON)
}
