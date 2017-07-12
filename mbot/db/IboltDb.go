package db

import (
	"github.com/boltdb/bolt"
	"time"
	"log"
	"fmt"
)


const dbBucket string = "marketingClient"

type BoltDb struct {
	dbPath string
}

func NewBoltDb(path string) (*BoltDb) {
	return &BoltDb{path}
}

func (b *BoltDb)Save(enc []byte){
	db, err := bolt.Open(b.dbPath,0600,nil)
	defer db.Close()
	if err != nil{
		log.Println(err)
	}

	err = db.Update(func(tx *bolt.Tx)error {
		req, err := tx.CreateBucketIfNotExists([]byte(dbBucket))

		if err != nil {
			return err
		}

		err = req.Put([]byte(time.Now().String()),enc)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil{
		log.Println(nil)
	}

}
func (b *BoltDb)GetAll(){
	db, err := bolt.Open(b.dbPath,0600,nil)
	if err != nil{
		log.Println(err)
		return
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		if b == nil{
			log.Println("file is empty")
			return nil
		}
		err = b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})
		if err != nil{
			log.Println(err)
		}
		return nil
	})
	if err != nil{
		log.Println(err)
		return
	}
}

func (b *BoltDb)DeleteAll(){
	db, err := bolt.Open(b.dbPath,0600,nil)
	if err != nil{
		log.Println(err)
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dbBucket))
		if b == nil{
			return nil
		}
		tx.DeleteBucket([]byte(dbBucket))
		return nil
	})
}
