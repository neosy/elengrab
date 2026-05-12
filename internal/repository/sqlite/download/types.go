package download

import (
	"github.com/google/uuid"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

type downloadRepositoryFilters struct {
	userID *uuid.UUID
	title  *string
}

func (f *downloadRepositoryFilters) copy() downloadRepositoryFilters {
	return downloadRepositoryFilters{
		userID: uptr.Copy(f.userID),
		title:  uptr.Copy(f.title),
	}
}
