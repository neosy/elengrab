package ddownload

import "github.com/google/uuid"

type DuplicateHashRow struct {
	Hash   string
	UserID *uuid.UUID
}
