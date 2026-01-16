package lock

type WriteLocker interface {
	Lock()
	Unlock()
}

type RWLocker interface {
	RLock()
	RUnlock()
	Lock()
	Unlock()
}
