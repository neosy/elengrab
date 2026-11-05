package mappers

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/downloader/dto"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/dto"
	ierrors "github.com/neosy/elengrab/internal/errors"
)

func (m *Mappers) MapExtractByURLResponseDomainToDto(respIn *ucdto.ExtractByURLResponse) (dto.ExtractByURLResponse, error) {
	if respIn == nil {
		return dto.ExtractByURLResponse{}, ierrors.ErrFuncParamNullPointer
	}

	items := make([]dto.ExtractByURLItem, 0, len(respIn.Items))

	for _, itemIn := range respIn.Items {
		items = append(
			items,
			dto.ExtractByURLItem{
				URL:   itemIn.URL,
				Title: itemIn.Title,
			},
		)
	}

	return dto.ExtractByURLResponse{
		Items: items,
	}, nil
}
