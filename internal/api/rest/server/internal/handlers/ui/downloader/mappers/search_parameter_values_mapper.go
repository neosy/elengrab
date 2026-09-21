package mappers

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/consts"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/types"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (m *Mappers) MapSearchParameterValuesToUsecaseQuery(
	values *types.SearchParameterValues,
) (udto.MediaDownloadQuery, error) {
	query := udto.MediaDownloadQuery{
		ViewMode:   dtypes.QueryMediaViewModeDefault,
		Limit:      consts.LoadHistoryLimit,
		LastRecord: values.LastRecord,
	}

	if values.ViewMode != "" {
		var err error
		query.ViewMode, err = dtypes.ParseQueryMediaViewMode(values.ViewMode)
		if err != nil {
			return query, nil
		}
	}

	filters, err := m.MapQueryFiltersToDomain(values.Filters)
	if err != nil {
		return query, nil
	}

	query.Filters = filters

	return query, nil
}
