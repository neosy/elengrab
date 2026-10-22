package dauth

import "uuid"

var AnonymousUserID = func() uuid.UUID { return uuid.Nil() }
