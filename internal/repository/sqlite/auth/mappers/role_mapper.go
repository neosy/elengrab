package mappers

import (
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	eauth "github.com/neosy/elengrab/internal/repository/sqlite/auth/entity"
)

func (m *Mappers) MapRoleDomainToEntity(role *dauth.Role) (*eauth.Role, error) {
	return &eauth.Role{
		RoleID:      role.RoleID,
		Name:        role.Name,
		Description: role.Description,
	}, nil
}

func (m *Mappers) MapRoleEntityToDomain(role *eauth.Role) (*dauth.Role, error) {
	return &dauth.Role{
		RoleID:      role.RoleID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}, nil
}
