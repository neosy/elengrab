package searchindex

import (
	"fmt"
	"strings"
	"uuid"

	"github.com/Masterminds/squirrel"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
	esearchindex "github.com/neosy/elengrab/internal/repository/sqlite/search_index/entity"
	"github.com/neosy/elengrab/internal/repository/sqlite/sqlutil"
)

func (r *MediaSourceIndexRepository) applyContentTypeFilter(conditions []squirrel.Sqlizer) []squirrel.Sqlizer {
	var eIndex esearchindex.MediaSourceIndex

	switch r.queryOptions.ContentType {
	case dtypes.QueryMediaContentTypeAll:
		conditions = append(conditions, squirrel.Eq{eIndex.FieldName(&eIndex.IsReady): true})
	case dtypes.QueryMediaContentTypeVideo:
		conditions = append(conditions, squirrel.Eq{eIndex.FieldName(&eIndex.IsShort): false})
		conditions = append(conditions, squirrel.Eq{eIndex.FieldName(&eIndex.IsAudio): false})
		conditions = append(conditions, squirrel.Eq{eIndex.FieldName(&eIndex.IsReady): true})
	case dtypes.QueryMediaContentTypeShorts:
		conditions = append(conditions, squirrel.Eq{eIndex.FieldName(&eIndex.IsShort): true})
		conditions = append(conditions, squirrel.Eq{eIndex.FieldName(&eIndex.IsReady): true})
	case dtypes.QueryMediaContentTypeAudio:
		conditions = append(conditions, squirrel.Eq{eIndex.FieldName(&eIndex.IsAudio): true})
		conditions = append(conditions, squirrel.Eq{eIndex.FieldName(&eIndex.IsReady): true})
	case dtypes.QueryMediaContentTypeUploads:
		conditions = append(conditions, squirrel.Eq{eIndex.FieldName(&eIndex.SourceType): dtypes.MediaDownloadSourceTypeUpload.String()})
		conditions = append(conditions, squirrel.Eq{eIndex.FieldName(&eIndex.IsReady): true})
	case dtypes.QueryMediaContentTypeInProgress:
		conditions = append(conditions, squirrel.Eq{eIndex.FieldName(&eIndex.IsReady): false})
	}

	return conditions
}

func (r *MediaSourceIndexRepository) applySearchFilter(conditions []squirrel.Sqlizer) []squirrel.Sqlizer {
	var eIndex esearchindex.MediaSourceIndex

	if text := uptr.Deref(r.queryOptions.SearchText); strings.TrimSpace(text) != "" {
		text = strings.ToLower(text)
		conditions = append(conditions, squirrel.Or{
			sqlutil.Like(eIndex.FieldName(&eIndex.TitleLower), text),
			sqlutil.Like(eIndex.FieldName(&eIndex.DescriptionLower), text),
		})
	}

	return conditions
}

func (r *MediaSourceIndexRepository) applyDeletedFilter(conditions []squirrel.Sqlizer) []squirrel.Sqlizer {
	if !r.queryOptions.includeDeleted {
		var eIndex esearchindex.MediaSourceIndex
		conditions = append(conditions, squirrel.Eq{eIndex.FieldName(&eIndex.DeletedAt): nil})
	}

	return conditions
}

// applyInProgressOwnerFilter limits in-progress downloads to the specified owner.
// Anonymous users can only see ready downloads.
func (r *MediaSourceIndexRepository) applyInProgressOwnerFilter(conditions []squirrel.Sqlizer) []squirrel.Sqlizer {
	if r.queryOptions.OwnerIDForInProgress != nil {
		var eIndex esearchindex.MediaSourceIndex

		if *r.queryOptions.OwnerIDForInProgress != dauth.AnonymousUserID() {
			conditions = append(conditions, squirrel.Or{
				squirrel.Eq{eIndex.FieldName(&eIndex.IsReady): true},
				squirrel.Eq{eIndex.FieldName(&eIndex.UserID): r.queryOptions.OwnerIDForInProgress.String()},
			})
		} else {
			conditions = append(conditions, squirrel.Eq{eIndex.FieldName(&eIndex.IsReady): true})
		}
	}

	return conditions
}

func (r *MediaSourceIndexRepository) applyVisibilityAndCheckUserFilter(
	conditions []squirrel.Sqlizer,
	filterUserID string,
) ([]squirrel.Sqlizer, bool) {
	var (
		eIndex        esearchindex.MediaSourceIndex
		hasUserFilter bool
	)

	if r.queryOptions.Visibility != nil && *r.queryOptions.Visibility > dtypes.QueryMediaVisibilityAll {
		queryVisibility := *r.queryOptions.Visibility

		var sqlOr squirrel.Or
		if queryVisibility > dtypes.QueryMediaVisibilityAll {
			sqlOr = squirrel.Or{
				squirrel.Eq{eIndex.FieldName(&eIndex.UserID): nil},
				squirrel.Eq{eIndex.FieldName(&eIndex.UserID): uuid.Nil().String()},
				squirrel.Eq{eIndex.FieldName(&eIndex.Visibility): dtypes.MediaVisibilityPublic.String()},
			}
		}

		if queryVisibility == dtypes.QueryMediaVisibilityAuthenticated {
			if !r.queryOptions.IsGuestRequest {
				sqlOr = append(sqlOr, squirrel.Eq{eIndex.FieldName(&eIndex.Visibility): dtypes.MediaVisibilityAuthenticated.String()})
			}
			if filterUserID != "" {
				sqlOr = append(sqlOr, squirrel.Eq{eIndex.FieldName(&eIndex.UserID): filterUserID})
				hasUserFilter = true
			}
		}

		if len(sqlOr) > 0 {
			conditions = append(conditions, sqlOr)
		}
	}

	return conditions, hasUserFilter
}

func (r *MediaSourceIndexRepository) applyVisibilityFilter(conditions []squirrel.Sqlizer, filterUserID string) []squirrel.Sqlizer {
	conditions, _ = r.applyVisibilityAndCheckUserFilter(conditions, filterUserID)
	return conditions
}

func (r *MediaSourceIndexRepository) applyFiltersWithoutUserID(conditions []squirrel.Sqlizer) ([]squirrel.Sqlizer, string, error) {
	var (
		eIndex       esearchindex.MediaSourceIndex
		filterUserID string
	)

	for name, filter := range r.queryOptions.Filters {
		switch name {
		case "":
			continue
		case eIndex.FieldName(&eIndex.TitleLower):
			title, ok := filter.Condition().Value().(string)
			if !ok {
				return nil, "", fmt.Errorf("expected string, got %T (%v)", filter.Condition().Value(), filter.Condition().Value())
			}
			conditions = append(conditions, sqlutil.Like(eIndex.FieldName(&eIndex.TitleLower), title))
		case eIndex.FieldName(&eIndex.UserID):
			userID, ok := filter.Condition().Value().(uuid.UUID)
			if !ok {
				return nil, "", fmt.Errorf("expected uuid.UUID, got %T", filter.Condition().Value())
			}
			filterUserID = userID.String()
		default:
			conditions = append(conditions, filter.PrepareForSql().SqlCondition())
		}
	}

	return conditions, filterUserID, nil
}

func (r *MediaSourceIndexRepository) applyFilters(conditions []squirrel.Sqlizer) ([]squirrel.Sqlizer, error) {
	conditions, filterUserID, err := r.applyFiltersWithoutUserID(conditions)
	if err != nil {
		return nil, err
	}

	if filterUserID != "" {
		var eIndex esearchindex.MediaSourceIndex
		conditions = append(conditions, squirrel.Eq{eIndex.FieldName(&eIndex.UserID): filterUserID})
	}

	return conditions, nil
}
