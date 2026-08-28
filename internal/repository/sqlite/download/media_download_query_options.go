package download

import (
	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/ports/persistence"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
)

type mediaDownloadQueryOptions struct {
	dtypes.QueryMediaOptions

	includeDeleted bool
	statuses       []dtypes.MediaDownloadStatus

	partialHash **string
}

func (r *MediaDownloadRepository) WithOptions(options dtypes.QueryMediaOptions) persistence.MediaDownloadRepository {
	if options.Before != nil {
		r.queryOptions.Before = options.Before
	}

	if options.Limit != nil {
		r.queryOptions.Limit = options.Limit
	}

	return r
}

func (r *MediaDownloadRepository) WithStatus(statuses ...dtypes.MediaDownloadStatus) persistence.MediaDownloadRepository {
	if len(statuses) == 0 {
		return r
	}

	r.queryOptions.statuses = statuses

	return r
}

func (r *MediaDownloadRepository) WithDeleted() persistence.MediaDownloadRepository {
	r.queryOptions.includeDeleted = true
	return r
}

func (r *MediaDownloadRepository) WithUser(userID uuid.UUID) persistence.MediaDownloadRepository {
	var eDownload edownload.MediaDownload

	if userID == uuid.Nil {
		r.queryOptions.Visibility = new(dtypes.QueryMediaVisibilityPublic)
	} else {
		r.filtersByName[eDownload.FieldName(&eDownload.UserID)] = userID
	}

	return r
}

func (r *MediaDownloadRepository) WithFilters(filters map[string]any) persistence.MediaDownloadRepository {
	var (
		eDownload edownload.MediaDownload

		fieldNameByAllowedFilter = map[string]string{
			dtypes.QueryFilterNameUserID: eDownload.FieldName(&eDownload.UserID),
			dtypes.QueryFilterNameTitle:  eDownload.FieldName(&eDownload.MediaTitleLower),
		}
	)

	for name, value := range filters {
		fieldName, exists := fieldNameByAllowedFilter[name]
		if exists {
			switch fieldName {
			case eDownload.FieldName(&eDownload.UserID):
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
