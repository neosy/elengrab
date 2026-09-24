package persistence

import (
	"context"

	"github.com/google/uuid"
	dmedia "github.com/neosy/elengrab/internal/domain/media"
	"github.com/neosy/elengrab/internal/pkg/cache/memory"
	memsimple "github.com/neosy/elengrab/internal/pkg/cache/memory/simple"
)

type ChannelRepositoryFactory func() ChannelRepository

type ChannelRepository interface {
	Transactional

	Insert(ctx context.Context, channel *dmedia.Channel) error
	Update(ctx context.Context, channel *dmedia.Channel) error
	Save(ctx context.Context, channel *dmedia.Channel) error
	UpdateChannelID(ctx context.Context, oldChannelID, newChannelID uuid.UUID) error

	FindByChannelID(ctx context.Context, channelID uuid.UUID) (*dmedia.Channel, error)
	ExistsByChannelID(ctx context.Context, channelID uuid.UUID) (bool, error)
	FindByExternalChannelID(
		ctx context.Context,
		externalID string,
		platform string,
	) (*dmedia.Channel, error)
	ExistsByExternalChannelID(
		ctx context.Context,
		externalID string,
		platform string,
	) (bool, error)

	GetIDs(ctx context.Context) ([]uuid.UUID, error)
	IterateAll(ctx context.Context, fn func(*dmedia.Channel) error) error
}

type ChannelCacheRepository interface {
	memory.CacheRepository

	Save(ctx context.Context, channel *dmedia.Channel) error
	SaveNegative(ctx context.Context, channelID uuid.UUID) error
	SaveNegativeByExternalChannelID(
		ctx context.Context,
		externalID string,
		platform string,
	) error
	Delete(ctx context.Context, channelID uuid.UUID) error

	FindByChannelID(
		ctx context.Context,
		channelID uuid.UUID,
	) (*dmedia.Channel, memsimple.CacheStatus, error)
	ExistsByChannelID(ctx context.Context, channelID uuid.UUID) (bool, error)
	FindByExternalChannelID(
		ctx context.Context,
		externalID string,
		platform string,
	) (*dmedia.Channel, memsimple.CacheStatus, error)
	ExistsByExternalChannelID(
		ctx context.Context,
		externalID string,
		platform string,
	) (bool, error)

	CleanExpired(ctx context.Context) error
}
