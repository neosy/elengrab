package sse

import (
	"bufio"
	"strconv"
	"strings"
	"time"

	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/domain/types/broadcast"
	dsse "github.com/neosy/elengrab/internal/domain/types/sse"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
	"github.com/valyala/fasthttp"
)

func (h *SSEHandler) StreamEvents(
	ctx *fasthttp.RequestCtx,
	w *bufio.Writer,
	authCtx dauth.AuthContext,
) {
	eventKey := authCtx.EventKey()
	clientKey := broadcast.BuildClientKey(eventKey)

	subscription := h.downloader.Broadcaster().Subscribe(clientKey, authCtx.UserID)
	defer h.downloader.Broadcaster().Unsubscribe(clientKey)

	connectionID := subscription.ConnectionID()

	h.logger.Debug(
		"SSE subscribed",
		"connectionId", connectionID,
		"eventKey", eventKey,
		"login", authCtx.Login,
	)
	defer h.logger.Debug(
		"SSE closed",
		"connectionId", connectionID,
		"eventKey", eventKey,
		"login", authCtx.Login,
	)

	h.sendConnected(ctx, w, authCtx, connectionID)

	if eventID := h.resolveEventsReplayPosition(ctx, authCtx); !eventID.IsZero() {
		_ = h.replayEvents(ctx, w, authCtx, eventID)
	}

	h.sendInitialEvents(ctx, w, authCtx)

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

// sendInitialEvents sends the initial events to the client.
func (h *SSEHandler) sendInitialEvents(
	ctx *fasthttp.RequestCtx,
	w *bufio.Writer,
	authCtx dauth.AuthContext,
) {
	event, err := h.downloader.GetDownloadsInProgressEvent(ctx, authCtx)
	if err == nil && !event.IsEmpty() {
		h.handleDownloadsInProgress(w, event)
	}
}

// resolveEventsReplayPosition reads the replay position from the request.
// A nil logID means there is nothing to replay.
func (h *SSEHandler) resolveEventsReplayPosition(
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
		return h.eventLogCacheRepo.GetPosition(ctx, authCtx.EventKey().UUID().String())
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
func (h *SSEHandler) replayEvents(
	ctx *fasthttp.RequestCtx,
	w *bufio.Writer,
	authCtx dauth.AuthContext,
	eventID dsse.EventID,
) error {
	if eventID.IsZero() {
		return nil
	}

	key := authCtx.EventKey().UUID().String()

	events, _, err := h.eventLogCacheRepo.GetEventsAfter(ctx, key, eventID)
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
			err := h.sendEventWithID(w, event.EventID(eventID.LogID), *event)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
