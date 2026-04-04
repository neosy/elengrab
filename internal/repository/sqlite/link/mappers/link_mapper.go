package mappers

import (
	"encoding/json"
	"fmt"

	dlink "github.com/neosy/elengrab/internal/domain/link"
	"github.com/neosy/elengrab/internal/pkg/fnx"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
	elink "github.com/neosy/elengrab/internal/repository/sqlite/link/entity"
)

func (m *Mappers) MapLinkDomainToEntity(link *dlink.Link) (*elink.Link, error) {
	var allowedUserIDs *string
	if len(link.AllowedUserIDs) > 0 {
		data, err := json.Marshal(link.AllowedUserIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal allowedUserIDs: %w", err)
		}
		allowedUserIDs = uptr.Any(string(data))
	}

	var allowedIPs *string
	if len(link.AllowedIPs) > 0 {
		data, err := json.Marshal(link.AllowedIPs)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal allowedIPs: %w", err)
		}
		allowedIPs = uptr.Any(string(data))
	}

	var maxClicks *int
	if link.MaxClicks != nil {
		maxClicks = uptr.Any(int(*link.MaxClicks))
	}

	return &elink.Link{
		LinkID:          link.LinkID,
		OriginalURL:     link.OriginalURL,
		ShortCode:       link.ShortCode,
		ShortURL:        link.ShortURL,
		IsMatchShortURL: fnx.Ternary(link.IsMatchShortURL, 1, 0),
		MaxClicks:       maxClicks,
		AllowedUserIDs:  allowedUserIDs,
		AllowedIPs:      allowedIPs,
		ExpiresAt:       link.ExpiresAt,
	}, nil
}

func (m *Mappers) MapLinkEntityToDomain(link *elink.Link) (*dlink.Link, error) {
	var allowedUserIDs []string
	if link.AllowedUserIDs != nil {
		err := json.Unmarshal([]byte(*link.AllowedUserIDs), &allowedUserIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal allowedUserIDs: %w", err)
		}
	}

	var allowedIPs []string
	if link.AllowedIPs != nil {
		err := json.Unmarshal([]byte(*link.AllowedIPs), &allowedIPs)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal allowedIPs: %w", err)
		}
	}

	var maxClicks *uint16
	if link.MaxClicks != nil {
		maxClicks = uptr.Any(uint16(*link.MaxClicks))
	}

	return &dlink.Link{
		LinkID:          link.LinkID,
		OriginalURL:     link.OriginalURL,
		ShortCode:       link.ShortCode,
		ShortURL:        link.ShortURL,
		IsMatchShortURL: fnx.Ternary(link.IsMatchShortURL == 1, true, false),
		MaxClicks:       maxClicks,
		AllowedUserIDs:  allowedUserIDs,
		AllowedIPs:      allowedIPs,
		ExpiresAt:       link.ExpiresAt,
		CreatedAt:       link.CreatedAt,
		UpdatedAt:       link.UpdatedAt,
		DeletedAt:       link.DeletedAt,
	}, nil
}
