package broadcaster

import (
	"sync"
	"time"

	"uuid"

	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/authz"
	"github.com/neosy/elengrab/internal/app/usecases/types"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/domain/types/broadcast"
	eventkey "github.com/neosy/elengrab/internal/domain/types/event_key"
	dsse "github.com/neosy/elengrab/internal/domain/types/sse"
)

const broadcastChannelBufferSize = 100

type Broadcaster struct {
	clients map[broadcast.ClientKey]subscription
	lock    sync.RWMutex
	authz   *authz.Authorization
}

func NewBroadcaster(authz *authz.Authorization) *Broadcaster {
	return &Broadcaster{
		clients: make(map[broadcast.ClientKey]subscription),
		authz:   authz,
	}
}

func (b *Broadcaster) Subscribe(clientKey broadcast.ClientKey, roles dtypes.UserRoleIDs) subscription {
	ch := make(chan dsse.BroadcastEvent, broadcastChannelBufferSize)

	subscription := subscription{
		connectionID: clientKey.ConnectionID,
		roles:        roles,
		eventCh:      ch,
	}

	b.lock.Lock()
	b.clients[clientKey] = subscription
	b.lock.Unlock()

	return subscription
}

func (b *Broadcaster) Unsubscribe(clientKey broadcast.ClientKey) {
	b.lock.Lock()
	defer b.lock.Unlock()

	data, ok := b.clients[clientKey]
	if !ok {
		return
	}

	delete(b.clients, clientKey)
	close(data.eventCh)
}

func (b *Broadcaster) Broadcast(eventType types.BroadcastEventType, data any) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	event := dsse.BroadcastEvent{
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now().UnixMicro(),
	}

	for _, subscription := range b.clients {
		select {
		case subscription.eventCh <- event:
			// sent successfully
		default:
			// skip if buffer is full
		}
	}
}

func (b *Broadcaster) BroadcastPublic(key eventkey.EventKey, eventType types.BroadcastEventType, data any) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	event := dsse.BroadcastEvent{
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now().UnixMicro(),
	}

	for client, subscription := range b.clients {
		allowed := client.EventKey.IsEqual(key)

		if !allowed {
			allowed = b.authz.HasPublicViewAccess(subscription.roles)
		}

		if !allowed {
			allowed = b.authz.HasViewAllAccess(subscription.roles)
		}

		if !allowed {
			continue
		}

		select {
		case subscription.eventCh <- event:
			// sent successfully
		default:
			// skip if buffer is full
		}
	}
}

// BroadcastToViewAllAccess broadcasts an event to clients with view-all access, excluding the specified event key.
func (b *Broadcaster) BroadcastToViewAllAccess(excludeKey eventkey.EventKey, eventType types.BroadcastEventType, data any) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	event := dsse.BroadcastEvent{
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now().UnixMicro(),
	}

	for client, subscription := range b.clients {
		if client.EventKey.IsEqual(excludeKey) {
			continue
		}

		if !b.authz.HasViewAllAccess(subscription.roles) {
			continue
		}

		select {
		case subscription.eventCh <- event:
			// sent successfully
		default:
			// skip if buffer is full
		}
	}
}

func (b *Broadcaster) BroadcastByKey(key eventkey.EventKey, eventType types.BroadcastEventType, data any) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	event := dsse.BroadcastEvent{
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now().UnixMicro(),
	}

	for client, subscription := range b.clients {
		allowed := client.EventKey.IsEqual(key)
		if !allowed {
			continue
		}

		select {
		case subscription.eventCh <- event:
			// sent successfully
		default:
			// skip if buffer is full
		}
	}
}

// BroadcastByAccess broadcast to all clients that have access to the event
func (b *Broadcaster) BroadcastByAccess(key eventkey.EventKey, eventType types.BroadcastEventType, data any) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	event := dsse.BroadcastEvent{
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now().UnixMicro(),
	}

	for client, subscription := range b.clients {
		allowed := client.EventKey.IsEqual(key)

		if !allowed {
			allowed = b.authz.HasViewAllAccess(subscription.roles)
		}

		if !allowed {
			continue
		}

		select {
		case subscription.eventCh <- event:
			// sent successfully
		default:
			// skip if buffer is full
		}
	}
}

// BroadcastToAuthorized sends the event to all clients that have the given event key
func (b *Broadcaster) BroadcastToAuthorized(authCtx dauth.AuthContext, eventType types.BroadcastEventType, data any) {
	if authCtx.IsAnonymousWithoutSession() {
		return
	}
	b.BroadcastByKey(authCtx.EventKey(), eventType, data)
}

func (b *Broadcaster) BroadcastToUser(userID uuid.UUID, eventType types.BroadcastEventType, data any) {
	if dauth.IsNotRegisteredUserID(userID) {
		return
	}
	b.BroadcastByKey(eventkey.NewEventKeyUserID(userID), eventType, data)
}

func (b *Broadcaster) BroadcastToSession(sessionID uuid.UUID, eventType types.BroadcastEventType, data any) {
	if sessionID == uuid.Nil() {
		return
	}
	b.BroadcastByKey(eventkey.NewEventKeySessionID(sessionID), eventType, data)
}

func (b *Broadcaster) BroadcastToUsersWithAccess(userID uuid.UUID, eventType types.BroadcastEventType, data any) {
	if dauth.IsNotRegisteredUserID(userID) {
		return
	}
	b.BroadcastByAccess(eventkey.NewEventKeyUserID(userID), eventType, data)
}
