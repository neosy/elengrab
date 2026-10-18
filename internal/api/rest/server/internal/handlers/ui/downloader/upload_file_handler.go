package downloader

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"os"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	"github.com/neosy/elengrab/internal/exceptions"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/fasthttpx"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

// UploadFileHandler receives a media file from the upload dialog (multipart/form-data, field "file"),
// saves it to the temporary upload directory, validates it and keeps it in the cache
// until the user confirms the upload.
func (h *DownloaderHandlers) UploadFileHandler(ctx *fasthttp.RequestCtx) {
	authCtx, err := policy.ResolveUserOrFallback(ctx, h.appMode)
	if err != nil {
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	if ctx.IsHead() {
		fasthttpx.WriteHTMLHeadResponse(ctx)
		return
	}

	boundary := ctx.Request.Header.MultipartFormBoundary()
	if len(boundary) == 0 {
		fasthttpx.WriteErrorx(ctx, errorx.NewMessage("Expected a multipart/form-data request.", exceptions.INVALID_REQUEST))
		return
	}

	// With StreamRequestBody enabled the body is read as a stream and is never buffered in memory.
	var body io.Reader
	if ctx.Request.IsBodyStream() {
		body = ctx.RequestBodyStream()
	} else {
		body = bytes.NewReader(ctx.Request.Body())
	}

	reader := multipart.NewReader(body, string(boundary))

	// Find the "file" part; the dialog sends nothing else in this request.
	var part *multipart.Part
	for {
		p, err := reader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fasthttpx.WriteErrorx(ctx, errorx.NewMessage("The upload request is malformed.", exceptions.INVALID_REQUEST, err))
			return
		}
		if p.FormName() == "file" && p.FileName() != "" {
			part = p
			break
		}
	}

	if part == nil {
		fasthttpx.WriteErrorx(ctx, errorx.NewMessage("No file was uploaded.", exceptions.INVALID_REQUEST))
		return
	}

	upload := h.uploadFiles.NewUpload(authCtx.UserID, part.FileName())

	// The temporary uploaded file must outlive this handler: it stays in the cache until
	// the user confirms the upload or the entry expires. Remove it only if the request fails.
	uploaded := false
	defer func() {
		if !uploaded {
			upload.Discard()
		}
	}()

	err = h.writePartUploadFile(upload.PartPath(), part, h.maxUploadSize)
	if err != nil {
		h.logger.Error("Failed to save the uploaded file", "path", upload.PartPath(), "error", err)
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	err = upload.Complete(ctx)
	if err != nil {
		h.logger.Warn("Uploaded file was rejected", "name", part.FileName(), "error", err)
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	err = h.uploadedFileCacheRepo.Save(ctx, upload)
	if err != nil {
		h.logger.Error("Failed to save the uploaded file to the cache", "name", part.FileName(), "error", err)
		fasthttpx.WriteErrorx(ctx, err)
		return
	}

	uploaded = true

	dataResponse := dto.UploadFileResponse{
		FileID: idcodec.EncodeUUIDBase64URL(upload.FileID()),
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	fasthttpx.WriteResponse(ctx, dataResponse)
}

// writePartUploadFile streams the part to a new file and fails if it is larger than maxSize.
func (h *DownloaderHandlers) writePartUploadFile(path string, part io.Reader, maxSize int64) (err error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return errorx.Errorf("failed to create part file %q: %w", path, err,
			errorx.WithErrorMessage("Failed to save the file. Please try again."),
			errorx.WithHttpStatus(fasthttp.StatusInternalServerError),
		)
	}
	defer func() {
		if closeErr := f.Close(); err == nil && closeErr != nil {
			err = errorx.Errorf("failed to close part file %q: %w", path, closeErr,
				errorx.WithErrorMessage("Failed to save the file. Please try again."),
				errorx.WithHttpStatus(fasthttp.StatusInternalServerError),
			)
		}
	}()

	// Read one byte more than allowed to tell "exactly maxSize" from "too large".
	n, err := io.Copy(f, io.LimitReader(part, maxSize+1))
	if err != nil {
		// The request body was cut off: the client cancelled the upload or lost the connection.
		if errors.Is(err, io.ErrUnexpectedEOF) {
			return errorx.Errorf("upload interrupted while writing part file %q: %w", path, err,
				errorx.WithErrorMessage("The upload was interrupted."),
				errorx.WithHttpStatus(fasthttp.StatusBadRequest),
			)
		}

		return errorx.Errorf("failed to write part file %q: %w", path, err,
			errorx.WithErrorMessage("Failed to save the file. Please try again."),
			errorx.WithHttpStatus(fasthttp.StatusInternalServerError),
		)
	}
	if n > maxSize {
		return errorx.Errorf("uploaded file exceeds the limit of %d bytes", maxSize,
			errorx.WithErrorMessage("The file is too large."),
			errorx.WithHttpStatus(fasthttp.StatusRequestEntityTooLarge),
		)
	}

	return nil
}
