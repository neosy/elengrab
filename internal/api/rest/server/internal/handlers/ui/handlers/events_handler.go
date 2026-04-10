package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	authmw "github.com/neosy/elengrab/internal/api/rest/server/internal/auth_middleware"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	uformat "github.com/neosy/elengrab/internal/pkg/utils/format"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) EventsHandler(ctx *fasthttp.RequestCtx) {
	ctxUser := authmw.UserFromContext(ctx)
	if ctxUser == nil {
		// anonymous
		ctxUser = dauth.UserContextAnonymous()
		ctxUser.UserID = uuid.New()
	}

	// SSE headers
	h.setupSSEHeaders(ctx)

	ctx.SetBodyStreamWriter(func(w *bufio.Writer) {
		h.streamEvents(ctx, w, ctxUser.UserID.String())
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
	chEvent := h.downloader.Broadcaster().Subscribe(userID)
	defer h.downloader.Broadcaster().Unsubscribe(userID, chEvent)

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
	// fmt.Fprintf(w, "event: connected\ndata: {}\n\n")

	jsonData, err := json.Marshal(struct {
		UserID string `json:"id"`
	}{uuid.NewString()})
	if err != nil {
		return
	}

	fmt.Fprintf(w, "event: connected\n")
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	w.Flush()
}

func (h *DownloaderHandlers) sendPing(w *bufio.Writer) {
	fmt.Fprintf(w, "event: ping\ndata: {}\n\n")
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
	case ucdto.BroadcastEventTypeSystemInfoUpdate:
		h.handleSystemInfoUpdate(w, event)
	case ucdto.BroadcastEventTypeNotification:
		h.handleNotification(w, event)
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
		FileID string `json:"fileId"`
		HTML   string `json:"html"`
	}{
		FileID: fileInfo.FileID.String(),
		HTML:   html,
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

	row := h.genRow(fileInfo, true)
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
	err = tmpl.ExecuteTemplate(&buf, row.templateName, row.data)
	if err != nil {
		return
	}

	html := strings.TrimSpace(buf.String())

	data := struct {
		FileID string `json:"fileId"`
		HTML   string `json:"html"`
	}{
		FileID: fileInfo.FileID.String(),
		HTML:   html,
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
		FileID string `json:"fileId"`
	}{
		FileID: fileID.String(),
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

	data := struct {
		Field  string `json:"field"`
		FileID string `json:"fileId"`
		Value  any    `json:"value"`
	}{
		Field:  "progress",
		FileID: progress.FileID.String(),
		Value:  progress.Percent,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "event: row-patch-field\n")
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
		DiskFree: uformat.BytesHuman(systemInfo.DiskFree),
		DiskUsed: uformat.BytesHuman(systemInfo.DiskUsed),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "event: system-info-update\n")
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

	fmt.Fprintf(w, "event: notification\n")
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	w.Flush()
}
