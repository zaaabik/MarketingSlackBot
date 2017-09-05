package db

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

const dbBucket = "marketingClient"

type BoltDb struct {
	db *bolt.DB
}

func NewBoltDb(path string) (*BoltDb, error) {
	db, err := bolt.Open(path, 0600, nil)
	return &BoltDb{db}, err
}

func (b *BoltDb) Save(m map[string]string) {
	enc, err := json.Marshal(m)
	if err != nil {
		log.Print(err)
		return
	}

	err = b.db.Update(func(tx *bolt.Tx) error {
		req, err := tx.CreateBucketIfNotExists([]byte(dbBucket))

		if err != nil {
			return err
		}

		err = req.Put([]byte(time.Now().String()), enc)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Println(nil)
	}

}
func (b *BoltDb) GetAll() {
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		if b == nil {
			log.Println("file is empty")
			return nil
		}
		err := b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return
	}
}

func (b *BoltDb) DeleteAll() {
	b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		if b == nil {
			return nil
		}
		tx.DeleteBucket([]byte(dbBucket))
		return nil
	})
}

func (b *BoltDb) CloseDataBase() {
	b.db.Close()
}
