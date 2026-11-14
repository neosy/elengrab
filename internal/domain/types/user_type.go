package dtypes

type UserType uint8

const (
	UserTypeAnonymous UserType = iota
	UserTypeGuest
	UserTypeUser
	UserTypeAdmin
)

// HasRequiredAccess reports whether the user type meets the required access level.
func (ut UserType) HasRequiredAccess(required UserType) bool {
	return ut >= required
}

// ResolveUserTypeByRoles determines the user type based on the provided role IDs.
func ResolveUserTypeByRoles(roleIDs UserRoleIDs) UserType {
	if len(roleIDs) == 0 {
		return UserTypeAnonymous
	}

	uType := UserTypeAnonymous

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
