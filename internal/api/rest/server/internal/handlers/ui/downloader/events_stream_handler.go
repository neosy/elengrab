package downloader

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"uuid"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/policy"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
	utypes "github.com/neosy/elengrab/internal/app/usecases/types"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/domain/types/broadcast"
	dsse "github.com/neosy/elengrab/internal/domain/types/sse"
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

// sendSSEEventWithID sends a single SSE event to the client.
func (h *DownloaderHandlers) sendSSEEventWithID(
	w *bufio.Writer,
	eventID dsse.EventID,
	event dsse.Event,
) error {
	_, err := fmt.Fprintf(w, "event: %s\n", event.Name)
	if err != nil {
		return err
	}

	if !eventID.IsZero() {
		_, err = fmt.Fprintf(w, "id: %s\n", eventID.String())
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(w, "data: %s\n\n", event.Data)
	if err != nil {
		return err
	}

	w.Flush()

	return nil
}

// sendSSEEvent sends a single SSE event to the client without an event ID.
func (h *DownloaderHandlers) sendSSEEvent(
	w *bufio.Writer,
	name string,
	data []byte,
) error {
	return h.sendSSEEventWithID(w, dsse.EventIDNil(), dsse.NewEvent(name, data))
}

// sendAndStoreSSEEvent stores the event in the replay log and sends it to the client.
func (h *DownloaderHandlers) sendAndStoreSSEEvent(
	ctx context.Context,
	w *bufio.Writer,
	authCtx dauth.AuthContext,
	event dsse.Event,
) error {
	eventID, _ := h.sseEventLogCacheRepo.Push(ctx, authCtx.EventKey().UUID().String(), &event)

	return h.sendSSEEventWithID(w, eventID, event)
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

	h.sendConnected(ctx, w, authCtx, subscription.ConnectionID())

	if eventID := h.resolveEventsReplayPosition(ctx, authCtx); !eventID.IsZero() {
		_ = h.replayEvents(ctx, w, authCtx, eventID)
	}

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

// resolveEventsReplayPosition reads the replay position from the request.
// A nil logID means there is nothing to replay.
func (h *DownloaderHandlers) resolveEventsReplayPosition(
	ctx *fasthttp.RequestCtx,
	authCtx dauth.AuthContext,
) dsse.EventID {
	var eventIDNil dsse.EventID

	lastEventId := string(ctx.QueryArgs().Peek(qkeys.SSELastEventID.String()))
	reconnect := ctx.QueryArgs().GetBool(qkeys.SSEReconnect.String())

	if !reconnect && lastEventId == "" {
		return eventIDNil
	}

	if reconnect && lastEventId == "" {
		return h.sseEventLogCacheRepo.GetPosition(ctx, authCtx.EventKey().UUID().String())
	}

	parts := strings.SplitN(lastEventId, ":", 2)
	if len(parts) != 2 {
		return eventIDNil
	}

	logID, err := idcodec.DecodeUUIDBase64URL(parts[0])
	if err != nil {
		return eventIDNil
	}

	lastNum, err := strconv.Atoi(parts[1])
	if err != nil {
		return eventIDNil
	}

	return dsse.EventID{
		LogID:  logID,
		LogNum: lastNum,
	}
}

// replayEvents sends events stored after lastEventId (the whole log if it is empty).
func (h *DownloaderHandlers) replayEvents(
	ctx *fasthttp.RequestCtx,
	w *bufio.Writer,
	authCtx dauth.AuthContext,
	eventID dsse.EventID,
) error {
	if eventID.IsZero() {
		return nil
	}

	key := authCtx.EventKey().UUID().String()

	events, _, err := h.sseEventLogCacheRepo.GetEventsAfter(ctx, key, eventID)
	if err != nil || len(events) == 0 {
		return nil
	}

	for _, event := range events {
		if event == nil {
			continue
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := h.sendSSEEventWithID(w, event.EventID(eventID.LogID), *event)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *DownloaderHandlers) sendConnected(
	ctx *fasthttp.RequestCtx,
	w *bufio.Writer,
	authCtx dauth.AuthContext,
	connectionID uuid.UUID,
) {
	data := struct {
		ConnectionID string `json:"id"`
	}{
		ConnectionID: connectionID.String(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	sseEvent := dsse.Event{
		Name: "connected",
		Data: jsonData,
	}

	eventID := h.sseEventLogCacheRepo.GetPosition(ctx, authCtx.EventKey().UUID().String())

	_ = h.sendSSEEventWithID(w, eventID, sseEvent)
}

func (h *DownloaderHandlers) sendPing(w *bufio.Writer) error {
	return h.sendSSEEvent(w, "ping", []byte("{}"))
}

func (h *DownloaderHandlers) handleEvent(
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
	case utypes.BroadcastEventTypeSystemInfoUpdate:
		if h.downloader.HasWriteAccess(authCtx) {
			h.handleSystemInfoUpdate(w, event)
		}
	case utypes.BroadcastEventTypeNotification:
		h.handleNotification(w, event)
	}
}

func (h *DownloaderHandlers) handleDownloadAdd(
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
		DownloadID     string `json:"itemId"`
		OwnerUserLogin string `json:"owner"`
		HTML           string `json:"html"`
	}{
		DownloadID:     idcodec.EncodeUUIDBase64URL(downloadInfo.DownloadID),
		OwnerUserLogin: downloadInfo.UserLogin,
		HTML:           html,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	sseEvent := dsse.Event{
		Name: event.Type.SSEEventName(),
		Data: jsonData,
	}

	_ = h.sendAndStoreSSEEvent(ctx, w, authCtx, sseEvent)
}

func (h *DownloaderHandlers) handleDownloadUpdate(
	ctx *fasthttp.RequestCtx,
	authCtx dauth.AuthContext,
	w *bufio.Writer,
	event dsse.BroadcastEvent,
) {
	downloadInfo, ok := event.Data.(*ucdto.MediaDownloadInfo)
	if !ok {
		return
	}

	row := h.renderMediaItemRow(
		ctx,
		renderMediaItemRowParams{
			downloadInfo:    downloadInfo,
			isDownloadEvent: true,
			AuthCtx:         authCtx,
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

	_ = h.sendAndStoreSSEEvent(ctx, w, authCtx, sseEvent)
}

func (h *DownloaderHandlers) handleDownloadPatch(
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

	data := h.mappers.MapMediaDownloadChangedToEventResponse(
		downloadChanged,
		h.buildVisibilityResponse,
		func(downloadID uuid.UUID) bool {
			return h.hasShareLink(ctx, downloadID)
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

	_ = h.sendAndStoreSSEEvent(ctx, w, authCtx, sseEvent)
}

func (h *DownloaderHandlers) handleDownloadDelete(
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

	_ = h.sendAndStoreSSEEvent(ctx, w, authCtx, sseEvent)
}

func (h *DownloaderHandlers) handleProgressUpdate(w *bufio.Writer, event dsse.BroadcastEvent) {
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

	_ = h.sendSSEEvent(w, event.Type.SSEEventName(), jsonData)
}

func (h *DownloaderHandlers) handleStartRefreshing(
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

	_ = h.sendAndStoreSSEEvent(ctx, w, authCtx, sseEvent)
}

func (h *DownloaderHandlers) handleSystemInfoUpdate(w *bufio.Writer, event dsse.BroadcastEvent) {
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

	_ = h.sendSSEEvent(w, event.Type.SSEEventName(), jsonData)
}

func (h *DownloaderHandlers) handleNotification(w *bufio.Writer, event dsse.BroadcastEvent) {
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

	_ = h.sendSSEEvent(w, event.Type.SSEEventName(), jsonData)
}
