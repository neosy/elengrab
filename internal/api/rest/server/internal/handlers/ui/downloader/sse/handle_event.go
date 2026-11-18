package sse

import (
	"bufio"
	"encoding/json"
	"strings"
	"uuid"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/helper"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
	utypes "github.com/neosy/elengrab/internal/app/usecases/types"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	dsse "github.com/neosy/elengrab/internal/domain/types/sse"
	"github.com/neosy/elengrab/internal/pkg/humanize"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

func (h *SSEHandler) handleEvent(
	ctx *fasthttp.RequestCtx,
	w *bufio.Writer,
	authCtx dauth.AuthContext,
	event dsse.BroadcastEvent,
) {
	switch event.Type {
	case utypes.BroadcastEventTypeDownloadAdd:
		h.handleDownloadAdd(ctx, authCtx, w, event)
	case utypes.BroadcastEventTypeDownloadUpdate:
		h.handleDownloadUpdate(ctx, authCtx, w, event)
	case utypes.BroadcastEventTypeDownloadPatch:
		h.handleDownloadPatch(ctx, authCtx, w, event)
	case utypes.BroadcastEventTypeDownloadDelete:
		h.handleDownloadDelete(ctx, authCtx, w, event)
	case utypes.BroadcastEventTypeProgressUpdate:
		h.handleProgressUpdate(w, event)
	case utypes.BroadcastEventTypeDownloadStartRefreshing:
		h.handleStartRefreshing(ctx, authCtx, w, event)
	case utypes.BroadcastEventTypeDownloadsInProgress:
		h.handleDownloadsInProgress(w, event)
	case utypes.BroadcastEventTypeSystemInfoUpdate:
		if h.downloader.HasWriteAccess(authCtx) {
			h.handleSystemInfoUpdate(w, event)
		}
	case utypes.BroadcastEventTypeNotification:
		h.handleNotification(w, event)
	}
}

func (h *SSEHandler) handleDownloadAdd(
	ctx *fasthttp.RequestCtx,
	authCtx dauth.AuthContext,
	w *bufio.Writer,
	event dsse.BroadcastEvent,
) {
	downloadInfo, ok := event.Data.(*ucdto.MediaDownloadInfo)
	if !ok {
		return
	}

	buf, err := h.renderMediaItemRowPlaceholder(ctx, authCtx, downloadInfo)
	if err != nil {
		return
	}

	html := strings.TrimSpace(buf.String())

	data := struct {
		DownloadID      string `json:"itemId"`
		OwnerUserLogin  string `json:"owner"`
		IsReady         bool   `json:"isReady"`
		CreatedAtUnixMs int64  `json:"createdAtUnixMs"`
		HTML            string `json:"html"`
	}{
		DownloadID:      idcodec.EncodeUUIDBase64URL(downloadInfo.DownloadID),
		OwnerUserLogin:  downloadInfo.UserLogin,
		IsReady:         downloadInfo.IsReady(),
		CreatedAtUnixMs: downloadInfo.CreatedAt.UnixMilli(),
		HTML:            html,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	sseEvent := dsse.Event{
		Name: event.Type.SSEEventName(),
		Data: jsonData,
	}

	_ = h.sendAndStoreEvent(ctx, w, authCtx, sseEvent)
}

func (h *SSEHandler) handleDownloadUpdate(
	ctx *fasthttp.RequestCtx,
	authCtx dauth.AuthContext,
	w *bufio.Writer,
	event dsse.BroadcastEvent,
) {
	downloadInfo, ok := event.Data.(*ucdto.MediaDownloadInfo)
	if !ok {
		return
	}

	buf, err := h.renderMediaItemRow(ctx, authCtx, downloadInfo)
	if err != nil || len(buf) == 0 {
		return
	}

	html := strings.TrimSpace(string(buf))

	data := struct {
		DownloadID string `json:"itemId"`
		HTML       string `json:"html"`
		Timestamp  int64  `json:"timestamp"`
	}{
		DownloadID: idcodec.EncodeUUIDBase64URL(downloadInfo.DownloadID),
		HTML:       html,
		Timestamp:  event.Timestamp,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	sseEvent := dsse.Event{
		Name: event.Type.SSEEventName(),
		Data: jsonData,
	}

	_ = h.sendAndStoreEvent(ctx, w, authCtx, sseEvent)
}

func (h *SSEHandler) handleDownloadPatch(
	ctx *fasthttp.RequestCtx,
	authCtx dauth.AuthContext,
	w *bufio.Writer,
	event dsse.BroadcastEvent) {
	downloadChanged, ok := event.Data.(*ucdto.MediaDownloadChanged)
	if !ok {
		return
	}

	if downloadChanged == nil {
		return
	}

	data := h.mapMediaDownloadChangedToEventResponse(
		downloadChanged,
		helper.BuildVisibilityResponse,
		func(downloadID uuid.UUID) bool {
			link, _ := h.linkWeb.ResolveURL(ctx, helper.BuildMediaWatchURL(ctx, downloadID))
			return link != nil
		},
	)

	data.Timestamp = event.Timestamp

	if !data.HasChanges() {
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	sseEvent := dsse.Event{
		Name: event.Type.SSEEventName(),
		Data: jsonData,
	}

	_ = h.sendAndStoreEvent(ctx, w, authCtx, sseEvent)
}

func (h *SSEHandler) handleDownloadDelete(
	ctx *fasthttp.RequestCtx,
	authCtx dauth.AuthContext,
	w *bufio.Writer,
	event dsse.BroadcastEvent,
) {
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

	sseEvent := dsse.Event{
		Name: event.Type.SSEEventName(),
		Data: jsonData,
	}

	_ = h.sendAndStoreEvent(ctx, w, authCtx, sseEvent)
}

func (h *SSEHandler) handleProgressUpdate(w *bufio.Writer, event dsse.BroadcastEvent) {
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

	_ = h.sendEvent(w, event.Type.SSEEventName(), jsonData)
}

func (h *SSEHandler) handleStartRefreshing(
	ctx *fasthttp.RequestCtx,
	authCtx dauth.AuthContext,
	w *bufio.Writer,
	event dsse.BroadcastEvent,
) {
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

	sseEvent := dsse.Event{
		Name: event.Type.SSEEventName(),
		Data: jsonData,
	}

	_ = h.sendAndStoreEvent(ctx, w, authCtx, sseEvent)
}

func (h *SSEHandler) handleDownloadsInProgress(w *bufio.Writer, event dsse.BroadcastEvent) {
	downloadsInProgress, ok := event.Data.(dtypes.IncompleteDownloadsCount)
	if !ok {
		return
	}

	jsonData, err := json.Marshal(downloadsInProgress)
	if err != nil {
		return
	}

	_ = h.sendEvent(w, event.Type.SSEEventName(), jsonData)
}

func (h *SSEHandler) handleSystemInfoUpdate(w *bufio.Writer, event dsse.BroadcastEvent) {
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

	_ = h.sendEvent(w, event.Type.SSEEventName(), jsonData)
}

func (h *SSEHandler) handleNotification(w *bufio.Writer, event dsse.BroadcastEvent) {
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

	_ = h.sendEvent(w, event.Type.SSEEventName(), jsonData)
}
