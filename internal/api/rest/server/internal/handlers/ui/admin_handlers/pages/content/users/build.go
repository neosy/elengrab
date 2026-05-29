package users

import (
	"context"

	"github.com/neosy/elengrab/internal/api/rest/server/internal/handlers/ui/composition/pages"
	httppaths "github.com/neosy/elengrab/internal/api/rest/server/internal/paths"
)

func (c *Content) Build(ctx context.Context) (any, error) {
	users, err := c.admin.GetAllUsersWithoutGuest(ctx)
	if err != nil {
		return nil, err
	}

	usersData := make([]pages.AdminUser, 0, len(users))
	userIcons := pages.NewAdminUserIcons()

	for _, user := range users {
		usersData = append(usersData, c.mappers.UserToUserOnPage(user, userIcons))
	}

	return pages.AdminContentUsersData{
		AdminContentData:    pages.NewAdminContentData(c.title),
		UserDetailQueryPath: httppaths.BuildAdminUserDetailTemplatePath(),
		Users:               usersData,
	}, nil
}
