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
func AddTask(newTask task.Task) error {
	db, err := bolt.Open(dbFilename, 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return nil
		}

		preTask := b.Get([]byte(newTask.Name))
		if preTask != nil {
			var de task.DuplicateErr
			de.Name = newTask.Name
			return &de
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}

		buf, err := json.Marshal(newTask)
		if err != nil {
			return err
		}

		err = b.Put([]byte(newTask.Name), buf)
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
