package dtypes

import "github.com/neosy/elengrab/internal/pkg/stringx"

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
)

var (
	userStatusByStatus = map[UserStatus]struct{}{
		UserStatusActive:   {},
		UserStatusInactive: {},
	}
)

// String returns the value as a string.
func (v UserStatus) String() string {
	return string(v)
}

// String returns the value as a string.
func (v UserStatus) Title() string {
	return stringx.Capitalize(v.String())
}

// Ptr returns the pointer.
func (v UserStatus) Ptr() *UserStatus {
	return &v
}

// Exists returns true if the VideoCodec is valid.
func (v UserStatus) Exists() bool {
	_, exists := userStatusByStatus[v]
	return exists
}
