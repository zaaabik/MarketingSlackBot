package db

type Store interface {
	Save(map[string]string) error
	DeleteAll()
	Close()
}
