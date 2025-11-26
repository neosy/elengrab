package nfasthttp

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/valyala/fasthttp"
)

func SendFile(ctx *fasthttp.RequestCtx, path string, downloadName string, contentType string) {
	// Open file
	f, err := os.Open(path)
	if err != nil {
		ctx.Error("failed to open file", fasthttp.StatusInternalServerError)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		ctx.Error("failed to stat file", fasthttp.StatusInternalServerError)
		return
	}

	size := fi.Size()

	// Base headers
	ctx.Response.Header.Set("Content-Type", contentType)
	ctx.Response.Header.Set("Accept-Ranges", "bytes")
	ctx.Response.Header.Set("Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s"`, downloadName))

	rangeHeader := string(ctx.Request.Header.Peek("Range"))

	// -------------------------
	// Handle Range request
	// -------------------------
	if strings.HasPrefix(rangeHeader, "bytes=") {
		var start, end int64
		end = size - 1

		_, err := fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
		if err != nil {
			ctx.Error("invalid range", fasthttp.StatusBadRequest)
			return
		}

		if start < 0 || start >= size {
			ctx.Error("invalid range start", fasthttp.StatusRequestedRangeNotSatisfiable)
			return
		}
		if end <= 0 || end >= size {
			end = size - 1
		}

		ctx.SetStatusCode(fasthttp.StatusPartialContent)
		ctx.Response.Header.Set("Content-Range",
			fmt.Sprintf("bytes %d-%d/%d", start, end, size))

		if _, err := f.Seek(start, io.SeekStart); err != nil {
			ctx.Error("seek error", fasthttp.StatusInternalServerError)
			return
		}

		_, err = io.CopyN(ctx, f, end-start+1)
		if err != nil && !errors.Is(err, io.EOF) {
			ctx.Error("stream error", fasthttp.StatusInternalServerError)
		}
		return
	}

	// -------------------------
	// Full file download
	// -------------------------
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.Set("Content-Length", fmt.Sprintf("%d", size))

	_, err = io.Copy(ctx, f)
	if err != nil && !errors.Is(err, io.EOF) {
		ctx.Error("stream error", fasthttp.StatusInternalServerError)
	}
}
