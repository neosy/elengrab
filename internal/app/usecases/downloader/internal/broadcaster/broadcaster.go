package broadcaster

import (
	"sync"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
)

const broadcastChannelBufferSize = 100

type client struct {
	userID string
	ch     chan dto.BroadcastEvent
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

func (b *Broadcaster) Subscribe(userID string) chan dto.BroadcastEvent {
	ch := make(chan dto.BroadcastEvent, broadcastChannelBufferSize)
	client := client{
		userID: userID,
		ch:     ch,
	}
	b.lock.Lock()
	b.clients[client] = struct{}{}
	b.lock.Unlock()
	return ch
}

func (b *Broadcaster) Unsubscribe(userID string, ch chan dto.BroadcastEvent) {
	client := client{
		userID: userID,
		ch:     ch,
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

func (b *Broadcaster) BroadcastToUser(userID uuid.UUID, eventType dto.BroadcastEventType, data any) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	for client := range b.clients {
		if client.userID != userID.String() {
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
