package dtypes

import (
	"slices"
)

// HasRoleID checks if a single role exists in the roles slice
func HasRoleID(roleIDs []UserRoleID, checkRole UserRoleID) bool {
	return slices.Contains(roleIDs, checkRole)
}

// HasAnyRoleID checks if any of the checkRoles exist in the roles slice
func HasAnyRoleID(roleIDs []UserRoleID, checkRoleIDs ...UserRoleID) bool {
	if len(roleIDs) == 0 || len(checkRoleIDs) == 0 {
		return false
	}

	for _, r := range roleIDs {
		if slices.Contains(checkRoleIDs, r) {
			return true
		}
	}
	return false
}

// HasAllRoleIDs checks if all of the checkRoles exist in the roles slice
func HasAllRoleIDs(roleIDs []UserRoleID, checkRoleIDs ...UserRoleID) bool {
	if len(checkRoleIDs) == 0 {
		return true
	}

	if len(checkRoleIDs) > len(roleIDs) {
		return false
	}

	for _, r := range checkRoleIDs {
		if !slices.Contains(roleIDs, r) {
			return false
		}
	}
	return true
}
