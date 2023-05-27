package model

import (
	"context"
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

func (m *Model) GetProductById(id int) (Product, error) {
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
		return product, err
	}

	return product, nil
}
