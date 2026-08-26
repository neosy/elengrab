package mappers

import (
	"net/http"

	"github.com/google/uuid"
	dto "github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/errorx"
)

func (m *Mappers) MapPatchMediaRequestToUsecase(
	downloadID uuid.UUID,
	req dto.PatchMediaByDownloadIDRequest,
) (ucdto.PatchMediaDownloadRequest, error) {
	var description *string

	if req.Description != "" {
		description = &req.Description
	}

	visibility, err := dtypes.ParseMediaVisibility(req.Visibility)
	if err != nil {
		return ucdto.PatchMediaDownloadRequest{}, errorx.NewHTTPMessage("Invalid visibility value", http.StatusBadRequest)
	}

	return ucdto.PatchMediaDownloadRequest{
		DownloadID: downloadID,

		MediaTitle:       &req.Title,
		MediaDescription: &description,
		Visibility:       &visibility,
	}, nil
}
