package mappers

import (
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	dlink "github.com/neosy/elengrab/internal/domain/link"
)

// MapShortLinkClickRequestToDomain converting dto ShortLinkClickRequest to domain LinkClick
func (m *Mappers) MapShortLinkClickRequestToDomain(req *dto.ShortLinkClickRequest) *dlink.LinkClick {
	return &dlink.LinkClick{
		ClickedBy: req.ClickedBy,
		IPAddress: req.IPAddress,
		ShortURL:  req.ShortURL,
		UserAgent: req.UserAgent,
		Referrer:  req.Referrer,
	}
}
