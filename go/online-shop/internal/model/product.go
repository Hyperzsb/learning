package model

import (
	"context"
	"database/sql"
	"time"
)

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	Inventory   int       `json:"inventory"`
	CreateTime  time.Time `json:"-"`
	UpdateTime  time.Time `json:"-"`
}

func (m *Model) CreateProduct(p Product) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	statement := `
		insert into products (name, description, price, inventory) 
		values (?, ?, ?, ?) 
	`
	result, err := m.db.ExecContext(ctx, statement,
		p.Name,
		p.Description,
		p.Price,
		p.Inventory,
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

func (m *Model) GetProduct(id int) (Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	statement := `
		select id, name, description, price, inventory, create_time, update_time
		from products
		where id = ?
	`
	row := m.db.QueryRowContext(ctx, statement, id)

	product := Product{}
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Inventory,
		&product.CreateTime,
		&product.UpdateTime,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return product, &EmptyQueryError{err.Error()}
		}

		return product, err
	}

	return product, nil
}

func (m *Model) UpdateProduct(id int, p Product) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	statement := `
		update products
		set name = ?, description = ?, price = ?, inventory = ?
		where id = ?
	`
	result, err := m.db.ExecContext(ctx, statement,
		p.Name,
		p.Description,
		p.Price,
		p.Inventory,
		id,
	)
	if err != nil {
		return -1, err
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(cnt), nil
}

func (m *Model) DeleteProduct(id int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	statement := `
		delete from products
		where id = ?
	`
	result, err := m.db.ExecContext(ctx, statement, id)
	if err != nil {
		return -1, err
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(cnt), err
}
