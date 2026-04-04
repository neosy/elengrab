package mappers

import (
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dlink "github.com/neosy/elengrab/internal/domain/link"
)

func (m *Mappers) MapLinkCreateRequestDtoToDomain(
	req *dto.LinkCreateRequest,
) *dlink.Link {
	return &dlink.Link{
		OriginalURL:     req.OriginalURL,
		IsMatchShortURL: req.IsMatchShortURL,
		MaxClicks:       req.MaxClicks,
		AllowedUserIDs:  req.AllowedUserIDs,
		AllowedIPs:      req.AllowedIPs,
		ExpiresAt:       req.ExpiresAt,
	}
}
