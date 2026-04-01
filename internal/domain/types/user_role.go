package dtypes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type UserRole string

const (
	// Full access to all system features and settings
	UserRoleAdmin UserRole = "admin"

	// Regular user with standard access
	UserRoleUser UserRole = "user"

	// Limited access, read-only or temporary sessions
	UserRoleGuest UserRole = "guest"

	// Can view all from all users
	UserRoleViewerAll UserRole = "viewer_all"
)

var (
	userRoleMap = map[UserRole]struct{}{
		UserRoleAdmin:     {},
		UserRoleUser:      {},
		UserRoleGuest:     {},
		UserRoleViewerAll: {},
	}
	parseUserRoleMap = map[string]UserRole{
		"admin":      UserRoleAdmin,
		"user":       UserRoleUser,
		"guest":      UserRoleGuest,
		"viewer_all": UserRoleViewerAll,
	}
)

// String returns the value as a string.
func (v UserRole) String() string {
	return string(v)
}

// Ptr returns the pointer.
func (v UserRole) Ptr() *UserRole {
	return &v
}

func (v UserRole) Login() Login {
	return NewLogin(v.String())
}

// Exists returns true if the UserRole is valid.
func (v UserRole) Exists() bool {
	_, exists := userRoleMap[v]
	return exists
}

// ParseUserRole converting string to UserRole
func ParseUserRole(s string) (UserRole, error) {
	userRole, exists := parseUserRoleMap[strings.ToLower(s)]
	if !exists {
		return "", errors.New("invalid value for UserRole")
	}
	return userRole, nil
}

// MustParseUserRole converting string to UserRole, ignoring any errors.
func MustParseUserRole(s string) UserRole {
	role, err := ParseUserRole(s)
	if err != nil {
		return UserRoleGuest
	}
	return role
}

// ValidateUserRole checks if the field value is a valid UserRole enum.
func ValidateUserRole(fl validator.FieldLevel) bool {
	return UserRole(fl.Field().String()).Exists()
}
