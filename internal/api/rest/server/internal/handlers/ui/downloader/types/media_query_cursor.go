package types

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type MediaQueryCursor struct {
	ViewMode     dtypes.QueryMediaViewMode
	LastID       uuid.UUID
	LastCreateAt time.Time
	LastViews    uint32
}

func (cursor *MediaQueryCursor) Encode() string {
	var lastID string
	if cursor.LastID != uuid.Nil {
		lastID = cursor.LastID.String()
	}

	var lastCreateAt string
	if !cursor.LastCreateAt.IsZero() {
		lastCreateAt = strconv.FormatInt(cursor.LastCreateAt.UTC().UnixMilli(), 10)
	}

	values := []string{
		cursor.ViewMode.String(),
		lastID,
		lastCreateAt,
		strconv.Itoa(int(cursor.LastViews)),
	}

	valuesString := strings.Join(values, ",")

	return valuesString
}

func DecodeMediaQueryCursor(value string) (*MediaQueryCursor, error) {
	values := strings.Split(value, ",")
	if len(values) != 4 {
		return nil, fmt.Errorf("invalid cursor format: expected 4 values, got %d", len(values))
	}

	var (
		viewMode     = dtypes.QueryMediaViewModeDefault
		lastID       uuid.UUID
		lastCreateAt time.Time
		lastViews    uint32
		err          error
	)

	if values[0] != "" {
		viewMode, err = dtypes.ParseQueryMediaViewMode(values[0])
		if err != nil {
			return nil, err
		}
	}

	if values[1] != "" {
		lastID, err = uuid.Parse(values[1])
		if err != nil {
			return nil, err
		}
	}

	if values[2] != "" {
		timestamp, err := strconv.ParseInt(values[2], 10, 64)
		if err != nil {
			return nil, err
		}
		lastCreateAt = time.UnixMilli(timestamp).UTC()
	}

	if values[3] != "" {
		views, err := strconv.Atoi(values[3])
		if err != nil {
			return nil, err
		}
		lastViews = uint32(views)
	}

	return &MediaQueryCursor{
		ViewMode:     viewMode,
		LastID:       lastID,
		LastCreateAt: lastCreateAt,
		LastViews:    lastViews,
	}, nil
}
