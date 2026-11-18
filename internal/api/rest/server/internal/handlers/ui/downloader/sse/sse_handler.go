package sse

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"github.com/neosy/elengrab/internal/app/usecases/downloader"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
	linkweb "github.com/neosy/elengrab/internal/app/usecases/link_web"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dsse "github.com/neosy/elengrab/internal/domain/types/sse"
	"github.com/neosy/elengrab/internal/ports/persistence"
	"github.com/valyala/fasthttp"
)

type (
	renderMediaItemRowPlaceholderFunc func(
		ctx *fasthttp.RequestCtx,
		authCtx dauth.AuthContext,
		downloadInfo *ucdto.MediaDownloadInfo,
	) (*bytes.Buffer, error)

	renderMediaItemRowFunc func(
		ctx *fasthttp.RequestCtx,
		authCtx dauth.AuthContext,
		downloadInfo *ucdto.MediaDownloadInfo,
	) ([]byte, error)
)

type SSEHandler struct {
	logger *slog.Logger

	eventLogCacheRepo persistence.SSEEventLogCacheRepository

	// Usecases
	downloader downloader.DownloaderAPI
	linkWeb    linkweb.LinkWeb

	renderMediaItemRowPlaceholder renderMediaItemRowPlaceholderFunc
	renderMediaItemRow            renderMediaItemRowFunc
}

func NewSSEHandler(
	logger *slog.Logger,

	eventLogCacheRepo persistence.SSEEventLogCacheRepository,
	downloader downloader.DownloaderAPI,
	linkWeb linkweb.LinkWeb,

	renderMediaItemRowPlaceholder renderMediaItemRowPlaceholderFunc,
	renderMediaItemRow renderMediaItemRowFunc,
) *SSEHandler {
	return &SSEHandler{
		logger: logger,

		eventLogCacheRepo: eventLogCacheRepo,
		downloader:        downloader,
		linkWeb:           linkWeb,

		renderMediaItemRowPlaceholder: renderMediaItemRowPlaceholder,
		renderMediaItemRow:            renderMediaItemRow,
	}
}

// SetupHeaders configures HTTP response headers for an SSE connection.
func (h *SSEHandler) SetupHeaders(ctx *fasthttp.RequestCtx) {
	header := &ctx.Response.Header
	header.Set("Content-Type", "text/event-stream")
	header.Set("Cache-Control", "no-cache")
	header.Set("Connection", "keep-alive")
}

// sendEventWithID sends a single SSE event to the client, optionally including its event ID.
func (h *SSEHandler) sendEventWithID(
	w *bufio.Writer,
	eventID dsse.EventID,
	event dsse.Event,
) error {
	if _, err := fmt.Fprintf(w, "event: %s\n", event.Name); err != nil {
		return err
	}

	if !eventID.IsZero() {
		if _, err := fmt.Fprintf(w, "id: %s\n", eventID.String()); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(w, "data: %s\n\n", event.Data); err != nil {
		return err
	}

	return w.Flush()
}

// sendEvent sends a single SSE event to the client without an event ID.
func (h *SSEHandler) sendEvent(w *bufio.Writer, name string, data []byte) error {
	return h.sendEventWithID(w, dsse.EventIDNil(), dsse.NewEvent(name, data))
}

// sendAndStoreEvent stores the event in the replay log and sends it to the client.
func (h *SSEHandler) sendAndStoreEvent(
	ctx context.Context,
	w *bufio.Writer,
	authCtx dauth.AuthContext,
	event dsse.Event,
) error {
	eventID, _ := h.eventLogCacheRepo.Push(ctx, authCtx.EventKey().UUID().String(), &event)

	return h.sendEventWithID(w, eventID, event)
}
