package fasthttpx

import (
	"encoding/json"
	"mime"

	"github.com/valyala/fasthttp"
)

// WriteResponse writes the given data as a JSON response to the provided fasthttp.RequestCtx.
func WriteResponse(ctx *fasthttp.RequestCtx, data any) error {
	if isNil(data) {
		ctx.SetStatusCode(fasthttp.StatusNoContent)
		return nil
	}

	response, err := json.Marshal(data)
	if err != nil {
		WriteErrorx(ctx, err)
		return err
	}

	statusCode := ctx.Response.StatusCode()
	if statusCode == 0 {
		statusCode = fasthttp.StatusOK
	}

	ctx.SetContentType(mime.TypeByExtension(".json"))
	ctx.SetStatusCode(statusCode)
	ctx.SetBody(response)

	return nil
}

// WriteHTMLHeadResponse writes a 200 OK response with an HTML Content-Type for HEAD requests.
func WriteHTMLHeadResponse(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType(mime.TypeByExtension(".html"))
	ctx.SetStatusCode(fasthttp.StatusOK)
}
