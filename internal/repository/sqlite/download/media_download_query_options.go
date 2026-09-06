package download

import (
	"strings"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/ports/persistence"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/types"
)

type mediaDownloadQueryOptions struct {
	types.QueryMediaOptions

	includeDeleted bool
	statuses       []dtypes.MediaDownloadStatus
	downloadIDs    []uuid.UUID

	partialHash **string
}

type queryArgs struct {
	Placeholder string
	Values      []any
}

func newMediaDownloadQueryOptions() mediaDownloadQueryOptions {
	return mediaDownloadQueryOptions{
		QueryMediaOptions: types.NewQueryMediaOptions(),
	}
}

func (o *mediaDownloadQueryOptions) downloadIDsQuery() queryArgs {
	placeholders := make([]string, len(o.downloadIDs))
	values := make([]any, len(o.downloadIDs))

	for i, id := range o.downloadIDs {
		placeholders[i] = "?"
		values[i] = id.String()
	}

	return queryArgs{
		Placeholder: strings.Join(placeholders, ", "),
		Values:      values,
	}
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
		r.queryOptions.Filters.Add(eDownload.FieldName(&eDownload.UserID), userID)
	}

	return r
}

func (r *MediaDownloadRepository) WithFilters(filters map[dtypes.QueryFilterName]any) persistence.MediaDownloadRepository {
	var (
		eDownload edownload.MediaDownload

		fieldNameByAllowedFilter = map[dtypes.QueryFilterName]string{
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
				r.queryOptions.Filters.Add(fieldName, value)
			default:
				r.queryOptions.Filters.Add(fieldName, value)
			}
			continue
		}

		switch name {
		case dtypes.QueryFilterNameDownloadIDs:
			ids, ok := value.([]uuid.UUID)
			if ok {
				r.queryOptions.downloadIDs = ids
			}
		}
	}

	return r
}
