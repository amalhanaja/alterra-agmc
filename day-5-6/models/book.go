package models

import "time"

type Book struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Isbn      string    `json:"isbn"`
	Writer    string    `json:"writer"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uint      `json:"-"`
}

type CreateBookPayload struct {
	Title  string `json:"title" validate:"required"`
	Isbn   string `json:"isbn" validate:"required"`
	Writer string `json:"writer" validate:"required"`
}

type UpdateBookPayload struct {
	Title  string `json:"title,omitempty" validate:"omitempty"`
	Isbn   string `json:"isbn,omitempty" validate:"omitempty"`
	Writer string `json:"writer,omitempty" validate:"omitempty"`
}
