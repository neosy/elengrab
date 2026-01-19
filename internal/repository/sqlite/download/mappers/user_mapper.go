package mappers

import (
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	edownload "github.com/neosy/elengrab/internal/repository/sqlite/download/entity"
)

func (m *Mappers) MapUserDomainToEntity(user *dauth.User) (*edownload.User, error) {
	var isGuest int = 0
	if user.IsGuest {
		isGuest = 1
	}

	var isActive int = 0
	if user.IsActive {
		isActive = 1
	}

	return &edownload.User{
		UserID:   user.UserID,
		Login:    user.Login,
		Email:    user.Email,
		IsGuest:  isGuest,
		IsActive: isActive,
	}, nil
}

func (m *Mappers) MapUserEntityToDomain(user *edownload.User) (*dauth.User, error) {
	var isGuest bool = false
	if user.IsGuest == 1 {
		isGuest = true
	}

	var isActive bool = false
	if user.IsActive == 1 {
		isActive = true
	}

	return &dauth.User{
		UserID:    user.UserID,
		Login:     user.Login,
		Email:     user.Email,
		IsGuest:   isGuest,
		IsActive:  isActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
