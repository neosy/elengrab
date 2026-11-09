package mappers

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	udto "github.com/neosy/elengrab/internal/app/usecases/dto"
	ddownload "github.com/neosy/elengrab/internal/domain/download"
	ierrors "github.com/neosy/elengrab/internal/errors"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
)

func (m *Mappers) MapQueueMediaDownloadListRequestToScheduleDownloadListRequest(
	req dto.QueueMediaDownloadListRequest,
) (*udto.ScheduleDownloadListRequest, error) {
	resultItems := make([]udto.ScheduleDownloadListRequestItem, 0, len(req.Items))

	for _, item := range req.Items {
		resultItems = append(resultItems,
			udto.ScheduleDownloadListRequestItem{
				URL:   item.URL,
				Title: item.Title,
			},
		)
	}

	return &udto.ScheduleDownloadListRequest{
		Items: resultItems,
		Options: ddownload.DownloadOptions{
			FormatType:      m.MapFormatType(req.Options.QualityCodec, req.Options.Format),
			VideoFormat:     m.MapVideoFormat(req.Options.QualityCodec, req.Options.Format),
			VideoCodec:      m.MapVideoCodec(req.Options.QualityCodec),
			VideoResolution: m.MapVideoResolution(req.Options.QualityResolution),
			AudioFormat:     m.MapAudioFormat(req.Options.Format),
		},
	}, nil
}

func (m *Mappers) MapScheduleDownloadListResponseToQueueMediaDownloadListResponse(
	respIn *udto.ScheduleDownloadListResponse,
) (*dto.QueueMediaDownloadListResponse, error) {
	if respIn == nil {
		return nil, ierrors.ErrFuncParamNullPointer
	}

	resultItems := make([]dto.QueueMediaDownloadListResponseItem, 0, len(respIn.Items))

	for _, item := range respIn.Items {
		resultItems = append(resultItems,
			dto.QueueMediaDownloadListResponseItem{
				DownloadID: idcodec.EncodeUUIDBase64URL(item.DownloadID),
				URL:        item.URL,
			},
		)
	}

	return &dto.QueueMediaDownloadListResponse{
		Items: resultItems,
		Count: len(resultItems),
	}, nil
}
