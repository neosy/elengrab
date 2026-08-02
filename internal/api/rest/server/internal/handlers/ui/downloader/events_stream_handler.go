package downloader

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/domain/types/broadcast"
	"github.com/neosy/elengrab/internal/pkg/humanize"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

func (h *DownloaderHandlers) EventsStreamHandler(ctx *fasthttp.RequestCtx) {
	authCtx := policy.ResolveUserOrAnonym(ctx)

	// SSE headers
	h.setupSSEHeaders(ctx)

	ctx.SetBodyStreamWriter(func(w *bufio.Writer) {
		h.streamEvents(ctx, w, authCtx)
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
	authCtx dauth.AuthContext,
) {
	eventKey := authCtx.EventKey()
	clientKey := broadcast.BuildClientKey(eventKey)

	subscription := h.downloader.Broadcaster().Subscribe(clientKey, authCtx.RoleIDs)
	defer h.downloader.Broadcaster().Unsubscribe(clientKey)

	h.sendConnected(w, subscription.ConnectionID())

	pingTicker := time.NewTicker(15 * time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-pingTicker.C:
			err := h.sendPing(w)
			if err != nil {
				return
			}
		case event, ok := <-subscription.EventCh():
			if !ok {
				return
			}
			h.handleEvent(ctx, w, authCtx, event)
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

func (h *DownloaderHandlers) sendPing(w *bufio.Writer) error {
	if _, err := fmt.Fprint(w, "event: ping\ndata: {}\n\n"); err != nil {
		return err
	}

	return w.Flush()
}

func (h *DownloaderHandlers) handleEvent(
	ctx context.Context,
	w *bufio.Writer,
	authCtx dauth.AuthContext,
	event ucdto.BroadcastEvent,
) {
	switch event.Type {
	case ucdto.BroadcastEventTypeDownloadAdd:
		h.handleDownloadAdd(w, event)
	case ucdto.BroadcastEventTypeDownloadUpdate:
		h.handleDownloadUpdate(ctx, w, event)
	case ucdto.BroadcastEventTypeDownloadDelete:
		h.handleDownloadDelete(w, event)
	case ucdto.BroadcastEventTypeProgressUpdate:
		h.handleProgressUpdate(w, event)
	case ucdto.BroadcastEventTypeDownloadStartRefreshing:
		h.handleStartRefreshing(w, event)
	case ucdto.BroadcastEventTypeSystemInfoUpdate:
		if h.downloader.HasWriteOperation(authCtx) {
			h.handleSystemInfoUpdate(w, event)
		}
	case ucdto.BroadcastEventTypeNotification:
		h.handleNotification(w, event)
	}
}

func (h *DownloaderHandlers) handleDownloadAdd(w *bufio.Writer, event ucdto.BroadcastEvent) {
	downloadInfo, ok := event.Data.(*ucdto.ScheduleDownloadResponse)
	if !ok {
		return
	}

	buf := h.renderMediaItemRowPlaceholder(downloadInfo, false)

	html := strings.TrimSpace(buf.String())

	data := struct {
		DownloadID string `json:"itemId"`
		HTML       string `json:"html"`
	}{
		DownloadID: idcodec.EncodeUUIDBase64URL(downloadInfo.DownloadID),
		HTML:       html,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "event: %s\n", event.Type.SSEEventName())
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	w.Flush()
}

func (h *DownloaderHandlers) handleDownloadUpdate(ctx context.Context, w *bufio.Writer, event ucdto.BroadcastEvent) {
	downloadInfo, ok := event.Data.(*ucdto.MediaDownloadInfo)
	if !ok {
		return
	}

	row := h.renderMediaItemRow(
		ctx,
		renderMediaItemRowParams{
			downloadInfo:    downloadInfo,
			isDownloadEvent: true,
		},
	)
	if row.err != nil {
		return
	}
	if row.httpStatus == fasthttp.StatusNoContent {
		return
	}

	var buf bytes.Buffer
	err := h.templates.Base.ExecuteTemplate(&buf, components.ResultRowStatusKey, row.data)
	if err != nil {
		return
	}

	html := strings.TrimSpace(buf.String())

	data := struct {
		DownloadID string `json:"itemId"`
		HTML       string `json:"html"`
	}{
		DownloadID: idcodec.EncodeUUIDBase64URL(downloadInfo.DownloadID),
		HTML:       html,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "event: %s\n", event.Type.SSEEventName())
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	w.Flush()
}

func (h *DownloaderHandlers) handleDownloadDelete(w *bufio.Writer, event ucdto.BroadcastEvent) {
	downloadID, ok := event.Data.(uuid.UUID)
	if !ok {
		return
	}

	data := struct {
		DownloadID string `json:"itemId"`
	}{
		DownloadID: idcodec.EncodeUUIDBase64URL(downloadID),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "event: %s\n", event.Type.SSEEventName())
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
		DownloadID string `json:"itemId"`
		Value      any    `json:"value"`
	}{
		Field:      "progress",
		DownloadID: idcodec.EncodeUUIDBase64URL(progress.DownloadID),
		Value:      progress.Percent,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "event: %s\n", event.Type.SSEEventName())
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	w.Flush()
}

func (h *DownloaderHandlers) handleStartRefreshing(w *bufio.Writer, event ucdto.BroadcastEvent) {
	downloadInfo, ok := event.Data.(*ucdto.MediaDownloadInfo)
	if !ok {
		return
	}

	data := struct {
		DownloadID string `json:"itemId"`
	}{
		DownloadID: idcodec.EncodeUUIDBase64URL(downloadInfo.DownloadID),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "event: %s\n", event.Type.SSEEventName())
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
		DiskFree: humanize.Bytes(int64(systemInfo.DiskFree)),
		DiskUsed: humanize.Bytes(int64(systemInfo.DiskUsed)),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "event: %s\n", event.Type.SSEEventName())
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

	fmt.Fprintf(w, "event: %s\n", event.Type.SSEEventName())
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	w.Flush()
}
