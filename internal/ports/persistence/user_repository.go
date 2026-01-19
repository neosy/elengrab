package persistence

import (
	"context"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

type UserRepository interface {
	Transactional
	Insert(ctx context.Context, user *dauth.User) error
	Update(ctx context.Context, user *dauth.User) error
	Save(ctx context.Context, user *dauth.User) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*dauth.User, error)
	ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error)
}
