package helper

import (
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/types"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func mapSearchQueryToSearchValues(queryItems *types.QueryFilters) (types.SearchValues, error) {
	if queryItems == nil {
		return types.SearchValues{}, nil
	}

	searchValues := types.NewSearchValues()

	filters := types.NewQueryFilters()
	parmFilters := types.NewQueryFilters()

	for _, item := range queryItems.List() {
		switch item.Key {
		case qkeys.SearchKey:
			fallthrough
		case qkeys.SearchQueryKey:
			searchValues.QueryText = item.Value
		case qkeys.UserIDKey:
			filters.Add(item.Key, item.Value)
		case qkeys.UserNameKey:
			filters.Add(qkeys.UserLoginKey, item.Value)
		case qkeys.UserLoginKey:
			filters.Add(item.Key, item.Value)
		case qkeys.ChannelIDKey:
			filters.Add(item.Key, item.Value)
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
			if qkeys.SearchParameterFilterKeys.ExistsByKey(item.Key) {
				parmFilters.Add(item.Key, item.Value)
			}
		}
	}

	if filters.Len() != 0 {
		searchValues.Filters = filters
	}

	if parmFilters.Len() != 0 {
		if searchValues.Parameters.Filters == nil {
			searchValues.Parameters.Filters = parmFilters
		} else {
			for _, filter := range parmFilters.List() {
				searchValues.Parameters.Filters.Append(filter)
			}
		}
	}

	return searchValues, nil
}
