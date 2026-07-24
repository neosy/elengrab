package mappers

import (
	"time"

	"github.com/google/uuid"
	dto "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (m *Mappers) MapMediaItemWatchTrackingRequestToUsecase(
	authCtx dauth.UserContext,
	downloadID uuid.UUID,
	req dto.MediaItemWatchTrackingRequest,
) ucdto.TrackMediaWatchEventRequest {
	var userID *uuid.UUID
	if authCtx.UserID != uuid.Nil {
		userID = &authCtx.UserID
	}

	return ucdto.TrackMediaWatchEventRequest{
		DownloadID: downloadID,
		UserID:     userID,
		Position:   time.Duration(req.PositionMs) * time.Millisecond,
		Interval:   time.Duration(req.IntervalMs) * time.Millisecond,
	}
}
