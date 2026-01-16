package useruc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (uc *User) Create(ctx context.Context, user *dauth.User) (uuid.UUID, error) {
	if user == nil {
		uc.logger.Warn("Nil pointer in function")
		return uuid.Nil, errors.New("function parameter is a null pointer")
	}

	if user.UserID == uuid.Nil {
		user.UserID = uuid.New()
	}

	err := uc.userRep.Insert(ctx, user)
	if err != nil {
		uc.logger.Warn(
			"Failed to insert record into repository",
			"error", err,
		)
		return uuid.Nil, err
	}

	return user.UserID, nil
}
