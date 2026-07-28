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

	login := dtypes.NewLogin(user.Login.String()).String()

	return &eauth.User{
		UserID:       user.UserID,
		Login:        login,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		IsActive:     isActive,
	}, nil
}

func (m *Mappers) MapUserEntityToDomain(user *eauth.User, rolesCSV string) (*dauth.User, error) {
	var isActive bool = false
	if user.IsActive == 1 {
		isActive = true
	}

	roleIds := strings.Split(rolesCSV, ",")

	var roleIDs = make([]string, 0, len(roleIds))
	for _, roleID := range roleIds {
		roleIDs = append(roleIDs, roleID)
	}

	return &dauth.User{
		UserID:            user.UserID,
		Login:             dtypes.NewLogin(user.Login),
		Email:             user.Email,
		PasswordHash:      user.PasswordHash,
		PasswordUpdatedAt: user.PasswordUpdatedAt,
		IsActive:          isActive,
		RoleIDs:           roleIDs,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
		DeletedAt:         user.DeletedAt,
	}, nil
}
