package fasthttpx

import (
	"encoding/json"
	"mime"
	"time"

	appenv "github.com/neosy/elengrab/internal/pkg/config/app_env"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/valyala/fasthttp"
)

// WriteErrorxResponse represents the structure of the error response returned by WriteErrorx.
type WriteErrorxResponse struct {
	// HTTPStatus is the HTTP status code of the response.
	HTTPStatus int `json:"status"`

	// Code is the numeric exception code.
	Code int `json:"code,omitempty"`

	// CodeText is the short textual code identifying the exception, e.g., "TOO_MANY_REQUESTS".
	// This field is included in the response only when the application runs in
	// local, develop, or test environments (AppEnvLocal, AppEnvDevelop, AppEnvTest).
	CodeText string `json:"code_text,omitempty"`

	// Message is a human-readable error message.
	Message string `json:"message"`

	// Errors contains additional error details, typically used for debugging purposes.
	// It is included in the response only when the application is running in debug mode.
	Errors string `json:"errors,omitempty"`

	// Timestamp is the time when the error occurred.
	Timestamp string `json:"timestamp,omitempty"`
}

// WriteErrorx writes an error response to the given fasthttp.RequestCtx based on the provided error.
// It extracts relevant information from the error, such as the HTTP status code, exception code, and message,
// and constructs a JSON response with this information. The response includes additional details in debug mode.
func WriteErrorx(ctx *fasthttp.RequestCtx, err error) {
	var (
		isDebugMode       bool
		message           string
		debugText         string
		httpStatus        int
		exceptionCode     int
		exceptionCodeText string
	)

	errx := errorx.OuterErrorx(err)

	userValue, ok := ctx.UserValue(AppConfigCtxKey).(appenv.AppEnv)
	if ok {
		switch userValue {
		case appenv.AppEnvDevelop, appenv.AppEnvLocal, appenv.AppEnvTest:
			isDebugMode = true
		}
	}

	httpStatus = fasthttp.StatusInternalServerError
	if errx != nil {
		httpStatus = errx.PublicHttpStatus()
		message = errx.PublicMessage()

		exception := errx.OuterException()
		if exception != nil {
			exceptionCode = int(exception.Num())
			if isDebugMode {
				exceptionCodeText = exception.Code()
			}
		}

		if isDebugMode {
			debugText = errx.Error()
		}
	} else if err != nil {
		if isDebugMode {
			debugText = err.Error()
		}
	}

	if message == "" {
		message = "Internal Server Error"
	}

	errorResponse := WriteErrorxResponse{
		HTTPStatus: httpStatus,
		Code:       exceptionCode,
		CodeText:   exceptionCodeText,
		Message:    message,
		Errors:     debugText,
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	ctx.SetStatusCode(httpStatus)
	ctx.SetContentType(mime.TypeByExtension(".json"))

	errorResponseJSON, err := json.Marshal(errorResponse)
	if err != nil {
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetBody(errorResponseJSON)
}
