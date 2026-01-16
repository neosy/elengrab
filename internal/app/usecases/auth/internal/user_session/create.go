package usersession

import (
	"context"
	"errors"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (uc *UserSession) Create(ctx context.Context, session *dauth.UserSession) (uuid.UUID, error) {
	if session == nil {
		uc.logger.Warn("Nil pointer in function")
		return uuid.Nil, errors.New("function parameter is a null pointer")
	}

	if session.SessionID == uuid.Nil {
		session.SessionID = uuid.New()
	}

	err := uc.userSessionRep.Insert(ctx, session)
	if err != nil {
		uc.logger.Warn(
			"Failed to insert record into repository",
			"error", err,
		)
		return uuid.Nil, err
	}

	return session.SessionID, nil
}
