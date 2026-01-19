package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (u *Auth) SessionNeedsRefresh(expiresAt time.Time) bool {
	return time.Until(expiresAt) <= sessionRefreshInterval
}

func (u *Auth) RefreshSession(ctx context.Context, sessionID uuid.UUID) (time.Time, error) {
	var retExpiresAt time.Time

	refresh := func(ctx context.Context) error {
		session, err := u.userSession.GetBySessionID(ctx, sessionID)
		if err != nil {
			return err
		}

		if !u.SessionNeedsRefresh(session.ExpiresAt) {
			retExpiresAt = session.ExpiresAt
			return nil
		}

		newExpiry := time.Now().Add(sessionTTL)
		session.ExpiresAt = newExpiry
		err = u.userSession.Update(ctx, session)
		if err != nil {
			return err
		}

		return nil
	}

	err := u.userSession.Tx(ctx, refresh)
	if err != nil {
		return time.Time{}, err
	}

	return retExpiresAt, nil
}
