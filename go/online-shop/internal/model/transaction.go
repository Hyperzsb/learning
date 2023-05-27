package model

import "time"

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
