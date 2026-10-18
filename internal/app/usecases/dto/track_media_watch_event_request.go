package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/consts"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type TrackMediaWatchEventRequest struct {
	DownloadID uuid.UUID

	UserID    *uuid.UUID
	SessionID *uuid.UUID

	Position time.Duration
	Interval time.Duration

	EventType dtypes.MediaWatchEventType
}

func (t *TrackMediaWatchEventRequest) AdjustForMediaDuration(mediaDuration time.Duration) {
	if t.Position > mediaDuration {
		t.Position = mediaDuration
	}

	untilEnd := mediaDuration - t.Position
	if untilEnd > 0 && untilEnd < consts.MinMediaWatchInterval {
		t.EventType = dtypes.MediaWatchEventTypeEnded
		t.Interval += untilEnd
	}

	if t.EventType == dtypes.MediaWatchEventTypeEnded {
		t.Position = mediaDuration
	}
}

func (t *TrackMediaWatchEventRequest) Validate() error {
	if t.Interval <= 0 {
		return fmt.Errorf("interval must be positive")
	}

	if t.EventType != dtypes.MediaWatchEventTypeEnded && t.Interval < consts.MinMediaWatchInterval {
		return fmt.Errorf("interval must be at least %s", consts.MinMediaWatchInterval)
	}

	if t.EventType == dtypes.MediaWatchEventTypeEnded {
		if t.Interval > consts.MaxMediaWatchEndedInterval {
			return fmt.Errorf("ended interval should be no more than %s", consts.MaxMediaWatchEndedInterval)
		}
	} else {
		if t.Interval > consts.MaxMediaWatchInterval {
			return fmt.Errorf("interval should be no more than %s", consts.MaxMediaWatchInterval)
		}
	}

	return nil
}
