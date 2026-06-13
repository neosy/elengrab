package mappers

import (
	"strings"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/common/composition/pages"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
	dauth "github.com/neosy/elengrab/internal/domain/auth"
)

func (m *Mappers) UserToUserOnPage(user *dauth.User, icons *pages.AdminUserIcons) pages.AdminUser {
	var email string
	if user.Email != nil {
		email = *user.Email
	}

	return pages.AdminUser{
		UserID: user.UserID.String(),
		Login:  user.Login.String(),
		Name:   user.Login.String(),
		Email:  email,
		Status: user.Status().Title(),
		Roles:  strings.Join(user.RoleIDs, ", "),

		LoadTableRowPath: httppaths.BuildAdminUserTableRowPath(user.UserID),

		Icons: icons,
	}
}
