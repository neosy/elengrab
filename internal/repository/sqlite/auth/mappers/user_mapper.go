package mappers

import (
	"strings"

	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
	eauth "github.com/neosy/elengrab/internal/repository/sqlite/auth/entity"
)

func (m *Mappers) MapUserDomainToEntity(user *dauth.User) (*eauth.User, error) {
	var isActive int = 0
	if user.IsActive {
		isActive = 1
	}

	return &eauth.User{
		UserID:   user.UserID,
		Login:    user.Login,
		Email:    user.Email,
		IsActive: isActive,
	}, nil
}

func (m *Mappers) MapUserEntityToDomain(user *eauth.User, rolesCSV string) (*dauth.User, error) {
	var isActive bool = false
	if user.IsActive == 1 {
		isActive = true
	}

	roleIds := strings.Split(rolesCSV, ",")

	var roles = make([]dtypes.UserRole, 0, len(roleIds))
	for _, roleID := range roleIds {
		role, err := dtypes.ParseUserRole(roleID)
		if err != nil {
			continue
		}
		roles = append(roles, role)
	}

	return &dauth.User{
		UserID:    user.UserID,
		Login:     user.Login,
		Email:     user.Email,
		IsActive:  isActive,
		Roles:     roles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
