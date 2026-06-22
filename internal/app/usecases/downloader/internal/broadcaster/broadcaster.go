package broadcaster

import (
	"sync"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/authz"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

const broadcastChannelBufferSize = 100

type Broadcaster struct {
	clients map[clientKey]subscription
	lock    sync.RWMutex
	authz   *authz.Authorization
}

func NewBroadcaster(authz *authz.Authorization) *Broadcaster {
	return &Broadcaster{
		clients: make(map[clientKey]subscription),
		authz:   authz,
	}
}

func (b *Broadcaster) Subscribe(key dtypes.EventKey, roles dtypes.UserRoleIDs) subscription {
	ch := make(chan dto.BroadcastEvent, broadcastChannelBufferSize)
	connectionID := uuid.New()

	clientID := clientKey{
		connectionID: connectionID,
		eventKey:     key,
	}

	subscription := subscription{
		connectionID: connectionID,
		roles:        roles,
		eventCh:      ch,
	}

	b.lock.Lock()
	b.clients[clientID] = subscription
	b.lock.Unlock()

	return subscription
}

func (b *Broadcaster) Unsubscribe(key dtypes.EventKey) {
	clientID := clientKey{
		eventKey: key,
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	data, ok := b.clients[clientID]
	if !ok {
		return
	}

	delete(b.clients, clientID)
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

func (b *Broadcaster) BroadcastPublic(key dtypes.EventKey, eventType dto.BroadcastEventType, data any) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	for client, subscription := range b.clients {
		allowed := client.eventKey.Type == key.Type && client.eventKey.ID == key.ID

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

func (b *Broadcaster) BroadcastByKey(key dtypes.EventKey, eventType dto.BroadcastEventType, data any) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	for client, subscription := range b.clients {
		allowed := client.eventKey.Type == key.Type && client.eventKey.ID == key.ID

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

func (b *Broadcaster) BroadcastByAccess(key dtypes.EventKey, eventType dto.BroadcastEventType, data any) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	for client, subscription := range b.clients {
		allowed := client.eventKey.Type == key.Type && client.eventKey.ID == key.ID

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

func (b *Broadcaster) BroadcastToUser(userID uuid.UUID, eventType dto.BroadcastEventType, data any) {
	b.BroadcastByKey(dtypes.NewEventKeyUserID(userID), eventType, data)
}

func (b *Broadcaster) BroadcastToUsersWithAccess(userID uuid.UUID, eventType dto.BroadcastEventType, data any) {
	b.BroadcastByAccess(dtypes.NewEventKeyUserID(userID), eventType, data)
}
