package fileuc

import (
	"context"

	"github.com/google/uuid"
)

func (uc *File) Delete(ctx context.Context, fileId uuid.UUID) error {
	return uc.FileRep.Delete(ctx, fileId)
}
