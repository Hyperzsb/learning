package model

import (
	"context"
	"database/sql"
	"time"
)

type TransactionStatus struct {
	ID         int       `json:"id"`
	Status     string    `json:"status"`
	CreateTime time.Time `json:"-"`
	UpdateTime time.Time `json:"-"`
}

type Transaction struct {
	ID            int       `json:"id"`
	PaymentIntent string    `json:"payment_intent"`
	PaymentMethod string    `json:"payment_method"`
	Currency      string    `json:"currency"`
	Amount        int       `json:"amount"`
	Card          string    `json:"card"`
	BankCode      string    `json:"bank_code"`
	StatusID      int       `json:"status_id"`
	CreateTime    time.Time `json:"-"`
	UpdateTime    time.Time `json:"-"`
}

func (m *Model) CreateTransaction(t Transaction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	statement := `
		insert into transactions (payment_intent, payment_method, currency, amount, card, bank_code, status_id)
		values (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := m.db.ExecContext(ctx, statement,
		t.PaymentIntent,
		t.PaymentMethod,
		t.Currency,
		t.Amount,
		t.Card,
		t.BankCode,
		1,
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

func (m *Model) GetTransaction(id int) (Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	statement := `
		select id, payment_intent, payment_method, currency, amount, card, bank_code, status_id
		from transactions
		where id = ?
	`
	row := m.db.QueryRowContext(ctx, statement, id)

	transaction := Transaction{}
	err := row.Scan(
		&transaction.ID,
		&transaction.PaymentIntent,
		&transaction.PaymentMethod,
		&transaction.Currency,
		&transaction.Amount,
		&transaction.Card,
		&transaction.BankCode,
		&transaction.StatusID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return transaction, &EmptyQueryError{err.Error()}
		}

		return transaction, err
	}

	return transaction, nil
}
