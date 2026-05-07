package tablenames

const (
	Users        = "users"
	Roles        = "roles"
	UserRoles    = "user_roles"
	UserSessions = "user_sessions"
)

var tableNames = []string{
	Users,
	Roles,
	UserRoles,
	UserSessions,
}

func TableNames() []string {
	return tableNames
}
