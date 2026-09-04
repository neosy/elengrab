package dto

import (
	"maps"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type MediaDownloadFilters struct {
	Search string
}

type MediaDownloadQuery struct {
	ViewMode dtypes.QueryMediaViewMode
	Limit    uint64

	LastRecord dtypes.QueryMediaDownloadCursor

	Filters dtypes.QueryFiltersByName
}

func MediaDownloadQueryDefault(limit uint64) MediaDownloadQuery {
	query := MediaDownloadQuery{
		ViewMode: dtypes.QueryMediaViewModeDefault,
		Limit:    limit,
	}

	query.Normalize()

	return query
}

func BuildMediaDownloadQuery(query MediaDownloadQuery) MediaDownloadQuery {
	return query.Copy()
}

func (q MediaDownloadQuery) Copy() MediaDownloadQuery {
	query := q

	query.Filters = maps.Clone(query.Filters)

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
