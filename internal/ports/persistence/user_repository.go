package persistence

import (
	"context"

	"github.com/google/uuid"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

type UserRepositoryFactory func() UserRepository

type UserRepository interface {
	Transactional

	Insert(ctx context.Context, user *dauth.User) error
	Update(ctx context.Context, user *dauth.User) error
	Delete(ctx context.Context, userID uuid.UUID, soft bool) error

	UpdatePassword(ctx context.Context, userID uuid.UUID, newPasswHash string) error

	FindByUserID(ctx context.Context, userID uuid.UUID) (*dauth.User, error)
	FindByLogin(ctx context.Context, login dtypes.Login) (*dauth.User, error)
	FindByEmail(ctx context.Context, email string) (*dauth.User, error)

	ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error)
	ExistsByLogin(ctx context.Context, login dtypes.Login) (bool, error)

	IterateGetAll(ctx context.Context, fn func(*dauth.User) error) error

	WithFilters(filtersByName map[string]any) UserRepository
	WithoutDeleted() UserRepository
	WithoutGuest() UserRepository
}
