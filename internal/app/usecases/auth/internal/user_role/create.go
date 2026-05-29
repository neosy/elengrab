package authuserrole

import (
	"context"
	"strings"

	"github.com/google/uuid"
	apperrors "github.com/neosy/elengrab/internal/app/errors"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
	"github.com/neosy/elengrab/internal/pkg/errorx"
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
)

func (uc *UserRole) Create(ctx context.Context, userRole *dauth.UserRole) error {
	if userRole == nil {
		uc.logger.Warn("Nil pointer in function")
		return apperrors.ErrFuncParamNullPointer
	}

	userRole.RoleID = strings.TrimSpace(userRole.RoleID)

	if userRole.UserID == uuid.Nil || userRole.RoleID == "" {
		return errorx.New("empty userID or roleID fields", exceptionx.ERROR)
	}

	err := uc.userRoleRep.Insert(ctx, userRole)
	if err != nil {
		uc.logger.Warn(
			"Failed to insert record into repository",
			"error", err,
		)
		return errorx.NewFromError(err, exceptionx.ERROR)
	}

	return nil
}
