package sse

import (
	"bufio"
	"encoding/json"
	"uuid"

	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dsse "github.com/neosy/elengrab/internal/domain/types/sse"
	"github.com/valyala/fasthttp"
)

func (h *SSEHandler) sendConnected(
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

	eventID := h.eventLogCacheRepo.GetPosition(ctx, authCtx.EventKey().UUID().String())

	_ = h.sendEventWithID(w, eventID, sseEvent)
}

func (h *SSEHandler) sendPing(w *bufio.Writer) error {
	return h.sendEvent(w, "ping", []byte("{}"))
}
