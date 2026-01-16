package mappers

import (
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
)

func (m *Mappers) MapUserDomainToEntity(user *dauth.User) (*edownload.User, error) {
	var isActive int
	if user.IsActive {
		isActive = 1
	} else {
		isActive = 0
	}

	return &edownload.User{
		UserID:   user.UserID,
		Login:    user.Login,
		Email:    user.Email,
		IsActive: isActive,
	}, nil
}

func (m *Mappers) MapUserEntityToDomain(user *edownload.User) (*dauth.User, error) {
	var isActive bool
	if user.IsActive == 0 {
		isActive = false
	} else {
		isActive = true
	}

	return &dauth.User{
		UserID:    user.UserID,
		Login:     user.Login,
		Email:     user.Email,
		IsActive:  isActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
