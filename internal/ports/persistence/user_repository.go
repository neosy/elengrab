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
	Delete(ctx context.Context, userID uuid.UUID, soft bool) error

	UpdatePassword(ctx context.Context, userID uuid.UUID, newPasswHash string) error

	FindByUserID(ctx context.Context, userID uuid.UUID) (*dauth.User, error)
	FindByLogin(ctx context.Context, login string) (*dauth.User, error)
	FindByEmail(ctx context.Context, email string) (*dauth.User, error)

	ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error)
	ExistsByLogin(ctx context.Context, login string) (bool, error)
}
