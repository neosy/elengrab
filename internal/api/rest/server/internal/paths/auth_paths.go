package httppaths

// UI
const (
	// Groups
	AuthGroup = "/auth"

	// Paths Auth
	RegisterPath = "/register"
	LoginPath    = "/login"
	LogoutPath   = "/logout"

	// Routes Auth
	AuthRegisterPath = AuthGroup + RegisterPath
	AuthLoginPath    = AuthGroup + LoginPath
	AuthLogoutPath   = AuthGroup + LogoutPath
)
