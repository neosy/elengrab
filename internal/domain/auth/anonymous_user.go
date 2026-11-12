package dauth

import "uuid"

// AnonymousUserID returns the ID of the anonymous user.
func AnonymousUserID() uuid.UUID {
	return uuid.Nil()
}

// IsRegisteredUserID checks whether the ID belongs to a registered user.
func IsRegisteredUserID(userID uuid.UUID) bool {
	return userID != AnonymousUserID() && userID != uuid.Nil()
}

// IsNotRegisteredUserID checks whether the ID does not belong to a registered user.
func IsNotRegisteredUserID(userID uuid.UUID) bool {
	return !IsRegisteredUserID(userID)
}
