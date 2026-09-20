package dauth

import "github.com/google/uuid"

var AnonymousUserID = func() uuid.UUID { return uuid.Nil }
