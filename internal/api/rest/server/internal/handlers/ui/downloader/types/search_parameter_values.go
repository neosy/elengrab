package types

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/consts"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type SearchParameterValues struct {
	ViewMode   dtypes.QueryMediaViewMode
	Filters    *QueryFilters
	LastRecord dtypes.QueryMediaDownloadCursor
}

func NewSearchParameterValues() SearchParameterValues {
	return SearchParameterValues{
		ViewMode: dtypes.QueryMediaViewModeDefault,
	}
}

func (v *SearchParameterValues) IsZero() bool {
	if v == nil {
		return true
	}

	return v.ViewMode == dtypes.QueryMediaViewModeNone &&
		v.Filters.Len() == 0
}

func (v SearchParameterValues) BuildUsecaseQuery(filters *dtypes.QueryFilters) udto.MediaDownloadQuery {
	viewMode := dtypes.QueryMediaViewModeDefault

	if v.ViewMode != dtypes.QueryMediaViewModeNone {
		viewMode = v.ViewMode
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

	for _, filter := range v.Filters.List() {
		if qkeys.SearchFilterKeys.ExistsByKey(filter.Key) {
			queryFilters[filter.Key] = filter.Value
		}
	}

	return queryFilters
}
