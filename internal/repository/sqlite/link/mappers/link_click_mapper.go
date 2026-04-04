package mappers

import (
	dlink "github.com/neosy/elengrab/internal/domain/link"
	elink "github.com/neosy/elengrab/internal/repository/sqlite/link/entity"
)

func (m *Mappers) MapLinkClickDomainToEntity(linkClick *dlink.LinkClick) (*elink.LinkClick, error) {
	return &elink.LinkClick{
		LinkClickID: linkClick.LinkClickID,
		LinkID:      linkClick.LinkID,
		IPAddress:   linkClick.IPAddress,
		ShortURL:    linkClick.ShortURL,
		ClickedBy:   linkClick.ClickedBy,
		ClickedAt:   linkClick.ClickedAt,
		UserAgent:   linkClick.UserAgent,
		Referrer:    linkClick.Referrer,
	}, nil
}

func (m *Mappers) MapLinkClickEntityToDomain(linkClick *elink.LinkClick) (*dlink.LinkClick, error) {
	return &dlink.LinkClick{
		LinkClickID: linkClick.LinkClickID,
		LinkID:      linkClick.LinkID,
		IPAddress:   linkClick.IPAddress,
		ShortURL:    linkClick.ShortURL,
		ClickedBy:   linkClick.ClickedBy,
		ClickedAt:   linkClick.ClickedAt,
		UserAgent:   linkClick.UserAgent,
		Referrer:    linkClick.Referrer,
		CreatedAt:   linkClick.CreatedAt,
	}, nil
}
