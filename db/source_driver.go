package database

import "github.com/golang-migrate/migrate/v4/source"

type sourceDriverWrapper func(source.Driver) (source.Driver, sourceCleanup, error)
type sourceCleanup func() error
