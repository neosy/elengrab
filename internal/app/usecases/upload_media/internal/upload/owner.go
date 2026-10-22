package upload

import "uuid"

// OwnerID returns the unique identifier of the user who uploaded the file.
func (u *Upload) OwnerID() uuid.UUID {
	return u.userID
}
