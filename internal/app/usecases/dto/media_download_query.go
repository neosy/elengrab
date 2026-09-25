package dto

import (
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type MediaDownloadFilters struct {
	Search string
}

type MediaDownloadQuery struct {
	ViewMode dtypes.QueryMediaViewMode
	Limit    uint64

	LastRecord dtypes.QueryMediaDownloadCursor

	Filters *dtypes.QueryFilters
}

func NewMediaDownloadQuery(
	viewMode dtypes.QueryMediaViewMode,
	limit uint64,
	lastRecord dtypes.QueryMediaDownloadCursor,
) MediaDownloadQuery {
	vMode := dtypes.QueryMediaViewModeDefault
	if viewMode != dtypes.QueryMediaViewModeNone {
		vMode = viewMode
	}

	return MediaDownloadQuery{
		ViewMode:   vMode,
		Limit:      limit,
		LastRecord: lastRecord,
	}
}

func MediaDownloadQueryDefault(limit uint64) MediaDownloadQuery {
	query := MediaDownloadQuery{
		ViewMode: dtypes.QueryMediaViewModeDefault,
		Limit:    limit,
	}

	query.Normalize()

	return query
}

func (q MediaDownloadQuery) Clone() MediaDownloadQuery {
	query := q

	query.Filters = query.Filters.Clone()

	query.Normalize()

	return query
}

func (q *MediaDownloadQuery) Normalize() {
	if !q.ViewMode.Exists() {
		q.ViewMode = dtypes.QueryMediaViewModeDefault
	}

	if q.Limit == 0 {
		q.Limit = 20
	}
}

func (q *MediaDownloadQuery) GetSearchQueryString() string {
	if q.Filters.Len() == 0 {
		return ""
	}

	return q.Filters.GetStringValue(dtypes.QueryFilterNameSearchQuery)
}
