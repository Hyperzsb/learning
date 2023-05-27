package model

import "time"

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
