package mappers

import (
	"github.com/google/uuid"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/types"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
)

func (m *Mappers) MapQueryFiltersDomainToFilters(filters *dtypes.QueryFilters) *types.QueryFilters {
	if filters == nil {
		return nil
	}

	newFilters := types.NewQueryFilters()

	for _, filter := range filters.List() {
		queryKey := qkeys.Keys.FindByFilterName(filter.Name)
		if queryKey.IsZero() {
			continue
		}

		switch v := filter.Value().(type) {
		case string:
			newFilters.Add(queryKey, v)
		case uuid.UUID:
			newFilters.Add(queryKey, idcodec.EncodeUUIDBase64URL(v))
		}
	}

	return newFilters
}

func (m *Mappers) MapQueryFiltersToDomain(filters *types.QueryFilters) (*dtypes.QueryFilters, error) {
	if filters == nil || filters.Len() == 0 {
		return nil, nil
	}

	newFilter := dtypes.NewQueryFilters()

	for _, filter := range filters.List() {
		filterName := filter.Key.FilterName()
		if filterName == dtypes.QueryFilterNameNone {
			continue
		}

		switch filter.Key {
		case qkeys.ChannelIDKey:
			id, err := idcodec.DecodeUUIDBase64URL(filter.Value)
			if err != nil {
				return nil, err
			}
			newFilter.Add(filterName, id)
		case qkeys.SearchKey:
			fallthrough
		case qkeys.SearchQueryKey:
			newFilter.Add(filterName, filter.Value)
		}
	}

	if newFilter.Len() == 0 {
		return nil, nil
	}

	return newFilter, nil
}
