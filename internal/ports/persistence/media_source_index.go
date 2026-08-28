package persistence

type MediaSourceIndexRepositoryFactory func() MediaSourceIndexRepository

type MediaSourceIndexRepository interface {
	Transactional
}
