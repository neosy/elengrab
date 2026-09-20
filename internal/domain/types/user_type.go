package dtypes

type UserType uint8

const (
	UserTypeAnonymous UserType = iota
	UserTypeGuest
	UserTypeUser
	UserTypeAdmin
)

func ResolveUserTypeByRoles(roleIDs UserRoleIDs) UserType {
	if len(roleIDs) == 0 {
		return UserTypeAnonymous
	}

	var uType UserType = UserTypeAnonymous

	for _, r := range roleIDs {
		switch r {
		case UserRoleAdmin:
			return UserTypeAdmin
		case UserRoleGuest:
			if uType < UserTypeGuest {
				uType = UserTypeGuest
			}
		default:
			uType = UserTypeUser
		}
	}

	return uType
}
