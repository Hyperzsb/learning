package model

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

type User struct {
	ID         int       `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	CreateTime time.Time `json:"-"`
	UpdateTime time.Time `json:"-"`
}

func (m *Model) CreateUser(u User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	statement := `
		insert into users (first_name, last_name, email, password)
		values (?, ?, ?, ?)
	`
	result, err := m.db.ExecContext(ctx, statement,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Password,
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

func (m *Model) GetUser(id int) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	statement := `
		select id, first_name, last_name, email, password
		from Users
		where id = ?
	`
	row := m.db.QueryRowContext(ctx, statement, id)

	user := User{}
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, &EmptyQueryError{err.Error()}
		}

		return user, err
	}

	return user, nil
}

func (m *Model) GetUserByEmail(email string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	email = strings.ToLower(email)
	statement := `
		select id, first_name, last_name, email, password
		from Users
		where email = ?
	`
	row := m.db.QueryRowContext(ctx, statement, email)

	user := User{}
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, &EmptyQueryError{err.Error()}
		}

		return user, err
	}

	return user, nil
}
