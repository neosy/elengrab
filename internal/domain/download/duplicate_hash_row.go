package ddownload

import "uuid"

type DuplicateHashRow struct {
	Hash   string
	UserID *uuid.UUID
}
