package mappers

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (m *Mappers) MapQueueUploadedFileRequestToScheduleUpload(req dto.QueueUploadedFileRequest) (udto.ScheduleUploadRequest, error) {
	visibility, err := dtypes.ParseMediaVisibility(req.Visibility)
	if err != nil {
		return udto.ScheduleUploadRequest{}, err
	}

	return udto.ScheduleUploadRequest{
		Title:       req.Title,
		Description: req.Description,
		Visibility:  visibility,
	}, nil
}
