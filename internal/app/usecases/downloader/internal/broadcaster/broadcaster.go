package broadcaster

import (
	"context"
	"sync"

	"uuid"

	"github.com/neosy/elengrab/internal/app/usecases/downloader/internal/authz"
	"github.com/neosy/elengrab/internal/app/usecases/types"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dconsts "github.com/neosy/elengrab/internal/domain/consts"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/domain/types/broadcast"
	eventkey "github.com/neosy/elengrab/internal/domain/types/event_key"
	dsse "github.com/neosy/elengrab/internal/domain/types/sse"
	pservices "github.com/neosy/elengrab/internal/ports/services"
)

type Broadcaster struct {
	clients map[broadcast.ClientKey]subscription
	lock    sync.RWMutex

	auth  pservices.AuthUsecase
	authz *authz.Authorization
}

func NewBroadcaster(auth pservices.AuthUsecase, authz *authz.Authorization) *Broadcaster {
	return &Broadcaster{
		clients: make(map[broadcast.ClientKey]subscription),
		auth:    auth,
		authz:   authz,
	}
}

func (b *Broadcaster) Subscribe(clientKey broadcast.ClientKey, userID uuid.UUID) subscription {
	ch := make(chan dsse.BroadcastEvent, dconsts.SSEBroadcastChannelBufferSize)

	subscription := subscription{
		connectionID: clientKey.ConnectionID,
		userID:       userID,
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

func (b *Broadcaster) roleIDsByUserID(ctx context.Context, userID uuid.UUID) (dtypes.UserRoleIDs, error) {
	if userID == dauth.AnonymousUserID() {
		return nil, nil
	}

	user, err := b.auth.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	roleIDs, err := user.ParseRoleIDs()
	if err != nil {
		return nil, err
	}

	return roleIDs, nil
}

func (b *Broadcaster) Broadcast(eventType types.BroadcastEventType, data any) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	event := dsse.NewBroadcastEvent(eventType, data)

	for _, subscription := range b.clients {
		select {
		case subscription.eventCh <- event:
			// sent successfully
		default:
			// skip if buffer is full
		}
	}
}

func (b *Broadcaster) BroadcastPublic(ctx context.Context, key eventkey.EventKey, eventType types.BroadcastEventType, data any) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	event := dsse.NewBroadcastEvent(eventType, data)

	for client, subscription := range b.clients {
		allowed := client.EventKey.IsEqual(key)

		roleIDs, err := b.roleIDsByUserID(ctx, subscription.userID)
		if err != nil {
			continue
		}

		if !allowed {
			allowed = b.authz.HasPublicViewAccess(roleIDs)
		}

		if !allowed {
			allowed = b.authz.HasViewAllAccess(roleIDs)
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
func (b *Broadcaster) BroadcastToViewAllAccess(
	ctx context.Context,
	excludeKey eventkey.EventKey,
	eventType types.BroadcastEventType,
	data any,
) {
	b.BroadcastToViewAllAccessWithDataBuilder(
		ctx, excludeKey, eventType,
		func(_ uuid.UUID, _ dtypes.UserRoleIDs) (any, bool) {
			return data, true
		},
	)
}

// BroadcastToViewAllAccessWithDataBuilder broadcasts an event to clients with view-all access,
// excluding the specified event key and building event data separately for each user.
func (b *Broadcaster) BroadcastToViewAllAccessWithDataBuilder(
	ctx context.Context,
	excludeKey eventkey.EventKey,
	eventType types.BroadcastEventType,
	dataBuilder func(userID uuid.UUID, roleIDs dtypes.UserRoleIDs) (any, bool),
) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	event := dsse.NewBroadcastEvent(eventType, nil)

	for client, subscription := range b.clients {
		if client.EventKey.IsEqual(excludeKey) {
			continue
		}

		roleIDs, err := b.roleIDsByUserID(ctx, subscription.userID)
		if err != nil {
			continue
		}

		if !b.authz.HasViewAllAccess(roleIDs) {
			continue
		}

		data, ok := dataBuilder(subscription.userID, roleIDs)
		if !ok {
			continue
		}

		event.Data = data

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

	event := dsse.NewBroadcastEvent(eventType, data)

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
func (b *Broadcaster) BroadcastByAccess(
	ctx context.Context,
	key eventkey.EventKey,
	eventType types.BroadcastEventType,
	data any,
) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	event := dsse.NewBroadcastEvent(eventType, data)

	for client, subscription := range b.clients {
		allowed := client.EventKey.IsEqual(key)

		roleIDs, err := b.roleIDsByUserID(ctx, subscription.userID)
		if err != nil {
			continue
		}

		if !allowed {
			allowed = b.authz.HasViewAllAccess(roleIDs)
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

func (b *Broadcaster) BroadcastToUsersWithAccess(
	ctx context.Context,
	userID uuid.UUID,
	eventType types.BroadcastEventType,
	data any,
) {
	if dauth.IsNotRegisteredUserID(userID) {
		return
	}
	b.BroadcastByAccess(ctx, eventkey.NewEventKeyUserID(userID), eventType, data)
}
