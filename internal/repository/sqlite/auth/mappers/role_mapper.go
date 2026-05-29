package mappers

import (
	"database/sql"

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

func (m *Mappers) MapRoleRowsToDomainRoles(rows *sql.Rows, fn func(*dauth.Role) error) error {
	var eRole eauth.Role

	for rows.Next() {
		err := rows.Scan(eRole.FieldPointers()...)
		if err != nil {
			return err
		}

		role, err := m.MapRoleEntityToDomain(&eRole)
		if err != nil {
			return err
		}

		err = fn(role)
		if err != nil {
			return err
		}
	}

	return nil
}
