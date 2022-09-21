package models

import "time"

type Book struct {
	ID        uint
	Title     string
	Isbn      string
	Writer    string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint
}
