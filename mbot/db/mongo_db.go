package db

import (
	"gopkg.in/mgo.v2"
	"time"
)

const collectionName = "requests"
const dbName = "logs"

type MongoDb struct {
	session *mgo.Session
}

func NewMongoDb(addres string) (*MongoDb, error) {
	if addres == "" {
		addres = "localhost"
	}
	db, err := mgo.Dial(addres)
	if db == nil {
		return nil, err
	}

	return &MongoDb{db}, err
}

func (mongo *MongoDb) Close() {
	mongo.session.Close()
}

func (mongo *MongoDb) Save(m map[string]string) error {
	collection := mongo.session.DB(dbName).C(collectionName)
	m["time"] = time.Now().Format(time.ANSIC)
	err := collection.Insert(&m)
	if err != nil {
		return err
	}
	return err
}

func (mongo *MongoDb) DeleteAll() {
	collection := mongo.session.DB(dbName).C(collectionName)
	collection.RemoveAll(nil)
}
