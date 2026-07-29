package fasthttpx

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/valyala/fasthttp"
)

// SendFileBuffered sends a file to the client over HTTP with support for partial downloads (Range requests).
// This is a simple, buffered implementation: the file is read through Go memory buffers.
// For large files, memory usage may be high, especially when multiple downloads occur simultaneously.
//
// Parameters:
//
//	ctx          - the fasthttp request context used to write the response.
//	path         - the path to the file on disk to be sent.
//	downloadName - the filename suggested to the client for saving the file.
//	contentType  - the MIME type of the file (e.g., "application/octet-stream", "image/png").
//
// Features:
//   - Supports full file downloads.
//   - Supports HTTP Range requests for partial downloads.
//   - Uses Go buffers to stream data, which may increase RAM usage for large files.
//   - Simple implementation without guaranteed zero-copy; not recommended for very large files.
func SendFileBuffered(ctx *fasthttp.RequestCtx, path string, downloadName string, contentType string) {
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
	ctx.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, downloadName))

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
		ctx.Response.Header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, size))
		ctx.Response.Header.Set("Content-Length", fmt.Sprintf("%d", end-start+1))

		if _, err := f.Seek(start, io.SeekStart); err != nil {
			ctx.Error("seek error", fasthttp.StatusInternalServerError)
			return
		}

		w := ctx.Response.BodyWriter()
		_, err = io.CopyN(w, f, end-start+1)
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

	w := ctx.Response.BodyWriter()
	_, err = io.Copy(w, f)
	if err != nil && !errors.Is(err, io.EOF) {
		ctx.Error("stream error", fasthttp.StatusInternalServerError)
	}
}

// SendFileDirect streams a file to the client as a downloadable attachment using fasthttp's built-in SendFile method.
// This implementation is optimized for large files: it uses zero-copy where possible, minimizing
// memory usage and avoiding buffering the entire file in Go.
//
// Parameters:
//
//	ctx          - the fasthttp request context used to write the response.
//	path         - the path to the file on disk to be sent.
//	downloadName - the filename suggested to the client for saving the file.
//	contentType  - the MIME type of the file (e.g., "application/octet-stream", "image/png").
//
// Features:
//   - Forces file download using Content-Disposition: attachment.
//   - Supports HTTP Range requests (partial downloads).
//   - Minimal RAM usage, even for gigabyte-scale files.
//   - Uses fasthttp.SendFile, which leverages zero-copy sendfile on supported operating systems.
//   - Suitable for production environments and large files.
func SendFileDirect(ctx *fasthttp.RequestCtx, path string, downloadName string, contentType string) {
	prepareFileResponse(ctx, contentType)

	ctx.Response.Header.Set(
		"Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s"`, downloadName),
	)

	// fasthttp.SendFile streams the file with zero-copy where possible
	// and handles Range requests automatically.
	ctx.SendFile(path)
}

// StreamFileDirect streams a file to the client for inline browser playback using fasthttp's built-in SendFile method.
// This implementation is optimized for media streaming: it allows the browser to play the file directly
// and supports seeking through HTTP Range requests.
//
// Parameters:
//
//	ctx         - the fasthttp request context used to write the response.
//	path        - the path to the file on disk to be streamed.
//	contentType - the MIME type of the file (e.g., "video/mp4", "audio/mpeg").
//
// Features:
//   - Allows inline playback in browsers using Content-Disposition: inline.
//   - Supports HTTP Range requests for seeking and partial loading.
//   - Minimal RAM usage, even for large media files.
//   - Uses fasthttp.SendFile, which leverages zero-copy sendfile on supported operating systems.
//   - Suitable for video and audio streaming.
func StreamFileDirect(ctx *fasthttp.RequestCtx, path string, contentType string) {
	prepareFileResponse(ctx, contentType)

	ctx.Response.Header.Set("Content-Disposition", "inline")

	// fasthttp.SendFile streams the file with zero-copy where possible
	// and handles Range requests automatically.
	ctx.SendFile(path)
}

// prepareFileResponse configures common HTTP headers for file responses.
func prepareFileResponse(ctx *fasthttp.RequestCtx, contentType string) {
	// Disable compression for this response.
	ctx.Response.Header.Del("Content-Encoding")

	ctx.Response.Header.Set("Content-Type", contentType)
	ctx.Response.Header.Set("Accept-Ranges", "bytes")
}
