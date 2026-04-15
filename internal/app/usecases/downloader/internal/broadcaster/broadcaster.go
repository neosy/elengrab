package broadcaster

import (
	"sync"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

const broadcastChannelBufferSize = 100

type client struct {
	key dtypes.EventKey
	ch  chan dto.BroadcastEvent
}

type Broadcaster struct {
	clients map[client]struct{}
	lock    sync.RWMutex
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients: make(map[client]struct{}),
	}
}

func (b *Broadcaster) Subscribe(key dtypes.EventKey) chan dto.BroadcastEvent {
	ch := make(chan dto.BroadcastEvent, broadcastChannelBufferSize)
	client := client{
		key: key,
		ch:  ch,
	}
	b.lock.Lock()
	b.clients[client] = struct{}{}
	b.lock.Unlock()
	return ch
}

func (b *Broadcaster) Unsubscribe(key dtypes.EventKey, ch chan dto.BroadcastEvent) {
	client := client{
		key: key,
		ch:  ch,
	}
	b.lock.Lock()
	delete(b.clients, client)
	close(ch)
	b.lock.Unlock()
}

func (b *Broadcaster) Broadcast(eventType dto.BroadcastEventType, data any) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	for client := range b.clients {
		select {
		case client.ch <- dto.BroadcastEvent{
			Type: eventType,
			Data: data,
		}:
			// sent successfully
		default:
			// skip if buffer is full
		}
	}
}

func (b *Broadcaster) BroadcastTo(key dtypes.EventKey, eventType dto.BroadcastEventType, data any) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	for client := range b.clients {
		if client.key.Type != key.Type || client.key.ID != key.ID {
			continue
		}
		select {
		case client.ch <- dto.BroadcastEvent{
			Type: eventType,
			Data: data,
		}:
			// sent successfully
		default:
			// skip if buffer is full
		}
	}
}

func (b *Broadcaster) BroadcastToUser(userID uuid.UUID, eventType dto.BroadcastEventType, data any) {
	b.BroadcastTo(dtypes.NewEventKeyUserID(userID), eventType, data)
}

func (b *Broadcaster) BroadcastToAnonym(sessionID uuid.UUID, eventType dto.BroadcastEventType, data any) {
	b.BroadcastTo(dtypes.NewEventKeyAnonSessionID(sessionID), eventType, data)
}
