package dtypes

import (
	"slices"
)

// HasRole checks if a single role exists in the roles slice
func HasRole(roles []UserRole, checkRole UserRole) bool {
	return slices.Contains(roles, checkRole)
}

// HasAnyRole checks if any of the checkRoles exist in the roles slice
func HasAnyRole(roles []UserRole, checkRoles ...UserRole) bool {
	for _, r := range roles {
		if slices.Contains(checkRoles, r) {
			return true
		}
	}
	return false
}

// HasAllRoles checks if all of the checkRoles exist in the roles slice
func HasAllRoles(roles []UserRole, checkRoles ...UserRole) bool {
	for _, r := range checkRoles {
		if !slices.Contains(roles, r) {
			return false
		}
	}
	return true
}
