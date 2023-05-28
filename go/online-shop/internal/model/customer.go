package model

import (
	"context"
	"database/sql"
	"time"
)

type Customer struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	CreateTime time.Time `json:"-"`
	UpdateTime time.Time `json:"-"`
}

func (m *Model) CreateCustomer(c Customer) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	statement := `
		insert into customers (name, email) 
		values (?, ?)
	`
	result, err := m.db.ExecContext(ctx, statement,
		c.Name,
		c.Email,
	)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *Model) GetCustomer(id int) (Customer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	statement := `
		select id, name, email
		from customers
		where id = ?
	`
	row := m.db.QueryRowContext(ctx, statement, id)

	customer := Customer{}
	err := row.Scan(
		&customer.ID,
		&customer.Name,
		&customer.Email,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return customer, &EmptyQueryError{err.Error()}
		}

		return customer, err
	}

	return customer, nil
}
