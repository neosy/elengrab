package download

import (
	"github.com/google/uuid"
	uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"
)

type fileRepositoryFilters struct {
	userID *uuid.UUID
	title  *string
}

func (f *fileRepositoryFilters) copy() fileRepositoryFilters {
	return fileRepositoryFilters{
		userID: uptr.Copy(f.userID),
		title:  uptr.Copy(f.title),
	}
}
