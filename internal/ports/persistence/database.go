package persistence

type Database interface {
	Backup(path string) error
	FlushWAL() error
}
