package searchindex

import (
	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/ports/persistence"
	esearchindex "github.com/neosy/elengrab/internal/repository/sqlite/search_index/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/types"
)

type queryOptions struct {
	types.QueryMediaOptions

	includeDeleted bool
}

func newQueryOptions() queryOptions {
	return queryOptions{
		QueryMediaOptions: types.NewQueryMediaOptions(),
	}
}

func (r *MediaSourceIndexRepository) WithOptions(options dtypes.QueryMediaOptions) persistence.MediaSourceIndexRepository {
	if options.Before != nil {
		r.queryOptions.Before = options.Before
	}

	if options.Offset != nil {
		r.queryOptions.Offset = options.Offset
	}

	if options.Limit != nil {
		r.queryOptions.Limit = options.Limit
	}

	if options.Visibility != nil {
		r.queryOptions.Visibility = options.Visibility
	}

	if len(options.OrderBys) > 0 {
		r.WithOrderBy(options.OrderBys...)
	}

	if len(options.Filters) > 0 {
		r.WithFilters(options.Filters...)
	}

	return r
}

func (r *MediaSourceIndexRepository) WithDeleted() persistence.MediaSourceIndexRepository {
	r.queryOptions.includeDeleted = true
	return r
}

func (r *MediaSourceIndexRepository) WithUser(userID uuid.UUID) persistence.MediaSourceIndexRepository {
	var eIndex esearchindex.MediaSourceIndex

	if userID == uuid.Nil {
		r.queryOptions.Visibility = new(dtypes.QueryMediaVisibilityPublic)
	} else {
		r.queryOptions.Filters.Add(eIndex.FieldName(&eIndex.UserID), userID)
	}

	return r
}

func (r *MediaSourceIndexRepository) WithFilters(filters ...dtypes.QueryFilter) persistence.MediaSourceIndexRepository {
	var (
		eIndex esearchindex.MediaSourceIndex

		fieldNameByAllowedFilter = eIndex.PaginationFieldNames()
	)

	for _, filter := range filters {
		fieldName, exists := fieldNameByAllowedFilter[filter.Name.String()]
		if !exists {
			continue
		}

		switch fieldName {
		case eIndex.FieldName(&eIndex.UserID):
			id, ok := filter.Value.(uuid.UUID)
			if !ok {
				continue
			}
			if id == uuid.Nil {
				r.queryOptions.Visibility = new(dtypes.QueryMediaVisibilityPublic)
				continue
			}
			r.queryOptions.Filters.Add(fieldName, filter.Value)
		default:
			r.queryOptions.Filters.Add(fieldName, filter.Value)
		}
	}

	return r
}

func (r *MediaSourceIndexRepository) WithOrderBy(list ...dtypes.QueryOrderBy) persistence.MediaSourceIndexRepository {
	var (
		eIndex esearchindex.MediaSourceIndex

		fieldNameByAllowedFilter = eIndex.PaginationFieldNames()
	)

	for _, orderBy := range list {
		fieldName, exists := fieldNameByAllowedFilter[orderBy.Field]
		if !exists {
			continue
		}

		r.queryOptions.OrderBys.Add(fieldName, orderBy.Direction)
	}

	return r
}
