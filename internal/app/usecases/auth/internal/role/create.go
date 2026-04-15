package authrole

import (
	"context"

	apperrors "github.com/neosy/elengrab/internal/app/errors"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *Role) Create(ctx context.Context, role *dauth.Role) (string, error) {
	if role == nil {
		uc.logger.Warn("Nil pointer in function")
		return "", apperrors.ErrFuncParamNullPointer
	}

	if role.RoleID == "" {
		return "", errorx.New("empty roleID field", exceptionx.ERROR)
	}

	err := uc.roleRep.Insert(ctx, role)
	if err != nil {
		uc.logger.Warn(
			"Failed to insert record into repository",
			"error", err,
		)
		return "", errorx.NewFromError(err, exceptionx.ERROR)
	}

	return role.RoleID, nil
}
