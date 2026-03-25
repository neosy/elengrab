package nfasthttp

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

func WriteResponse(ctx *fasthttp.RequestCtx, t interface{}) error {
	response, err := json.Marshal(t)
	if err != nil {
		WriteError(ctx, err, fasthttp.StatusInternalServerError)
		return err
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(response)

	return nil
}
