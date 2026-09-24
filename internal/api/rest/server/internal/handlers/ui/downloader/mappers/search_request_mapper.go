package mappers

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/types"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (m *Mappers) MapSearchRequestToSearchValues(
	req dto.SearchRequest,
) (types.SearchValues, error) {
	viewMode := dtypes.QueryMediaViewModeDefault
	if req.ViewMode != "" {
		var err error
		viewMode, err = dtypes.ParseQueryMediaViewMode(req.ViewMode)
		if err != nil {
			return types.SearchValues{}, err
		}
	}

	searchValues := types.NewSearchValues()
	searchValues.Parameters.ViewMode = viewMode

	filters := types.NewQueryFilters()

	if queryText := dtypes.SearchText(req.QueryText).Normalize(); queryText.IsLongEnough() {
		if err := queryText.Validate(); err != nil {
			return types.SearchValues{}, err
		}
		searchValues.QueryText = queryText.String()
	}

	if req.ChannelID != "" {
		filters.Add(qkeys.ChannelIDKey, req.ChannelID)
	}

	if filters.Len() != 0 {
		searchValues.Parameters.Filters = filters
	}

	return searchValues, nil
}

func (m *Mappers) MapSearchRequestToUsecaseQuery(req dto.SearchRequest) (udto.MediaDownloadQuery, error) {
	values, err := m.MapSearchRequestToSearchValues(req)
	if err != nil {
		return udto.MediaDownloadQuery{}, err
	}

	return m.MapSearchValuesToUsecaseQuery(values)
}
