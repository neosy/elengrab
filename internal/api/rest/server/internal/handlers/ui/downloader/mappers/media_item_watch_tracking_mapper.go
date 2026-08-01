package mappers

import (
	"time"

	"github.com/google/uuid"
	dto "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (m *Mappers) MapMediaItemWatchTrackingRequestToUsecase(
	authCtx dauth.AuthContext,
	downloadID uuid.UUID,
	req dto.MediaItemWatchTrackingRequest,
) (ucdto.TrackMediaWatchEventRequest, error) {
	var userID *uuid.UUID
	if authCtx.UserID != uuid.Nil {
		userID = &authCtx.UserID
	}

	var sessionID *uuid.UUID
	if userID == nil && authCtx.AnonSessionID != uuid.Nil {
		sessionID = &authCtx.AnonSessionID
	}

	eventType, err := dtypes.ParseMediaWatchEventType(req.EventType)
	if err != nil {
		return ucdto.TrackMediaWatchEventRequest{}, err
	}

	return ucdto.TrackMediaWatchEventRequest{
		DownloadID: downloadID,
		UserID:     userID,
		SessionID:  sessionID,
		Position:   time.Duration(req.PositionMs) * time.Millisecond,
		Interval:   time.Duration(req.IntervalMs) * time.Millisecond,
		EventType:  eventType,
	}, nil
}
