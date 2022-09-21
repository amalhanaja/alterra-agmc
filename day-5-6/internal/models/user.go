package models

import "time"

type User struct {
	Name      string
	Email     string
	Password  string
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
