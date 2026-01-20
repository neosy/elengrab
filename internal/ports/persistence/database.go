package persistence

type DBName string

const (
	DBMainName DBName = "elengrab"
	DBAuthName DBName = "auth"
)

func (n DBName) String() string {
	return string(n)
}

type Database interface {
	Backup(dbName DBName, path string) error
	FlushWAL() error
	DBFileName(dbName DBName) string
	GetDBNames() []DBName
}
