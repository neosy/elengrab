package ucdownloader

import (
	"context"

	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/app/usecases/dto"
)

func (uc *YouTubeDownloader) UpdateFileInfo(ctx context.Context, fileId uuid.UUID, patch *dto.FileInfoPatch) error {
	return uc.file.Patch(ctx, fileId, patch)
}
