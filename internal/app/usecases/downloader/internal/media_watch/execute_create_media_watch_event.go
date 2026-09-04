package mediawatch

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
	"github.com/neosy/elengrab/internal/pkg/errorx"
)

func (uc *MediaWatch) ExecuteCreateMediaWatchEvent(
	ctx context.Context,
	workerID uint64,
	req *dto.CreateMediaWatchEventRequest,
) error {
	if req == nil || req.Event == nil {
		return apperrors.ErrFuncParamNullPointer
	}

	event := req.Event

	position := uc.mappers.MapUserWatchEventToWatchPosition(event, req.MediaDuration)
	if position.Validate() == nil {
		uc.userPosition.Write(ctx, position)
		uc.onWatchUserPositionUpdated(ctx, req.Event.AuthCtx(), req.Event.DownloadID)
	}

	create := func(ctx context.Context) error {
		err := uc.event.Create(ctx, event)
		if err != nil {
			return err
		}

		chunks := uc.eventToChunks(event, req.MediaDuration)
		if len(chunks) == 0 {
			return errorx.NewHTTPMessage(
				"No media watch chunks to process",
				http.StatusConflict,
			)
		}

		return uc.userChunk.AddChunkQtyBatch(ctx, chunks)
	}

	err := uc.event.Tx(ctx, create)
	if err != nil {
		return err
	}

	var userID uuid.UUID
	if event.UserID != nil {
		userID = *event.UserID
	}

	var sessionID uuid.UUID
	if userID == uuid.Nil && event.SessionID != nil {
		sessionID = *event.SessionID
	}

	uc.pendingStats.add(event.DownloadID, userID, sessionID, req.MediaDuration)

	return nil
}
