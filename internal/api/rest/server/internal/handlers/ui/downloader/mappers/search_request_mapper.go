package mappers

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/types"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (m *Mappers) MapSearchRequestToSearchParameterValues(req dto.SearchRequest) (*types.SearchParameterValues, error) {
	viewMode := dtypes.QueryMediaViewModeDefault
	if req.ViewMode != "" {
		var err error
		viewMode, err = dtypes.ParseQueryMediaViewMode(req.ViewMode)
		if err != nil {
			return nil, err
		}
	}

	paramValues := &types.SearchParameterValues{
		ViewMode: viewMode.String(),
	}

	filters := types.NewQueryFilters()

	if queryText := dtypes.SearchText(req.Query).Normalize(); queryText.IsLongEnough() {
		if err := queryText.Validate(); err != nil {
			return nil, err
		}
		filters.Add(qkeys.SearchQueryKey, queryText.String())
	}

	if req.ChannelID != "" {
		filters.Add(qkeys.ChannelIDKey, req.ChannelID)
	}

	if filters.Len() != 0 {
		paramValues.Filters = filters
	}

	return paramValues, nil
}

func (m *Mappers) MapSearchRequestToUsecaseQuery(req dto.SearchRequest) (udto.MediaDownloadQuery, error) {
	parmValues, err := m.MapSearchRequestToSearchParameterValues(req)
	if err != nil {
		return udto.MediaDownloadQuery{}, err
	}

	return m.MapSearchParameterValuesToUsecaseQuery(parmValues)
}
