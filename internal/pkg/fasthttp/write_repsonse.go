package nfasthttp

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

func WriteResponse(ctx *fasthttp.RequestCtx, t any) error {
	if isNil(t) {
		ctx.SetStatusCode(fasthttp.StatusNoContent)
		return nil
	}

	response, err := json.Marshal(t)
	if err != nil {
		WriteErrorx(ctx, err)
		return err
	}

	statusCode := ctx.Response.StatusCode()
	if statusCode == 0 {
		statusCode = fasthttp.StatusOK
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(statusCode)
	ctx.SetBody(response)

	return nil
}
