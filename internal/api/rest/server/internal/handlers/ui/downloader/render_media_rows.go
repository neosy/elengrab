package downloader

import (
	"bytes"
	"context"
	"html/template"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
	qkeys "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/query_keys.go"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/types"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (h *DownloaderHandlers) renderDownloadRows(
	ctx context.Context,
	buf *bytes.Buffer,
	authCtx dauth.AuthContext,
	query udto.MediaDownloadQuery,
) error {
	var rowsBuf bytes.Buffer
	err := h.renderDownloadItemList(ctx, &rowsBuf, authCtx, query)
	if err != nil {
		return err
	}

	queryFilters := h.mappers.MapQueryFiltersDomainToFilters(query.Filters)

	searchQueryFilters, err := h.prepareSearchQueryFilters(ctx, queryFilters)
	if err != nil {
		return err
	}

	searchParam := types.NewSearchParameters()
	searchParam.AddValues(query.ViewMode, searchQueryFilters, dtypes.QueryMediaDownloadCursor{})

	searchFilters := searchQueryFilters.FilterByKeys(qkeys.SearchFilterKeys)
	searchParamFilters := searchParam.QueryFilters()

	pageData := pages.RowsFragmentData{
		BasePaths: paths.NewHttpPaths(),
		Values: &pages.RowsFragmentValues{
			HasSearchFilters:     searchFilters.Len() != 0 || searchParamFilters.Len() != 0,
			SearchFiltersJSON:    string(searchFilters.BuildJSON()),
			SearchParametersJSON: string(searchParam.BuildJSON()),
			SearchParameters:     template.HTML(searchParam.EncodeShortQueryValue()),

			ResultNoRows:   rowsBuf.Len() == 0,
			ResultRowsHTML: template.HTML(rowsBuf.String()),
		},
	}

	if err := h.templates.Base.ExecuteTemplate(buf, components.ResultRowsKey, pageData); err != nil {
		return err
	}

	return nil
}
