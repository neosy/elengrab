package dbexec

import (
	"database/sql"
)

type dbEntry interface {
	DBName() string
	DB() *sql.DB
	Locker() WriteLocker
}
