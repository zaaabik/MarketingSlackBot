package db

type Store interface {
	Save(map[string]string)
	GetAll()
	DeleteAll()
}
