package model

import (
	"database/sql"
	"onlineshop/internal/driver"
)

type Model struct {
	db *sql.DB
}

func New(dsn string) (*Model, error) {
	db, err := driver.NewDB(dsn)
	if err != nil {
		return nil, err
	}

	return &Model{
		db: db,
	}, err
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
