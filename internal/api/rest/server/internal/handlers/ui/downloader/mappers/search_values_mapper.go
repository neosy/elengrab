package mappers

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/consts"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/types"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (m *Mappers) MapSearchValuesToUsecaseQuery(values types.SearchValues) (udto.MediaDownloadQuery, error) {
	query := udto.NewMediaDownloadQuery(
		values.Parameters.ViewMode,
		consts.LoadHistoryLimit,
		values.Parameters.LastRecord,
	)

	qFilters := values.Parameters.Filters.Clone()

	if qFilters == nil {
		qFilters = types.NewQueryFilters()
	}

	if values.QueryText != "" {
		qFilters.Add(qkeys.SearchQueryKey, values.QueryText)
	}

	filters, err := m.MapQueryFiltersToDomain(qFilters)
	if err != nil {
		return query, err
	}

	query.Filters = filters

	return query, nil
}

func (m *Mappers) MapSearchQueryToSearchValues(queryItems *types.QueryFilters) (types.SearchValues, error) {
	if queryItems == nil {
		return types.SearchValues{}, nil
	}

	searchValues := types.NewSearchValues()
	filters := types.NewQueryFilters()

	for _, item := range queryItems.List() {
		switch item.Key {
		case qkeys.SearchKey:
			fallthrough
		case qkeys.SearchQueryKey:
			searchValues.QueryText = item.Value
		case qkeys.SearchParametersKey:
			parameters, err := types.ParseSearchEncodeQueryStringToParameters(item.Value)
			if err != nil {
				return types.SearchValues{}, err
			}

			parmValues, err := parameters.ParseValues()
			if err != nil {
				return types.SearchValues{}, err
			}

			searchValues.Parameters = parmValues
		case qkeys.ViewModeKey:
			viewMode, err := dtypes.ParseQueryMediaViewMode(item.Value)
			if err != nil {
				return types.SearchValues{}, err
			}
			searchValues.Parameters.ViewMode = viewMode
		default:
			if qkeys.SearchFilterKeys.ExistsByKey(item.Key) {
				filters.Add(item.Key, item.Value)
			}
		}
	}

	if filters.Len() != 0 {
		if searchValues.Parameters.Filters == nil {
			searchValues.Parameters.Filters = filters
		} else {
			for _, filter := range filters.List() {
				searchValues.Parameters.Filters.Append(filter)
			}
		}
	}

	return searchValues, nil
}
