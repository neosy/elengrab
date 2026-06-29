package mappers

import (
	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/admin/dto"
	ucdto "github.com/neosy/elengrab/internal/app/usecases/admin/dto"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	"github.com/neosy/elengrab/internal/pkg/idcodec"
)

func (m *Mappers) MapSetUserRolesRequestToUsecase(req dto.SetUserRolesRequest) (ucdto.SetUserRolesRequest, error) {
	roleIDs, err := dtypes.ParseUserRoleIDs(req.RoleIDs)
	if err != nil {
		return ucdto.SetUserRolesRequest{}, err
	}

	userID, err := idcodec.DecodeUUIDBase64URL(req.UserID)
	if err != nil {
		return ucdto.SetUserRolesRequest{}, err
	}

	return ucdto.SetUserRolesRequest{
		UserID:  userID,
		RoleIDs: roleIDs,
	}, nil
}
