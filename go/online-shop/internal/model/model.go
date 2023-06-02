package model

import (
	"database/sql"
	"onlineshop/internal/driver"
)

type Model struct {
	db *sql.DB
}

func New(dsn string) (*Model, error) {
	model := &Model{
		db: nil,
	}

	db, err := driver.NewDB(dsn)
	if err != nil {
		return model, err
	}

	model.db = db

	return model, nil
}

func (m *Model) DB() *sql.DB {
	return m.db
}

func (m *Model) Close() error {
	if m.db == nil {
		return nil
	}

	err := m.db.Close()
	if err != nil {
		return err
	}

	return nil
}
