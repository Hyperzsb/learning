package db

import (
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"gophercises/taskmanager/task"
	"time"
)

const (
	dbFilename = ".data/task.db"
	bucketName = "TaskBucket"
)

// AddTask adds a new task to db, with an auto-increasing Task.ID
func AddTask(t task.Task) error {
	db, err := bolt.Open(dbFilename, 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return nil
		}

		buf := b.Get([]byte(t.Name))
		if buf != nil {
			var de task.DuplicateErr
			de.Name = t.Name
			return de
		}

		buf, err := json.Marshal(t)
		if err != nil {
			return err
		}

		err = b.Put([]byte(t.Name), buf)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func RemoveTask(name string) error {
	db, err := bolt.Open(dbFilename, 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return nil
		}

		buf := b.Get([]byte(name))
		if buf == nil {
			var nte task.NotFoundErr
			nte.Name = name
			return nte
		}

		if err := b.Delete([]byte(name)); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func UpdateTask(name string, status task.Status) error {
	db, err := bolt.Open(dbFilename, 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return nil
		}

		buf := b.Get([]byte(name))
		if buf == nil {
			var nte task.NotFoundErr
			nte.Name = name
			return nte
		}

		var t task.Task
		err = json.Unmarshal(buf, &t)
		if err != nil {
			return err
		}

		t.Status = status

		buf, err = json.Marshal(t)
		if err != nil {
			return err
		}

		err = b.Put([]byte(name), buf)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// ListTask lists all validate tasks with different Task.Status
func ListTask() ([]task.Task, error) {
	db, err := bolt.Open(dbFilename, 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var tasks []task.Task
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return nil
		}
		c := b.Cursor()

		var t task.Task
		for key, val := c.First(); key != nil; key, val = c.Next() {
			err := json.Unmarshal(val, &t)
			if err != nil {
				return err
			}

			tasks = append(tasks, t)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
