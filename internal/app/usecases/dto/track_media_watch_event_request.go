package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

const (
	minMediaWatchInterval      = 2 * time.Second
	maxMediaWatchInterval      = 15500 * time.Millisecond
	maxMediaWatchEndedInterval = maxMediaWatchInterval + minMediaWatchInterval
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
	if untilEnd > 0 && untilEnd < minMediaWatchInterval {
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

	if t.EventType != dtypes.MediaWatchEventTypeEnded && t.Interval < minMediaWatchInterval {
		return fmt.Errorf("interval must be at least %s", minMediaWatchInterval)
	}

	if t.EventType == dtypes.MediaWatchEventTypeEnded {
		if t.Interval > maxMediaWatchEndedInterval {
			return fmt.Errorf("ended interval should be no more than %s", maxMediaWatchEndedInterval)
		}
	} else {
		if t.Interval > maxMediaWatchInterval {
			return fmt.Errorf("interval should be no more than %s", maxMediaWatchInterval)
		}
	}

	return nil
}
