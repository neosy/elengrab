package file

import (
	"context"

	"github.com/google/uuid"
)

func (uc *File) Delete(ctx context.Context, fileId uuid.UUID) error {
	return uc.fileRep.Delete(ctx, fileId)
}
