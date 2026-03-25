package mappers

import (
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	eauth "github.com/neosy/elengrab/internal/repository/sqlite/auth/entity"
)

func (m *Mappers) MapUserRoleDomainToEntity(userRole *dauth.UserRole) (*eauth.UserRole, error) {
	return &eauth.UserRole{
		UserID: userRole.UserID,
		RoleID: userRole.RoleID,
	}, nil
}

func (m *Mappers) MapUserRoleEntityToDomain(userRole *eauth.UserRole) (*dauth.UserRole, error) {
	return &dauth.UserRole{
		UserID:    userRole.UserID,
		RoleID:    userRole.RoleID,
		CreatedAt: userRole.CreatedAt,
	}, nil
}
