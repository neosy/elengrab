package types

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/consts"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type SearchParameterValues struct {
	ViewMode   string
	Filters    *QueryFilters
	LastRecord dtypes.QueryMediaDownloadCursor
}

func (v SearchParameterValues) BuildUsecaseQuery(filters *dtypes.QueryFilters) udto.MediaDownloadQuery {
	viewMode := dtypes.QueryMediaViewModeDefault

	mode, err := dtypes.ParseQueryMediaViewMode(v.ViewMode)
	if err == nil {
		viewMode = mode
	}

	return udto.MediaDownloadQuery{
		ViewMode:   viewMode,
		Limit:      consts.LoadHistoryLimit,
		LastRecord: v.LastRecord,
		Filters:    filters.Clone(),
	}
}

func (v SearchParameterValues) FilterValuesByKey() map[qkeys.QueryKey]string {
	if v.Filters.Len() == 0 {
		return nil
	}

	queryFilters := make(map[qkeys.QueryKey]string)

	searchKeys := qkeys.SearchKeys()

	for _, filter := range v.Filters.List() {
		if searchKeys.ExistsByKey(filter.Key) {
			queryFilters[filter.Key] = filter.Value
		}
	}

	return queryFilters
}
