package pservices

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dlink "github.com/neosy/elengrab/internal/domain/link"
)

type ShortLinkService interface {
	Create(ctx context.Context, req *dto.LinkCreateRequest) (*dlink.Link, error)
	SoftDelete(ctx context.Context, linkID uuid.UUID) error

	FindLastByShortCode(ctx context.Context, shortCode string) (*dlink.Link, error)
	GetLastByShortCode(ctx context.Context, shortCode string) (*dlink.Link, error)

	Click(ctx context.Context, req *dto.ShortLinkClickRequest) (*dlink.Link, error)

	GenerateShortCode(url string, length uint8, deterministic bool) string
	ResolveShortCode(ctx context.Context, shortCode string) (*dlink.Link, error)
}
