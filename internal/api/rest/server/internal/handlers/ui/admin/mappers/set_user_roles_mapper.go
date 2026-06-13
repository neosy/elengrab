package mappers

import (
	"github.com/google/uuid"
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin/dto"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/admin/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

func (m *Mappers) MapSetUserRolesRequestToUsecase(req dto.SetUserRolesRequest) (ucdto.SetUserRolesRequest, error) {
	roleIDs, err := dtypes.ParseUserRoleIDs(req.RoleIDs)
	if err != nil {
		return ucdto.SetUserRolesRequest{}, err
	}

	return ucdto.SetUserRolesRequest{
		UserID:  uuid.MustParse(req.UserID),
		RoleIDs: roleIDs,
	}, nil
}
