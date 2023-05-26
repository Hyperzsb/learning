package model

import "time"

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	Inventory   int       `json:"inventory"`
	CreateTime  time.Time `json:"-"`
	UpdateTime  time.Time `json:"-"`
}

type TransactionStatus struct {
	ID         int       `json:"id"`
	Status     string    `json:"status"`
	CreateTime time.Time `json:"-"`
	UpdateTime time.Time `json:"-"`
}

type Transaction struct {
	ID         int       `json:"id"`
	Currency   string    `json:"currency"`
	Amount     int       `json:"amount"`
	Card       string    `json:"card"`
	BankCode   string    `json:"bank_code"`
	StatusID   int       `json:"status_id"`
	CreateTime time.Time `json:"-"`
	UpdateTime time.Time `json:"-"`
}

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
	StatusID      int       `json:"status_id"`
	Quantity      int       `json:"quantity"`
	Amount        int       `json:"amount"`
	CreateTime    time.Time `json:"-"`
	UpdateTime    time.Time `json:"-"`
}

type User struct {
	ID         int       `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	CreateTime time.Time `json:"-"`
	UpdateTime time.Time `json:"-"`
}
