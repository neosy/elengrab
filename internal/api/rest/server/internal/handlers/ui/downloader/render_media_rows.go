package downloader

import (
	"bytes"
	"context"
	"html/template"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/components"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/paths"
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

	searchParameters := types.NewSearchParameters()
	searchParameters.AddValues(query.ViewMode.String(), queryFilters, dtypes.QueryMediaDownloadCursor{})

	channelHeaderPageData := pages.ChannelHeader{}
	if channelID := searchParameters.FindChannelID(); channelID != uuid.Nil {
		channelHeaderPageData = h.buildChannelHeaderPageData(ctx, channelID)
	}

	pageData := pages.RowsFragmentData{
		BasePaths: paths.NewHttpPaths(),
		Values: &pages.RowsFragmentValues{
			HasSearchFilters:     searchParameters.QueryFilters().Len() != 0,
			SearchParametersJSON: string(searchParameters.BuildJSON()),
			SearchParameters:     template.HTML(searchParameters.EncodeShortQueryValue()),

			ChannelJSON: string(channelHeaderPageData.JSON()),

			ResultNoRows:   rowsBuf.Len() == 0,
			ResultRowsHTML: template.HTML(rowsBuf.String()),
		},
	}

	if err := h.templates.Base.ExecuteTemplate(buf, components.ResultRowsKey, pageData); err != nil {
		return err
	}

	return nil
}
