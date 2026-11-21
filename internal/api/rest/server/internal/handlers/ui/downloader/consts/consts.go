package consts

import (
	"time"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

const (
	DateFormate = time.RFC3339
)

const (
	defaultLoadHistoryLimit    = 20
	inProgressLoadHistoryLimit = 500

	preloadHistoryOffset = 1
)

func LoadHistoryLimit(contentType dtypes.QueryMediaContentType) uint64 {
	switch contentType {
	case dtypes.QueryMediaContentTypeInProgress:
		return inProgressLoadHistoryLimit
	default:
		return defaultLoadHistoryLimit
	}
}

func PreloadHistoryAfter(contentType dtypes.QueryMediaContentType) uint64 {
	return LoadHistoryLimit(contentType) - preloadHistoryOffset
}
