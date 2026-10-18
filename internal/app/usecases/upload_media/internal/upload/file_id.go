package upload

import "github.com/google/uuid"

func newFileID() uuid.UUID {
	return uuid.New()
}

// FileID returns the unique identifier for the uploaded file.
func (u *Upload) FileID() uuid.UUID {
	return u.fileID
}
