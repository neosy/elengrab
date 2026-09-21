package types

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type MediaQueryCursor dtypes.QueryMediaDownloadCursor

func (cursor MediaQueryCursor) IsZero() bool {
	return cursor.ID == uuid.Nil &&
		cursor.CreatedAt.IsZero() &&
		cursor.Views == 0
}

func (cursor MediaQueryCursor) Encode() string {
	var lastID string
	if cursor.ID != uuid.Nil {
		lastID = cursor.ID.String()
	}

	var lastCreateAt string
	if !cursor.CreatedAt.IsZero() {
		lastCreateAt = strconv.FormatInt(cursor.CreatedAt.UTC().UnixMilli(), 10)
	}

	values := []string{
		lastID,
		lastCreateAt,
		strconv.Itoa(int(cursor.Views)),
	}

	valuesString := strings.Join(values, ",")

	return valuesString
}

func DecodeMediaQueryCursor(value string) (*MediaQueryCursor, error) {
	values := strings.Split(value, ",")
	if len(values) != 3 {
		return nil, fmt.Errorf("invalid cursor format: expected 4 values, got %d", len(values))
	}

	var (
		lastID       uuid.UUID
		lastCreateAt time.Time
		lastViews    uint32
		err          error
	)

	if value := values[0]; value != "" {
		lastID, err = uuid.Parse(value)
		if err != nil {
			return nil, err
		}
	}

	if value := values[1]; value != "" {
		timestamp, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}
		lastCreateAt = time.UnixMilli(timestamp).UTC()
	}

	if value := values[2]; value != "" {
		views, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		lastViews = uint32(views)
	}

	return &MediaQueryCursor{
		ID:        lastID,
		CreatedAt: lastCreateAt,
		Views:     lastViews,
	}, nil
}

func (cursor MediaQueryCursor) DomainCursor() dtypes.QueryMediaDownloadCursor {
	return dtypes.QueryMediaDownloadCursor(cursor)
}
