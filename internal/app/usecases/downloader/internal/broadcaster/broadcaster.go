package broadcaster

import (
	"sync"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/authz"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/domain/types/broadcast"
	eventkey "github.com/neosy/elengrab/internal/domain/types/event_key"
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
	ch := make(chan dto.BroadcastEvent, broadcastChannelBufferSize)

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

func (b *Broadcaster) Broadcast(eventType dto.BroadcastEventType, data any) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	for _, subscription := range b.clients {
		select {
		case subscription.eventCh <- dto.BroadcastEvent{
			Type: eventType,
			Data: data,
		}:
			// sent successfully
		default:
			// skip if buffer is full
		}
	}
}

func (b *Broadcaster) BroadcastPublic(key eventkey.EventKey, eventType dto.BroadcastEventType, data any) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	for client, subscription := range b.clients {
		allowed := client.EventKey.Type == key.Type && client.EventKey.ID == key.ID

		if !allowed {
			allowed = b.authz.HasPublicViewAccess(subscription.roles)
		}

		if !allowed {
			allowed = b.authz.HasViewAllAccess(subscription.roles)
		}

		if !allowed {
			continue
		}

		event := dto.BroadcastEvent{
			Type: eventType,
			Data: data,
		}

		select {
		case subscription.eventCh <- event:
			// sent successfully
		default:
			// skip if buffer is full
		}
	}
}

func (b *Broadcaster) BroadcastByKey(key eventkey.EventKey, eventType dto.BroadcastEventType, data any) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	for client, subscription := range b.clients {
		allowed := client.EventKey.Type == key.Type && client.EventKey.ID == key.ID

		if !allowed {
			continue
		}

		event := dto.BroadcastEvent{
			Type: eventType,
			Data: data,
		}

		select {
		case subscription.eventCh <- event:
			// sent successfully
		default:
			// skip if buffer is full
		}
	}
}

func (b *Broadcaster) BroadcastByAccess(key eventkey.EventKey, eventType dto.BroadcastEventType, data any) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	for client, subscription := range b.clients {
		allowed := client.EventKey.Type == key.Type && client.EventKey.ID == key.ID

		if !allowed {
			allowed = b.authz.HasViewAllAccess(subscription.roles)
		}

		if !allowed {
			continue
		}

		event := dto.BroadcastEvent{
			Type: eventType,
			Data: data,
		}

		select {
		case subscription.eventCh <- event:
			// sent successfully
		default:
			// skip if buffer is full
		}
	}
}

func (b *Broadcaster) BroadcastToAuth(authCtx dauth.AuthContext, eventType dto.BroadcastEventType, data any) {
	if authCtx.UserID == uuid.Nil && authCtx.AnonSessionID == uuid.Nil {
		return
	}
	b.BroadcastByKey(authCtx.EventKey(), eventType, data)
}

func (b *Broadcaster) BroadcastToUser(userID uuid.UUID, eventType dto.BroadcastEventType, data any) {
	if userID == uuid.Nil {
		return
	}
	b.BroadcastByKey(eventkey.NewEventKeyUserID(userID), eventType, data)
}

func (b *Broadcaster) BroadcastToSession(sessionID uuid.UUID, eventType dto.BroadcastEventType, data any) {
	if sessionID == uuid.Nil {
		return
	}
	b.BroadcastByKey(eventkey.NewEventKeySessionID(sessionID), eventType, data)
}

func (b *Broadcaster) BroadcastToUsersWithAccess(userID uuid.UUID, eventType dto.BroadcastEventType, data any) {
	if userID == uuid.Nil {
		return
	}
	b.BroadcastByAccess(eventkey.NewEventKeyUserID(userID), eventType, data)
}
