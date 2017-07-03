package DB

import (
	"github.com/boltdb/bolt"
	"time"
	"log"
	"fmt"
)


const dbName string = "requestHistory.db"
const dbBucket string = "request"

func Save(enc []byte){

	db, err := bolt.Open(dbName,0600,nil)
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
func ShowDb(){
	db, err := bolt.Open(dbName,0600,nil)
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

func DeleteDb(){
	db, err := bolt.Open(dbName,0600,nil)
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
