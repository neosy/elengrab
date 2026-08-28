package searchindex

import (
	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/ports/persistence"
	esearchindex "github.com/neosy/elengrab/internal/repository/sqlite/search_index/entity"
)

type queryOptions struct {
	dtypes.QueryMediaOptions

	includeDeleted bool
}

func (r *MediaSourceIndexRepository) WithOptions(options dtypes.QueryMediaOptions) persistence.MediaSourceIndexRepository {
	if options.Before != nil {
		r.queryOptions.Before = options.Before
	}

	if options.Limit != nil {
		r.queryOptions.Limit = options.Limit
	}

	if options.Visibility != nil {
		r.queryOptions.Visibility = options.Visibility
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
		r.filtersByName[eIndex.FieldName(&eIndex.UserID)] = userID
	}

	return r
}

func (r *MediaSourceIndexRepository) WithFilters(filters map[string]any) persistence.MediaSourceIndexRepository {
	var (
		eIndex esearchindex.MediaSourceIndex

		fieldNameByAllowedFilter = map[string]string{
			dtypes.QueryFilterNameUserID: eIndex.FieldName(&eIndex.UserID),
			dtypes.QueryFilterNameTitle:  eIndex.FieldName(&eIndex.TitleLower),
		}
	)

	for name, value := range filters {
		fieldName, exists := fieldNameByAllowedFilter[name]
		if exists {
			switch fieldName {
			case eIndex.FieldName(&eIndex.UserID):
				id, ok := value.(uuid.UUID)
				if !ok {
					continue
				}
				if id == uuid.Nil {
					r.queryOptions.Visibility = new(dtypes.QueryMediaVisibilityPublic)
					continue
				}
				r.filtersByName[fieldName] = value
			default:
				r.filtersByName[fieldName] = value
			}

		}
	}

	return r
}
