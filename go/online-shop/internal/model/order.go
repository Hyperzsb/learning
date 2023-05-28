package model

import (
	"context"
	"database/sql"
	"time"
)

type OrderStatus struct {
	ID         int       `json:"id"`
	Status     string    `json:"status"`
	CreateTime time.Time `json:"-"`
	UpdateTime time.Time `json:"-"`
}

type Order struct {
	ID            int       `json:"id"`
	ProductID     int       `json:"product_id"`
	TransactionID int       `json:"transaction_id"`
	CustomerID    int       `json:"customer_id"`
	StatusID      int       `json:"status_id"`
	Quantity      int       `json:"quantity"`
	Amount        int       `json:"amount"`
	CreateTime    time.Time `json:"-"`
	UpdateTime    time.Time `json:"-"`
}

func (m *Model) CreateOrder(o Order) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	statement := `
		insert into orders (product_id, transaction_id, customer_id, status_id, quantity, amount) 
		values (?, ?, ?, ?, ?, ?)
	`
	result, err := m.db.ExecContext(ctx, statement,
		o.ProductID,
		o.TransactionID,
		o.CustomerID,
		o.StatusID,
		o.Quantity,
		o.Amount,
	)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

func (m *Model) GetOrder(id int) (Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	statement := `
		select id, product_id, transaction_id, customer_id, status_id, quantity, amount
		from orders
		where id = ?
	`
	row := m.db.QueryRowContext(ctx, statement, id)

	order := Order{}
	err := row.Scan(
		&order.ID,
		&order.ProductID,
		&order.TransactionID,
		&order.CustomerID,
		&order.StatusID,
		&order.Quantity,
		&order.Amount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return order, &EmptyQueryError{err.Error()}
		}

		return order, err
	}

	return order, nil
}
