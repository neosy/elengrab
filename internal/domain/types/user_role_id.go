package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type (
	UserRoleID  string
	UserRoleIDs []UserRoleID
)

const (
	// Full access to all system features and settings
	UserRoleAdmin UserRoleID = "admin"

	// Regular user with standard access
	UserRoleUser UserRoleID = "user"

	// Limited access, read-only or temporary sessions
	UserRoleGuest UserRoleID = "guest"

	// Can view all from all users
	UserRoleViewerAll UserRoleID = "viewer_all"
)

var (
	userRoleMap = map[UserRoleID]struct{}{
		UserRoleAdmin:     {},
		UserRoleUser:      {},
		UserRoleGuest:     {},
		UserRoleViewerAll: {},
	}
	parseUserRoleMap = map[string]UserRoleID{
		"admin":      UserRoleAdmin,
		"user":       UserRoleUser,
		"guest":      UserRoleGuest,
		"viewer_all": UserRoleViewerAll,
	}
)

// String returns the value as a string.
func (v UserRoleID) String() string {
	return string(v)
}

// Ptr returns the pointer.
func (v UserRoleID) Ptr() *UserRoleID {
	return &v
}

func (v UserRoleID) Login() Login {
	return NewLogin(v.String())
}

// Exists returns true if the UserRoleID is valid.
func (v UserRoleID) Exists() bool {
	_, exists := userRoleMap[v]
	return exists
}

// ParseUserRole converting string to UserRoleID
func ParseUserRole(s string) (UserRoleID, error) {
	userRole, exists := parseUserRoleMap[strings.ToLower(s)]
	if !exists {
		return "", errors.New("invalid value for UserRoleID")
	}
	return userRole, nil
}

// ParseUserRoleIDs converting strings to UserRoleIDs
func ParseUserRoleIDs(roleIDs []string) (UserRoleIDs, error) {
	ids := make(UserRoleIDs, 0, len(roleIDs))
	for _, r := range roleIDs {
		roleID, err := ParseUserRole(r)
		if err != nil {
			return nil, err
		}
		ids = append(ids, roleID)
	}
	return ids, nil
}

// MustParseUserRole converting string to UserRoleID, ignoring any errors.
func MustParseUserRole(s string) UserRoleID {
	role, err := ParseUserRole(s)
	if err != nil {
		return UserRoleGuest
	}
	return role
}

// ValidateUserRole checks if the field value is a valid UserRoleID enum.
func ValidateUserRole(fl validator.FieldLevel) bool {
	return UserRoleID(fl.Field().String()).Exists()
}

// Strings returns a slice of string representations of the UserRoles.
func (roles UserRoleIDs) Strings() []string {
	values := make([]string, 0, len(roles))
	for _, role := range roles {
		values = append(values, role.String())
	}
	return values
}

// Join concatenates the string representations of the UserRoles with the specified separator.
func (roles UserRoleIDs) Join(sep string) string {
	return strings.Join(roles.Strings(), sep)
}
