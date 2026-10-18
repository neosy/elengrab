package fasthttpx

import (
	"encoding/json"
	"fmt"
	"mime"

	"github.com/valyala/fasthttp"
)

type responseError struct {
	// http status code
	StatusCode int `json:"statusCode"`
	// error message
	Msg string `json:"message,omitempty"`
}

// WriteErrorx writes an error response in JSON format to the provided fasthttp.RequestCtx.
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
	ctx.SetContentType(mime.TypeByExtension(".json"))
	ctx.SetBody(respBody)

	return nil
}
