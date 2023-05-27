package model

import "time"

type Customer struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	CreateTime time.Time `json:"-"`
	UpdateTime time.Time `json:"-"`
}
