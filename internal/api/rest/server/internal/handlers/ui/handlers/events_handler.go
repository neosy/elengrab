package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/handlers/policy"
	uivalues "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/values"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) EventsHandler(ctx *fasthttp.RequestCtx) {
	ctxUser := policy.ResolveUserOrAnonym(ctx)
	connectionID := uuid.New()

	// SSE headers
	h.setupSSEHeaders(ctx)

	ctx.SetBodyStreamWriter(func(w *bufio.Writer) {
		h.streamEvents(ctx, w, connectionID, ctxUser.EventKey())
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
	connectionID uuid.UUID,
	eventKey dtypes.EventKey,
) {
	chEvent := h.downloader.Broadcaster().Subscribe(eventKey)
	defer h.downloader.Broadcaster().Unsubscribe(eventKey, chEvent)

	h.sendConnected(w, connectionID)

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

func (h *DownloaderHandlers) sendConnected(w *bufio.Writer, connectionID uuid.UUID) {
	// fmt.Fprintf(w, "event: connected\ndata: {}\n\n")

	jsonData, err := json.Marshal(struct {
		ConnectionID string `json:"id"`
	}{connectionID.String()})
	if err != nil {
		return
	}

	fmt.Fprint(w, "event: connected\n")
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	w.Flush()
}

func (h *DownloaderHandlers) sendPing(w *bufio.Writer) {
	fmt.Fprint(w, "event: ping\ndata: {}\n\n")
	w.Flush()
}

func (h *DownloaderHandlers) handleEvent(w *bufio.Writer, event ucdto.BroadcastEvent) {
	switch event.Type {
	case ucdto.BroadcastEventTypeDownloadAdd:
		h.handleDownloadAdd(w, event)
	case ucdto.BroadcastEventTypeDownloadUpdate:
		h.handleDownloadUpdate(w, event)
	case ucdto.BroadcastEventTypeDownloadDelete:
		h.handleDownloadDelete(w, event)
	case ucdto.BroadcastEventTypeProgressUpdate:
		h.handleProgressUpdate(w, event)
	case ucdto.BroadcastEventTypeSystemInfoUpdate:
		h.handleSystemInfoUpdate(w, event)
	case ucdto.BroadcastEventTypeNotification:
		h.handleNotification(w, event)
	}
}

func (h *DownloaderHandlers) handleDownloadAdd(w *bufio.Writer, event ucdto.BroadcastEvent) {
	downloadInfo, ok := event.Data.(*ucdto.ScheduleDownloadResponse)
	if !ok {
		return
	}

	buf := h.genNewRow(downloadInfo, false)

	html := strings.TrimSpace(buf.String())

	data := struct {
		DownloadID string `json:"downloadID"`
		HTML       string `json:"html"`
	}{
		DownloadID: downloadInfo.DownloadID.String(),
		HTML:       html,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	fmt.Fprint(w, "event: row-add\n")
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	w.Flush()
}

func (h *DownloaderHandlers) handleDownloadUpdate(w *bufio.Writer, event ucdto.BroadcastEvent) {
	downloadInfo, ok := event.Data.(*ucdto.GetMediaDownloadInfoResponse)
	if !ok {
		return
	}

	row := h.genRow(downloadInfo, true)
	if row.err != nil {
		return
	}
	if row.httpStatus == fasthttp.StatusNoContent {
		return
	}

	// Load template
	tmpl, err := h.templates.Clone()
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, uivalues.ComponentResultRowStatusKey, row.data)
	if err != nil {
		return
	}

	html := strings.TrimSpace(buf.String())

	data := struct {
		DownloadID string `json:"downloadID"`
		HTML       string `json:"html"`
	}{
		DownloadID: downloadInfo.DownloadID.String(),
		HTML:       html,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	fmt.Fprint(w, "event: row-update\n")
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	w.Flush()
}

func (h *DownloaderHandlers) handleDownloadDelete(w *bufio.Writer, event ucdto.BroadcastEvent) {
	downloadID, ok := event.Data.(uuid.UUID)
	if !ok {
		return
	}

	data := struct {
		DownloadID string `json:"downloadID"`
	}{
		DownloadID: downloadID.String(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	fmt.Fprint(w, "event: row-delete\n")
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	w.Flush()
}

func (h *DownloaderHandlers) handleProgressUpdate(w *bufio.Writer, event ucdto.BroadcastEvent) {
	type dataProgress struct {
		ProgressID              string
		IsDownloadProgressEvent bool
		ProgressPercent         int
	}

	progress, ok := event.Data.(ucdto.MediaDownloadProgressResponse)
	if !ok {
		return
	}

	data := struct {
		Field      string `json:"field"`
		DownloadID string `json:"downloadID"`
		Value      any    `json:"value"`
	}{
		Field:      "progress",
		DownloadID: progress.DownloadID.String(),
		Value:      progress.Percent,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	fmt.Fprint(w, "event: row-patch-field\n")
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	w.Flush()
}

func (h *DownloaderHandlers) handleSystemInfoUpdate(w *bufio.Writer, event ucdto.BroadcastEvent) {
	systemInfo, ok := event.Data.(ucdto.SystemInfoResponse)
	if !ok {
		return
	}

	data := struct {
		DiskFree string `json:"diskFree"`
		DiskUsed string `json:"diskUsed"`
	}{
		DiskFree: uformat.BytesHuman(int64(systemInfo.DiskFree)),
		DiskUsed: uformat.BytesHuman(int64(systemInfo.DiskUsed)),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	fmt.Fprint(w, "event: system-info-update\n")
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	w.Flush()
}

func (h *DownloaderHandlers) handleNotification(w *bufio.Writer, event ucdto.BroadcastEvent) {
	notification, ok := event.Data.(ucdto.BroadcastNotification)
	if !ok {
		return
	}

	data := struct {
		Module  string `json:"module"`
		Type    string `json:"type"`
		Message string `json:"message"`
	}{
		Module:  notification.Module.String(),
		Type:    notification.Type.String(),
		Message: notification.Message,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	fmt.Fprint(w, "event: notification\n")
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	w.Flush()
}
