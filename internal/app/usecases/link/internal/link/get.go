package linklink

import (
	"context"

	"github.com/google/uuid"
	dlink "github.com/neosy/elengrab/internal/domain/link"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (u *Link) FindByLinkID(ctx context.Context, linkID uuid.UUID) (*dlink.Link, error) {
	if linkID == uuid.Nil {
		return nil, nil
	}

	link, err := u.linkRep.Find(ctx, linkID)
	if err != nil {
		u.logger.Warn("Failed get link", "error", err)
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	return link, nil
}

func (u *Link) GetByLinkID(ctx context.Context, linkID uuid.UUID) (*dlink.Link, error) {
	link, err := u.FindByLinkID(ctx, linkID)
	if err != nil {
		return nil, err
	}

	if link == nil {
		u.logger.Debug("Link not found", "linkID", linkID)
		return nil, errorx.New("link not found", exceptionx.NOT_FOUND)
	}

	return link, nil
}

func (u *Link) FindLastByShortCode(ctx context.Context, shortCode string) (*dlink.Link, error) {
	if shortCode == "" {
		return nil, nil
	}

	link, err := u.linkRep.FindLastByShortCode(ctx, shortCode)
	if err != nil {
		u.logger.Warn("Failed get link", "shortCode", shortCode, "error", err)
		return nil, errorx.NewFromError(err, exceptionx.ERROR)
	}

	return link, nil
}

func (u *Link) GetLastByShortCode(ctx context.Context, shortCode string) (*dlink.Link, error) {
	link, err := u.FindLastByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	if link == nil {
		u.logger.Debug("Link not found", "shortCode", shortCode)
		return nil, errorx.New("link not found", exceptionx.NOT_FOUND)
	}

	return link, nil
}
