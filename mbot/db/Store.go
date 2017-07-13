package db

type Store interface {
	Save([]byte)
	GetAll()
	DeleteAll()
}
