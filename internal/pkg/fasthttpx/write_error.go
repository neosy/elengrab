package fasthttpx

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

type responseError struct {
	// http status code
	StatusCode int `json:"statusCode"`
	// error message
	Msg string `json:"message,omitempty"`
}

func WriteError(ctx *fasthttp.RequestCtx, err error, statusCode int) error {
	resp := responseError{
		StatusCode: statusCode,
		Msg:        err.Error(),
	}
	respBody, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to make response body: %w", err)
	}
	ctx.SetStatusCode(statusCode)
	ctx.Response.Header.Add("Content-Type", "application/json")
	ctx.SetBody(respBody)

	return nil
}
