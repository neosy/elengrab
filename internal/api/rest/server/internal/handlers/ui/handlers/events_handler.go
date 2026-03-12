package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) EventsHandler(ctx *fasthttp.RequestCtx) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetBodyString(fmt.Sprintf("Authorization error: %v", err))
		return
	}

	// SSE headers
	h.setupSSEHeaders(ctx)

	ctx.SetBodyStreamWriter(func(w *bufio.Writer) {
		h.streamEvents(ctx, w, userID.String())
	})
}

func (h *DownloaderHandlers) setupSSEHeaders(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "text/event-stream")
	ctx.Response.Header.Set("Cache-Control", "no-cache")
	ctx.Response.Header.Set("Connection", "keep-alive")
}

func (h *DownloaderHandlers) streamEvents(
	ctx *fasthttp.RequestCtx,
	w *bufio.Writer,
	userID string,
) {
	chEvent := h.usecases.Downloader.Broadcaster().Subscribe(userID)
	defer h.usecases.Downloader.Broadcaster().Unsubscribe(userID, chEvent)

	h.sendConnected(w)

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			h.sendPing(w)

		case event := <-chEvent:
			h.handleEvent(w, event)
		}
	}
}

func (h *DownloaderHandlers) sendConnected(w *bufio.Writer) {
	fmt.Fprintf(w, ": connected\n\n")
	w.Flush()
}

func (h *DownloaderHandlers) sendPing(w *bufio.Writer) {
	fmt.Fprintf(w, ": ping\n\n")
	w.Flush()
}

func (h *DownloaderHandlers) handleEvent(w *bufio.Writer, event ucdto.BroadcastEvent) {
	switch event.Type {
	case ucdto.BroadcastEventTypeFileAdd:
		h.handleFileAdd(w, event)
	case ucdto.BroadcastEventTypeFileUpdate:
		h.handleFileUpdate(w, event)
	case ucdto.BroadcastEventTypeFileDelete:
		h.handleFileDelete(w, event)
	case ucdto.BroadcastEventTypeProgressUpdate:
		h.handleProgressUpdate(w, event)
	}
}

func (h *DownloaderHandlers) handleFileAdd(w *bufio.Writer, event ucdto.BroadcastEvent) {
	fileInfo, ok := event.Data.(*ucdto.ScheduleDownloadResponse)
	if !ok {
		return
	}

	buf := h.genNewRow(fileInfo, false)

	html := strings.TrimSpace(buf.String())

	data := struct {
		ID   string `json:"id"`
		HTML string `json:"html"`
	}{
		ID:   "row-" + fileInfo.FileID.String(),
		HTML: html,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "event: row-add\n")
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	w.Flush()
}

func (h *DownloaderHandlers) handleFileUpdate(w *bufio.Writer, event ucdto.BroadcastEvent) {
	fileInfo, ok := event.Data.(*ucdto.GetFileInfoResponse)
	if !ok {
		return
	}

	buf, httpStatus, err := h.genRow(fileInfo, true)
	if err != nil {
		return
	}
	if httpStatus == fasthttp.StatusNoContent {
		return
	}

	html := strings.TrimSpace(buf.String())

	data := struct {
		ID   string `json:"id"`
		HTML string `json:"html"`
	}{
		ID:   "row-" + fileInfo.FileID.String(),
		HTML: html,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "event: row-update\n")
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	w.Flush()
}

func (h *DownloaderHandlers) handleFileDelete(w *bufio.Writer, event ucdto.BroadcastEvent) {
	fileID, ok := event.Data.(uuid.UUID)
	if !ok {
		return
	}

	data := struct {
		RowID string `json:"id"`
	}{
		RowID: "row-" + fileID.String(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "event: row-delete\n")
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	w.Flush()
}

func (h *DownloaderHandlers) handleProgressUpdate(w *bufio.Writer, event ucdto.BroadcastEvent) {
	type dataProgress struct {
		ProgressID          string
		IsFileProgressEvent bool
		ProgressPercent     int
	}

	progress, ok := event.Data.(ucdto.FileProgressResponse)
	if !ok {
		return
	}
	data := dataProgress{
		ProgressID:          fmt.Sprintf("progress-%s", progress.FileID),
		IsFileProgressEvent: true,
		ProgressPercent:     progress.Percent,
	}

	var buf bytes.Buffer
	err := h.templates.ExecuteTemplate(&buf, uivalues.GrabResultProgressHtmlFileName, data)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "event: row-patch\n")
	fmt.Fprintf(w, "data: %s\n\n", buf.String())
	w.Flush()
}
