package demo

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"time"
)

type People struct {
	ID        int
	FirstName string
	LastName  string
}

func BBoltDemo() error {
	const (
		dbFilename = ".data/demo.db"
	)

	db, err := bolt.Open(dbFilename, 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("DemoBucket"))
		if err != nil {
			return err
		}

		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		people := People{int(id), "Harry", "Potter"}

		buf, err := json.Marshal(people)
		if err != nil {
			return err
		}

		err = b.Put(itob(people.ID), buf)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DemoBucket"))
		c := b.Cursor()

		var people People
		for key, val := c.First(); key != nil; key, val = c.Next() {
			err := json.Unmarshal(val, &people)
			if err != nil {
				return err
			}

			fmt.Println(people)
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte("DemoBucket"))
		if err != nil {
			return err
		}

		return nil
	})

	return nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
