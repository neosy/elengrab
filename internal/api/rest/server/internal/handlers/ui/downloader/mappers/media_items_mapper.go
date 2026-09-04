package mappers

import (
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/consts"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
)

func (m *Mappers) MapMediaItemsRequestToQuery(req dto.MediaItemsRequest) (udto.MediaDownloadQuery, error) {
	viewMode, err := dtypes.ParseQueryMediaViewMode(req.ViewMode)
	if err != nil {
		return udto.MediaDownloadQuery{}, err
	}

	var lastID uuid.UUID
	if req.LastID != "" {
		lastID, err = idcodec.DecodeUUIDBase64URL(req.LastID)
		if err != nil {
			return udto.MediaDownloadQuery{}, err
		}
	}

	var lastCreatedAt time.Time
	if req.LastCreatedAt != "" {
		lastCreatedAt, err = time.Parse(consts.DateFormate, req.LastCreatedAt)
		if err != nil {
			return udto.MediaDownloadQuery{}, err
		}
	}

	filters := make(dtypes.QueryFiltersByName)
	if req.Search != "" {
		filters.Add(dtypes.QueryFilterNameSearch, req.Search)
	}

	query := udto.MediaDownloadQuery{
		ViewMode: viewMode,
		Limit:    consts.LoadHistoryLimit,

		LastRecord: dtypes.QueryMediaDownloadCursor{
			ID:        lastID,
			CreatedAt: lastCreatedAt,
			Views:     uint32(req.LastViews),
		},

		Filters: filters,
	}

	return udto.BuildMediaDownloadQuery(query), nil
}
